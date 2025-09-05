package main

import (
	"fmt"
	"math"
	"time"

	"github.com/shopspring/decimal"
)

// TradingStrategy 交易策略
type TradingStrategy struct {
	dataLoader *StockDataLoader
	config     *Config
}

// NewTradingStrategy 创建新的交易策略
func NewTradingStrategy(dataLoader *StockDataLoader, config *Config) *TradingStrategy {
	return &TradingStrategy{
		dataLoader: dataLoader,
		config:     config,
	}
}

// ExecuteStrategy 执行交易策略
func (strategy *TradingStrategy) ExecuteStrategy() ([]*MonthlyReport, error) {
	var reports []*MonthlyReport
	cash := decimal.NewFromFloat(strategy.config.InitialCapital)
	portfolio := &Portfolio{
		Cash:      cash,
		Positions: make(map[string]*Position),
		Value:     cash,
	}

	// 按月遍历交易周期
	currentDate := strategy.config.StartDate
	monthIndex := 0

	for currentDate.Before(strategy.config.EndDate) || currentDate.Equal(strategy.config.EndDate) {
		fmt.Printf("处理 %d年%d月...\n", currentDate.Year(), int(currentDate.Month()))
		
		// 处理当月交易
		report, err := strategy.processMonth(currentDate, portfolio, monthIndex)
		if err != nil {
			fmt.Printf("警告: 处理 %d年%d月 时出错: %v\n", currentDate.Year(), int(currentDate.Month()), err)
			// 继续处理下一个月
		} else {
			reports = append(reports, report)
		}
		
		// 移动到下一个月
		currentDate = time.Date(currentDate.Year(), time.Month(int(currentDate.Month())+1), 1, 0, 0, 0, 0, time.UTC)
		monthIndex++
	}

	return reports, nil
}

// processMonth 处理单个月的交易
func (strategy *TradingStrategy) processMonth(date time.Time, portfolio *Portfolio, monthIndex int) (*MonthlyReport, error) {
	
	// 加载交易信号
	signals, err := strategy.dataLoader.LoadTradeSignals(date)
	if err != nil {
		return nil, fmt.Errorf("加载交易信号失败: %v", err)
	}

	var tradingActions []TradingAction
	previousValue := portfolio.Value

	// 1. 处理剔除股票
	for _, signal := range signals {
		if signal.Status == "剔除" {
			if position, exists := portfolio.Positions[signal.Symbol]; exists {
				action, err := strategy.sellStock(signal.Symbol, position, portfolio, date)
				if err != nil {
					fmt.Printf("警告: 卖出股票 %s 失败: %v\n", signal.Symbol, err)
					continue
				}
				tradingActions = append(tradingActions, *action)
			}
		}
	}

	// 2. 计算渐进式建仓比例
	allocationRatio := strategy.calculateAllocationRatio(monthIndex)
	
	// 3. 处理纳入股票
	var stocksToBuy []*TradeSignal
	for _, signal := range signals {
		if signal.Status == "纳入" {
			// 检查是否已持有该股票
			if _, exists := portfolio.Positions[signal.Symbol]; !exists {
				stocksToBuy = append(stocksToBuy, signal)
			}
		}
	}

	// 4. 执行买入操作
	if len(stocksToBuy) > 0 {
		buyActions, err := strategy.buyStocks(stocksToBuy, portfolio, date, allocationRatio)
		if err != nil {
			fmt.Printf("警告: 买入股票失败: %v\n", err)
		} else {
			tradingActions = append(tradingActions, buyActions...)
		}
	}

	// 5. 更新投资组合价值
	err = strategy.updatePortfolioValue(portfolio, date)
	if err != nil {
		return nil, fmt.Errorf("更新投资组合价值失败: %v", err)
	}

	// 6. 计算收益率
	monthlyReturn := decimal.Zero
	cumulativeReturn := decimal.Zero
	if !previousValue.IsZero() {
		monthlyReturn = portfolio.Value.Sub(previousValue).Div(previousValue)
	}
	if strategy.config.InitialCapital > 0 {
		cumulativeReturn = portfolio.Value.Sub(decimal.NewFromFloat(strategy.config.InitialCapital)).Div(decimal.NewFromFloat(strategy.config.InitialCapital))
	}

	// 7. 生成月度报告
	report := &MonthlyReport{
		Date:             date,
		TotalValue:       portfolio.Value,
		Cash:             portfolio.Cash,
		StockValue:       portfolio.Value.Sub(portfolio.Cash),
		MonthlyReturn:    monthlyReturn,
		CumulativeReturn: cumulativeReturn,
		Positions:        copyPositions(portfolio.Positions),
		TradingActions:   tradingActions,
	}

	return report, nil
}

