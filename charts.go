package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/shopspring/decimal"
)

// Helper function to create Bool pointer
func boolPtr(b bool) types.Bool {
	return &b
}

// ChartGenerator 图表生成器
type ChartGenerator struct {
	config *Config
}

// NewChartGenerator 创建新的图表生成器
func NewChartGenerator(config *Config) *ChartGenerator {
	return &ChartGenerator{
		config: config,
	}
}

// GenerateAllCharts 生成所有图表
func (cg *ChartGenerator) GenerateAllCharts(reports []*MonthlyReport) error {
	if len(reports) == 0 {
		return fmt.Errorf("no report data available")
	}

	// 创建图表输出目录
	chartDir := filepath.Join(cg.config.OutputDir, "charts")
	err := os.MkdirAll(chartDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create chart directory: %v", err)
	}

	// 生成投资组合价值趋势图
	err = cg.generatePortfolioValueChart(reports, chartDir)
	if err != nil {
		return fmt.Errorf("failed to generate portfolio value chart: %v", err)
	}

	// 生成收益率趋势图
	err = cg.generateReturnChart(reports, chartDir)
	if err != nil {
		return fmt.Errorf("failed to generate return chart: %v", err)
	}

	// 生成资产配置饼图
	err = cg.generateAssetAllocationChart(reports[len(reports)-1], chartDir)
	if err != nil {
		return fmt.Errorf("failed to generate asset allocation chart: %v", err)
	}

	// 生成持仓分布图
	err = cg.generatePositionDistributionChart(reports[len(reports)-1], chartDir)
	if err != nil {
		return fmt.Errorf("failed to generate position distribution chart: %v", err)
	}

	// 生成月度交易活动图
	err = cg.generateTradingActivityChart(reports, chartDir)
	if err != nil {
		return fmt.Errorf("failed to generate trading activity chart: %v", err)
	}

	fmt.Printf("All charts generated successfully in: %s\n", chartDir)
	return nil
}

// generatePortfolioValueChart 生成投资组合价值趋势图
func (cg *ChartGenerator) generatePortfolioValueChart(reports []*MonthlyReport, outputDir string) error {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Portfolio Value Trend",
			Subtitle: "Monthly Portfolio Value Over Time",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Date",
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Value (USD)",
			Type: "value",
		}),
		charts.WithLegendOpts(opts.Legend{Show: boolPtr(true)}),
		charts.WithTooltipOpts(opts.Tooltip{Show: boolPtr(true)}),
	)

	// 准备数据
	var xAxis []string
	var totalValues []opts.LineData
	var cashValues []opts.LineData
	var stockValues []opts.LineData

	for _, report := range reports {
		xAxis = append(xAxis, report.Date.Format("2006-01"))
		totalValueFloat, _ := report.TotalValue.Round(2).Float64()
		cashFloat, _ := report.Cash.Round(2).Float64()
		stockValueFloat, _ := report.StockValue.Round(2).Float64()
		totalValues = append(totalValues, opts.LineData{Value: totalValueFloat})
		cashValues = append(cashValues, opts.LineData{Value: cashFloat})
		stockValues = append(stockValues, opts.LineData{Value: stockValueFloat})
	}

	line.SetXAxis(xAxis).
		AddSeries("Total Value", totalValues).
		AddSeries("Cash", cashValues).
		AddSeries("Stock Value", stockValues).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: boolPtr(true)}),
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name: "Maximum",
				Type: "max",
			}),
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name: "Minimum",
				Type: "min",
			}),
		)

	// 保存图表
	filePath := filepath.Join(outputDir, "portfolio_value_trend.html")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return line.Render(f)
}

