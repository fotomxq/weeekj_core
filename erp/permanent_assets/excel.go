package ERPPermanentAssets

import (
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2Excel "github.com/fotomxq/weeekj_core/v5/router2/excel"
	"time"
)

// ExcelSortChangeAnalysis 固定资产变动情况表
func ExcelSortChangeAnalysis(c any, logErr string, orgID int64, endAt time.Time) {
	//分析参数
	endAtCarbon := CoreFilter.GetCarbonByTime(endAt)
	//预先加载
	excelObj := Router2Excel.ExcelQuick{
		C:                 c,
		LogErr:            logErr,
		FileParams:        fmt.Sprint("erp_permanent_assets_", orgID, "_sort_change_analysis_end_", endAtCarbon.Time.Format("2006-01-02")),
		FileName:          fmt.Sprint("固定资产", endAtCarbon.Time.Format("2006"), "年度变动情况.xlsx"),
		TemplatePath:      fmt.Sprint("erp", CoreFile.Sep, "permanent_assets", CoreFile.Sep, "sort_change_analysis.xlsx"),
		NeedReplaceStyle:  true,
		ReplaceStyleSheet: "",
		ReplaceStyleRef:   "A1",
		ReplaceStyleStart: "A1",
		ReplaceStyleEnd:   "",
		ExcelObj:          nil,
		CacheSaveTime:     300,
	}
	//初始化
	if excelObj.InitCache() {
		return
	}
	//开始位置
	rowStep := 2
	//预先写入数据
	insertData := map[string]string{}
	//列总数
	var colB, colC, colD, colE, colF, colG, colH, colI, colJ, colK, colL, colM, colN, colO int64
	//获取固定资产产品分类
	sortList, _ := Sort.GetAll(orgID, 0)
	for _, vSort := range sortList {
		//购进价值
		var vAllBuyPrice int64 = 0
		//年初价值
		var vBeginPrice int64 = 0
		//期末价值
		var vEndPrice int64 = 0
		//本年度购进价值
		var vNewBuyPrice int64 = 0
		//本年度减少价值
		var vReducePrice int64 = 0
		//本年度计划折旧
		var vPlanDepreciation int64 = 0
		//本年度报废折旧价值
		var vScrapDepreciation int64 = 0
		//获取分类下产品，计算总价值
		vProductList := GetProductListBySortID(orgID, vSort.ID)
		for _, vProduct := range vProductList {
			//叠加购进价值
			vAllBuyPrice += vProduct.BuyAllPrice
			//获取年初价值
			vProductBeginLogData := getLogFirstInBetween(vProduct.ID, "check", endAtCarbon.StartOfYear().Time, endAtCarbon.EndOfYear().Time)
			if vProductBeginLogData.ID > 0 {
				vBeginPrice += vProductBeginLogData.AllPrice
			} else {
				vBeginPrice += vProduct.NowAllPrice
			}
			//获取期末价值
			vProductEndLogData := getLogLastInBetween(vProduct.ID, "check", endAtCarbon.StartOfYear().Time, endAtCarbon.EndOfYear().Time)
			if vProductEndLogData.ID > 0 {
				vEndPrice += vProductEndLogData.AllPrice
			} else {
				vEndPrice += vProduct.NowAllPrice
			}
			//如果产品创建时间在本年度内，则计入本年度购进价值
			if vProduct.CreateAt.After(endAtCarbon.StartOfYear().Time) && vProduct.CreateAt.Before(endAtCarbon.EndOfYear().Time) {
				vNewBuyPrice += vProduct.BuyAllPrice
			}
			//获取本年度销毁数据总和
			vProductScrapLogData := getLogSUMInBetween(vProduct.ID, "delete", endAtCarbon.StartOfYear().Time, endAtCarbon.EndOfYear().Time)
			if vProductScrapLogData.ID > 0 {
				vScrapDepreciation += vProductScrapLogData.AllPrice
			}
		}
		//计算本年度减少值
		vReducePrice = vBeginPrice - vEndPrice
		//构建数据
		insertData[fmt.Sprint("A", rowStep)] = vSort.Name
		insertData[fmt.Sprint("B", rowStep)] = fmt.Sprint(float64(vBeginPrice) / 100)
		colB += vBeginPrice
		insertData[fmt.Sprint("C", rowStep)] = fmt.Sprint(float64(vNewBuyPrice) / 100)
		colC += vNewBuyPrice
		insertData[fmt.Sprint("D", rowStep)] = fmt.Sprint(float64(vReducePrice) / 100)
		colD += vReducePrice
		insertData[fmt.Sprint("E", rowStep)] = fmt.Sprint(float64(vBeginPrice+vNewBuyPrice-vReducePrice) / 100)
		colE += vBeginPrice + vNewBuyPrice - vReducePrice
		insertData[fmt.Sprint("F", rowStep)] = fmt.Sprint(float64(vReducePrice) / 100)
		colF += vReducePrice
		insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(vPlanDepreciation) / 100)
		colG += vPlanDepreciation
		insertData[fmt.Sprint("H", rowStep)] = fmt.Sprint(float64(vScrapDepreciation) / 100)
		colH += vScrapDepreciation
		insertData[fmt.Sprint("I", rowStep)] = fmt.Sprint(float64(vReducePrice+vPlanDepreciation-vScrapDepreciation) / 100)
		colI += vReducePrice + vPlanDepreciation - vScrapDepreciation
		insertData[fmt.Sprint("N", rowStep)] = fmt.Sprint(float64(vBeginPrice-vReducePrice) / 100)
		colN += vBeginPrice - vReducePrice
		insertData[fmt.Sprint("O", rowStep)] = fmt.Sprint(float64((vBeginPrice+vNewBuyPrice-vReducePrice)-(vReducePrice+vPlanDepreciation-vScrapDepreciation)) / 100)
		colO += (vBeginPrice + vNewBuyPrice - vReducePrice) - (vReducePrice + vPlanDepreciation - vScrapDepreciation)
		//叠加一行
		rowStep++
	}
	//计算总数
	insertData[fmt.Sprint("A", rowStep)] = "合计"
	insertData[fmt.Sprint("B", rowStep)] = fmt.Sprint(float64(colB) / 100)
	insertData[fmt.Sprint("C", rowStep)] = fmt.Sprint(float64(colC) / 100)
	insertData[fmt.Sprint("D", rowStep)] = fmt.Sprint(float64(colD) / 100)
	insertData[fmt.Sprint("E", rowStep)] = fmt.Sprint(float64(colE) / 100)
	insertData[fmt.Sprint("F", rowStep)] = fmt.Sprint(float64(colF) / 100)
	insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(colG) / 100)
	insertData[fmt.Sprint("H", rowStep)] = fmt.Sprint(float64(colH) / 100)
	insertData[fmt.Sprint("I", rowStep)] = fmt.Sprint(float64(colI) / 100)
	insertData[fmt.Sprint("J", rowStep)] = fmt.Sprint(float64(colJ) / 100)
	insertData[fmt.Sprint("K", rowStep)] = fmt.Sprint(float64(colK) / 100)
	insertData[fmt.Sprint("L", rowStep)] = fmt.Sprint(float64(colL) / 100)
	insertData[fmt.Sprint("M", rowStep)] = fmt.Sprint(float64(colM) / 100)
	insertData[fmt.Sprint("N", rowStep)] = fmt.Sprint(float64(colN) / 100)
	insertData[fmt.Sprint("O", rowStep)] = fmt.Sprint(float64(colO) / 100)
	//写入数据
	excelObj.QuickInsertCol("", insertData)
	//设置样式覆盖
	excelObj.ReplaceStyleEnd = fmt.Sprint("O", rowStep)
	//结束处理
	_ = excelObj.Done()
	//反馈成功
	return
}

