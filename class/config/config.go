package ClassConfig

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// Config 通用对象内嵌配置设计
type Config struct {
	//表名称
	TableName string
	//默认配置项
	Default ConfigDefault
}

// ArgsGetConfigList 读取组织的配置数据参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//绑定ID
	// 必须填写
	BindID int64 `json:"bindID" check:"id"`
	//访问的渠道
	// public / self / admin
	VisitType string `json:"visitType" check:"mark"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 读取组织的配置数据
func (t *Config) GetConfigList(args *ArgsGetConfigList) (configList []FieldsConfigDefault, dataList []FieldsConfig, dataCount int64, err error) {
	//组合权限
	needAllowPublic := false
	allowPublic := false
	needAllowSelfView := false
	allowSelfView := false
	switch args.VisitType {
	case "public":
		needAllowPublic = true
		allowPublic = true
	case "self":
		needAllowSelfView = true
		allowSelfView = true
	case "admin":
	default:
		err = errors.New("visit type error")
		return
	}
	//获取全部权限
	configList, dataCount, err = t.Default.GetConfigDefaultList(&ArgsGetConfigDefaultList{
		Pages:             args.Pages,
		NeedAllowPublic:   needAllowPublic,
		AllowPublic:       allowPublic,
		NeedAllowSelfView: needAllowSelfView,
		AllowSelfView:     allowSelfView,
		NeedAllowSelfSet:  false,
		AllowSelfSet:      false,
		Search:            args.Search,
	})
	if err != nil {
		return
	}
	//根据查询的数据，遍历查询所有mark
	// 注意，本方法不构建虚假值，如果需完整填充，请使用后续的获取指定集合mark方法
	var marks pq.StringArray
	for _, v := range configList {
		marks = append(marks, v.Mark)
	}
	where := "bind_id = :bind_id AND mark = ANY(:marks)"
	maps := map[string]interface{}{
		"bind_id": args.BindID,
		"marks":   marks,
	}
	var rawList []FieldsConfig
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		t.TableName,
		"id",
		"SELECT mark "+"FROM "+t.TableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "mark"},
	)
	if err != nil {
		err = nil
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := t.getConfig(args.BindID, v.Mark)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// DataConfig 获取配置列表数据
type DataConfig struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//配置标识码
	Mark string `db:"mark" json:"mark"`
	//值
	Val string `db:"val" json:"val"`
}

// GetConfigMarge 获取配置列表
// 对配置进行融合处理
func (t *Config) GetConfigMarge(args *ArgsGetConfigList) (dataList []DataConfig, dataCount int64, err error) {
	var defaultList []FieldsConfigDefault
	var configList []FieldsConfig
	defaultList, configList, dataCount, err = t.GetConfigList(args)
	if err != nil {
		return
	}
	for _, v := range configList {
		dataList = append(dataList, DataConfig{
			CreateAt: v.CreateAt,
			UpdateAt: v.UpdateAt,
			Mark:     v.Mark,
			Val:      v.Val,
		})
	}
	for _, v := range defaultList {
		isFind := false
		for _, v2 := range dataList {
			if v.Mark == v2.Mark {
				isFind = true
				break
			}
		}
		if !isFind {
			dataList = append(dataList, DataConfig{
				CreateAt: v.CreateAt,
				UpdateAt: v.UpdateAt,
				Mark:     v.Mark,
				Val:      v.ValueDefault,
			})
		}
	}
	return
}

// ArgsGetConfigByMarks 获取指定一组mark的方法参数
type ArgsGetConfigByMarks struct {
	//指定绑定ID
	BindID int64
	//一组标识码
	Marks []string
	//访问的渠道
	// public / self / admin
	VisitType string
}

// GetConfigByMarks 获取指定一组mark的方法
// 本方法将自动补全数据，即如果组织尚未定义数据，则依赖于全局配置设置
func (t *Config) GetConfigByMarks(args *ArgsGetConfigByMarks) (dataList []FieldsConfig, err error) {
	//组合权限
	needAllowPublic := false
	allowPublic := false
	needAllowSelfView := false
	allowSelfView := false
	switch args.VisitType {
	case "public":
		needAllowPublic = true
		allowPublic = true
	case "self":
		needAllowSelfView = true
		allowSelfView = true
	case "admin":
	default:
		err = errors.New("visit type error")
		return
	}
	//获取默认配置
	var configList []FieldsConfigDefault
	configList, err = t.Default.GetConfigDefaultByMarks(&ArgsGetConfigDefaultByMarks{
		Marks:             args.Marks,
		NeedAllowPublic:   needAllowPublic,
		AllowPublic:       allowPublic,
		NeedAllowSelfView: needAllowSelfView,
		AllowSelfView:     allowSelfView,
		NeedAllowSelfSet:  false,
		AllowSelfSet:      false,
	})
	if err != nil {
		return
	}
	//获取组织配置
	var marks pq.StringArray
	for _, v := range configList {
		marks = append(marks, v.Mark)
	}
	where := "bind_id = :bind_id AND mark = ANY(:marks)"
	maps := map[string]interface{}{
		"bind_id": args.BindID,
		"marks":   marks,
	}
	var rawList []FieldsConfig
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"SELECT mark "+"FROM "+t.TableName+" WHERE "+where,
		maps,
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := t.getConfig(args.BindID, v.Mark)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//补全数据集合
	for _, v := range configList {
		isFind := false
		for _, v2 := range dataList {
			if v.Mark == v2.Mark {
				isFind = true
				break
			}
		}
		if !isFind {
			dataList = append(dataList, FieldsConfig{
				ID:       0,
				CreateAt: time.Time{},
				UpdateAt: time.Time{},
				BindID:   args.BindID,
				Mark:     v.Mark,
				Val:      v.ValueDefault,
			})
		}
	}
	return
}

// GetConfigByMarksMerge 获取指定的组织配置
func (t *Config) GetConfigByMarksMerge(args *ArgsGetConfigByMarks) (dataList []DataConfig, err error) {
	//组合权限
	needAllowPublic := false
	allowPublic := false
	needAllowSelfView := false
	allowSelfView := false
	switch args.VisitType {
	case "public":
		needAllowPublic = true
		allowPublic = true
	case "self":
		needAllowSelfView = true
		allowSelfView = true
	case "admin":
	default:
		err = errors.New("visit type error")
		return
	}
	var defaultList []FieldsConfigDefault
	defaultList, err = t.Default.GetConfigDefaultByMarks(&ArgsGetConfigDefaultByMarks{
		Marks:             args.Marks,
		NeedAllowPublic:   needAllowPublic,
		AllowPublic:       allowPublic,
		NeedAllowSelfView: needAllowSelfView,
		AllowSelfView:     allowSelfView,
		NeedAllowSelfSet:  false,
		AllowSelfSet:      false,
	})
	if err != nil {
		return
	}
	var configList []FieldsConfig
	configList, err = t.GetConfigByMarks(args)
	if err != nil {
		err = nil
	} else {
		for _, v := range configList {
			dataList = append(dataList, DataConfig{
				CreateAt: v.CreateAt,
				UpdateAt: v.UpdateAt,
				Mark:     v.Mark,
				Val:      v.Val,
			})
		}
	}
	for _, v := range defaultList {
		isFind := false
		for _, v2 := range dataList {
			if v.Mark == v2.Mark {
				isFind = true
				break
			}
		}
		if !isFind {
			dataList = append(dataList, DataConfig{
				CreateAt: v.CreateAt,
				UpdateAt: v.UpdateAt,
				Mark:     v.Mark,
				Val:      v.ValueDefault,
			})
		}
	}
	return
}

// ArgsGetConfig 获取指定的配置数据参数
type ArgsGetConfig struct {
	//绑定ID
	BindID int64
	//标识码
	Mark string
	//访问的渠道
	// public / self / admin
	VisitType string
}

// GetConfig 获取指定的配置数据
func (t *Config) GetConfig(args *ArgsGetConfig) (configData FieldsConfigDefault, data FieldsConfig, err error) {
	configData, err = t.Default.GetConfigDefaultMark(args.Mark)
	if err != nil {
		err = errors.New("get config by mark, " + err.Error())
		return
	}
	switch args.VisitType {
	case "public":
		if !configData.AllowPublic {
			err = errors.New("cannot view config")
			return
		}
	case "self":
		if !configData.AllowSelfView {
			err = errors.New("cannot view config")
			return
		}
	case "admin":
	default:
		err = errors.New("visit type error")
		return
	}
	data = t.getConfig(args.BindID, args.Mark)
	if data.ID < 1 {
		data = FieldsConfig{
			ID:       0,
			CreateAt: configData.CreateAt,
			UpdateAt: configData.UpdateAt,
			BindID:   0,
			Mark:     configData.Mark,
			Val:      configData.ValueDefault,
		}
		err = nil
	}
	return
}

func (t *Config) getConfig(bindID int64, mark string) (data FieldsConfig) {
	cacheMark := t.getConfigCacheMark(mark, bindID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, fmt.Sprint("SELECT id, create_at, update_at, bind_id, mark, val ", "FROM ", t.TableName, " WHERE bind_id = $1 AND mark = $2;"), bindID, mark)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 2592000)
	return
}

// GetConfigVal 获取配置
// 只获取指定数据的值
func (t *Config) GetConfigVal(args *ArgsGetConfig) (value string, err error) {
	var data FieldsConfig
	_, data, err = t.GetConfig(args)
	if err != nil {
		return
	}
	value = data.Val
	return
}

func (t *Config) GetConfigValNoErr(bindID int64, mark string) (value string) {
	value, _ = t.GetConfigVal(&ArgsGetConfig{
		BindID:    bindID,
		Mark:      mark,
		VisitType: "admin",
	})
	return
}

// GetConfigValInt 封装对配置的几个转化方法
func (t *Config) GetConfigValInt(args *ArgsGetConfig) (value int, err error) {
	var data FieldsConfig
	_, data, err = t.GetConfig(args)
	if err != nil {
		return
	}
	value, err = CoreFilter.GetIntByString(data.Val)
	return
}

func (t *Config) GetConfigValIntNoErr(bindID int64, mark string) (value int) {
	value, _ = t.GetConfigValInt(&ArgsGetConfig{
		BindID:    bindID,
		Mark:      mark,
		VisitType: "admin",
	})
	return
}

func (t *Config) GetConfigValInt64(args *ArgsGetConfig) (value int64, err error) {
	var data FieldsConfig
	_, data, err = t.GetConfig(args)
	if err != nil {
		return
	}
	value, err = CoreFilter.GetInt64ByString(data.Val)
	return
}

func (t *Config) GetConfigValInt64NoErr(bindID int64, mark string) (value int64) {
	value, _ = t.GetConfigValInt64(&ArgsGetConfig{
		BindID:    bindID,
		Mark:      mark,
		VisitType: "admin",
	})
	return
}

func (t *Config) GetConfigValFloat64(args *ArgsGetConfig) (value float64, err error) {
	var data FieldsConfig
	_, data, err = t.GetConfig(args)
	if err != nil {
		return
	}
	value, err = CoreFilter.GetFloat64ByString(data.Val)
	return
}

func (t *Config) GetConfigValBool(args *ArgsGetConfig) (value bool, err error) {
	var data FieldsConfig
	_, data, err = t.GetConfig(args)
	if err != nil {
		return
	}
	value = data.Val == "true"
	return
}

func (t *Config) GetConfigValBoolNoErr(bindID int64, mark string) (value bool) {
	var data FieldsConfig
	var err error
	_, data, err = t.GetConfig(&ArgsGetConfig{
		BindID:    bindID,
		Mark:      mark,
		VisitType: "admin",
	})
	if err != nil {
		return
	}
	value = data.Val == "true"
	return
}

// ArgsSetConfig 修改配置参数
type ArgsSetConfig struct {
	//绑定ID
	BindID int64 `json:"bindID" check:"id"`
	//标识码
	Mark string `json:"mark" check:"mark"`
	//访问的渠道
	// self / admin
	VisitType string `json:"visitType" check:"mark"`
	//值
	Val string `json:"val"`
}

// SetConfig 修改配置
func (t *Config) SetConfig(args *ArgsSetConfig) (err error) {
	var configData FieldsConfigDefault
	configData, err = t.Default.GetConfigDefaultMark(args.Mark)
	if err != nil {
		return
	}
	switch args.VisitType {
	case "self":
		if !configData.AllowSelfSet {
			err = errors.New(fmt.Sprint("org cannot set config, mark is ", args.Mark, ", bind id: ", args.BindID))
			return
		}
	case "admin":
	default:
		err = errors.New("visit type error")
		return
	}
	var data FieldsConfig
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id "+"FROM "+t.TableName+" WHERE bind_id = $1 AND mark = $2;", args.BindID, configData.Mark)
	if err == nil {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE "+t.TableName+" SET val = :val WHERE id = :id", map[string]interface{}{
			"id":  data.ID,
			"val": args.Val,
		})
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT "+"INTO "+t.TableName+" (bind_id, mark, val) VALUES (:bind_id, :mark, :val)", map[string]interface{}{
			"bind_id": args.BindID,
			"mark":    args.Mark,
			"val":     args.Val,
		})
	}
	if err != nil {
		return
	}
	t.deleteConfigCache(args.Mark, args.BindID)
	return
}

// SetConfigValSimple 内部快速修改配置
func (t *Config) SetConfigValSimple(bindID int64, mark string, val string) (err error) {
	err = t.SetConfig(&ArgsSetConfig{
		BindID:    bindID,
		Mark:      mark,
		VisitType: "admin",
		Val:       val,
	})
	return
}

// ArgsDeleteConfig 清空配置参数
type ArgsDeleteConfig struct {
	//标识码
	// 如果单独给mark，将清理所有配置
	Mark string `db:"mark" json:"mark"`
	//绑定ID
	// 如果单独给绑定ID，则将清空数据
	BindID int64 `db:"bind_id" json:"bindID"`
}

// DeleteConfig 清空配置
// 两个参数必须至少给予一个
func (t *Config) DeleteConfig(args *ArgsDeleteConfig) (err error) {
	if args.Mark == "" && args.BindID < 1 {
		err = errors.New("mark or bind is empty")
		return
	}
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, t.TableName, "(mark = :mark OR :mark = '') AND (bind_id = :bind_id OR :bind_id < 1)", args)
	if err != nil {
		return
	}
	t.deleteConfigCache(args.Mark, args.BindID)
	return
}