// generateReturnChart 生成收益率趋势图
func (cg *ChartGenerator) generateReturnChart(reports []*MonthlyReport, outputDir string) error {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Return Trend",
			Subtitle: "Monthly and Cumulative Returns",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Date",
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Return (%)",
			Type: "value",
		}),
		charts.WithLegendOpts(opts.Legend{Show: boolPtr(true)}),
		charts.WithTooltipOpts(opts.Tooltip{Show: boolPtr(true)}),
	)

	// 准备数据
	var xAxis []string
	var monthlyReturns []opts.LineData
	var cumulativeReturns []opts.LineData

	for _, report := range reports {
		xAxis = append(xAxis, report.Date.Format("2006-01"))
		monthlyReturnFloat, _ := report.MonthlyReturn.Mul(decimal.NewFromInt(100)).Round(2).Float64()
		cumulativeReturnFloat, _ := report.CumulativeReturn.Mul(decimal.NewFromInt(100)).Round(2).Float64()
		monthlyReturns = append(monthlyReturns, opts.LineData{
			Value: monthlyReturnFloat,
		})
		cumulativeReturns = append(cumulativeReturns, opts.LineData{
			Value: cumulativeReturnFloat,
		})
	}

	line.SetXAxis(xAxis).
		AddSeries("Monthly Return (%)", monthlyReturns).
		AddSeries("Cumulative Return (%)", cumulativeReturns).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: boolPtr(true)}),
		)

	// 保存图表
	filePath := filepath.Join(outputDir, "return_trend.html")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return line.Render(f)
}

// generateAssetAllocationChart 生成资产配置饼图
func (cg *ChartGenerator) generateAssetAllocationChart(report *MonthlyReport, outputDir string) error {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Asset Allocation",
			Subtitle: fmt.Sprintf("As of %s", report.Date.Format("2006-01-02")),
		}),
		charts.WithLegendOpts(opts.Legend{Show: boolPtr(true), Orient: "vertical", Left: "left"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: boolPtr(true)}),
	)

	// 准备数据
	var items []opts.PieData
	
	// 添加现金
	cashPercent := report.Cash.Div(report.TotalValue).Mul(decimal.NewFromInt(100))
	cashPercentFloat, _ := cashPercent.Round(2).Float64()
	items = append(items, opts.PieData{
		Name:  "Cash",
		Value: cashPercentFloat,
	})

	// 添加股票
	stockPercent := report.StockValue.Div(report.TotalValue).Mul(decimal.NewFromInt(100))
	stockPercentFloat, _ := stockPercent.Round(2).Float64()
	items = append(items, opts.PieData{
		Name:  "Stocks",
		Value: stockPercentFloat,
	})

	pie.AddSeries("allocation", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      boolPtr(true),
				Formatter: "{b}: {c}%",
			}),
		)

	// 保存图表
	filePath := filepath.Join(outputDir, "asset_allocation.html")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return pie.Render(f)
}

// generatePositionDistributionChart 生成持仓分布图
func (cg *ChartGenerator) generatePositionDistributionChart(report *MonthlyReport, outputDir string) error {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Position Distribution",
			Subtitle: "Top Holdings by Market Value",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Symbol",
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Market Value (USD)",
			Type: "value",
		}),
		charts.WithLegendOpts(opts.Legend{Show: boolPtr(true)}),
		charts.WithTooltipOpts(opts.Tooltip{Show: boolPtr(true)}),
	)

	// 按市值排序持仓
	var positions []*Position
	for _, position := range report.Positions {
		positions = append(positions, position)
	}
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].MarketValue.GreaterThan(positions[j].MarketValue)
	})

	// 取前15个持仓
	maxPositions := 15
	if len(positions) < maxPositions {
		maxPositions = len(positions)
	}

	// 准备数据
	var xAxis []string
	var marketValues []opts.BarData
	var pnlValues []opts.BarData

	for i := 0; i < maxPositions; i++ {
		position := positions[i]
		xAxis = append(xAxis, position.Symbol)
		marketValueFloat, _ := position.MarketValue.Round(2).Float64()
		pnlFloat, _ := position.PnL.Round(2).Float64()
		marketValues = append(marketValues, opts.BarData{Value: marketValueFloat})
		pnlValues = append(pnlValues, opts.BarData{Value: pnlFloat})
	}

	bar.SetXAxis(xAxis).
		AddSeries("Market Value", marketValues).
		AddSeries("P&L", pnlValues)

	// 保存图表
	filePath := filepath.Join(outputDir, "position_distribution.html")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return bar.Render(f)
}

