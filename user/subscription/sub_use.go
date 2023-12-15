package UserSubscription

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsUseSub 使用目标订阅参数
type ArgsUseSub struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//使用来源
	UseFrom     string `db:"use_from" json:"useFrom"`
	UseFromName string `db:"use_from_name" json:"useFromName"`
}

// UseSub 使用目标订阅
func UseSub(args *ArgsUseSub) (err error) {
	//获取数据
	var data FieldsSub
	data, err = GetSub(&ArgsGetSub{
		ConfigID: args.ConfigID,
		UserID:   args.UserID,
	})
	if err != nil {
		return
	}
	//检查过期状况
	if data.DeleteAt.Unix() > 1000000 {
		err = errors.New("sub not exist")
		return
	}
	if data.ExpireAt.Unix() < CoreFilter.GetNowTime().Unix() {
		err = errors.New("sub not exist")
		return
	}
	//获取配置
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    data.ConfigID,
		OrgID: data.OrgID,
	})
	if err != nil {
		return
	}
	//检查限制
	if len(configData.Limits) > 0 {
		type countDataType struct {
			//数量
			Count int `db:"count" json:"count"`
		}
		for _, v := range configData.Limits {
			beforeAt := CoreFilter.GetNowTimeCarbon()
			switch v.TimeType {
			case 0:
				beforeAt = beforeAt.SubHour()
			case 1:
				beforeAt = beforeAt.SubDay().StartOfDay()
			case 2:
				beforeAt = beforeAt.SubWeek().StartOfWeek().StartOfDay()
			case 3:
				beforeAt = beforeAt.SubMonth().StartOfMonth().StartOfDay()
			case 4:
				beforeAt = beforeAt.SubYear().StartOfYear().StartOfMonth().StartOfDay()
			default:
				continue
			}
			var countData countDataType
			if err = Router2SystemConfig.MainDB.Get(&countData, "SELECT COUNT(id) FROM user_sub_log WHERE config_id = $1 AND create_at >= $2", configData.ID, beforeAt.Time); err != nil {
				err = nil
				continue
			}
			if countData.Count < v.Count {
				continue
			}
			err = errors.New("sub limit")
			return
		}
	}
	//通过后反馈
	// 记录日志
	err = appendLog(data.OrgID, data.ConfigID, data.UserID, args.UseFrom, fmt.Sprint("[", args.UseFromName, "]使用[", configData.Title, "]订阅"))
	if err != nil {
		return
	}
	return
}
