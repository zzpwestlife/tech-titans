package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// StockDataLoader 股价数据加载器
type StockDataLoader struct {
	stockPriceDir string
	historyDir    string
}

// NewStockDataLoader 创建新的数据加载器
func NewStockDataLoader(stockPriceDir, historyDir string) *StockDataLoader {
	return &StockDataLoader{
		stockPriceDir: stockPriceDir,
		historyDir:    historyDir,
	}
}

// LoadStockPrice 加载指定股票的价格数据
func (loader *StockDataLoader) LoadStockPrice(symbol string) (map[string]*StockPrice, error) {
	filePath := filepath.Join(loader.stockPriceDir, symbol+".csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开股价文件 %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// 设置CSV reader不严格检查字段数量，因为Volume字段可能包含逗号
	reader.FieldsPerRecord = -1
	header, err := reader.Read() // 读取表头
	if err != nil {
		return nil, fmt.Errorf("读取CSV表头失败: %v", err)
	}

	// 查找列索引
	columnIndex := make(map[string]int)
	for i, col := range header {
		columnIndex[col] = i
	}

	prices := make(map[string]*StockPrice)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取CSV记录失败: %v", err)
		}

		// 解析日期
		dateStr := strings.Trim(record[columnIndex["Date"]], `"`)
		date, err := parseDate(dateStr)
		if err != nil {
			fmt.Printf("警告: 无法解析日期 %s: %v\n", dateStr, err)
			continue
		}

		// 处理Volume字段可能被逗号分割的情况
		volumeIdx := columnIndex["Volume"]
		if len(record) > len(header) {
			// 如果字段数量超过表头，说明Volume字段被分割了
			extraFields := len(record) - len(header)
			// 将Volume及其后续字段重新组合
			volumeStr := ""
			for i := 0; i <= extraFields; i++ {
				if volumeIdx+i < len(record) {
					volumeStr += record[volumeIdx+i]
				}
			}
			// 创建修正后的记录
			correctedRecord := make([]string, len(header))
			copy(correctedRecord, record[:volumeIdx])
			correctedRecord[volumeIdx] = volumeStr
			record = correctedRecord
		}

		// 检查记录字段数量是否匹配表头
		if len(record) != len(header) {
			fmt.Printf("警告: 第%d行字段数量不匹配，跳过该行\n", len(prices)+2)
			continue
		}

		// 解析价格数据
		open, _ := parseDecimal(record[columnIndex["Open"]])
		high, _ := parseDecimal(record[columnIndex["High"]])
		low, _ := parseDecimal(record[columnIndex["Low"]])
		close, _ := parseDecimal(record[columnIndex["Close"]])
		
		// 处理调整后收盘价
		adjClose := close // 默认使用 Close 的值
		if adjCloseIdx, exists := columnIndex["Adj Close"]; exists && adjCloseIdx < len(record) {
			// 如果存在 Adj Close 字段，则使用该字段的值
			if adjCloseVal, err := parseDecimal(record[adjCloseIdx]); err == nil && !adjCloseVal.IsZero() {
				adjClose = adjCloseVal
			}
		}
		// 注意：对于没有 Adj Close 字段的文件（如 UCTT.csv），adjClose 会使用 close 的值

		// 解析成交量
		volume, _ := parseVolume(record[columnIndex["Volume"]])

		// 跳过无效数据行（如分红等情况）
		if close.IsZero() || adjClose.IsZero() {
			continue
		}

		stockPrice := &StockPrice{
			Date:     date,
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			AdjClose: adjClose,
			Volume:   volume,
		}

		dateKey := date.Format("20060102")
		prices[dateKey] = stockPrice
	}

	return prices, nil
}

// LoadTradeSignals 加载指定日期的交易信号
func (loader *StockDataLoader) LoadTradeSignals(date time.Time) ([]*TradeSignal, error) {
	year := date.Format("2006")
	dateStr := date.Format("20060102")
	filePath := filepath.Join(loader.historyDir, year, dateStr+".csv")
	
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开交易信号文件 %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read() // 读取表头
	if err != nil {
		return nil, fmt.Errorf("读取CSV表头失败: %v", err)
	}

	// 查找列索引
	columnIndex := make(map[string]int)
	for i, col := range header {
		columnIndex[col] = i
	}

	var signals []*TradeSignal

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取CSV记录失败: %v", err)
		}

		signal := &TradeSignal{
			Symbol: record[columnIndex["symbol"]],
			Name:   record[columnIndex["name"]],
			Price:  record[columnIndex["price"]],
			PL:     record[columnIndex["pl"]],
			Status: record[columnIndex["status"]],
		}

		signals = append(signals, signal)
	}

	return signals, nil
}

// GetFirstTradingDay 获取指定月份的第一个交易日
func (loader *StockDataLoader) GetFirstTradingDay(year, month int, stockPrices map[string]*StockPrice) (time.Time, error) {
	// 从月初开始查找第一个有数据的交易日
	for day := 1; day <= 31; day++ {
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		if date.Month() != time.Month(month) {
			break // 超出当月范围
		}
		
		dateKey := date.Format("20060102")
		if _, exists := stockPrices[dateKey]; exists {
			return date, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("未找到 %d年%d月 的交易日", year, month)
}

// parseDate 解析日期字符串
func parseDate(dateStr string) (time.Time, error) {
	// 支持多种日期格式
	formats := []string{
		"Jan 2, 2006",
		"2006-01-02",
		"20060102",
		"01/02/2006",
	}
	
	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("无法解析日期格式: %s", dateStr)
}

// parseDecimal 解析十进制数字
func parseDecimal(str string) (decimal.Decimal, error) {
	str = strings.Trim(str, `"`)
	str = strings.ReplaceAll(str, ",", "") // 移除千位分隔符
	if str == "" {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(str)
}

// parseVolume 解析成交量
func parseVolume(str string) (int64, error) {
	str = strings.Trim(str, `"`)
	str = strings.ReplaceAll(str, ",", "") // 移除千位分隔符
	if str == "" {
		return 0, nil
	}
	return strconv.ParseInt(str, 10, 64)
}