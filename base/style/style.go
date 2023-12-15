package BaseStyle

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// 获取样式库列表
type ArgsGetStyleList struct {
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

func GetStyleList(args *ArgsGetStyleList) (dataList []FieldsStyle, dataCount int64, err error) {
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
		where = where + " AND (name ILIKE '%' || :search || '%' OR title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_style",
		"id",
		"SELECT id, create_at, update_at, delete_at, name, mark, system_mark, components, title, cover_file_id, des_files, sort_id, tags FROM core_style WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// 获取样式库
type ArgsGetStyleByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

func GetStyleByID(args *ArgsGetStyleByID) (data FieldsStyle, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, name, mark, system_mark, components, title, des, cover_file_id, des_files, sort_id, tags, params FROM core_style WHERE id = $1 AND delete_at < to_timestamp(1000000)", args.ID)
	return
}

// 通过mark获取样式哭
type ArgsGetStyleByMark struct {
	//样式mark
	Mark string `db:"mark" json:"mark" check:"mark"`
}

func GetStyleByMark(args *ArgsGetStyleByMark) (data FieldsStyle, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, name, mark, system_mark, components, title, des, cover_file_id, des_files, sort_id, tags, params FROM core_style WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	return
}

// 获取一组IDs
type ArgsGetStyleMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

func GetStyleMore(args *ArgsGetStyleMore) (dataList []FieldsStyle, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "core_style", "id, create_at, update_at, delete_at, name, mark, system_mark, components, title, des, cover_file_id, des_files, sort_id, tags, params", args.IDs, args.HaveRemove)
	return
}

func GetStyleMoreMap(args *ArgsGetStyleMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsNameAndDelete("core_style", args.IDs, args.HaveRemove)
	return
}

// 创建新的样式库
type ArgsCreateStyle struct {
	//样式库名称
	Name string `db:"name" json:"name" check:"name"`
	//关联标识码
	// 用于识别代码片段
	Mark string `db:"mark" json:"mark" check:"mark"`
	//样式使用渠道
	// app APP；wxx 小程序等，可以任意定义，模块内不做限制
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark"`
	//分栏样式结构设计
	Components pq.Int64Array `db:"components" json:"components" check:"ids" empty:"true"`
	//默认标题
	// 标题是展示给用户的，样式库名称和该标题不是一个
	Title string `db:"title" json:"title" check:"name" empty:"true"`
	//默认描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//默认封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//默认描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func CreateStyle(args *ArgsCreateStyle) (data FieldsStyle, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at FROM core_style WHERE mark = $1 AND delete_at > to_timestamp(1000000)", args.Mark)
	if err == nil {
		if data.ID > 0 {
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_style SET update_at = NOW(), delete_at = to_timestamp(0), name = :name, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, sort_id = :sort_id, tags = :tags, params = :params WHERE id = :id", map[string]interface{}{
				"id":            data.ID,
				"name":          args.Name,
				"system_mark":   args.SystemMark,
				"components":    args.Components,
				"title":         args.Title,
				"des":           args.Des,
				"cover_file_id": args.CoverFileID,
				"des_files":     args.DesFiles,
				"sort_id":       args.SortID,
				"tags":          args.Tags,
				"params":        args.Params,
			})
			if err == nil {
				data = FieldsStyle{
					ID:          data.ID,
					CreateAt:    data.CreateAt,
					UpdateAt:    data.UpdateAt,
					DeleteAt:    data.DeleteAt,
					Name:        args.Name,
					Mark:        args.Mark,
					SystemMark:  args.SystemMark,
					Components:  args.Components,
					Title:       args.Title,
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
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_style", "INSERT INTO core_style (name, mark, system_mark, components, title, des, cover_file_id, des_files, sort_id, tags, params) VALUES (:name, :mark, :system_mark, :components, :title, :des, :cover_file_id, :des_files, :sort_id, :tags, :params)", args, &data)
	return
}

// 修改样式库
type ArgsUpdateStyle struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//样式库名称
	Name string `db:"name" json:"name" check:"name"`
	//关联标识码
	// 用于识别代码片段
	Mark string `db:"mark" json:"mark" check:"mark"`
	//样式使用渠道
	// app APP；wxx 小程序等，可以任意定义，模块内不做限制
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark"`
	//分栏样式结构设计
	Components pq.Int64Array `db:"components" json:"components" check:"ids" empty:"true"`
	//默认标题
	// 标题是展示给用户的，样式库名称和该标题不是一个
	Title string `db:"title" json:"title" check:"name" empty:"true"`
	//默认描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//默认封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//默认描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func UpdateStyle(args *ArgsUpdateStyle) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_style SET update_at = NOW(), name = :name, mark = :mark, system_mark = :system_mark, components = :components, title = :title, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, sort_id = :sort_id, tags = :tags, params = :params WHERE id = :id", args)
	return
}

// ArgsDeleteStyle 删除样式库参数
type ArgsDeleteStyle struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteStyle 删除样式库
func DeleteStyle(args *ArgsDeleteStyle) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "core_style", "id", args)
	return
}
