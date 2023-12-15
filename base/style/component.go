package BaseStyle

import (
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
	//关联标识码
	// 必填
	// 页面内独特的代码片段，声明后可直接定义该组件的默认参数形式
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//分类ID
	Sort int64 `db:"sort" json:"sort" check:"id" empty:"true"`
	//标签组
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
	if args.Mark != "" {
		where = where + " AND mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.Sort > 0 {
		where = where + " AND sort = :sort"
		maps["sort"] = args.Sort
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_style_component",
		"id",
		"SELECT id, create_at, update_at, delete_at, mark, name, cover_file_id, des_files, sort_id, tags FROM core_style_component WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetComponentByID 获取指定组件参数
type ArgsGetComponentByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetComponentByID 获取指定组件
func GetComponentByID(args *ArgsGetComponentByID) (data FieldsComponent, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, name, des, cover_file_id, des_files, sort_id, tags, params FROM core_style_component WHERE id = $1 AND delete_at < to_timestamp(1000000)", args.ID)
	return
}

type ArgsGetComponentByMark struct {
	//组件mark
	Mark string `db:"mark" json:"mark" check:"mark"`
}

func GetComponentByMark(args *ArgsGetComponentByMark) (data FieldsComponent, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, name, des, cover_file_id, des_files, sort_id, tags, params FROM core_style_component WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	return
}

// ArgsGetComponentMore 获取一组IDs参数
type ArgsGetComponentMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetComponentMore 获取一组IDs
func GetComponentMore(args *ArgsGetComponentMore) (dataList []FieldsComponent, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "core_style_component", "id, create_at, update_at, delete_at, mark, name, des, cover_file_id, des_files, sort_id, tags, params", args.IDs, args.HaveRemove)
	return
}

func GetComponentMoreMap(args *ArgsGetComponentMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsNameAndDelete("core_style_component", args.IDs, args.HaveRemove)
	return
}

// 创建新的组件
type ArgsCreateComponent struct {
	//关联标识码
	// 必填
	// 页面内独特的代码片段，声明后可直接定义该组件的默认参数形式
	Mark string `db:"mark" json:"mark" check:"mark"`
	//组件名称
	Name string `db:"name" json:"name" check:"name"`
	//组件介绍
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//组件封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//组件描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func CreateComponent(args *ArgsCreateComponent) (data FieldsComponent, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at FROM core_style_component WHERE mark = $1 AND delete_at > to_timestamp(1000000)", args.Mark)
	if err == nil {
		if data.ID > 0 {
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_style_component SET update_at = NOW(), delete_at = to_timestamp(0), name = :name, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, sort_id = :sort_id, tags = :tags, params = :params WHERE id = :id", map[string]interface{}{
				"id":            data.ID,
				"name":          args.Name,
				"des":           args.Des,
				"cover_file_id": args.CoverFileID,
				"des_files":     args.DesFiles,
				"sort_id":       args.SortID,
				"tags":          args.Tags,
				"params":        args.Params,
			})
			if err == nil {
				data = FieldsComponent{
					ID:          data.ID,
					CreateAt:    data.CreateAt,
					UpdateAt:    data.UpdateAt,
					DeleteAt:    data.DeleteAt,
					Mark:        args.Mark,
					Name:        args.Name,
					Des:         args.Des,
					CoverFileID: args.CoverFileID,
					DesFiles:    args.DesFiles,
					SortID:      args.SortID,
					Tags:        args.Tags,
					Params:      args.Params,
				}
			}
			return
		}
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_style_component", "INSERT INTO core_style_component (mark, name, des, cover_file_id, des_files, sort_id, tags, params) VALUES (:mark, :name, :des, :cover_file_id, :des_files, :sort_id, :tags, :params)", args, &data)
	return
}

// 修改组件
type ArgsUpdateComponent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//关联标识码
	// 必填
	// 页面内独特的代码片段，声明后可直接定义该组件的默认参数形式
	Mark string `db:"mark" json:"mark" check:"mark"`
	//组件名称
	Name string `db:"name" json:"name" check:"name"`
	//组件介绍
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//组件封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//组件描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func UpdateComponent(args *ArgsUpdateComponent) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_style_component SET update_at = NOW(), mark = :mark, name = :name, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, sort_id = :sort_id, tags = :tags, params = :params WHERE id = :id", args)
	return
}

// 删除组件
type ArgsDeleteComponent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

func DeleteComponent(args *ArgsDeleteComponent) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "core_style_component", "id", args)
	return
}
