package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// 命令行参数
	var (
		initialCapital = flag.Float64("capital", 100000, "Initial capital in USD")
		startDate      = flag.String("start", "20230101", "Start date (YYYYMMDD)")
		endDate        = flag.String("end", "20250831", "End date (YYYYMMDD)")
		stockPriceDir  = flag.String("stock-dir", "stock_price", "Stock price data directory")
		historyDir     = flag.String("history-dir", "history", "Trading history directory")
		outputDir      = flag.String("output-dir", "output", "Output directory")
	)
	flag.Parse()

	// 解析日期
	startTime, err := time.Parse("20060102", *startDate)
	if err != nil {
		log.Fatalf("Invalid start date format: %v", err)
	}
	endTime, err := time.Parse("20060102", *endDate)
	if err != nil {
		log.Fatalf("Invalid end date format: %v", err)
	}

	// 创建配置
	config := &Config{
		InitialCapital: *initialCapital,
		StartDate:      startTime,
		EndDate:        endTime,
		StockPriceDir:  *stockPriceDir,
		HistoryDir:     *historyDir,
		OutputDir:      *outputDir,
		ChartsDir:      filepath.Join(*outputDir, "charts"),
		ReportsDir:     filepath.Join(*outputDir, "reports"),
	}

	// 创建输出目录
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fmt.Printf("=== Tech Titans Quantitative Investment Strategy Analysis ===\n")
	fmt.Printf("Initial Capital: $%.2f\n", config.InitialCapital)
	fmt.Printf("Analysis Period: %s - %s\n", config.StartDate, config.EndDate)
	fmt.Printf("Stock Price Directory: %s\n", config.StockPriceDir)
	fmt.Printf("History Directory: %s\n", config.HistoryDir)
	fmt.Printf("Output Directory: %s\n", config.OutputDir)
	fmt.Println()

	// 初始化数据加载器
	dataLoader := NewStockDataLoader(config.StockPriceDir, config.HistoryDir)

	// 初始化交易策略
	strategy := NewTradingStrategy(dataLoader, config)

	// 执行策略
	fmt.Println("Executing trading strategy...")
	start := time.Now()
	reports, err := strategy.ExecuteStrategy()
	if err != nil {
		log.Fatalf("Strategy execution failed: %v", err)
	}
	executionTime := time.Since(start)
	fmt.Printf("Strategy execution completed in %v\n\n", executionTime)

	// 生成报告
	fmt.Println("Generating reports...")
	reportGenerator := NewReportGenerator(config)

	// 生成月度报告
	for _, report := range reports {
		if err := reportGenerator.GenerateMonthlyReport(report); err != nil {
			log.Printf("Failed to generate monthly report for %s: %v", report.Date.Format("2006-01"), err)
		}
	}

	// 生成最终报告
	if len(reports) > 0 {
		finalReport := reports[len(reports)-1]
		if err := reportGenerator.generateFinalPositionReport([]*MonthlyReport{finalReport}); err != nil {
			log.Printf("Failed to generate final position report: %v", err)
		}

		if err := reportGenerator.generatePerformanceSummary(reports); err != nil {
			log.Printf("Failed to generate performance summary: %v", err)
		}

		// 打印控制台摘要
		reportGenerator.PrintSummary(reports)
	}

	// 生成图表
	fmt.Println("\nGenerating charts...")
	chartGenerator := NewChartGenerator(config)
	if err := chartGenerator.GenerateAllCharts(reports); err != nil {
		log.Printf("Failed to generate charts: %v", err)
	} else {
		fmt.Println("Charts generated successfully")
	}

	fmt.Printf("\n=== Analysis Complete ===\n")
	fmt.Printf("Total execution time: %v\n", time.Since(start))
	fmt.Printf("Reports and charts saved to: %s\n", config.OutputDir)

	// 列出生成的文件
	fmt.Println("\nGenerated files:")
	filepath.Walk(config.OutputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(config.OutputDir, path)
			fmt.Printf("  - %s\n", relPath)
		}
		return nil
	})
}