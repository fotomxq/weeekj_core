package ServiceUserInfo

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"strings"
)

// ArgsGetDocList 获取文件列表参数
type ArgsGetDocList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//采用模版
	TemplateID int64 `db:"template_id" json:"templateID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetDocList 获取文件列表
func GetDocList(args *ArgsGetDocList) (dataList []FieldsDoc, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.InfoID > -1 {
		where = where + " AND info_id = :info_id"
		maps["info_id"] = args.InfoID
	}
	if args.TemplateID > -1 {
		where = where + " AND template_id = :template_id"
		maps["template_id"] = args.TemplateID
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
		args.Search = strings.ReplaceAll(args.Search, " ", "")
		where = where + " AND (REPLACE(title, ' ', '') ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_user_info_doc"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, info_id, title, sort_id, tags, template_id, params "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetDocID 获取ID参数
type ArgsGetDocID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetDocID 获取ID
func GetDocID(args *ArgsGetDocID) (data FieldsDoc, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, info_id, title, sort_id, tags, template_id, file_data, params FROM service_user_info_doc WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	return
}

// ArgsGetDocMore 获取多个信息ID参数
type ArgsGetDocMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetDocMore 获取多个信息ID
func GetDocMore(args *ArgsGetDocMore) (dataList []FieldsDoc, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "service_user_info_doc", "id, create_at, update_at, delete_at, org_id, info_id, title, sort_id, tags, template_id, params", args.IDs, args.HaveRemove)
	return
}

func GetDocMoreNames(args *ArgsGetDocMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsTitleAndDelete("service_user_info_doc", args.IDs, args.HaveRemove)
	return
}

// ArgsGetOrgDocMore 获取多个信息ID带组织参数
type ArgsGetOrgDocMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetOrgDocMore 获取多个信息ID带组织
func GetOrgDocMore(args *ArgsGetOrgDocMore) (dataList []FieldsDoc, err error) {
	err = CoreSQLIDs.GetIDsOrgAndDelete(&dataList, "service_user_info_doc", "id, create_at, update_at, delete_at, org_id, info_id, title, sort_id, tags, template_id, file_data, params", args.IDs, args.OrgID, args.HaveRemove)
	return
}

func GetOrgDocMoreNames(args *ArgsGetOrgDocMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsOrgTitleAndDelete("service_user_info_doc", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// ArgsCheckDoc 检查数据有效性参数
type ArgsCheckDoc struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// CheckDoc 检查数据有效性
func CheckDoc(args *ArgsCheckDoc) (err error) {
	var data FieldsDoc
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_user_info_doc WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	if err != nil || data.ID < 1 {
		err = errors.New(fmt.Sprint("data not exist, ", err))
		return
	}
	return
}

// ArgsCreateDoc 创建新的文件参数
type ArgsCreateDoc struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//关联档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//模版ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//文件数据
	FileData string `db:"file_data" json:"fileData"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateDoc 创建新的文件
func CreateDoc(args *ArgsCreateDoc) (data FieldsDoc, err error) {
	if args.InfoID > 0 {
		//获取信息档案
		var infoData FieldsInfo
		err = Router2SystemConfig.MainDB.Get(&infoData, "SELECT id FROM service_user_info WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", args.InfoID, args.OrgID)
		if err != nil || infoData.ID < 1 {
			err = errors.New(fmt.Sprint("bind not exist, ", err))
			return
		}
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_user_info_doc", "INSERT INTO service_user_info_doc (org_id, info_id, title, sort_id, tags, template_id, file_data, params) VALUES (:org_id,:info_id,:title,:sort_id,:tags,:template_id,:file_data,:params)", args, &data)
	if err != nil {
		err = errors.New(fmt.Sprint("create new file, ", err))
		return
	}
	return
}

// ArgsUpdateDoc 修改文件信息参数
type ArgsUpdateDoc struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//模版ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id" empty:"true"`
	//文件数据
	FileData string `db:"file_data" json:"fileData"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateDoc 修改文件信息
func UpdateDoc(args *ArgsUpdateDoc) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_user_info_doc SET update_at = NOW(), title = :title, sort_id = :sort_id, tags = :tags, template_id = :template_id, file_data = :file_data, params = :params WHERE id = :id AND org_id = :org_id", args)
	return
}

// ArgsDeleteDoc 删除文档
type ArgsDeleteDoc struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteDoc 删除文档
func DeleteDoc(args *ArgsDeleteDoc) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_user_info_doc", "id = :id AND org_id = :org_id", args)
	return
}
