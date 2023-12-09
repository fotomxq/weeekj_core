package ERPDocument

import (
	"errors"
	CoreCache "gitee.com/weeekj/weeekj_core/v5/core/cache"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ERPCore "gitee.com/weeekj/weeekj_core/v5/erp/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"strings"
)

// ArgsGetDocList 获取配置列表参数
type ArgsGetDocList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetDocList 获取配置列表参
func GetDocList(args *ArgsGetDocList) (dataList []FieldsDoc, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ConfigID > -1 {
		//检查是否发布
		if !checkConfigPublish(args.ConfigID, args.OrgID) {
			err = errors.New("config no publish")
			return
		}
		//组装条件
		where = where + " AND config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.Search != "" {
		args.Search = strings.ReplaceAll(args.Search, " ", "")
		where = where + " AND (REPLACE(name, ' ', '') ILIKE '%' || :search || '%' OR REPLACE(des, ' ', '') ILIKE '%' || :search || '%' OR REPLACE(search_des, ' ', '') ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "erp_document_doc"
	var rawList []FieldsDoc
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getDocByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetDoc 获取指定的配置
func GetDoc(id int64, orgID int64) (data FieldsDoc) {
	data = getDocByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		data = FieldsDoc{}
		return
	}
	return
}

// GetDocAllVal 获取文档组件列
func GetDocAllVal(docID int64) (dataList []ERPCore.FieldsComponentVal) {
	return docComponentValObj.GetAllVal(docID)
}

// GetDocName 获取配置名称
//func GetDocName(id int64) (name string) {
//	data := getDocByID(id)
//	return data.Name
//}

// getDocByID 获取配置
func getDocByID(id int64) (data FieldsDoc) {
	cacheMark := getDocCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, config_id, org_id, name, des, cover_file_id, params FROM erp_document_doc WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime3Day)
	return
}
