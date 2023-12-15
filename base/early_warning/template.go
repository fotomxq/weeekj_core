package BaseEarlyWarning

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//通知模版

// 查看模版列表
type ArgsGetTemplateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//搜索
	Search string
}

func GetTemplateList(args *ArgsGetTemplateList) (dataList []FieldsTemplateType, dataCount int64, err error) {
	where := "(mark ILIKE '%' || :search || '%' OR name ILIKE '%' || :search || '%' OR title ILIKE '%' || :search || '%' OR content ILIKE '%')"
	maps := map[string]interface{}{
		"search": args.Search,
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_ew_template",
		"id",
		"SELECT id, create_at, update_at, mark, name, default_expire_time, title, content, template_id, bind_data FROM core_ew_template WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "default_expire_time"},
	)
	return
}

// 查看模版
type ArgsGetTemplateByID struct {
	//ID
	ID int64
}

func GetTemplateByID(args *ArgsGetTemplateByID) (data FieldsTemplateType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, mark, name, default_expire_time, title, content, template_id, bind_data FROM core_ew_template WHERE id = $1", args.ID)
	return
}

type ArgsGetTemplateByMark struct {
	//Mark
	Mark string
}

func GetTemplateByMark(args *ArgsGetTemplateByMark) (data FieldsTemplateType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, mark, name, default_expire_time, title, content, template_id, bind_data FROM core_ew_template WHERE mark = $1", args.Mark)
	return
}

// 创建新的模版
type ArgsCreateTemplate struct {
	//标识码
	Mark string `db:"mark"`
	//名称
	Name string `db:"name"`
	//默认过期时间
	DefaultExpireTime string `db:"default_expire_time"`
	//标题
	Title string `db:"title"`
	//内容
	Content string `db:"content"`
	//短消息模版ID
	TemplateID string `db:"template_id"`
	//短消息模版变量
	BindData []string `db:"bind_data"`
}

func CreateTemplate(args *ArgsCreateTemplate) (data FieldsTemplateType, err error) {
	data, err = GetTemplateByMark(&ArgsGetTemplateByMark{
		Mark: args.Mark,
	})
	if err == nil {
		return
	}
	var lastID int64
	lastID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_ew_template (mark, name, default_expire_time, title, content, template_id, bind_data) VALUES (:mark,:name,:default_expire_time,:title,:content,:template_id,:bind_data)", map[string]interface{}{
		"mark":                args.Mark,
		"name":                args.Name,
		"default_expire_time": args.DefaultExpireTime,
		"title":               args.Title,
		"content":             args.Content,
		"template_id":         args.TemplateID,
		"bind_data":           FieldsTemplateBindData(args.BindData),
	})
	if err == nil {
		data, err = GetTemplateByID(&ArgsGetTemplateByID{
			ID: lastID,
		})
	}
	return
}

// 修改模版
type ArgsUpdateTemplate struct {
	//ID
	ID int64 `db:"id"`
	//标识码
	Mark string `db:"mark"`
	//名称
	Name string `db:"name"`
	//默认过期时间
	DefaultExpireTime string `db:"default_expire_time"`
	//标题
	Title string `db:"title"`
	//内容
	Content string `db:"content"`
	//短消息模版ID
	TemplateID string `db:"template_id"`
	//短消息模版变量
	BindData []string `db:"bind_data"`
}

func UpdateTemplate(args *ArgsUpdateTemplate) (err error) {
	var data FieldsTemplateType
	data, err = GetTemplateByID(&ArgsGetTemplateByID{
		ID: args.ID,
	})
	if err != nil {
		err = errors.New("cannot find id, " + err.Error())
		return
	}
	if data.Mark != args.Mark {
		_, err = GetTemplateByMark(&ArgsGetTemplateByMark{
			Mark: args.Mark,
		})
		if err == nil {
			return
		}
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_ew_template SET update_at = NOW(), mark = :mark, name = :name, default_expire_time = :default_expire_time, title = :title, content = :content, template_id = :template_id, bind_data = :bind_data WHERE id = :id", map[string]interface{}{
		"id":                  args.ID,
		"mark":                args.Mark,
		"name":                args.Name,
		"default_expire_time": args.DefaultExpireTime,
		"title":               args.Title,
		"content":             args.Content,
		"template_id":         args.TemplateID,
		"bind_data":           FieldsTemplateBindData(args.BindData),
	})
	return
}

// 删除模块
type ArgsDeleteTemplate struct {
	//ID
	ID int64 `db:"id"`
}

func DeleteTemplate(args *ArgsDeleteTemplate) (err error) {
	//解绑模版
	//如果失败则继续执行后续
	_ = SetUnBind(&ArgsSetUnBind{
		TemplateID: args.ID,
	})
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "DELETE FROM core_ew_template WHERE id = :id", args)
	return
}
