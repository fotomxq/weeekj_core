package PedometerCore

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// 获取数据
func GetData(args CoreSQLFrom.FieldsFrom) (data FieldsPedometerType, err error) {
	err = args.GetFromOne(Router2SystemConfig.MainDB.DB, "core_pedometer", "id, create_at, update_at, expire_at, create_info, count", "create_info", &data)
	return
}

// 进一位
// 自动协调进一步或退一步
func NextData(args CoreSQLFrom.FieldsFrom) (int, error) {
	count, err := saveData(args, 1, false, false)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 回退一步
func PrevData(args CoreSQLFrom.FieldsFrom) (int, error) {
	count, err := saveData(args, -1, false, false)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 归位数据
func ReturnData(args CoreSQLFrom.FieldsFrom) (int, error) {
	count, err := saveData(args, 0, true, false)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 更改指定的位数
type ArgsSetData struct {
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//个数
	Count int
}

func SetData(args *ArgsSetData) (int, error) {
	count, err := saveData(args.CreateInfo, args.Count, false, true)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 清理某个系统下所有内容
type ArgsClearData struct {
	//来源系统
	System string `db:"create_info_system"`
}

func ClearData(args *ArgsClearData) error {
	creatInfo := CoreSQLFrom.FieldsFrom{
		System: args.System,
		ID:     0,
		Mark:   "",
		Name:   "",
	}
	maps, err := creatInfo.GetMaps("create_info", nil)
	if err != nil {
		return err
	}
	if _, err := CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_pedometer", "create_info @> :create_info", maps); err != nil {
		return err
	}
	return nil
}

// 检查当前有多少个
func GetCount(args CoreSQLFrom.FieldsFrom) int {
	data, err := GetData(args)
	if err != nil {
		return 0
	}
	return data.Count
}

// 达到max/min时反馈true
func CheckData(args CoreSQLFrom.FieldsFrom) bool {
	data, err := GetData(args)
	if err != nil {
		return false
	}
	configData, err := getConfig(args.Mark)
	if err != nil {
		//如果配置不存在则自动放行，避免异常
		return false
	}
	if configData.IsAdd {
		if data.Count >= configData.MaxCount {
			return true
		}
	} else {
		if data.Count <= configData.MinCount {
			return true
		}
	}
	return false
}

/*
*
内部统一保存和修改函数
@params needReturn 是否需要回归原始值
@params needSet 是否为直接指定一个值，而不是递增或递减
*/
func saveData(createInfo CoreSQLFrom.FieldsFrom, stepCount int, needReturn bool, needSet bool) (int, error) {
	//获取配置
	configData, err := getConfig(createInfo.System)
	if err != nil {
		return 0, errors.New("get default by from system, " + err.Error())
	}
	//生成新的过期时间
	expireTime, err := CoreFilter.GetTimeByAdd(configData.DefaultExpire)
	if err != nil {
		expireTime = CoreFilter.GetNowTime().Add(time.Minute * 30)
	}
	//尝试获取数据
	data, err := GetData(createInfo)
	//写入数据
	if err == nil {
		var count int
		if needReturn {
			count = configData.Count
		} else {
			if needSet {
				//特殊处理方案，上一级已经指定了特定数字，不需要修改
				count = stepCount
			} else {
				if configData.IsAdd {
					data.Count = data.Count + stepCount
				} else {
					data.Count = data.Count - stepCount
				}
				if data.Count < configData.MinCount {
					data.Count = configData.MinCount
				}
				if data.Count > configData.MaxCount {
					data.Count = configData.MaxCount
				}
				count = data.Count
			}
		}
		if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_pedometer SET expire_at=:expire_at, update_at=NOW(), count = :count WHERE id=:id", map[string]interface{}{
			"expire_at": expireTime,
			"count":     count,
			"id":        data.ID,
		}); err != nil {
			return 0, errors.New("update data, " + err.Error())
		}
		return count, nil
	} else {
		var count int
		if needSet {
			count = stepCount
		} else {
			count = configData.Count
		}
		var createData string
		createData, err = createInfo.GetRaw()
		if err != nil {
			return 0, err
		}
		if _, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_pedometer(expire_at, create_info, count) VALUES(:expire_at, :create_info, :count)", map[string]interface{}{
			"expire_at":   expireTime,
			"create_info": createData,
			"count":       count,
		}); err != nil {
			return 0, errors.New("insert data, " + err.Error())
		}
		return count, nil
	}
}
