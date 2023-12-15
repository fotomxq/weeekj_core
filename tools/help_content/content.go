package ToolsHelpContent

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLMarks "github.com/fotomxq/weeekj_core/v5/core/sql/marks"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// 获取列表
type ArgsGetContentList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//文本唯一标识码
	// none为预留值，指定后可以重复且该值将无效
	// 用于在不同页面使用
	// 删除后，将不会占用该mark设置，一个mark可以指定一个正常数据和多个已删除数据
	Mark string `db:"mark" json:"mark"`
	//分类ID
	// > -1 为包含；否则不包含。0为没有设定
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetContentList(args *ArgsGetContentList) (dataList []FieldsContent, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
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
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"tools_help_content",
		"id",
		"SELECT id, create_at, update_at, delete_at, mark, is_public, sort_id, tags, title, cover_file_id, bind_ids, bind_marks, params FROM tools_help_content WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetContentByID 获取ID参数
type ArgsGetContentByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//是否需要public
	NeedPublic bool `db:"need_public" json:"needPublic" check:"bool"`
	//是否公开
	IsPublic bool `db:"is_public" json:"isPublic" check:"bool"`
}

// GetContentByID 获取ID
func GetContentByID(args *ArgsGetContentByID) (data FieldsContent, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, is_public, sort_id, tags, title, cover_file_id, des, bind_ids, bind_marks, params FROM tools_help_content WHERE id = $1 AND delete_at < to_timestamp(1000000) AND ($2 = FALSE OR ($2 = TRUE AND is_public = $3))", args.ID, args.NeedPublic, args.IsPublic)
	return
}

// 获取指定mark
type ArgsGetContentByMark struct {
	//文本唯一标识码
	// none为预留值，指定后可以重复且该值将无效
	// 用于在不同页面使用
	// 删除后，将不会占用该mark设置，一个mark可以指定一个正常数据和多个已删除数据
	Mark string `db:"mark" json:"mark"`
	//是否需要public
	NeedPublic bool `db:"need_public" json:"needPublic" check:"bool"`
	//是否公开
	IsPublic bool `db:"is_public" json:"isPublic" check:"bool"`
}

func GetContentByMark(args *ArgsGetContentByMark) (data FieldsContent, err error) {
	//禁止访问none和空数据
	if args.Mark == "" || args.Mark == "none" {
		err = errors.New("mark not support")
		return
	}
	//查询数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, is_public, sort_id, tags, title, cover_file_id, des, bind_ids, bind_marks, params FROM tools_help_content WHERE mark = $1 AND delete_at < to_timestamp(1000000) AND ($2 = FALSE OR ($2 = TRUE AND is_public = $3))", args.Mark, args.NeedPublic, args.IsPublic)
	return
}

// 获取一组IDs
type ArgsGetContentMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

func GetContentMore(args *ArgsGetContentMore) (dataList []FieldsContent, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "tools_help_content", "id, create_at, update_at, delete_at, mark, is_public, sort_id, tags, title, cover_file_id, bind_ids, bind_marks, params", args.IDs, args.HaveRemove)
	return
}

func GetContentMoreMap(args *ArgsGetContentMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsTitleAndDelete("tools_help_content", args.IDs, args.HaveRemove)
	return
}

// 获取一组Marks
type ArgsGetContentMoreByMark struct {
	//Mark列
	Marks pq.StringArray `json:"marks" check:"marks"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

func GetContentMoreByMark(args *ArgsGetContentMoreByMark) (dataList []FieldsContent, err error) {
	err = CoreSQLMarks.GetMarksAndDelete(&dataList, "tools_help_content", "id, create_at, update_at, delete_at, mark, is_public, sort_id, tags, title, cover_file_id, bind_ids, bind_marks, params", args.Marks, args.HaveRemove)
	return
}

func GetContentMoreByMarkMap(args *ArgsGetContentMoreByMark) (data map[string]string, err error) {
	data, err = CoreSQLMarks.GetMarksTitleAndDelete("tools_help_content", args.Marks, args.HaveRemove)
	return
}

// 创建新的词条
type ArgsCreateContent struct {
	//文本唯一标识码
	// none为预留值，指定后可以重复且该值将无效
	// 用于在不同页面使用
	// 删除后，将不会占用该mark设置，一个mark可以指定一个正常数据和多个已删除数据
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否公开
	// 非公开数据将作为草稿或私有数据存在，只有管理员可以看到
	IsPublic bool `db:"is_public" json:"isPublic" check:"bool" empty:"true"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//封面文件
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//关联阅读引导ID
	BindIDs pq.Int64Array `db:"bind_ids" json:"bindIDs" check:"ids" empty:"true"`
	//关联阅读引导mark
	BindMarks pq.StringArray `db:"bind_marks" json:"bindMarks" check:"marks" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func CreateContent(args *ArgsCreateContent) (data FieldsContent, err error) {
	//如果mark可用，则检查是否存在数据，如果存在则拒绝
	if args.Mark != "" && args.Mark != "none" {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_help_content WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
		if err == nil {
			if data.ID > 0 {
				err = errors.New("mark exist")
				return
			}
		}
	}
	//构建新的数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tools_help_content", "INSERT INTO tools_help_content (mark, is_public, sort_id, tags, title, cover_file_id, des, bind_ids, bind_marks, params) VALUES (:mark, :is_public, :sort_id, :tags, :title, :cover_file_id, :des, :bind_ids, :bind_marks, :params)", args, &data)
	return
}

// 修改词条
type ArgsUpdateContent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" check:"id"`
	//文本唯一标识码
	// none为预留值，指定后可以重复且该值将无效
	// 用于在不同页面使用
	// 删除后，将不会占用该mark设置，一个mark可以指定一个正常数据和多个已删除数据
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否公开
	// 非公开数据将作为草稿或私有数据存在，只有管理员可以看到
	IsPublic bool `db:"is_public" json:"isPublic" check:"bool" empty:"true"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//封面文件
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//关联阅读引导ID
	BindIDs pq.Int64Array `db:"bind_ids" json:"bindIDs" check:"ids" empty:"true"`
	//关联阅读引导mark
	BindMarks pq.StringArray `db:"bind_marks" json:"bindMarks" check:"marks" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func UpdateContent(args *ArgsUpdateContent) (err error) {
	//检查mark，如果存在则检查是否一个数据，否则退出
	if args.Mark != "" && args.Mark != "none" {
		var data FieldsContent
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_help_content WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
		if err == nil {
			if data.ID > 0 && data.ID != args.ID {
				err = errors.New("mark exist")
				return
			}
		}
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tools_help_content SET update_at = NOW(), mark = :mark, is_public = :is_public, sort_id = :sort_id, tags = :tags, title = :title, cover_file_id = :cover_file_id, des = :des, bind_ids = :bind_ids, bind_marks = :bind_marks, params = :params WHERE id = :id", args)
	return
}

// 删除词条
type ArgsDeleteContent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

func DeleteContent(args *ArgsDeleteContent) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "tools_help_content", "id = :id", args)
	return
}
