package IOTDevice

import (
	"errors"
	"fmt"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAutoInfoTemplateList 查询模版列表参数
type ArgsGetAutoInfoTemplateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//任务动作
	SendAction string `db:"send_action" json:"sendAction" check:"mark" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `db:"search" json:"search" check:"search" empty:"true"`
}

// GetAutoInfoTemplateList 查询模版列表
func GetAutoInfoTemplateList(args *ArgsGetAutoInfoTemplateList) (dataList []FieldsAutoInfoTemplate, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.Mark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.SendAction != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "send_action = :send_action"
		maps["send_action"] = args.SendAction
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
	tableName := "iot_core_auto_info_template"
	var rawList []FieldsAutoInfoTemplate
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getAutoInfoTemplateByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetAutoInfoTemplate 获取指定ID的数据参数
type ArgsGetAutoInfoTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetAutoInfoTemplate 获取指定ID的数据
func GetAutoInfoTemplate(args *ArgsGetAutoInfoTemplate) (data FieldsAutoInfoTemplate, err error) {
	data = getAutoInfoTemplateByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateAutoInfoTemplate 创建模版参数
type ArgsCreateAutoInfoTemplate struct {
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//冷却时间
	WaitTime int64 `db:"wait_time" json:"waitTime" check:"int64Than0"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark" check:"mark"`
	// 等式
	// 0 等于; 1 小于; 2 大于; 3 不等于
	Eq int `db:"eq" json:"eq" check:"intThan0" empty:"true"`
	//值
	Val string `db:"val" json:"val"`
	//发送任务指令
	// 留空则发送触发条件的数据包
	SendAction string `db:"send_action" json:"sendAction" check:"mark" empty:"true"`
	//发送参数
	ParamsData []byte `db:"params_data" json:"paramsData"`
}

// CreateAutoInfoTemplate 创建模版
func CreateAutoInfoTemplate(args *ArgsCreateAutoInfoTemplate) (data FieldsAutoInfoTemplate, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "iot_core_auto_info_template", "INSERT INTO iot_core_auto_info_template (name, wait_time, mark, eq, val, send_action, params_data) VALUES (:name,:wait_time,:mark,:eq,:val,:send_action,:params_data)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateAutoInfoTemplate 修改模版参数
type ArgsUpdateAutoInfoTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//冷却时间
	WaitTime int64 `db:"wait_time" json:"waitTime" check:"int64Than0"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark" check:"mark"`
	// 等式
	// 0 等于; 1 小于; 2 大于; 3 不等于
	Eq int `db:"eq" json:"eq" check:"intThan0" empty:"true"`
	//值
	Val string `db:"val" json:"val"`
	//发送任务指令
	// 留空则发送触发条件的数据包
	SendAction string `db:"send_action" json:"sendAction" check:"mark" empty:"true"`
	//发送参数
	ParamsData []byte `db:"params_data" json:"paramsData"`
}

// UpdateAutoInfoTemplate 修改模版
func UpdateAutoInfoTemplate(args *ArgsUpdateAutoInfoTemplate) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_auto_info_template SET update_at = NOW(), name = :name, wait_time = :wait_time, mark = :mark, eq = :eq, val = :val, send_action = :send_action, params_data = :params_data WHERE id = :id", args)
	if err != nil {
		return
	}
	deleteAutoInfoTemplateCache(args.ID)
	return
}

// ArgsDeleteAutoInfoTemplate 删除模版参数
type ArgsDeleteAutoInfoTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteAutoInfoTemplate 删除模版
func DeleteAutoInfoTemplate(args *ArgsDeleteAutoInfoTemplate) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "iot_core_auto_info_template", "id", args)
	if err != nil {
		return
	}
	deleteAutoInfoTemplateCache(args.ID)
	CoreNats.PushDataNoErr("/iot/device/auto_info_template", "delete", args.ID, "", nil)
	return
}

// 获取指定ID
func getAutoInfoTemplateByID(id int64) (data FieldsAutoInfoTemplate) {
	cacheMark := getAutoInfoTemplateCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, name, wait_time, mark, eq, val, send_action, params_data FROM iot_core_auto_info_template WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getAutoInfoTemplateCacheMark(id int64) string {
	return fmt.Sprint("iot:device:auto:template:id:", id)
}

func deleteAutoInfoTemplateCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getAutoInfoTemplateCacheMark(id))
}