// calculateAllocationRatio 计算满仓建仓比例
func (strategy *TradingStrategy) calculateAllocationRatio(monthIndex int) decimal.Decimal {
	// 初始满仓策略：直接使用90%资金建仓（保留10%现金）
	maxRatio := 0.9 // 始终保持90%资金投资，10%现金
	
	return decimal.NewFromFloat(maxRatio)
}

// sellStock 卖出股票
func (strategy *TradingStrategy) sellStock(symbol string, position *Position, portfolio *Portfolio, date time.Time) (*TradingAction, error) {
	// 获取当前股价
	stockPrices, err := strategy.dataLoader.LoadStockPrice(symbol)
	if err != nil {
		return nil, fmt.Errorf("加载股价数据失败: %v", err)
	}

	// 获取当月第一个交易日
	tradingDay, err := strategy.dataLoader.GetFirstTradingDay(date.Year(), int(date.Month()), stockPrices)
	if err != nil {
		return nil, fmt.Errorf("获取交易日失败: %v", err)
	}

	tradingDayKey := tradingDay.Format("20060102")
	stockPrice, exists := stockPrices[tradingDayKey]
	if !exists {
		return nil, fmt.Errorf("未找到 %s 在 %s 的股价数据", symbol, tradingDayKey)
	}

	// 计算卖出金额
	sellAmount := stockPrice.Close.Mul(decimal.NewFromInt(int64(position.Shares)))
	
	// 更新现金和持仓
	portfolio.Cash = portfolio.Cash.Add(sellAmount)
	delete(portfolio.Positions, symbol)

	// 创建交易记录
	action := &TradingAction{
		Date:   tradingDay,
		Symbol: symbol,
		Action: "SELL",
		Shares: position.Shares,
		Price:  stockPrice.Close,
		Amount: sellAmount,
		Reason: "股票被剔除",
	}

	fmt.Printf("卖出: %s, 股数: %d, 价格: %s, 金额: %s\n", 
		symbol, position.Shares, stockPrice.Close.String(), sellAmount.String())

	return action, nil
}

// buyStocks 买入股票
func (strategy *TradingStrategy) buyStocks(stocksToBuy []*TradeSignal, portfolio *Portfolio, date time.Time, allocationRatio decimal.Decimal) ([]TradingAction, error) {
	var actions []TradingAction
	
	if len(stocksToBuy) == 0 {
		return actions, nil
	}

	// 计算可用于投资的资金
	totalValue := portfolio.Value
	availableCash := totalValue.Mul(allocationRatio)
	
	// 如果可用资金超过当前现金，使用当前现金
	if availableCash.GreaterThan(portfolio.Cash) {
		availableCash = portfolio.Cash
	}

	// 等分分配给所有要买入的股票
	cashPerStock := availableCash.Div(decimal.NewFromInt(int64(len(stocksToBuy))))
	
	fmt.Printf("可用资金: %s, 每只股票分配: %s\n", availableCash.String(), cashPerStock.String())

	for _, signal := range stocksToBuy {
		action, err := strategy.buyStock(signal.Symbol, cashPerStock, portfolio, date)
		if err != nil {
			fmt.Printf("警告: 买入股票 %s 失败: %v\n", signal.Symbol, err)
			continue
		}
		actions = append(actions, *action)
	}

	return actions, nil
}

