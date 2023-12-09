package ServiceUserInfo

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//修改的位置
	ChangeMark string `db:"change_mark" json:"changeMark" check:"mark" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.InfoID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "info_id = :info_id"
		maps["info_id"] = args.InfoID
	}
	if args.ChangeMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "change_mark = :change_mark"
		maps["change_mark"] = args.ChangeMark
	}
	if args.Search != "" {
		where = where + " AND (change_des ILIKE '%' || :search || '%' OR old_des ILIKE '%' || :search || '%' OR new_des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_user_info_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, info_id, change_mark, change_des, old_des, new_des FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// 添加日志
type argsAppendLog struct {
	//档案ID
	InfoID int64 `db:"info_id" json:"infoID"`
	//组织ID
	// 允许平台方的0数据，该数据可能来源于其他领域
	OrgID int64 `db:"org_id" json:"orgID"`
	//修改的位置
	// 1. 字段
	// 2. 或扩展参数指定的内容，例如params.[mark]
	// 3. 其他内容采用.形式跨越记录
	// 4. room.in 入驻房间变更
	ChangeMark string `db:"change_mark" json:"changeMark"`
	ChangeDes  string `db:"change_des" json:"changeDes"`
	//修改前描述
	OldDes string `db:"old_des" json:"oldDes"`
	//修改后描述
	NewDes string `db:"new_des" json:"newDes"`
}

func appendLog(args *argsAppendLog) {
	if args.OrgID < 1 {
		infoData := getInfoID(args.InfoID)
		args.OrgID = infoData.OrgID
	}
	_, err := CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_user_info_log (org_id, info_id, change_mark, change_des, old_des, new_des) VALUES (:org_id, :info_id, :change_mark, :change_des, :old_des, :new_des)", args)
	if err != nil {
		CoreLog.Error("service user info append log, ", err)
	}
}
