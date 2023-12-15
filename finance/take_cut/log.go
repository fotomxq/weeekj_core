package FinanceTakeCut

import (
	"errors"
	"fmt"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取日志参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
}

// GetLogList 获取日志
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if where == "" {
		where = "true"
	}
	tableName := "finance_take_cut_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, order_price, cut_price, cut_price_proportion, order_id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetLogByOrderID 查询指定订单抽成参数
type ArgsGetLogByOrderID struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联的订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
}

// GetLogByOrderID 查询指定订单抽成
func GetLogByOrderID(args *ArgsGetLogByOrderID) (data FieldsLog, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, order_price, cut_price, cut_price_proportion, order_id FROM finance_take_cut_log WHERE ($1 < 1 OR org_id = $1) AND order_id = $2", args.OrgID, args.OrderID)
	if err == nil && data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsAddLog 添加抽成参数
type ArgsAddLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//针对订单的系统来源
	// eg: user_sub / org_sub / mall
	OrderSystem string `db:"order_system" json:"orderSystem" check:"mark"`
	//订单金额
	OrderPrice int64 `db:"order_price" json:"orderPrice" check:"price"`
	//关联的订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
}

// AddLog 添加抽成
// 本模块同时将完成抽成处理
// price 反馈抽取的金额
func AddLog(args *ArgsAddLog) (price int64, err error) {
	//检查金额
	if args.OrderPrice < 1 {
		return
	}
	//检查订单是否发生过转账
	var logID int64
	err = Router2SystemConfig.MainDB.Get(&logID, "SELECT id FROM finance_take_cut_log WHERE order_id = $1", args.OrderID)
	if err == nil && logID > 0 {
		err = errors.New("replace data")
		return
	}
	//获取抽成设计
	var configData FieldsConfig
	err = Router2SystemConfig.MainDB.Get(&configData, "SELECT id, cut_price_proportion FROM finance_take_cut_config WHERE org_id = $1 AND order_system = $2", args.OrgID, args.OrderSystem)
	if err != nil || configData.ID < 1 {
		var orgData OrgCoreCore.FieldsOrg
		orgData, err = OrgCoreCore.GetOrg(&OrgCoreCore.ArgsGetOrg{
			ID: args.OrgID,
		})
		if err != nil {
			return
		}
		err = Router2SystemConfig.MainDB.Get(&configData, "SELECT id, cut_price_proportion FROM finance_take_cut_config WHERE sort_id = $1 AND order_system = $2", orgData.SortID, args.OrderSystem)
		if err != nil || configData.ID < 1 {
			err = errors.New("no config data")
			return
		}
	}
	//计算扣除金额
	cutPrice := int64(float64(args.OrderPrice) * (float64(configData.CutPriceProportion) / 10000000))
	//fmt.Println(fmt.Sprint("cutPrice: ", cutPrice, ", args.OrderPrice: ", args.OrderPrice, ", configData.CutPriceProportion: ", configData.CutPriceProportion))
	if cutPrice < 1 {
		return
	}
	//扣除商户等值金额
	var takeChannelMark string
	takeChannelMark, err = OrgCoreCore.Config.GetConfigVal(&ClassConfig.ArgsGetConfig{
		BindID:    args.OrgID,
		Mark:      "FinanceDepositDefaultMark",
		VisitType: "admin",
	})
	if err != nil {
		err = errors.New("get org deposit mark config, " + err.Error())
		return
	}
	_, _, err = FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
		UpdateHash: "",
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     args.OrgID,
			Mark:   "",
			Name:   "",
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     args.OrgID,
			Mark:   "",
			Name:   "",
		},
		ConfigMark:      takeChannelMark,
		AppendSavePrice: 0 - cutPrice,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("set finance deposit by org, ", err))
		return
	}
	//扣除商户的储蓄金额
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_take_cut_log (org_id, order_price, cut_price, cut_price_proportion, order_id) VALUES (:org_id,:order_price,:cut_price,:cut_price_proportion,:order_id)", map[string]interface{}{
		"org_id":               args.OrgID,
		"order_price":          args.OrderPrice,
		"cut_price":            cutPrice,
		"cut_price_proportion": configData.CutPriceProportion,
		"order_id":             args.OrderID,
	})
	//反馈
	return
}
