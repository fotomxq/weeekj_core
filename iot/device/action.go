package IOTDevice

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetActionList 获取动作列表参数
type ArgsGetActionList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetActionList 获取动作列表
func GetActionList(args *ArgsGetActionList) (dataList []FieldsAction, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "iot_core_action"
	var rawList []FieldsAction
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
		vData := getActionByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetActionMore 获取指定的一组动作参数
type ArgsGetActionMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetActionMore 获取指定的一组动作
func GetActionMore(args *ArgsGetActionMore) (dataList []FieldsAction, err error) {
	for _, v := range args.IDs {
		vData := getActionByID(v)
		if vData.ID < 1 {
			continue
		}
		if !args.HaveRemove && CoreSQL.CheckTimeHaveData(vData.DeleteAt) {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetActionMoreMap 获取指定的一组动作名称
func GetActionMoreMap(args *ArgsGetActionMore) (data map[int64]string, err error) {
	data = map[int64]string{}
	for _, v := range args.IDs {
		vData := getActionByID(v)
		if vData.ID < 1 {
			continue
		}
		if !args.HaveRemove && CoreSQL.CheckTimeHaveData(vData.DeleteAt) {
			continue
		}
		data[vData.ID] = vData.Name
	}
	return
}

// ArgsCreateAction 添加新的动作参数
type ArgsCreateAction struct {
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//动作对应任务的默认过期时间
	ExpireTime int64 `db:"expire_time" json:"expireTime" check:"int64Than0"`
	//连接方式
	// mqtt_client 与设备直接连接，用于标准物联网设计
	// mqtt_group 设备分组与设备进行mqtt广播，可用于app通告方法等
	// none 交给业务模块进行处理，任务终端不做任何广播处理
	// 本系统默认支持的是mqtt，tcp建议采用微服务跨应用或组件方式构建，以避免系统级阻塞
	ConnectType string `db:"connect_type" json:"connectType" check:"mark"`
	//扩展参数
	Configs CoreSQLConfig.FieldsConfigsType `db:"configs" json:"configs"`
}

// CreateAction 添加新的动作参数
func CreateAction(args *ArgsCreateAction) (data FieldsAction, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM iot_core_action WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	if err == nil && data.ID > 0 {
		err = errors.New("mark is replace")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "iot_core_action", "INSERT INTO iot_core_action (mark, name, des, expire_time, connect_type, configs) VALUES (:mark,:name,:des,:expire_time,:connect_type,:configs)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateAction 修改动作参数
type ArgsUpdateAction struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//动作对应任务的默认过期时间
	ExpireTime int64 `db:"expire_time" json:"expireTime" check:"int64Than0"`
	//连接方式
	// mqtt_client 与设备直接连接，用于标准物联网设计
	// mqtt_group 设备分组与设备进行mqtt广播，可用于app通告方法等
	// none 交给业务模块进行处理，任务终端不做任何广播处理
	// 本系统默认支持的是mqtt，tcp建议采用微服务跨应用或组件方式构建，以避免系统级阻塞
	ConnectType string `db:"connect_type" json:"connectType" check:"mark"`
	//扩展参数
	Configs CoreSQLConfig.FieldsConfigsType `db:"configs" json:"configs"`
}

// UpdateAction 修改动作
func UpdateAction(args *ArgsUpdateAction) (err error) {
	var data FieldsAction
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM iot_core_action WHERE mark = $1 AND id != $2 AND delete_at < to_timestamp(1000000)", args.Mark, args.ID)
	if err == nil && data.ID > 0 {
		err = errors.New("mark is replace")
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_action SET update_at = NOW(), mark = :mark, name = :name, des = :des, expire_time = :expire_time, connect_type = :connect_type, configs = :configs WHERE id = :id", args)
	if err != nil {
		return
	}
	deleteActionCache(args.ID)
	return
}

// ArgsDeleteAction 删除动作参数
type ArgsDeleteAction struct {
	//ID
	ID int64 `json:"id" check:"id"`
}

// DeleteAction 删除动作
func DeleteAction(args *ArgsDeleteAction) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "iot_core_action", "id", args)
	if err != nil {
		return
	}
	deleteActionCache(args.ID)
	return
}

// 获取指定ID
func getActionByID(id int64) (data FieldsAction) {
	cacheMark := getActionCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, name, des, expire_time, connect_type, configs FROM iot_core_action WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getActionCacheMark(id int64) string {
	return fmt.Sprint("iot:device:action:id:", id)
}

func deleteActionCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getActionCacheMark(id))
}
