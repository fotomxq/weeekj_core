package UserSubscription

import (
	"errors"
	"fmt"
	BaseExpireTip "gitee.com/weeekj/weeekj_core/v5/base/expire_tip"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	OrgUserMod "gitee.com/weeekj/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "gitee.com/weeekj/weeekj_core/v5/service/order/mod"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	UserMessageMod "gitee.com/weeekj/weeekj_core/v5/user/message/mod"
	"github.com/golang-module/carbon"
	"time"
)

// ArgsSetSub 设置订阅信息参数
type ArgsSetSub struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//新的到期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime"`
	//是否为继续订阅
	// 否则将覆盖过期时间
	HaveExpire bool `db:"have_expire" json:"haveExpire" check:"bool"`
	//使用来源
	UseFrom     string `db:"use_from" json:"useFrom"`
	UseFromName string `db:"use_from_name" json:"useFromName"`
}

// SetSub 设置订阅信息
func SetSub(args *ArgsSetSub) (err error) {
	logAppend := "user sub set sub, "
	//获取配置
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    args.ConfigID,
		OrgID: args.OrgID,
	})
	if err != nil || configData.ID < 1 {
		err = errors.New(fmt.Sprint("config not exist, ", err))
		return
	}
	//获取数据
	var subData FieldsSub
	subData, err = GetSub(&ArgsGetSub{
		ConfigID: args.ConfigID,
		UserID:   args.UserID,
	})
	if err == nil && subData.ID > 0 {
		//计算新的过期时间
		if args.HaveExpire {
			var newExpire carbon.Carbon
			if subData.ExpireAt.Unix() < CoreFilter.GetNowTime().Unix() {
				newExpire = CoreFilter.GetNowTimeCarbon()
			} else {
				newExpire = CoreFilter.GetCarbonByTime(subData.ExpireAt)
			}
			newExpire = newExpire.AddSeconds(int(args.ExpireAt.Unix() - CoreFilter.GetNowTime().Unix()))
			args.ExpireAt = newExpire.Time
		}
		//更新数据
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_sub SET update_at = NOW(), expire_at = :expire_at, params = :params, delete_at = to_timestamp(0) WHERE id = :id", map[string]interface{}{
			"expire_at": args.ExpireAt,
			"id":        subData.ID,
			"params":    subData.Params,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("insert failed, ", err))
			return
		}
		// 记录日志
		err = appendLog(args.OrgID, args.ConfigID, args.UserID, args.UseFrom, fmt.Sprint("[", args.UseFromName, "]给与[", configData.Title, "]订阅新到期时间: ", args.ExpireAt))
		if err != nil {
			err = errors.New(fmt.Sprint("insert log failed, ", err))
			return
		}
	} else {
		//创建数据
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_sub (expire_at, org_id, config_id, user_id, params) VALUES (:expire_at,:org_id,:config_id,:user_id,:params)", map[string]interface{}{
			"expire_at": args.ExpireAt,
			"org_id":    configData.OrgID,
			"config_id": args.ConfigID,
			"user_id":   args.UserID,
			"params":    subData.Params,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("insert failed, ", err))
			return
		}
		// 记录日志
		err = appendLog(args.OrgID, args.ConfigID, args.UserID, args.UseFrom, fmt.Sprint("[", args.UseFromName, "]新增[", configData.Title, "]订阅，到期时间: ", args.ExpireAt))
		if err != nil {
			err = errors.New(fmt.Sprint("insert log failed, ", err))
			return
		}
	}
	//获取会员数据包
	subData, err = GetSub(&ArgsGetSub{
		ConfigID: args.ConfigID,
		UserID:   args.UserID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get sub data, ", err))
		return
	}
	//如果没有过期
	if subData.ExpireAt.Unix() > CoreFilter.GetNowTime().Unix() {
		//授权用户组
		if len(configData.UserGroups) > 0 {
			var userData UserCore.FieldsUserType
			userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
				ID:    args.UserID,
				OrgID: -1,
			})
			if err == nil {
				for _, v := range configData.UserGroups {
					newExpireAt := CoreFilter.GetNowTimeCarbon()
					for _, v2 := range userData.Groups {
						if v == v2.GroupID && v2.ExpireAt.Unix() > newExpireAt.Time.Unix() {
							newExpireAt = newExpireAt.CreateFromGoTime(v2.ExpireAt)
						}
					}
					if args.ExpireAt.Unix() > CoreFilter.GetNowTime().Unix() {
						newExpireAt = newExpireAt.AddSeconds(int(args.ExpireAt.Unix() - CoreFilter.GetNowTime().Unix()))
					}
					err = UserCore.UpdateUserGroupByID(&UserCore.ArgsUpdateUserGroupByID{
						ID:       args.UserID,
						OrgID:    configData.OrgID,
						GroupID:  v,
						ExpireAt: newExpireAt.Time,
						IsRemove: false,
					})
					if err != nil {
						CoreLog.Error(logAppend, "update user group failed, user id: ", args.UserID, ", err: ", err)
						err = nil
					}
				}
			} else {
				CoreLog.Error(logAppend, "get user data, user id: ", subData.UserID, ", err: ", err)
				err = nil
			}
		}
	} else {
		//收尾工作
		expireSubLast(&configData, &subData)
	}
	//设置过期通知
	err = BaseExpireTip.AppendTip(&BaseExpireTip.ArgsAppendTip{
		OrgID:      subData.OrgID,
		UserID:     subData.UserID,
		SystemMark: "user_sub",
		BindID:     subData.ConfigID,
		Hash:       getSubHash(&subData),
		ExpireAt:   subData.ExpireAt,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("append tip, ", err))
		return
	}
	//强制更新组织用户数据
	if args.OrgID > 0 {
		OrgUserMod.PushUpdateUserData(args.OrgID, args.UserID)
	}
	//推送用户消息
	if CoreSQL.CheckTimeThanNow(subData.ExpireAt) {
		UserMessageMod.CreateSystemToUser(time.Time{}, subData.UserID, "会员订阅成功", fmt.Sprint("您已经成功订阅了(", configData.Title, ")会员，到期时间为", CoreFilter.GetTimeToDefaultTime(subData.ExpireAt), "。"), nil, nil)
	}
	//反馈
	return
}

