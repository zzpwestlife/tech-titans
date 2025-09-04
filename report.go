package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/shopspring/decimal"
)

// ReportGenerator 报告生成器
type ReportGenerator struct {
	config *Config
}

// NewReportGenerator 创建新的报告生成器
func NewReportGenerator(config *Config) *ReportGenerator {
	return &ReportGenerator{
		config: config,
	}
}

// GenerateMonthlyReport 生成月度报告
func (rg *ReportGenerator) GenerateMonthlyReport(report *MonthlyReport) error {
	// 创建输出目录
	outputDir := filepath.Join(rg.config.OutputDir, "monthly_reports")
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 生成文件名
	filename := fmt.Sprintf("%d%02d_monthly_report.csv", 
		report.Date.Year(), int(report.Date.Month()))
	filePath := filepath.Join(outputDir, filename)

	// 创建CSV文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建报告文件失败: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入标题
	headers := []string{
		"Symbol", "Buy Date", "Buy Price", "Current Price", 
		"Shares", "Market Value", "Cost Basis", "P&L", 
		"P&L %", "Weight %", "Trading Actions",
	}
	err = writer.Write(headers)
	if err != nil {
		return fmt.Errorf("写入标题失败: %v", err)
	}

	// 写入持仓数据
	for _, position := range report.Positions {
		row := []string{
			position.Symbol,
			position.BuyDate.Format("2006-01-02"),
			position.BuyPrice.StringFixed(2),
			position.CurrentPrice.StringFixed(2),
			strconv.Itoa(position.Shares),
			position.MarketValue.StringFixed(2),
			position.CostBasis.StringFixed(2),
			position.PnL.StringFixed(2),
			position.PnLPercent.Mul(decimal.NewFromInt(100)).StringFixed(2),
			position.Weight.Mul(decimal.NewFromInt(100)).StringFixed(2),
			rg.formatTradingActions(report.TradingActions, position.Symbol),
		}
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("写入持仓数据失败: %v", err)
		}
	}

	// 写入汇总信息
	summaryRows := [][]string{
		{"", "", "", "", "", "", "", "", "", "", ""},
		{"Summary", "", "", "", "", "", "", "", "", "", ""},
		{"Total Value", report.TotalValue.StringFixed(2), "", "", "", "", "", "", "", "", ""},
		{"Cash", report.Cash.StringFixed(2), "", "", "", "", "", "", "", "", ""},
		{"Stock Value", report.StockValue.StringFixed(2), "", "", "", "", "", "", "", "", ""},
		{"Monthly Return %", report.MonthlyReturn.Mul(decimal.NewFromInt(100)).StringFixed(2), "", "", "", "", "", "", "", "", ""},
		{"Cumulative Return %", report.CumulativeReturn.Mul(decimal.NewFromInt(100)).StringFixed(2), "", "", "", "", "", "", "", "", ""},
	}

	for _, row := range summaryRows {
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("写入汇总信息失败: %v", err)
		}
	}

	fmt.Printf("月度报告已生成: %s\n", filePath)
	return nil
}

// GenerateFinalReport 生成最终持仓报告
func (rg *ReportGenerator) GenerateFinalReport(reports []*MonthlyReport) error {
	if len(reports) == 0 {
		return fmt.Errorf("没有报告数据")
	}

	// 创建输出目录
	err := os.MkdirAll(rg.config.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 生成最终持仓报告
	err = rg.generateFinalPositionReport(reports)
	if err != nil {
		return fmt.Errorf("生成最终持仓报告失败: %v", err)
	}

	// 生成业绩汇总报告
	err = rg.generatePerformanceSummary(reports)
	if err != nil {
		return fmt.Errorf("生成业绩汇总报告失败: %v", err)
	}

	return nil
}

// generateFinalPositionReport 生成最终持仓报告
func (rg *ReportGenerator) generateFinalPositionReport(reports []*MonthlyReport) error {
	lastReport := reports[len(reports)-1]
	filePath := filepath.Join(rg.config.OutputDir, "final_position_report.csv")

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建最终持仓报告文件失败: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入标题
	headers := []string{
		"Symbol", "Buy Date", "Buy Price", "Final Price", 
		"Shares", "Final Market Value", "Cost Basis", "Total P&L", 
		"Total P&L %", "Final Weight %",
	}
	err = writer.Write(headers)
	if err != nil {
		return fmt.Errorf("写入标题失败: %v", err)
	}

	// 按市值排序持仓
	var positions []*Position
	for _, position := range lastReport.Positions {
		positions = append(positions, position)
	}
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].MarketValue.GreaterThan(positions[j].MarketValue)
	})

	// 写入持仓数据
	for _, position := range positions {
		row := []string{
			position.Symbol,
			position.BuyDate.Format("2006-01-02"),
			position.BuyPrice.StringFixed(2),
			position.CurrentPrice.StringFixed(2),
			strconv.Itoa(position.Shares),
			position.MarketValue.StringFixed(2),
			position.CostBasis.StringFixed(2),
			position.PnL.StringFixed(2),
			position.PnLPercent.Mul(decimal.NewFromInt(100)).StringFixed(2),
			position.Weight.Mul(decimal.NewFromInt(100)).StringFixed(2),
		}
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("写入持仓数据失败: %v", err)
		}
	}

	// 写入汇总信息
	summaryRows := [][]string{
		{"", "", "", "", "", "", "", "", "", ""},
		{"Final Summary", "", "", "", "", "", "", "", "", ""},
		{"Total Portfolio Value", lastReport.TotalValue.StringFixed(2), "", "", "", "", "", "", "", ""},
		{"Cash Balance", lastReport.Cash.StringFixed(2), "", "", "", "", "", "", "", ""},
		{"Total Stock Value", lastReport.StockValue.StringFixed(2), "", "", "", "", "", "", "", ""},
		{"Total Return %", lastReport.CumulativeReturn.Mul(decimal.NewFromInt(100)).StringFixed(2), "", "", "", "", "", "", "", ""},
		{"Initial Capital", decimal.NewFromFloat(rg.config.InitialCapital).StringFixed(2), "", "", "", "", "", "", "", ""},
	}

	for _, row := range summaryRows {
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("写入汇总信息失败: %v", err)
		}
	}

	fmt.Printf("最终持仓报告已生成: %s\n", filePath)
	return nil
}