// ExcelProductData 产品记录表
func ExcelProductData(c any, logErr string, orgID int64, startAt time.Time, endAt time.Time) {
	//分析参数
	startAtCarbon := CoreFilter.GetCarbonByTime(startAt)
	endAtCarbon := CoreFilter.GetCarbonByTime(endAt)
	//预先加载
	excelObj := Router2Excel.ExcelQuick{
		C:                 c,
		LogErr:            logErr,
		FileParams:        fmt.Sprint("erp_permanent_assets_", orgID, "_product_data_", startAtCarbon.Time.Format("2006-01-02"), ".", endAtCarbon.Time.Format("2006-01-02")),
		FileName:          fmt.Sprint("产品卡片表.", startAtCarbon.Time.Format("2006-01-02"), ".", endAtCarbon.Time.Format("2006-01-02"), "期间数据.xlsx"),
		TemplatePath:      fmt.Sprint("erp", CoreFile.Sep, "permanent_assets", CoreFile.Sep, "product_data.xlsx"),
		NeedReplaceStyle:  true,
		ReplaceStyleSheet: "",
		ReplaceStyleRef:   "A1",
		ReplaceStyleStart: "A1",
		ReplaceStyleEnd:   "",
		ExcelObj:          nil,
		CacheSaveTime:     600,
	}
	//初始化
	if excelObj.InitCache() {
		return
	}
	//开始位置
	rowStep := 2
	//预先写入数据
	insertData := map[string]string{}
	//准备合计数据
	var colA, colG, colH, colI, colJ int64
	//获取产品数据集合
	productList := GetProductListByCreateAt(orgID, startAt, endAt)
	for _, vProduct := range productList {
		//写入产品数据
		insertData[fmt.Sprint("A", rowStep)] = vProduct.Code
		colA += 1
		insertData[fmt.Sprint("B", rowStep)] = vProduct.Name
		insertData[fmt.Sprint("C", rowStep)] = Sort.GetNameNoErr(vProduct.SortID)
		insertData[fmt.Sprint("E", rowStep)] = vProduct.CreateAt.Format("2006.01.02")
		insertData[fmt.Sprint("F", rowStep)] = fmt.Sprint(vProduct.UseExpireMonth)
		//获取本月折旧，先获取上月月末数据，然后和本月月末数据对比
		beforeMonthLogData := getLogLastInBetween(vProduct.ID, "check", endAtCarbon.SubMonth().StartOfMonth().Time, endAtCarbon.SubMonth().EndOfMonth().Time)
		nowMonthLogData := getLogLastInBetween(vProduct.ID, "check", endAtCarbon.StartOfMonth().Time, endAtCarbon.EndOfMonth().Time)
		if nowMonthLogData.ID > 0 {
			if beforeMonthLogData.ID > 0 {
				insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(beforeMonthLogData.AllPrice-nowMonthLogData.AllPrice) / 100)
				colG += beforeMonthLogData.AllPrice - nowMonthLogData.AllPrice
			} else {
				beforeMonthLogData = getLogLastIDByProductAndBefore(vProduct.ID, "check", nowMonthLogData.ID)
				insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(beforeMonthLogData.AllPrice-nowMonthLogData.AllPrice) / 100)
				colG += beforeMonthLogData.AllPrice - nowMonthLogData.AllPrice
			}
		} else {
			insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(0.00)
			colG += 0
		}
		//获取最近的盘点记录
		lastCheckLogData := getLogLastByProduct(vProduct.ID, "check")
		if lastCheckLogData.ID > 0 {
			insertData[fmt.Sprint("H", rowStep)] = fmt.Sprint(float64(vProduct.BuyAllPrice-lastCheckLogData.AllPrice) / 100)
			colH += vProduct.BuyAllPrice - lastCheckLogData.AllPrice
		} else {
			insertData[fmt.Sprint("H", rowStep)] = fmt.Sprint(0.00)
			colH += 0
		}
		//继续写入数据
		insertData[fmt.Sprint("I", rowStep)] = fmt.Sprint(float64(vProduct.BuyAllPrice) / 100)
		colI += vProduct.BuyAllPrice
		insertData[fmt.Sprint("J", rowStep)] = fmt.Sprint(float64(vProduct.NowAllPrice) / 100)
		colJ += vProduct.NowAllPrice
		insertData[fmt.Sprint("K", rowStep)] = fmt.Sprint(getProductResidualRate(vProduct.ID))
		insertData[fmt.Sprint("L", rowStep)] = fmt.Sprint(vProduct.PlanUseOrgGroupName)
		//获取类别数据
		vSortData := Sort.GetByIDNoErr(vProduct.SortID, vProduct.OrgID)
		if vSortData.ID > 0 {
			insertData[fmt.Sprint("M", rowStep)] = fmt.Sprint(vSortData.Params.GetValNoErr("Code"))
			insertData[fmt.Sprint("N", rowStep)] = fmt.Sprint(vSortData.Params.GetValNoErr("DepreciationSortName"))
			insertData[fmt.Sprint("O", rowStep)] = fmt.Sprint(vSortData.Params.GetValNoErr("DepreciationSortCode"))
		}
		//继续写入数据
		insertData[fmt.Sprint("P", rowStep)] = vProduct.SavePlace
	}
	//写入合计数据
	insertData[fmt.Sprint("A", rowStep)] = fmt.Sprint("合计：(共计卡片", colA, "张)")
	insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(colG) / 100)
	insertData[fmt.Sprint("H", rowStep)] = fmt.Sprint(float64(colH) / 100)
	insertData[fmt.Sprint("I", rowStep)] = fmt.Sprint(float64(colI) / 100)
	insertData[fmt.Sprint("J", rowStep)] = fmt.Sprint(float64(colJ) / 100)
	//写入数据
	excelObj.QuickInsertCol("", insertData)
	//设置样式覆盖
	excelObj.ReplaceStyleEnd = fmt.Sprint("Q", rowStep)
	//结束处理
	_ = excelObj.Done()
	//反馈成功
	return
}