// generateTradingActivityChart 生成月度交易活动图
func (cg *ChartGenerator) generateTradingActivityChart(reports []*MonthlyReport, outputDir string) error {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Monthly Trading Activity",
			Subtitle: "Number of Buy and Sell Transactions",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Date",
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Number of Transactions",
			Type: "value",
		}),
		charts.WithLegendOpts(opts.Legend{Show: boolPtr(true)}),
		charts.WithTooltipOpts(opts.Tooltip{Show: boolPtr(true)}),
	)

	// 准备数据
	var xAxis []string
	var buyActions []opts.BarData
	var sellActions []opts.BarData

	for _, report := range reports {
		xAxis = append(xAxis, report.Date.Format("2006-01"))
		
		buyCount := 0
		sellCount := 0
		for _, action := range report.TradingActions {
			if action.Action == "BUY" {
				buyCount++
			} else if action.Action == "SELL" {
				sellCount++
			}
		}
		
		buyActions = append(buyActions, opts.BarData{Value: buyCount})
		sellActions = append(sellActions, opts.BarData{Value: sellCount})
	}

	bar.SetXAxis(xAxis).
		AddSeries("Buy Transactions", buyActions).
		AddSeries("Sell Transactions", sellActions)

	// 保存图表
	filePath := filepath.Join(outputDir, "trading_activity.html")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return bar.Render(f)
}

// GeneratePositionPieChart 生成持仓权重饼图
func (cg *ChartGenerator) GeneratePositionPieChart(report *MonthlyReport, outputDir string) error {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Position Weight Distribution",
			Subtitle: fmt.Sprintf("As of %s", report.Date.Format("2006-01-02")),
		}),
		charts.WithLegendOpts(opts.Legend{Show: boolPtr(true), Orient: "vertical", Left: "left"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: boolPtr(true)}),
	)

	// 按权重排序持仓
	var positions []*Position
	for _, position := range report.Positions {
		positions = append(positions, position)
	}
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].Weight.GreaterThan(positions[j].Weight)
	})

	// 准备数据
	var items []opts.PieData
	otherWeight := decimal.Zero
	maxShow := 10 // 显示前10个持仓

	for i, position := range positions {
		weightPercent := position.Weight.Mul(decimal.NewFromInt(100))
		if i < maxShow {
			weightPercentFloat, _ := weightPercent.Round(2).Float64()
			items = append(items, opts.PieData{
				Name:  position.Symbol,
				Value: weightPercentFloat,
			})
		} else {
			otherWeight = otherWeight.Add(position.Weight)
		}
	}

	// 添加其他持仓
	if !otherWeight.IsZero() {
		otherWeightFloat, _ := otherWeight.Mul(decimal.NewFromInt(100)).Round(2).Float64()
		items = append(items, opts.PieData{
			Name:  "Others",
			Value: otherWeightFloat,
		})
	}

	// 添加现金
	cashWeight := report.Cash.Div(report.TotalValue).Mul(decimal.NewFromInt(100))
	cashWeightFloat, _ := cashWeight.Round(2).Float64()
	items = append(items, opts.PieData{
		Name:  "Cash",
		Value: cashWeightFloat,
	})

	pie.AddSeries("Position Weight", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      boolPtr(true),
				Formatter: "{b}: {c}%",
			}),
		)

	// 保存图表
	filePath := filepath.Join(outputDir, "position_weight_distribution.html")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return pie.Render(f)
}