// AddUserSubAny 为任意一个会员叠加时间
// 将从平台或组织下，抽取任意一个配置来给用户叠加会员
// 该设计为过度设计，v2会员系统将没有会员配置得定义，所以考虑到此设计将不对多个配置进行处理，只识别一个配置即可
func AddUserSubAny(orgID int64, userID int64, addSec int64) (err error) {
	var configData FieldsConfig
	err = Router2SystemConfig.MainDB.Get(&configData, "SELECT id FROM user_sub_config WHERE org_id = $1 AND delete_at < to_timestamp(1000000) LIMIT 1", orgID)
	if err != nil || configData.ID < 1 {
		return
	}
	err = SetSub(&ArgsSetSub{
		OrgID:       orgID,
		ConfigID:    configData.ID,
		UserID:      userID,
		ExpireAt:    CoreFilter.GetNowTimeCarbon().AddSeconds(int(addSec)).Time,
		HaveExpire:  true,
		UseFrom:     "market_giving",
		UseFromName: "营销赠送",
	})
	if err != nil {
		return
	}
	return
}

// 叠加会员时间
type argsSetSubAdd struct {
	ConfigID int64 `json:"configID"`
	UserID   int64 `json:"userID"`
	Unit     int   `json:"unit"`
	OrderID  int64 `json:"orderID"`
}

func setSubAdd(args *argsSetSubAdd) (err error) {
	//获取配置
	var configData FieldsConfig
	configData, err = getConfigByID(args.ConfigID)
	if err != nil || configData.ID < 1 {
		err = errors.New(fmt.Sprint("get config, ", err))
		return
	}
	//新的到期时间
	expireAt := CoreFilter.GetNowTimeCarbon()
	//获取已有的续订
	var subData FieldsSub
	subData, err = GetSub(&ArgsGetSub{
		ConfigID: configData.ID,
		UserID:   args.UserID,
	})
	if err == nil {
		//重构初始化时间
		if subData.ExpireAt.Unix() >= CoreFilter.GetNowTime().Unix() {
			expireAt = CoreFilter.GetNowTimeCarbon().CreateFromGoTime(subData.ExpireAt)
		}
	}
	//根据单位时间，续订时间
	switch configData.TimeType {
	case 0:
		expireAt = expireAt.AddHours(configData.TimeN * args.Unit)
	case 1:
		expireAt = expireAt.AddDays(configData.TimeN * args.Unit)
	case 2:
		expireAt = expireAt.AddWeeks(configData.TimeN * args.Unit)
	case 3:
		expireAt = expireAt.AddMonths(configData.TimeN * args.Unit)
	case 4:
		expireAt = expireAt.AddYears(configData.TimeN * args.Unit)
	}
	//根据配置的时间设计，设置订阅过期时间
	if err = SetSub(&ArgsSetSub{
		OrgID:       configData.OrgID,
		ConfigID:    configData.ID,
		UserID:      args.UserID,
		ExpireAt:    expireAt.Time,
		HaveExpire:  false,
		UseFrom:     "buy",
		UseFromName: "购买订阅",
	}); err != nil {
		err = errors.New(fmt.Sprint("set sub, err: ", err))
		return
	}
	//完成订单
	ServiceOrderMod.UpdateFinish(args.OrderID, "用户订阅自动完成订单")
	//强制更新组织用户数据
	OrgUserMod.PushUpdateUserData(configData.OrgID, args.UserID)
	//反馈
	return
}

// 用户订阅过期的处理收尾工作
func expireSubLast(configData *FieldsConfig, subData *FieldsSub) {
	logAppend := "user sub expire sub last, "
	var err error
	//处理用户组
	if len(configData.UserGroups) > 0 {
		var userData UserCore.FieldsUserType
		userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
			ID:    subData.UserID,
			OrgID: -1,
		})
		if err == nil {
			for _, vGroup2 := range userData.Groups {
				isFind := false
				for _, vGroup := range configData.UserGroups {
					if vGroup == vGroup2.GroupID {
						isFind = true
						break
					}
				}
				if isFind {
					if err = UserCore.UpdateUserGroupByID(&UserCore.ArgsUpdateUserGroupByID{
						ID:       subData.UserID,
						OrgID:    subData.OrgID,
						GroupID:  vGroup2.GroupID,
						ExpireAt: vGroup2.ExpireAt,
						IsRemove: true,
					}); err != nil {
						CoreLog.Error(logAppend, "update user group by id, ", err)
						err = nil
						continue
					}
				}
			}
		} else {
			CoreLog.Error(logAppend, "get user data, user id: ", subData.UserID, ", err: ", err)
			err = nil
		}
	}
}