// ExcelFixedAssetsStatistics 固定资产统计表
func ExcelFixedAssetsStatistics(c any, logErr string, orgID int64) {
	//预先加载
	excelObj := Router2Excel.ExcelQuick{
		C:                 c,
		LogErr:            logErr,
		FileParams:        fmt.Sprint("erp_permanent_assets_", orgID, "_fixed_assets_statistics"),
		FileName:          fmt.Sprint("固定资产统计表.xlsx"),
		TemplatePath:      fmt.Sprint("erp", CoreFile.Sep, "permanent_assets", CoreFile.Sep, "fixed_assets_statistics.xlsx"),
		NeedReplaceStyle:  true,
		ReplaceStyleSheet: "",
		ReplaceStyleRef:   "A1",
		ReplaceStyleStart: "A1",
		ReplaceStyleEnd:   "",
		ExcelObj:          nil,
		CacheSaveTime:     300,
	}
	//初始化
	if excelObj.InitCache() {
		return
	}
	//预备数据
	type sortDataChildData struct {
		//数量
		Num int64
		//使用年限月
		UseMonth int64
		//购买原价总额
		BuyAllPrice int64
		//累计折旧金额
		DepreciationAllPrice int64
		//目前价值总额
		NowAllPrice int64
		//净残值
		ResidualPrice int64
		//本月计提折旧额
		DepreciationPrice int64
	}
	type sortDataChildSort struct {
		//分类ID
		SortID int64
		//分类名称
		SortName string
		//数据集合
		Data sortDataChildData
		//单位
		Unit string
	}
	type sortData struct {
		//部门名称
		OrgGroupName string
		//数据集合
		Data sortDataChildData
		//子分类数据集合
		ChildSort []sortDataChildSort
	}
	var waitData []sortData
	//合计数据
	var colC, colD, colF, colG, colI, colK int64
	//当前时间
	nowAt := CoreFilter.GetNowTimeCarbon()
	//获取所有分类
	sortList, _ := Sort.GetAll(orgID, 0)
	//遍历分类
	for _, vSort := range sortList {
		//获取分类下所有产品
		productList := GetProductListBySortID(orgID, vSort.ID)
		//遍历产品
		for _, vProduct := range productList {
			//找到归属部门
			findGroupKey := -1
			for vGroupKey, vGroup := range waitData {
				if vGroup.OrgGroupName == vProduct.PlanUseOrgGroupName {
					findGroupKey = vGroupKey
					break
				}
			}
			if findGroupKey < 0 {
				waitData = append(waitData, sortData{
					OrgGroupName: vProduct.PlanUseOrgGroupName,
					Data: sortDataChildData{
						Num:                  0,
						UseMonth:             0,
						BuyAllPrice:          0,
						DepreciationAllPrice: 0,
						NowAllPrice:          0,
						ResidualPrice:        0,
						DepreciationPrice:    0,
					},
					ChildSort: []sortDataChildSort{},
				})
				findGroupKey = len(waitData) - 1
			}
			if findGroupKey < 0 {
				continue
			}
			//找到分类
			findSortKey := 0
			for vSortKey, vSort := range waitData[findGroupKey].ChildSort {
				if vSort.SortID == vProduct.SortID {
					findSortKey = vSortKey
					break
				}
			}
			if findSortKey < 0 {
				waitData[findGroupKey].ChildSort = append(waitData[findGroupKey].ChildSort, sortDataChildSort{
					SortID:   vProduct.SortID,
					SortName: vSort.Name,
					Data: sortDataChildData{
						Num:                  0,
						UseMonth:             0,
						BuyAllPrice:          0,
						DepreciationAllPrice: 0,
						NowAllPrice:          0,
						ResidualPrice:        0,
						DepreciationPrice:    0,
					},
					Unit: vSort.Params.GetValNoErr("UnitName"),
				})
				findSortKey = len(waitData[findGroupKey].ChildSort) - 1
			}
			if findSortKey < 0 {
				continue
			}
			//获取本月价值变动日志
			vProductMonthLogData := getLogLastInBetween(vProduct.ID, "check", nowAt.StartOfMonth().Time, nowAt.EndOfMonth().Time)
			vProductMonthBeforeLogData := getLogLastIDByProductAndBefore(vProduct.ID, "check", vProductMonthLogData.ID)
			if vProductMonthBeforeLogData.ID < 1 {
				vProductMonthBeforeLogData.AllPrice = vProduct.BuyAllPrice
			}
			//给部门写入总数据
			waitData[findGroupKey].Data.Num += vProduct.Count
			waitData[findGroupKey].Data.UseMonth += vProduct.UseExpireMonth
			waitData[findGroupKey].Data.BuyAllPrice += vProduct.BuyAllPrice
			waitData[findGroupKey].Data.DepreciationAllPrice += vProduct.BuyAllPrice - vProduct.NowAllPrice
			waitData[findGroupKey].Data.NowAllPrice += vProduct.NowAllPrice
			waitData[findGroupKey].Data.ResidualPrice += int64(float64(vProduct.BuyAllPrice) * 0.05)
			waitData[findGroupKey].Data.DepreciationPrice += vProductMonthBeforeLogData.AllPrice - vProductMonthLogData.AllPrice
			//写入分类数据
			waitData[findGroupKey].ChildSort[findSortKey].Data.Num += vProduct.Count
			waitData[findGroupKey].ChildSort[findSortKey].Data.UseMonth += vProduct.UseExpireMonth
			waitData[findGroupKey].ChildSort[findSortKey].Data.BuyAllPrice += vProduct.BuyAllPrice
			waitData[findGroupKey].ChildSort[findSortKey].Data.DepreciationAllPrice += vProduct.BuyAllPrice - vProduct.NowAllPrice
			waitData[findGroupKey].ChildSort[findSortKey].Data.NowAllPrice += vProduct.NowAllPrice
			waitData[findGroupKey].ChildSort[findSortKey].Data.ResidualPrice += int64(float64(vProduct.BuyAllPrice) * 0.05)
			waitData[findGroupKey].ChildSort[findSortKey].Data.DepreciationPrice += vProductMonthBeforeLogData.AllPrice - vProductMonthLogData.AllPrice
			//总统计数据
			colC += vProduct.Count
			colD += vProduct.UseExpireMonth
			colF += vProduct.BuyAllPrice
			colG += vProduct.BuyAllPrice - vProduct.NowAllPrice
			colI += vProduct.NowAllPrice
			colK += int64(float64(vProduct.BuyAllPrice) * 0.05)
		}
	}
	//开始位置
	rowStep := 2
	//预先写入数据
	insertData := map[string]string{}
	//遍历构建数据
	for _, vGroup := range waitData {
		//写入部门基本数据
		insertData[fmt.Sprint("A", rowStep)] = vGroup.OrgGroupName
		insertData[fmt.Sprint("C", rowStep)] = fmt.Sprint(vGroup.Data.Num)
		insertData[fmt.Sprint("F", rowStep)] = fmt.Sprint(float64(vGroup.Data.BuyAllPrice) / 100)
		insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(vGroup.Data.DepreciationAllPrice) / 100)
		insertData[fmt.Sprint("H", rowStep)] = fmt.Sprint("0.00")
		insertData[fmt.Sprint("I", rowStep)] = fmt.Sprint(float64(vGroup.Data.NowAllPrice) / 100)
		insertData[fmt.Sprint("J", rowStep)] = fmt.Sprint("0.00")
		insertData[fmt.Sprint("K", rowStep)] = fmt.Sprint(float64(vGroup.Data.ResidualPrice) / 100)
		insertData[fmt.Sprint("L", rowStep)] = fmt.Sprint((float64(vGroup.Data.BuyAllPrice-vGroup.Data.ResidualPrice) / 5 / 12) / 100)
		insertData[fmt.Sprint("M", rowStep)] = fmt.Sprint("0.00")
		insertData[fmt.Sprint("N", rowStep)] = fmt.Sprint("0.00")
		insertData[fmt.Sprint("O", rowStep)] = fmt.Sprint("0.00")
		insertData[fmt.Sprint("P", rowStep)] = fmt.Sprint("0.00")
		insertData[fmt.Sprint("Q", rowStep)] = fmt.Sprint("0.00")
		insertData[fmt.Sprint("R", rowStep)] = fmt.Sprint("0.00")
		//进一行
		rowStep++
		//遍历分类数据
		for _, vSort := range vGroup.ChildSort {
			//写入分类数据
			insertData[fmt.Sprint("B", rowStep)] = vSort.SortName
			insertData[fmt.Sprint("C", rowStep)] = fmt.Sprint(vSort.Data.Num)
			insertData[fmt.Sprint("D", rowStep)] = fmt.Sprint(vSort.Data.UseMonth)
			insertData[fmt.Sprint("E", rowStep)] = vSort.Unit
			insertData[fmt.Sprint("F", rowStep)] = fmt.Sprint(float64(vSort.Data.BuyAllPrice) / 100)
			insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(vSort.Data.DepreciationAllPrice) / 100)
			insertData[fmt.Sprint("I", rowStep)] = fmt.Sprint(float64(vSort.Data.NowAllPrice) / 100)
			insertData[fmt.Sprint("J", rowStep)] = fmt.Sprint("0.00")
			insertData[fmt.Sprint("K", rowStep)] = fmt.Sprint(float64(vSort.Data.ResidualPrice) / 100)
			insertData[fmt.Sprint("L", rowStep)] = fmt.Sprint((float64(vSort.Data.BuyAllPrice-vSort.Data.ResidualPrice) / 5 / 12) / 100)
			insertData[fmt.Sprint("M", rowStep)] = fmt.Sprint("0.00")
			insertData[fmt.Sprint("N", rowStep)] = fmt.Sprint("0.00")
			insertData[fmt.Sprint("O", rowStep)] = fmt.Sprint("0.00")
			insertData[fmt.Sprint("P", rowStep)] = fmt.Sprint("0.00")
			insertData[fmt.Sprint("Q", rowStep)] = fmt.Sprint("0.00")
			insertData[fmt.Sprint("R", rowStep)] = fmt.Sprint("0.00")
			//进一行
			rowStep++
		}
		if len(vGroup.ChildSort) < 1 {
			//进一行
			rowStep++
		}
	}
	//构建总数据
	insertData[fmt.Sprint("A", rowStep)] = "总计"
	insertData[fmt.Sprint("C", rowStep)] = fmt.Sprint(colC)
	insertData[fmt.Sprint("D", rowStep)] = fmt.Sprint(colD)
	insertData[fmt.Sprint("F", rowStep)] = fmt.Sprint(float64(colF) / 100)
	insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(colG) / 100)
	insertData[fmt.Sprint("H", rowStep)] = fmt.Sprint("0.00")
	insertData[fmt.Sprint("I", rowStep)] = fmt.Sprint(float64(colI) / 100)
	insertData[fmt.Sprint("J", rowStep)] = fmt.Sprint("0.00")
	insertData[fmt.Sprint("K", rowStep)] = fmt.Sprint(float64(colK) / 100)
	insertData[fmt.Sprint("L", rowStep)] = fmt.Sprint((float64(colF-colK) / 5 / 12) / 100)
	insertData[fmt.Sprint("M", rowStep)] = fmt.Sprint("0.00")
	insertData[fmt.Sprint("N", rowStep)] = fmt.Sprint("0.00")
	insertData[fmt.Sprint("O", rowStep)] = fmt.Sprint("0.00")
	insertData[fmt.Sprint("P", rowStep)] = fmt.Sprint("0.00")
	insertData[fmt.Sprint("Q", rowStep)] = fmt.Sprint("0.00")
	insertData[fmt.Sprint("R", rowStep)] = fmt.Sprint("0.00")
	//写入数据
	excelObj.QuickInsertCol("", insertData)
	//设置样式覆盖
	excelObj.ReplaceStyleEnd = fmt.Sprint("R", rowStep)
	//结束处理
	_ = excelObj.Done()
	//反馈成功
	return
}

