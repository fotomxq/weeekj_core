package FinanceAssets

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

//操作记录

// 查看日志列表
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//实际操作人，组织绑定成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//是否为历史
	IsHistory bool `json:"isHistory"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.UserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.ProductID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "product_id = :product_id"
		maps["product_id"] = args.ProductID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if args.IsHistory {
		dataCount, err = CoreSQL.GetListPageAndCount(
			Router2SystemConfig.MainDB.DB,
			&dataList,
			"finance_assets_log_history",
			"id",
			"SELECT id, create_at, org_id, bind_id, user_id, product_id, count, des FROM finance_assets_log_history WHERE "+where,
			where,
			maps,
			&args.Pages,
			[]string{"id", "create_at", "count"},
		)
	} else {
		dataCount, err = CoreSQL.GetListPageAndCount(
			Router2SystemConfig.MainDB.DB,
			&dataList,
			"finance_assets_log",
			"id",
			"SELECT id, create_at, org_id, bind_id, user_id, product_id, count, des FROM finance_assets_log WHERE "+where,
			where,
			maps,
			&args.Pages,
			[]string{"id", "create_at", "count"},
		)
	}
	return
}

// 创建新的日志
type argsAppendLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//实际操作人，组织绑定成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//变动数量
	// 可以是正负数
	Count int64 `db:"count" json:"count"`
	//变动原因
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
}

func appendLog(args *argsAppendLog) (err error) {
	//检查产品是否是该组织的？
	err = checkProductAndOrg(args.ProductID, args.OrgID)
	if err != nil {
		return
	}
	//创建申请
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_assets_log (org_id, bind_id, user_id, product_id, count, des) VALUES (:org_id, :bind_id, :user_id, :product_id, :count, :des)", args)
	return
}