// generatePerformanceSummary 生成业绩汇总报告
func (rg *ReportGenerator) generatePerformanceSummary(reports []*MonthlyReport) error {
	filePath := filepath.Join(rg.config.OutputDir, "performance_summary.csv")

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建业绩汇总报告文件失败: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入标题
	headers := []string{
		"Date", "Total Value", "Cash", "Stock Value", 
		"Monthly Return %", "Cumulative Return %", "Number of Positions",
	}
	err = writer.Write(headers)
	if err != nil {
		return fmt.Errorf("写入标题失败: %v", err)
	}

	// 写入每月业绩数据
	for _, report := range reports {
		row := []string{
			report.Date.Format("2006-01-02"),
			report.TotalValue.StringFixed(2),
			report.Cash.StringFixed(2),
			report.StockValue.StringFixed(2),
			report.MonthlyReturn.Mul(decimal.NewFromInt(100)).StringFixed(2),
			report.CumulativeReturn.Mul(decimal.NewFromInt(100)).StringFixed(2),
			strconv.Itoa(len(report.Positions)),
		}
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("写入业绩数据失败: %v", err)
		}
	}

	fmt.Printf("业绩汇总报告已生成: %s\n", filePath)
	return nil
}

// formatTradingActions 格式化交易行为
func (rg *ReportGenerator) formatTradingActions(actions []TradingAction, symbol string) string {
	var result string
	for _, action := range actions {
		if action.Symbol == symbol {
			if result != "" {
				result += "; "
			}
			result += fmt.Sprintf("%s %d shares at $%s (%s)", 
				action.Action, action.Shares, action.Price.StringFixed(2), action.Reason)
		}
	}
	return result
}

// PrintSummary 打印汇总信息到控制台
func (rg *ReportGenerator) PrintSummary(reports []*MonthlyReport) {
	if len(reports) == 0 {
		fmt.Println("没有报告数据")
		return
	}

	lastReport := reports[len(reports)-1]
	initialCapital := decimal.NewFromFloat(rg.config.InitialCapital)

	fmt.Println("\n=== 投资策略执行汇总 ===")
	fmt.Printf("执行周期: %s 至 %s\n", 
		rg.config.StartDate.Format("2006-01-02"), 
		rg.config.EndDate.Format("2006-01-02"))
	fmt.Printf("初始资金: $%s\n", initialCapital.StringFixed(2))
	fmt.Printf("最终价值: $%s\n", lastReport.TotalValue.StringFixed(2))
	fmt.Printf("现金余额: $%s\n", lastReport.Cash.StringFixed(2))
	fmt.Printf("股票市值: $%s\n", lastReport.StockValue.StringFixed(2))
	fmt.Printf("总收益率: %s%%\n", lastReport.CumulativeReturn.Mul(decimal.NewFromInt(100)).StringFixed(2))
	fmt.Printf("持仓数量: %d\n", len(lastReport.Positions))
	fmt.Printf("总交易月数: %d\n", len(reports))

	// 计算年化收益率
	if len(reports) > 0 {
		months := decimal.NewFromInt(int64(len(reports)))
		years := months.Div(decimal.NewFromInt(12))
		if !years.IsZero() && !lastReport.CumulativeReturn.IsZero() {
			// 年化收益率 = (1 + 总收益率)^(1/年数) - 1
			annualizedReturn := lastReport.CumulativeReturn.Add(decimal.NewFromInt(1))
			annualizedReturn = decimal.NewFromFloat(math.Pow(annualizedReturn.InexactFloat64(), 1.0/years.InexactFloat64()))
			annualizedReturn = annualizedReturn.Sub(decimal.NewFromInt(1))
			fmt.Printf("年化收益率: %s%%\n", annualizedReturn.Mul(decimal.NewFromInt(100)).StringFixed(2))
		}
	}

	fmt.Println("\n=== 前10大持仓 ===")
	var positions []*Position
	for _, position := range lastReport.Positions {
		positions = append(positions, position)
	}
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].MarketValue.GreaterThan(positions[j].MarketValue)
	})

	for i, position := range positions {
		if i >= 10 {
			break
		}
		fmt.Printf("%d. %s: $%s (%s%%)\n", 
			i+1, position.Symbol, position.MarketValue.StringFixed(2), 
			position.Weight.Mul(decimal.NewFromInt(100)).StringFixed(2))
	}
}