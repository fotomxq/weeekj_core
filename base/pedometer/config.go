package PedometerCore

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 设置特定模块的默认值
// 本模块仅识别system，用于匹配
type ArgsSetConfig struct {
	//标识码
	// 对应createInfo.system
	Mark string
	//计数
	Count int
	//默认过期时间
	ExpireAdd string
	//最小数
	MinCount int
	//最大数
	MaxCount int
	//是否增加
	IsAdd bool
}

func SetConfig(args *ArgsSetConfig) (err error) {
	_, err = getConfig(args.Mark)
	if err == nil {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_pedometer_config SET count=:count, default_expire=:default_expire, min_count=:min_count, max_count=:max_count, is_add=:is_add WHERE mark = :mark;", map[string]interface{}{
			"mark":           args.Mark,
			"count":          args.Count,
			"default_expire": args.ExpireAdd,
			"min_count":      args.MinCount,
			"max_count":      args.MaxCount,
			"is_add":         args.IsAdd,
		})
		if err != nil {
			err = errors.New("update data, " + err.Error())
			return
		}
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_pedometer_config(mark, count, default_expire, min_count, max_count, is_add) VALUES(:mark, :count, :default_expire, :min_count, :max_count, :is_add)", map[string]interface{}{
			"mark":           args.Mark,
			"count":          args.Count,
			"default_expire": args.ExpireAdd,
			"min_count":      args.MinCount,
			"max_count":      args.MaxCount,
			"is_add":         args.IsAdd,
		})
		if err != nil {
			return errors.New("insert data, " + err.Error())
		}
	}
	return
}

// 根据from system获取默认值设定
// 其他值将自动抹去
func getConfig(mark string) (data FieldsPedometerConfigType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT mark, count, default_expire, min_count, max_count, is_add FROM core_pedometer_config WHERE mark = $1;", mark)
	return
}
