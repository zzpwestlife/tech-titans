package main

import (
	"time"
	"github.com/shopspring/decimal"
)

// StockPrice 股价数据结构
type StockPrice struct {
	Date     time.Time       `csv:"Date"`
	Open     decimal.Decimal `csv:"Open"`
	High     decimal.Decimal `csv:"High"`
	Low      decimal.Decimal `csv:"Low"`
	Close    decimal.Decimal `csv:"Close"`
	AdjClose decimal.Decimal `csv:"Adj Close"`
	Volume   int64           `csv:"Volume"`
}

// TradeSignal 交易信号数据结构
type TradeSignal struct {
	Symbol string `csv:"symbol"`
	Name   string `csv:"name"`
	Price  string `csv:"price"`
	PL     string `csv:"pl"`
	Status string `csv:"status"` // "纳入" 或 "剔除"
}

// Position 持仓信息
type Position struct {
	Symbol       string          // 股票代码
	Shares       int             // 持股数量
	BuyPrice     decimal.Decimal // 买入价格
	BuyDate      time.Time       // 买入日期
	CurrentPrice decimal.Decimal // 当前价格
	MarketValue  decimal.Decimal // 市值
	CostBasis    decimal.Decimal // 成本基础
	PnL          decimal.Decimal // 盈亏
	PnLPercent   decimal.Decimal // 盈亏百分比
	Weight       decimal.Decimal // 持仓占比
}

// Portfolio 投资组合
type Portfolio struct {
	Cash      decimal.Decimal            // 现金余额
	Positions map[string]*Position       // 持仓映射
	Value     decimal.Decimal            // 总价值
	Date      time.Time                  // 日期
}

// MonthlyReport 月度报告
type MonthlyReport struct {
	Date           time.Time                // 报告日期
	TotalValue     decimal.Decimal          // 总价值
	Cash           decimal.Decimal          // 现金
	StockValue     decimal.Decimal          // 股票市值
	MonthlyReturn  decimal.Decimal          // 月度收益率
	CumulativeReturn decimal.Decimal        // 累计收益率
	Positions      map[string]*Position     // 持仓详情
	TradingActions []TradingAction          // 交易行为
}

// TradingAction 交易行为
type TradingAction struct {
	Date   time.Time       // 交易日期
	Symbol string          // 股票代码
	Action string          // 行为类型："BUY" 或 "SELL"
	Shares int             // 股数
	Price  decimal.Decimal // 价格
	Amount decimal.Decimal // 金额
	Reason string          // 交易原因
}

// PerformanceMetrics 绩效指标
type PerformanceMetrics struct {
	TotalReturn       decimal.Decimal // 总收益率
	AnnualizedReturn  decimal.Decimal // 年化收益率
	MaxDrawdown       decimal.Decimal // 最大回撤
	SharpeRatio       decimal.Decimal // 夏普比率
	Volatility        decimal.Decimal // 波动率
	WinRate           decimal.Decimal // 胜率
	AverageReturn     decimal.Decimal // 平均收益率
	TotalTrades       int             // 总交易次数
}