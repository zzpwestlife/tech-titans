作为专业的美股量化投资策略分析师，协助我评估该策略的收益表现。以下是策略执行的具体要求：

## 策略参数
- 初始资金：${initial_capital} 美元（可配置）
- 初始状态：空仓
- 统计周期：20230101 至 20250831（可配置）


## 股价数据来源: stock_price 目录内的 csv 文件
每个文件的文件名格式为: {symbol}.csv, 例如: AAPL.csv, MSFT.csv 等, 每个文件的内容为:
- Date: 日期
- Open: 开盘价 (单位: 美元, 不参与计算)
- High: 最高价 (单位: 美元, 不参与计算)
- Low: 最低价 (单位: 美元, 不参与计算)
- Close: 收盘价 (单位: 美元)
- Adj Close: 调整后的收盘价 (如果有该列, 取数时 Close 取该列的值, 否则取 Close 列的值, 如果单行数据为空, 说明为分红等情况, 忽略当前行) (单位: 美元)
- Volume: 成交量 (可选, 取数时不使用) (单位: 股)

交易数据来源: history 目录内的 csv 文件, 每个文件的文件名格式为: {YYYYMMDD}.csv, 例如: 20230101.csv 等, 每个文件的内容为:
- symbol: 股票代码
- status: 股票状态 (剔除/纳入)
其他字段不关心, 不参与计算

交易规则：
1. 股票剔除：
   - 若持仓中存在状态标记为 "剔除" 的股票，在当月首个交易日以收盘价全部卖出
   - 示例：SYNA 标记为剔除时，以 20230103 (Jan 3, 2023) 收盘价清仓, 资金回到现金

2. 股票纳入：
   - 渐进式建仓策略：30%→40%→50%→...→90%（最少保留 10% 现金）
   - 新纳入股票等分分配可用资金，按当月首个交易日收盘价买入整数股
   - 示例： 20230101 当月的 10 万美金分配 30% 给 2 支新股票 (MU, UCTT)，各买入 15,000 美元 对应的整数股数

数据要求：
- 输出数据：
  a) 月度持仓明细：
     - 股票代码、买入日期、买入价格
     - 当前价格、市值、成本
     - 持仓占比、盈亏比例
  b) 最终持仓报告：
     - 总股票市值
     - 总现金余额

实现要求：
- 使用 Golang 开发
- 可视化图表使用英文标签

执行流程：
1. 按月遍历交易周期
2. 处理剔除股票（如有, 比如 ANSS 的数据找不到）
3. 计算新纳入股票的资金分配
4. 执行买卖操作
5. 生成月度报告
6. 最终生成汇总报告和可视化图表


当前所有股票标的
AAPL, ACLS, ACMR, ADBE, ADI, ADSK, AGYS, AKAM, ALAB, ALGM, ALRM, AMAT, AMKR, ANSS, APP, ASGN, ATEN, BDC, BILL, BL, CDNS, CDW, CGNX, CHKP, COHU, CORZ, CRUS, CRWD, CVLT, DBX, DDOG, DIOD, DLB, DOCN, DOCU, DV, DXC, ENPH, EPAM, EXTR, FFIV, FORM, FTNT, GFS, HLIT, IBM, ICHR, IDCC, INTC, INTU, IPGP, IT, JBL, KEYS, KLAC, LITE, LSCC, MANH, MARA, MCHP, MEI, MKSI, MPWR, MRVL, MSFT, MSI, MSTR, MU, MXL, NOW, NSIT, NSSC, NTAP, NTNX, NVDA, OLED, ON, ORCL, OTEX, PANW, PENG, PI, PLAB, PLUS, PRF, PRO, PSTG, QCOM, QLYS, QRVO, RMBS, RNG, ROG, ROP, RPD, SANM, SCSC, SMCI, SMTC, SNOW, SNPS, SPNS, SPSC, SWKS, SYNA, TDC, TER, TRMB, TTMI, TXN, UCTT, VECO, VIAV, VRNT, VSAT, VYX, WDAY, WDC, XRX, YOU, ZBRA, ZM, ZS


ANSS 找不到
https://cn.investing.com/pro/propicks/tech-titans

https://finance.yahoo.com/quote/AAPL/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ACLS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ACMR/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ADBE/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ADI/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ADSK/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/AGYS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/AKAM/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ALAB/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ALGM/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ALRM/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/AMAT/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/AMKR/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ANSS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/APP/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ASGN/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ATEN/history/?period1=1599194504&period2=1756960898


https://finance.yahoo.com/quote/BDC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/BILL/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/BL/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/CDNS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/CDW/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/CGNX/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/CHKP/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/COHU/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/CORZ/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/CRUS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/CRWD/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/CVLT/history/?period1=1599194504&period2=1756960898


https://finance.yahoo.com/quote/DBX/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/DDOG/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/DIOD/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/DLB/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/DOCN/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/DOCU/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/DV/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/DXC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ENPH/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/EPAM/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/EXTR/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/FFIV/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/FORM/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/FTNT/history/?period1=1599194504&period2=1756960898


https://finance.yahoo.com/quote/GFS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/HLIT/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/IBM/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ICHR/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/IDCC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/INTC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/INTU/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/IPGP/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/IT/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/JBL/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/KEYS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/KLAC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/LITE/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/LSCC/history/?period1=1599194504&period2=1756960898


https://finance.yahoo.com/quote/MANH/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MARA/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MCHP/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MEI/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MKSI/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MPWR/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MRVL/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MSFT/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MSI/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MSTR/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MU/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/MXL/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/NOW/history/?period1=1599194504&period2=1756960898


https://finance.yahoo.com/quote/NSIT/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/NSSC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/NTAP/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/NTNX/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/NVDA/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/OLED/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ON/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ORCL/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/OTEX/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/PANW/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/PENG/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/PI/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/PLAB/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/PLUS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/PRF/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/PRO/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/PSTG/history/?period1=1599194504&period2=1756960898


https://finance.yahoo.com/quote/QCOM/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/QLYS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/QRVO/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/RMBS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/RNG/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ROG/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ROP/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/RPD/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SANM/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SCSC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SMCI/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SMTC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SNOW/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SNPS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SPNS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SPSC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SWKS/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/SYNA/history/?period1=1599194504&period2=1756960898


https://finance.yahoo.com/quote/TDC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/TER/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/TRMB/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/TTMI/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/TXN/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/UCTT/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/VECO/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/VIAV/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/VRNT/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/VSAT/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/VYX/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/WDAY/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/WDC/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/XRX/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/YOU/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ZBRA/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ZM/history/?period1=1599194504&period2=1756960898
https://finance.yahoo.com/quote/ZS/history/?period1=1599194504&period2=1756960898

