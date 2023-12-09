package ERPPermanentAssets

import (
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL2 "gitee.com/weeekj/weeekj_core/v5/core/sql2"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	"math"
)

// 自动折旧定时任务
func subNatsAutoExpire() {
	//日志
	logAppend := "sub nats erp permanent assets auto expire, "
	//捕捉异常
	defer func() {
		//收尾处理
		runAutoExpireSysM.Bind.UpdateNextAtFutureHour(5, 0, 0)
		//跳出处理
		if r := recover(); r != nil {
			runAutoExpireSysM.Update(fmt.Sprint("发生错误: ", r), "run.error", 0)
			CoreLog.Error(logAppend, r)
		}
	}()
	//获取总产品数量
	productCount := productSQL.Analysis().Count("delete_at < to_timestamp(1000000)")
	//跟踪器
	runAutoExpireSysM.Start("开始计算折旧", "start", productCount)
	//遍历所有组织
	var page int64 = 1
	for {
		orgList := OrgCore.GetOrgListStep(page, 100, []string{"only"})
		if len(orgList) < 1 {
			break
		}
		for _, vOrg := range orgList {
			vOrgConfigAutoExpire := OrgCore.Config.GetConfigValBoolNoErr(vOrg.ID, "ERPPermanentAssetsAutoExpire")
			if !vOrgConfigAutoExpire {
				continue
			}
			//获取所有产品
			var pageProduct int64 = 1
			var productRawList []FieldsProduct
			for {
				_ = productSQL.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id"}).SetPages(CoreSQL2.ArgsPages{
					Page: pageProduct,
					Max:  1000,
					Sort: "id",
					Desc: false,
				}).SelectList("delete_at < to_timestamp(1000000) AND org_id = $1 AND wait_check_at < NOW()", vOrg.ID).Result(&productRawList)
				if len(productRawList) < 1 {
					break
				}
				for _, vProduct := range productRawList {
					subNatsAutoExpireChildProduct(logAppend, vOrg.ID, vProduct.ID, true, true, true)
				}
				pageProduct += 1
			}
		}
		page += 1
	}
	//收尾处理
	//runAutoExpireSysM.Bind.UpdateNextAt(CoreFilter.GetNowTimeCarbon().AddSeconds(2).Time)
	runAutoExpireSysM.Finish()
}

func subNatsAutoExpireChildProduct(logAppend string, orgID int64, productID int64, isAutoRun bool, isCheckMonth bool, isCheckLog bool) {
	if isAutoRun {
		//进度跟踪
		runAutoExpireSysM.Update("整理组织下所有库存产品", fmt.Sprint("org.", orgID, ".product.", productID), 1)
	}
	//获取该产品的扩展参数，折旧率
	vOrgConfigAutoExpirePFloat64ByProduct := getProductResidualRate(productID)
	//检查本月是否存在盘点记录？
	if isCheckLog {
		vLog := getLogLastByProduct(productID, "check")
		if vLog.ID > 0 && CoreFilter.GetCarbonByTime(vLog.CreateAt).Gt(CoreFilter.GetNowTimeCarbon().SubMonth()) {
			return
		}
	}
	//获取元数据
	vProductData := getProductByID(productID)
	//当月创建的产品不进行折旧
	if isCheckMonth {
		if vProductData.CreateAt.Unix() >= CoreFilter.GetNowTimeCarbon().StartOfMonth().Time.Unix() {
			return
		}
	}
	//客户提供的计算方式：（原值 - 残值）/使用年限/12月=每月计提折旧
	// 其中，残值=5%
	// 本计划方案采用晋中市保障房公司提供的计算形式，未来可以增加其他模式计算，需拆分处理。
	// 转化为程序的思路：
	/**
	1. 获取当前使用年限
	2. 按照公式，计算“原值 - 残值”部分；然后分别除尽即可得出每月计提折旧
	3. 在当前产品每个的价值上递减1次即可
	4. 计算产品的总价值
	*/
	//检查使用年限，必须大于1，否则拒绝计算
	if vProductData.UseExpireYear < 1 {
		return
	}
	//计算（原值 - 残值）
	// 原值=产品当前价值
	// 残值=原值x5%
	vD1 := float64(vProductData.NowPerPrice)
	vD2 := float64(vProductData.NowPerPrice) * vOrgConfigAutoExpirePFloat64ByProduct
	vD3 := vD1 - vD2
	//计算每月计提折旧
	vDResult := vD3 / float64(vProductData.UseExpireYear) / 12
	//计算新的价值
	vResultPer := vProductData.NowPerPrice - int64(math.Round(vDResult))
	vResultAll := vResultPer * vProductData.Count
	//赋予新的值
	_, err := CreateCheck(&ArgsCreateCheck{
		CreateAt:  CoreFilter.GetNowTime(),
		OrgID:     vProductData.OrgID,
		OrgBindID: 0,
		ProductID: vProductData.ID,
		AllPrice:  vResultAll,
		PerPrice:  vResultPer,
		Count:     vProductData.Count,
		SavePlace: vProductData.SavePlace,
		Des:       "自动盘点资产",
		Params:    nil,
	})
	if err != nil {
		CoreLog.Warn(logAppend, "create check product, id: ", vProductData.ID, ", err: ", err)
		return
	}
}