// ExcelDepreciationAllocation 折旧分配表
func ExcelDepreciationAllocation(c any, logErr string, orgID int64, startAt time.Time, endAt time.Time) {
	//分析时间
	startAtCarbon := CoreFilter.GetCarbonByTime(startAt)
	endAtCarbon := CoreFilter.GetCarbonByTime(endAt)
	//预先加载
	excelObj := Router2Excel.ExcelQuick{
		C:                 c,
		LogErr:            logErr,
		FileParams:        fmt.Sprint("erp_permanent_assets_", orgID, "_depreciation_allocation_", startAtCarbon.Time.Format("2006-01-02"), ".", endAtCarbon.Time.Format("2006-01-02")),
		FileName:          fmt.Sprint("折旧分配表.", startAtCarbon.Time.Format("20060102"), "_", endAtCarbon.Time.Format("20060102"), ".xlsx"),
		TemplatePath:      fmt.Sprint("erp", CoreFile.Sep, "permanent_assets", CoreFile.Sep, "depreciation_allocation.xlsx"),
		NeedReplaceStyle:  true,
		ReplaceStyleSheet: "",
		ReplaceStyleRef:   "",
		ReplaceStyleStart: "A1",
		ReplaceStyleEnd:   "",
		ExcelObj:          nil,
		CacheSaveTime:     300,
	}
	//初始化
	if excelObj.InitCache() {
		return
	}
	//构建数据
	type groupData struct {
		//部门名称
		OrgGroupName string
		//科目编码
		SubjectCode string
		//科目名称
		SubjectName string
		//折旧额度
		DepreciationPrice int64
	}
	var waitData []groupData
	//记录遍历过的产品
	type productData struct {
		//产品ID
		ProductID int64
		//创建时间
		CreatedAt time.Time
		//上次记录递减价值
		LastDepreciationPrice int64
	}
	var waitProductData []productData
	//获取本月发生折旧的记录
	logList := getLogListByOrgAndMode(orgID, "check", startAt, endAt)
	for _, vLog := range logList {
		//获取对应产品数据
		vProductData := getProductByID(vLog.ProductID)
		if vProductData.ID < 1 {
			continue
		}
		//获取上一次记录
		vBeforeLogData := getLogLastIDByProductAndBefore(vProductData.ID, "check", vLog.ID)
		if vBeforeLogData.ID < 1 {
			vBeforeLogData.AllPrice = vProductData.BuyAllPrice
		}
		//折旧情况
		vDepreciationPrice := vBeforeLogData.AllPrice - vLog.AllPrice
		//找到对应的部门
		findOrgGroupKey := -1
		for vGroupKey, vGroup := range waitData {
			if vGroup.OrgGroupName == vProductData.PlanUseOrgGroupName {
				findOrgGroupKey = vGroupKey
				break
			}
		}
		if findOrgGroupKey > 0 {
			//检查是否存在相同的产品被记录过，并检查时间是否超过上次记录的值
			isFindProduct := false
			for _, vProduct := range waitProductData {
				if vProduct.ProductID == vLog.ProductID {
					if vProduct.CreatedAt.After(vLog.CreateAt) {
						//递减掉金额记录
						waitData[findOrgGroupKey].DepreciationPrice -= vProduct.LastDepreciationPrice
						waitData[findOrgGroupKey].DepreciationPrice += vDepreciationPrice
					}
					isFindProduct = true
					break
				}
			}
			if isFindProduct {
				continue
			}
			//叠加折旧金额
			waitData[findOrgGroupKey].DepreciationPrice += vDepreciationPrice
		}
		//如果不存在，则新增记录
		if findOrgGroupKey < 1 {
			vSortData := Sort.GetByIDNoErr(vProductData.SortID, vProductData.OrgID)
			waitData = append(waitData, groupData{
				OrgGroupName:      vProductData.PlanUseOrgGroupName,
				SubjectCode:       vSortData.Params.GetValNoErr("DepreciationSortCode"),
				SubjectName:       vSortData.Params.GetValNoErr("DepreciationSortName"),
				DepreciationPrice: vDepreciationPrice,
			})
			findOrgGroupKey = len(waitData) - 1
		}
	}
	//开始位置
	rowStep := 2
	//预先写入数据
	insertData := map[string]string{}
	//合计数据
	var colG int64
	//写入数据
	for _, vGroup := range waitData {
		//写入数据
		insertData[fmt.Sprint("A", rowStep)] = ""
		insertData[fmt.Sprint("B", rowStep)] = vGroup.OrgGroupName
		insertData[fmt.Sprint("C", rowStep)] = ""
		insertData[fmt.Sprint("D", rowStep)] = ""
		insertData[fmt.Sprint("E", rowStep)] = vGroup.SubjectCode
		insertData[fmt.Sprint("F", rowStep)] = vGroup.SubjectName
		insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(vGroup.DepreciationPrice) / 100)
		//合计数据
		colG += vGroup.DepreciationPrice
		//进一步
		rowStep++
	}
	//写入合计数
	insertData[fmt.Sprint("A", rowStep)] = "合计"
	insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(float64(colG) / 100)
	//写入数据
	excelObj.QuickInsertCol("", insertData)
	//设置样式覆盖
	excelObj.ReplaceStyleEnd = fmt.Sprint("G", rowStep)
	//结束处理
	_ = excelObj.Done()
	//反馈成功
	return
}
