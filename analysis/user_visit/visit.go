package AnalysisUserVisit

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetVisitList 获取列表参数
type ArgsGetVisitList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//数据来源
	// 来自哪个模块
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//关联的用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//挖掘的电话号码
	Phone string `db:"phone" json:"phone" check:"phone" empty:"true"`
	//IP地址
	IP string `db:"ip" json:"ip" check:"ip" empty:"true"`
	//行为标记
	// insert 进入; out 离开; move 移动
	Action string `db:"action" json:"action" check:"mark" empty:"true"`
	//浏览器标识
	// 或设备标识
	Mark string `db:"mark" json:"mark"`
}

// GetVisitList 获取列表
func GetVisitList(args *ArgsGetVisitList) (dataList []FieldsVisit, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	where, maps, err = args.CreateInfo.GetListAnd("create_info", "create_info", where, maps)
	if err != nil {
		return
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.Country > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "country = :country"
		maps["country"] = args.Country
	}
	if args.Phone != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "phone = :phone"
		maps["phone"] = args.Phone
	}
	if args.IP != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ip = :ip"
		maps["ip"] = args.IP
	}
	if args.Mark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.Action != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "action = :action"
		maps["action"] = args.Action
	}
	if where == "" {
		where = "true"
	}
	tableName := "analysis_user_visit"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, create_info, user_id, country, phone, ip, mark, action, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsVisitCreate 创建新的访问数据参数
type ArgsVisitCreate struct {
	//组织ID
	// 如果存在数据，则表明该数据隶属于指定组织
	// 组织依可查看该数据
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//数据来源
	// 来自哪个模块
	// system: public 公共渠道，非模块内部创建数据
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo" check:"createInfo"`
	//关联的用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//挖掘的电话号码
	Phone string `db:"phone" json:"phone" check:"phone" empty:"true"`
	//IP地址
	IP string `db:"ip" json:"ip" check:"ip"`
	//浏览器标识
	// 或设备标识
	Mark string `db:"mark" json:"mark"`
	//行为标记
	// insert 进入; out 离开; move 移动
	Action string `db:"action" json:"action" check:"mark"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params"`
}

// VisitCreate 创建新的访问数据
func VisitCreate(args *ArgsVisitCreate) (err error) {
	if len(args.Mark) > 255 {
		args.Mark = args.Mark[0:254]
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_user_visit (org_id, create_info, user_id, country, phone, ip, mark, action, params) VALUES (:org_id,:create_info,:user_id,:country,:phone,:ip,:mark,:action,:params)", args)
	return
}

// ArgsDeleteVisitByUser 删除指定用户的数据参数
type ArgsDeleteVisitByUser struct {
	//关联的用户
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// DeleteVisitByUser 删除指定用户的数据
func DeleteVisitByUser(args *ArgsDeleteVisitByUser) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "analysis_user_visit", "user_id = :user_id", args)
	return
}
