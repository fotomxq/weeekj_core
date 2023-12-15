package ServiceUserInfo

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetTemplateList 获取模板列表参数
type ArgsGetTemplateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTemplateList 获取模板列表
func GetTemplateList(args *ArgsGetTemplateList) (dataList []FieldsTemplate, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.SortID > 0 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_user_info_template"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, des, cover_files, sort_id, tags, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetTemplateMore 获取指定ID组模板参数
type ArgsGetTemplateMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetTemplateMore 获取指定ID组模板
func GetTemplateMore(args *ArgsGetTemplateMore) (dataList []FieldsTemplate, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "service_user_info_template", "id, create_at, update_at, delete_at, org_id, name, des, cover_files, sort_id, tags, params", args.IDs, args.HaveRemove)
	return
}

// ArgsGetTemplateID 获取指定ID模板参数
type ArgsGetTemplateID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetTemplateID 获取指定ID模板
func GetTemplateID(args *ArgsGetTemplateID) (data FieldsTemplate, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, name, des, cover_files, sort_id, tags, file_data, params FROM service_user_info_template where id = $1 AND delete_at < to_timestamp(1000000) AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsCreateTemplate 创建新的模板参数
type ArgsCreateTemplate struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//封面图
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//文件数据
	FileData string `db:"file_data" json:"fileData"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateTemplate 创建新的模板
func CreateTemplate(args *ArgsCreateTemplate) (data FieldsTemplate, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_user_info_template", "INSERT INTO service_user_info_template (org_id, name, des, cover_files, sort_id, tags, file_data, params) VALUES (:org_id,:name,:des,:cover_files,:sort_id,:tags,:file_data,:params)", args, &data)
	return
}

// ArgsUpdateTemplate 修改模板参数
type ArgsUpdateTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//封面图
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//文件数据
	FileData string `db:"file_data" json:"fileData"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateTemplate 修改模板
func UpdateTemplate(args *ArgsUpdateTemplate) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_user_info_template SET name = :name, des = :des, cover_files = :cover_files, sort_id = :sort_id, tags = :tags, file_data = :file_data, params = :params WHERE id = :id AND org_id = :org_id", args)
	return
}

// ArgsDeleteTemplate 删除模板参数
type ArgsDeleteTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteTemplate 删除模板
func DeleteTemplate(args *ArgsDeleteTemplate) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_user_info_template", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
