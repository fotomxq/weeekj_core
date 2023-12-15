package BasePageStyle

import (
	"encoding/json"
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetComponentList 获取组件列表参数
type ArgsGetComponentList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//系统
	System string `db:"system" json:"system" check:"mark" empty:"true"`
	//关联标识码
	// 全局唯一
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetComponentList 获取组件列表
func GetComponentList(args *ArgsGetComponentList) (dataList []FieldsComponent, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.System != "" {
		where = where + " AND system = :system"
		maps["system"] = args.System
	}
	if args.Mark != "" {
		where = where + " AND mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.Search != "" {
		where = where + " AND (mark ILIKE '%' || :search || '%' OR name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "core_page_style_component"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, system, mark, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetComponentIDs 通过IDs获取一组组件参数
type ArgsGetComponentIDs struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetComponentIDs 通过IDs获取一组组件
func GetComponentIDs(args *ArgsGetComponentIDs) (dataList []FieldsComponent, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "core_page_style_component", "id, create_at, update_at, delete_at, system, mark, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list, data, params", args.IDs, args.HaveRemove)
	return
}

// ArgsGetComponentMarks 通过一组标识码获取一组组件参数
type ArgsGetComponentMarks struct {
	//系统
	System string `db:"system" json:"system" check:"mark"`
	//标识码列
	Marks pq.StringArray `json:"marks" check:"marks"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetComponentMarks 通过一组标识码获取一组组件
func GetComponentMarks(args *ArgsGetComponentMarks) (dataList []FieldsComponent, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM core_page_style_component WHERE system = $1 AND mark = ANY($2) AND delete_at < to_timestamp(1000000)", args.System, args.Marks)
	return
}

// ArgsCreateComponent 创建组件参数
type ArgsCreateComponent struct {
	//系统
	System string `db:"system" json:"system" check:"mark"`
	//关联标识码
	// 全局唯一
	Mark string `db:"mark" json:"mark" check:"mark"`
	//组件名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"150"`
	//组件介绍
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//组件封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//商户订阅
	// 必须存在商户订阅配置的订阅，才能使用该组件
	OrgSubConfigID pq.Int64Array `db:"org_sub_config_id" json:"orgSubConfigID" check:"ids" empty:"true"`
	//商户功能
	// 只有开通相关功能，才能使用使用该组件
	OrgFuncList pq.StringArray `db:"org_func_list" json:"orgFuncList" check:"marks" empty:"true"`
	//样式结构内容
	Data string `db:"data" json:"data"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateComponent 创建组件
func CreateComponent(args *ArgsCreateComponent) (err error) {
	//检查mark是否存在
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM core_page_style_component WHERE system = $1 AND mark = $2 AND delete_at < to_timestamp(1000000)", args.System, args.Mark)
	if err == nil || id > 0 {
		return errors.New("system mark is exist")
	}
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_page_style_component (system, mark, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list, data, params) VALUES (:system,:mark,:name,:des,:cover_file_id,:sort_id,:tags,:org_sub_config_id,:org_func_list,:data,:params)", args)
	return
}

// ArgsUpdateComponent 修改组件参数
type ArgsUpdateComponent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组件名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"150"`
	//组件介绍
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//组件封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//商户订阅
	// 必须存在商户订阅配置的订阅，才能使用该组件
	OrgSubConfigID pq.Int64Array `db:"org_sub_config_id" json:"orgSubConfigID" check:"ids" empty:"true"`
	//商户功能
	// 只有开通相关功能，才能使用使用该组件
	OrgFuncList pq.StringArray `db:"org_func_list" json:"orgFuncList" check:"marks" empty:"true"`
	//样式结构内容
	Data string `db:"data" json:"data"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateComponent 修改组件
func UpdateComponent(args *ArgsUpdateComponent) (err error) {
	//写入数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_page_style_component SET update_at = NOW(), name = :name, des = :des, cover_file_id = :cover_file_id, sort_id = :sort_id, tags = :tags, org_sub_config_id = :org_sub_config_id, org_func_list = :org_func_list, data = :data, params = :params WHERE id = :id", args)
	return
}

// ArgsDeleteComponent 删除组件参数
type ArgsDeleteComponent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteComponent 删除组件
func DeleteComponent(args *ArgsDeleteComponent) (err error) {
	//检查模版是否使用了组件
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM core_page_style_template WHERE component_ids <@ $1 AND delete_at < to_timestamp(1000000)", args.ID)
	if err == nil && id > 0 {
		err = errors.New("template have use component")
		return
	}
	//删除组件
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "core_page_style_component", "id", args)
	return
}

// OutputComponent 导出组件数据
func OutputComponent() (data string, err error) {
	var dataList []FieldsComponent
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, system, mark, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list, data, params FROM core_page_style_component WHERE delete_at < to_timestamp(1000000)")
	if err != nil || len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	var dataByte []byte
	dataByte, err = json.Marshal(dataList)
	if err != nil {
		return
	}
	data = string(dataByte)
	return
}

// ArgsImportComponent 导入数据参数
type ArgsImportComponent struct {
	//数据
	Data string `json:"data"`
}

// ImportComponent 导入数据
func ImportComponent(args *ArgsImportComponent) (err error) {
	var dataList []FieldsComponent
	err = json.Unmarshal([]byte(args.Data), &dataList)
	if err != nil {
		return
	}
	//准备写入数据包
	var insertData []interface{}
	for _, v := range dataList {
		//查询mark是否已经创建
		var id int64
		err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM core_page_style_component WHERE system = $1 AND mark = $2 AND delete_at < to_timestamp(1000000)", v.System, v.Mark)
		if err == nil && id > 0 {
			continue
		}
		//写入数据
		insertData = append(insertData, ArgsCreateComponent{
			System:         v.System,
			Mark:           v.Mark,
			Name:           v.Name,
			Des:            v.Des,
			CoverFileID:    v.CoverFileID,
			SortID:         v.SortID,
			Tags:           v.Tags,
			OrgSubConfigID: v.OrgSubConfigID,
			OrgFuncList:    v.OrgFuncList,
			Data:           v.Data,
			Params:         v.Params,
		})
	}
	err = CoreSQL.CreateMore(Router2SystemConfig.MainDB.DB, "INSERT INTO core_page_style_component (system, mark, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list, data, params) VALUES (:system,:mark,:name,:des,:cover_file_id,:sort_id,:tags,:org_sub_config_id,:org_func_list,:data,:params)", insertData)
	return
}
