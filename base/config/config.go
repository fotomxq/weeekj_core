package BaseConfig

import (
	"errors"
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

var (
	//最后1次更新时间
	lastUpdateTime = CoreFilter.GetNowTime().Unix() - 1
	//分组配置数据
	groupData DataGroup
	//RootDir 根目录
	RootDir = CoreFile.BaseDir()
	//缓存时间
	cacheTime = 2592000
)

// Init 初始化
func Init() (err error) {
	//加载数据
	err = loadGroupData()
	return
}

// GetLastUpdateTime 获取最后更新时间
// 用于外部缓冲、其他模块定期获取数据做判定处理
func GetLastUpdateTime() int64 {
	return lastUpdateTime
}

// GetAll 获取所有配置
func GetAll() (dataList []FieldsConfigType, err error) {
	//获取数据
	var rawList []FieldsConfigType
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT mark FROM core_config")
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	//覆盖数据
	for _, v := range rawList {
		vData, _ := getByMark(v.Mark)
		if vData.Mark == "" {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// ArgsGetAllByGroupMark 获取某个组的所有配置
// 注意该方法从缓冲获取，必须全局经过初始化才可以使用
type ArgsGetAllByGroupMark struct {
	//分组标识码
	GroupMark string
}

func GetAllByGroupMark(args *ArgsGetAllByGroupMark) (dataList []FieldsConfigType, err error) {
	//获取数据
	var configs []FieldsConfigType
	configs, err = GetAll()
	if err != nil {
		return
	}
	for _, v := range configs {
		if v.GroupMark != args.GroupMark {
			continue
		}
		dataList = append(dataList, v)
	}
	//反馈
	return
}

// ArgsGetByMark 获取配置参数
type ArgsGetByMark struct {
	//标识码
	Mark string `json:"mark" check:"mark"`
}

// GetByMark 获取配置
func GetByMark(args *ArgsGetByMark) (data FieldsConfigType, err error) {
	//获取数据
	data, err = getByMark(args.Mark)
	if err != nil {
		return
	}
	//反馈
	return
}

// 获取指定mark数据
func getByMark(mark string) (data FieldsConfigType, err error) {
	//获取缓冲
	cacheMark := fmt.Sprint(getConfigCacheMark(mark))
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.Mark != "" {
		return
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT mark, create_at, update_at, allow_public, update_hash, name, group_mark, des, value_type, value FROM core_config WHERE mark = $1 LIMIT 1", mark)
	if err != nil {
		return
	}
	if data.Mark == "" {
		err = errors.New("no data")
		return
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTime)
	//反馈
	return
}

// ArgsGetData 获取配置参数
type ArgsGetData struct {
	//标识码
	Mark string
}

// GetData 获取配置
func GetData(args *ArgsGetData) (string, error) {
	data, err := getByMark(args.Mark)
	if err != nil {
		return "", err
	}
	return data.Value, nil
}

// ArgsCreate 创建新的数据参数
// 如果已经存在，则忽略
type ArgsCreate struct {
	//标识码
	Mark string
	//是否公开
	AllowPublic bool
	//名称
	Name string
	//值类型
	ValueType string
	//值
	Value string
	//分组标识码
	GroupMark string
	//描述
	Des string
}

// Create 创建新的数据
func Create(args *ArgsCreate) (err error) {
	//不存在则创建新的
	data := FieldsConfigType{
		Mark:        args.Mark,
		AllowPublic: args.AllowPublic,
		UpdateHash:  getUpdateHash(),
		Name:        args.Name,
		GroupMark:   args.GroupMark,
		Des:         args.Des,
		ValueType:   args.ValueType,
		Value:       args.Value,
	}
	_, err = CoreSQL.CreateOne(
		Router2SystemConfig.MainDB.DB,
		"INSERT INTO core_config(mark, update_at, allow_public, update_hash, name, group_mark, des, value_type, value) VALUES(:mark, now(), :allow_public, :update_hash, :name, :group_mark, :des, :value_type, :value);",
		&data,
	)
	if err != nil {
		return
	}
	//更新最后1次更新时间
	lastUpdateTime = CoreFilter.GetNowTime().Unix()
	//清理缓冲
	deleteConfigCache(args.Mark)
	//反馈
	return
}

// ArgsUpdateByMark 写入配置参数
type ArgsUpdateByMark struct {
	//hash
	UpdateHash string
	//标识码
	Mark string
	//新的值
	Value string
}

// UpdateByMark 写入配置
func UpdateByMark(args *ArgsUpdateByMark) (err error) {
	newHash := getUpdateHash()
	if args.UpdateHash == "" {
		_, err = CoreSQL.UpdateOne(
			Router2SystemConfig.MainDB.DB,
			"UPDATE core_config SET value=:value, update_at=NOW(), update_hash=:new_update_hash WHERE mark=:mark",
			map[string]interface{}{
				"value":           args.Value,
				"new_update_hash": newHash,
				"mark":            args.Mark,
			},
		)
	} else {
		_, err = CoreSQL.UpdateOne(
			Router2SystemConfig.MainDB.DB,
			"UPDATE core_config SET value=:value, update_at=NOW(), update_hash=:new_update_hash WHERE mark=:mark",
			map[string]interface{}{
				"value":           args.Value,
				"new_update_hash": newHash,
				"mark":            args.Mark,
				"update_hash":     args.Value,
			},
		)
	}
	if err != nil {
		return
	}
	//更新最后1次更新时间
	lastUpdateTime = CoreFilter.GetNowTime().Unix()
	//清理缓冲
	deleteConfigCache(args.Mark)
	//反馈
	return
}

// ArgsDeleteByMark 删除配置参数
type ArgsDeleteByMark struct {
	//标识码
	Mark string `db:"mark"`
}

// DeleteByMark 删除配置
func DeleteByMark(args *ArgsDeleteByMark) (err error) {
	//执行操作
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_config", "mark", args)
	if err != nil {
		return
	}
	//更新最后1次更新时间
	lastUpdateTime = CoreFilter.GetNowTime().Unix()
	//清理缓冲
	deleteConfigCache(args.Mark)
	//反馈
	return
}

// ArgsUpdateInfo 修改配置的基本信息参数
type ArgsUpdateInfo struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//是否可以公开
	AllowPublic bool `db:"allow_public" json:"allowPublic"`
	//名称
	Name string `db:"name" json:"name"`
	//分组
	GroupMark string `db:"group_mark" json:"groupMark"`
	//描述
	Des string `db:"des" json:"des"`
	//结构
	// string / string_md / bool / int / int64 / float64
	// 结构也可用于前端判定某个特殊的样式，如时间样式、过期时间样式等，程序内不做任何限定，只是标记
	ValueType string `db:"value_type" json:"valueType"`
}

// UpdateInfo 修改配置的基本信息
func UpdateInfo(args *ArgsUpdateInfo) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_config SET update_at = NOW(), allow_public = :allow_public, name = :name, group_mark = :group_mark, des = :des, value_type = :value_type WHERE mark = :mark", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteConfigCache(args.Mark)
	//反馈
	return
}

// 获取随机吗
func getUpdateHash() string {
	str, err := CoreFilter.GetRandStr3(10)
	if err != nil {
		return "1"
	}
	return str
}
