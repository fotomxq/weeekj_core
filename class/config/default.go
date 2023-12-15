package ClassConfig

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ConfigDefault 默认配置项
type ConfigDefault struct {
	//表名称
	TableName string
}

// ArgsGetConfigDefaultList 获取默认配置列表参数
type ArgsGetConfigDefaultList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//是否需要公共限制参数
	NeedAllowPublic bool `json:"needAllowPublic" check:"bool"`
	AllowPublic     bool `json:"allowPublic" check:"bool"`
	//是否需要行政限制参数
	NeedAllowSelfView bool `json:"needAllowSelfView" check:"bool"`
	AllowSelfView     bool `json:"allowSelfView" check:"bool"`
	NeedAllowSelfSet  bool `json:"needAllowSelfSet" check:"bool"`
	AllowSelfSet      bool `json:"allowSelfSet" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigDefaultList 获取默认配置列表
func (t *ConfigDefault) GetConfigDefaultList(args *ArgsGetConfigDefaultList) (dataList []FieldsConfigDefault, dataCount int64, err error) {
	var where string
	maps := map[string]interface{}{}
	if args.NeedAllowPublic {
		where = where + "allow_public = :allow_public"
		maps["allow_public"] = args.AllowPublic
	}
	if args.NeedAllowSelfView {
		if where != "" {
			where = where + " AND "
		}
		where = where + "allow_self_view = :allow_self_view"
		maps["allow_self_view"] = args.AllowSelfView
	}
	if args.NeedAllowSelfSet {
		if where != "" {
			where = where + " AND "
		}
		where = where + "allow_self_set = :allow_self_set"
		maps["allow_self_set"] = args.AllowSelfSet
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	var rawList []FieldsConfigDefault
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		t.TableName,
		"id",
		"SELECT mark FROM "+t.TableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "mark", "name", "create_at", "update_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData, _ := t.GetConfigDefaultMark(v.Mark)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsSetConfigDefault 创建或修改指定的数据参数
type ArgsSetConfigDefault struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//是否可以公开
	AllowPublic bool `db:"allow_public" json:"allowPublic"`
	//是否允许组织查看该配置
	AllowSelfView bool `db:"allow_self_view" json:"allowSelfView"`
	//是否允许组织自己修改
	AllowSelfSet bool `db:"allow_self_set" json:"allowSelfSet"`
	//结构
	// 0 string / 1 bool / 2 int / 3 int64 / 4 float64
	// 5 time 时间 / 6 daytime 带有日期的时间 / 7 unix 时间戳
	// 8 fileID 文件ID / 9 fileIDList 文件ID列
	// 10 userID 用户ID / 11 userIDList 用户ID列
	// 结构也可用于前端判定某个特殊的样式，如时间样式、过期时间样式等，程序内不做任何限定，只是标记
	ValueType int `db:"value_type" json:"valueType"`
	//正则表达式
	ValueCheck string `db:"value_check" json:"valueCheck"`
	//默认值
	ValueDefault string `db:"value_default" json:"valueDefault"`
}

// SetConfigDefault 创建或修改指定的数据
func (t *ConfigDefault) SetConfigDefault(args *ArgsSetConfigDefault) (err error) {
	var data FieldsConfigDefault
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM "+t.TableName+" WHERE mark = $1;", args.Mark)
	if err == nil {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE "+t.TableName+" SET name = :name, update_at = NOW(), allow_public = :allow_public, allow_self_view = :allow_self_view, allow_self_set = :allow_self_set, value_type = :value_type, value_check = :value_check, value_default = :value_default WHERE id = :id", map[string]interface{}{
			"id":              data.ID,
			"name":            args.Name,
			"allow_public":    args.AllowPublic,
			"allow_self_view": args.AllowSelfView,
			"allow_self_set":  args.AllowSelfSet,
			"value_type":      args.ValueType,
			"value_check":     args.ValueCheck,
			"value_default":   args.ValueDefault,
		})
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO "+t.TableName+"(mark, name, allow_public, allow_self_view, allow_self_set, value_type, value_check, value_default) VALUES (:mark, :name, :allow_public, :allow_self_view, :allow_self_set, :value_type, :value_check, :value_default)", args)
	}
	if err != nil {
		return
	}
	t.deleteDefaultCache(args.Mark)
	return
}

// ArgsDeleteConfigDefault 删除配置参数
type ArgsDeleteConfigDefault struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
}

// DeleteConfigDefault 删除配置
func (t *ConfigDefault) DeleteConfigDefault(args *ArgsDeleteConfigDefault) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, t.TableName, "mark", args)
	if err != nil {
		return
	}
	t.deleteDefaultCache(args.Mark)
	return
}

// ArgsGetConfigDefaultByMarks 获取指定一组mark数据参数
type ArgsGetConfigDefaultByMarks struct {
	//要查询的marks
	Marks pq.StringArray
	//是否需要公共限制参数
	NeedAllowPublic bool
	AllowPublic     bool
	//是否需要行政限制参数
	NeedAllowSelfView bool
	AllowSelfView     bool
	NeedAllowSelfSet  bool
	AllowSelfSet      bool
}

// GetConfigDefaultByMarks 获取指定一组mark数据
func (t *ConfigDefault) GetConfigDefaultByMarks(args *ArgsGetConfigDefaultByMarks) (dataList []FieldsConfigDefault, err error) {
	where := "mark = ANY(:marks)"
	maps := map[string]interface{}{
		"marks": args.Marks,
	}
	if args.NeedAllowPublic {
		where = where + " AND allow_public = :allow_public"
		maps["allow_public"] = args.AllowPublic
	}
	if args.NeedAllowSelfView {
		if where != "" {
			where = where + " AND "
		}
		where = where + "allow_self_view = :allow_self_view"
		maps["allow_self_view"] = args.AllowSelfView
	}
	if args.NeedAllowSelfSet {
		if where != "" {
			where = where + " AND "
		}
		where = where + "allow_self_set = :allow_self_set"
		maps["allow_self_set"] = args.AllowSelfSet
	}
	var rawList []FieldsConfigDefault
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"SELECT mark FROM "+t.TableName+" WHERE "+where,
		maps,
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData, _ := t.GetConfigDefaultMark(v.Mark)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetConfigDefaultAll 内部获取所有配置的方法
// 同于快速构建新组织所需的配置数据组
func (t *ConfigDefault) GetConfigDefaultAll() (dataList []FieldsConfigDefault, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT mark FROM "+t.TableName+";")
	if err != nil {
		return
	}
	return
}

// GetConfigDefaultMark 内部获取指定的配置
func (t *ConfigDefault) GetConfigDefaultMark(mark string) (data FieldsConfigDefault, err error) {
	cacheMark := t.getDefaultCacheMark(mark)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, mark, name, create_at, update_at, allow_public, allow_self_view, allow_self_set, value_type, value_default FROM "+t.TableName+" WHERE mark = $1;", mark)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 2592000)
	return
}