// buyStock 买入单只股票
func (strategy *TradingStrategy) buyStock(symbol string, cashAmount decimal.Decimal, portfolio *Portfolio, date time.Time) (*TradingAction, error) {
	// 获取股价数据
	stockPrices, err := strategy.dataLoader.LoadStockPrice(symbol)
	if err != nil {
		return nil, fmt.Errorf("加载股价数据失败: %v", err)
	}

	// 获取当月第一个交易日
	tradingDay, err := strategy.dataLoader.GetFirstTradingDay(date.Year(), int(date.Month()), stockPrices)
	if err != nil {
		return nil, fmt.Errorf("获取交易日失败: %v", err)
	}

	tradingDayKey := tradingDay.Format("20060102")
	stockPrice, exists := stockPrices[tradingDayKey]
	if !exists {
		return nil, fmt.Errorf("未找到 %s 在 %s 的股价数据", symbol, tradingDayKey)
	}

	// 计算可买入的整数股数
	shares := int(math.Floor(cashAmount.Div(stockPrice.Close).InexactFloat64()))
	if shares <= 0 {
		return nil, fmt.Errorf("资金不足以买入 %s", symbol)
	}

	// 计算实际花费金额
	actualAmount := stockPrice.Close.Mul(decimal.NewFromInt(int64(shares)))
	
	// 检查现金是否足够
	if actualAmount.GreaterThan(portfolio.Cash) {
		return nil, fmt.Errorf("现金不足以买入 %s", symbol)
	}

	// 更新现金和持仓
	portfolio.Cash = portfolio.Cash.Sub(actualAmount)
	
	position := &Position{
		Symbol:       symbol,
		Shares:       shares,
		BuyPrice:     stockPrice.Close,
		BuyDate:      tradingDay,
		CurrentPrice: stockPrice.Close,
		MarketValue:  actualAmount,
		CostBasis:    actualAmount,
		PnL:          decimal.Zero,
		PnLPercent:   decimal.Zero,
	}
	
	portfolio.Positions[symbol] = position

	// 创建交易记录
	action := &TradingAction{
		Date:   tradingDay,
		Symbol: symbol,
		Action: "BUY",
		Shares: shares,
		Price:  stockPrice.Close,
		Amount: actualAmount,
		Reason: "股票被纳入",
	}

	fmt.Printf("买入: %s, 股数: %d, 价格: %s, 金额: %s\n", 
		symbol, shares, stockPrice.Close.String(), actualAmount.String())

	return action, nil
}

// updatePortfolioValue 更新投资组合价值
func (strategy *TradingStrategy) updatePortfolioValue(portfolio *Portfolio, date time.Time) error {
	totalStockValue := decimal.Zero
	
	for symbol, position := range portfolio.Positions {
		// 获取当前股价
		stockPrices, err := strategy.dataLoader.LoadStockPrice(symbol)
		if err != nil {
			fmt.Printf("警告: 无法加载 %s 的股价数据: %v\n", symbol, err)
			continue
		}

		// 获取当月第一个交易日的价格
		tradingDay, err := strategy.dataLoader.GetFirstTradingDay(date.Year(), int(date.Month()), stockPrices)
		if err != nil {
			fmt.Printf("警告: 无法获取 %s 的交易日: %v\n", symbol, err)
			continue
		}

		tradingDayKey := tradingDay.Format("20060102")
		stockPrice, exists := stockPrices[tradingDayKey]
		if !exists {
			fmt.Printf("警告: 未找到 %s 在 %s 的股价数据\n", symbol, tradingDayKey)
			continue
		}

		// 更新持仓信息
		position.CurrentPrice = stockPrice.Close
		position.MarketValue = stockPrice.Close.Mul(decimal.NewFromInt(int64(position.Shares)))
		position.PnL = position.MarketValue.Sub(position.CostBasis)
		if !position.CostBasis.IsZero() {
			position.PnLPercent = position.PnL.Div(position.CostBasis)
		}
		
		totalStockValue = totalStockValue.Add(position.MarketValue)
	}

	// 更新总价值
	portfolio.Value = portfolio.Cash.Add(totalStockValue)
	portfolio.Date = date

	// 计算持仓占比
	for _, position := range portfolio.Positions {
		if !portfolio.Value.IsZero() {
			position.Weight = position.MarketValue.Div(portfolio.Value)
		}
	}

	return nil
}

// copyPositions 复制持仓信息
func copyPositions(positions map[string]*Position) map[string]*Position {
	copy := make(map[string]*Position)
	for symbol, position := range positions {
		copy[symbol] = &Position{
			Symbol:       position.Symbol,
			Shares:       position.Shares,
			BuyPrice:     position.BuyPrice,
			BuyDate:      position.BuyDate,
			CurrentPrice: position.CurrentPrice,
			MarketValue:  position.MarketValue,
			CostBasis:    position.CostBasis,
			PnL:          position.PnL,
			PnLPercent:   position.PnLPercent,
			Weight:       position.Weight,
		}
	}
	return copy
}