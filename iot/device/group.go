package IOTDevice

import (
	"errors"
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"

	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	"github.com/lib/pq"
)

// ArgsGetGroupList 获取分组列表参数
type ArgsGetGroupList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetGroupList 获取分组列表
func GetGroupList(args *ArgsGetGroupList) (dataList []FieldsGroup, dataCount int64, err error) {
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
	tableName := "iot_core_group"
	var rawList []FieldsGroup
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
		vData := getGroupByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetGroupByID 获取指定分组ID参数
type ArgsGetGroupByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
}

// GetGroupByID 获取指定分组ID
func GetGroupByID(args *ArgsGetGroupByID) (data FieldsGroup, err error) {
	data = getGroupByID(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetGroupByMark 获取指定分组Mark参数
type ArgsGetGroupByMark struct {
	//Mark
	Mark string `json:"mark" check:"mark"`
}

// GetGroupByMark 获取指定分组Mark
func GetGroupByMark(args *ArgsGetGroupByMark) (data FieldsGroup, err error) {
	data = getGroupByMark(args.Mark)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetGroupMore 获取一组分组参数
type ArgsGetGroupMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetGroupMore 获取一组分组
func GetGroupMore(args *ArgsGetGroupMore) (dataList []FieldsGroup, err error) {
	for _, v := range args.IDs {
		vData := getGroupByID(v)
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

func GetGroupMoreMap(args *ArgsGetGroupMore) (data map[int64]string, err error) {
	data = map[int64]string{}
	for _, v := range args.IDs {
		vData := getGroupByID(v)
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

// ArgsCreateGroup 创建分组参数
type ArgsCreateGroup struct {
	//分区标识码
	// 全局必须唯一
	Mark string `db:"mark" json:"mark" check:"mark"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面列
	// 第一张作为封面
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//支持动作ID组
	Action pq.Int64Array `db:"action" json:"action" check:"ids" empty:"true"`
	//心跳超时时间
	// 超出时间没有通讯则判定掉线
	// 单位: 秒
	ExpireTime int64 `db:"expire_time" json:"expireTime" check:"int64Than0" empty:"true"`
	//设备的预计使用场景
	// 0 public 公共设备 / 1 private 私有设备
	// 如果>1 则为自定义设置，具体由设备驱动识别处理
	UseType int `db:"use_type" json:"useType" check:"intThan0" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateGroup 创建分组
func CreateGroup(args *ArgsCreateGroup) (data FieldsGroup, err error) {
	if len(args.Action) > 0 {
		var actionList []FieldsAction
		actionList, err = GetActionMore(&ArgsGetActionMore{
			IDs:        args.Action,
			HaveRemove: false,
		})
		if err != nil {
			err = errors.New("action not exist, " + err.Error())
			return
		}
		if len(actionList) != len(args.Action) {
			err = errors.New("action not exist")
			return
		}
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM iot_core_group WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	if err == nil && data.ID > 0 {
		err = errors.New("mark is replace")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "iot_core_group", "INSERT INTO iot_core_group (mark, name, des, cover_files, action, expire_time, use_type, params) VALUES (:mark,:name,:des,:cover_files,:action,:expire_time,:use_type,:params)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateGroup 修改分组参数
type ArgsUpdateGroup struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//分区标识码
	// 全局必须唯一
	Mark string `db:"mark" json:"mark" check:"mark"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面列
	// 第一张作为封面
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//支持动作ID组
	Action pq.Int64Array `db:"action" json:"action" check:"ids" empty:"true"`
	//心跳超时时间
	// 超出时间没有通讯则判定掉线
	// 单位: 秒
	ExpireTime int64 `db:"expire_time" json:"expireTime" check:"int64Than0" empty:"true"`
	//设备的预计使用场景
	// 0 public 公共设备 / 1 private 私有设备
	// 如果>1 则为自定义设置，具体由设备驱动识别处理
	UseType int `db:"use_type" json:"useType" check:"intThan0" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateGroup 修改分组
func UpdateGroup(args *ArgsUpdateGroup) (err error) {
	var data FieldsGroup
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM iot_core_group WHERE mark = $1 AND id != $2 AND delete_at < to_timestamp(1000000)", args.Mark, args.ID)
	if err == nil && data.ID > 0 {
		err = errors.New("mark is replace")
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_group SET update_at = NOW(), mark = :mark, name = :name, des = :des, cover_files = :cover_files, action = :action, expire_time = :expire_time, use_type = :use_type, params = :params WHERE id = :id", args)
	if err != nil {
		return
	}
	deleteGroupCache(args.ID)
	return
}

// ArgsDeleteGroup 删除分组参数
type ArgsDeleteGroup struct {
	//ID
	ID int64 `json:"id" check:"id"`
}

// DeleteGroup 删除分组
func DeleteGroup(args *ArgsDeleteGroup) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "iot_core_group", "id", args)
	if err != nil {
		return
	}
	deleteGroupCache(args.ID)
	return
}

// 获取指定ID
func getGroupByID(id int64) (data FieldsGroup) {
	cacheMark := getGroupCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, name, des, cover_files, action, expire_time, use_type, params FROM iot_core_group WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

func getGroupByMark(mark string) (data FieldsGroup) {
	cacheMark := getGroupMarkCacheMark(mark)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, name, des, cover_files, action, expire_time, use_type, params FROM iot_core_group WHERE mark = $1", mark)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getGroupCacheMark(id int64) string {
	return fmt.Sprint("iot:device:group:id:", id)
}

func getGroupMarkCacheMark(mark string) string {
	return fmt.Sprint("iot:device:group:mark:", mark)
}

func deleteGroupCache(id int64) {
	data := getGroupByID(id)
	Router2SystemConfig.MainCache.DeleteMark(getGroupCacheMark(id))
	if data.ID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getGroupMarkCacheMark(data.Mark))
	}
}
