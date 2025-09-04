package main

import (
	"time"
)

// Config 系统配置
type Config struct {
	InitialCapital float64   // 初始资金
	StartDate      time.Time // 开始日期
	EndDate        time.Time // 结束日期
	StockPriceDir  string    // 股价数据目录
	HistoryDir     string    // 交易信号数据目录
	OutputDir      string    // 输出目录
	ChartsDir      string    // 图表输出目录
	ReportsDir     string    // 报告输出目录
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	startDate, _ := time.Parse("20060102", "20230101")
	endDate, _ := time.Parse("20060102", "20250831")
	
	return &Config{
		InitialCapital: 100000.0, // 默认10万美元
		StartDate:      startDate,
		EndDate:        endDate,
		StockPriceDir:  "stock_price",
		HistoryDir:     "history",
		OutputDir:      "output",
		ChartsDir:      "output/charts",
		ReportsDir:     "output/reports",
	}
}