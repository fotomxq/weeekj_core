package OrgSubscription

import (
	"errors"
	"fmt"
	BaseExpireTip "gitee.com/weeekj/weeekj_core/v5/base/expire_tip"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "gitee.com/weeekj/weeekj_core/v5/service/order/mod"
	"github.com/golang-module/carbon"
)

// ArgsSetSub 设置组织的订阅参数
type ArgsSetSub struct {
	//配置单位
	ConfigUnit int `db:"config_unit" json:"configUnit" check:"intThan0"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//开通配置
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetSub 设置组织的订阅
func SetSub(args *ArgsSetSub) (errCode string, err error) {
	//获取配置
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID: args.ConfigID,
	})
	if err != nil {
		errCode = "config_not_exist"
		return
	}
	//获取组织数据
	var orgData OrgCoreCore.FieldsOrg
	orgData, err = OrgCoreCore.GetOrg(&OrgCoreCore.ArgsGetOrg{
		ID: args.OrgID,
	})
	if err != nil {
		errCode = "org_not_exist"
		return
	}
	//构建过期时间
	expireAt := CoreFilter.GetNowTimeCarbon()
	var subData FieldsSub
	err = Router2SystemConfig.MainDB.Get(&subData, "SELECT id, expire_at FROM org_sub WHERE org_id = $1 AND config_id = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.ConfigID)
	if err == nil && subData.ID > 0 {
		expireAt = carbon.CreateFromTimestamp(subData.ExpireAt.Unix())
	}
	switch configData.TimeType {
	case 0:
		//小时
		expireAt = expireAt.AddHours(configData.TimeN * args.ConfigUnit)
	case 1:
		//天
		expireAt = expireAt.AddDays(configData.TimeN * args.ConfigUnit)
	case 2:
		//周
		expireAt = expireAt.AddWeeks(configData.TimeN * args.ConfigUnit)
	case 3:
		//月份
		expireAt = expireAt.AddMonths(configData.TimeN * args.ConfigUnit)
	case 4:
		//年
		expireAt = expireAt.AddYears(configData.TimeN * args.ConfigUnit)
	default:
		//无法识别的时间类型，跳出
		err = errors.New("time type error")
		return
	}
	if err == nil {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_sub SET update_at = NOW(), expire_at = :expire_at, params = :params WHERE id = :id", map[string]interface{}{
			"id":        subData.ID,
			"expire_at": expireAt.Time,
			"params":    args.Params,
		})
		if err != nil {
			errCode = "update"
		}
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_sub (expire_at, org_id, config_id, params) VALUES (:expire_at, :org_id, :config_id, :params)", map[string]interface{}{
			"expire_at": expireAt.Time,
			"org_id":    args.OrgID,
			"config_id": args.ConfigID,
			"params":    args.Params,
		})
		if err != nil {
			errCode = "insert"
		}
	}
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&subData, "SELECT id, org_id, config_id, expire_at FROM org_sub WHERE org_id = $1 AND config_id = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.ConfigID)
	if err != nil {
		return
	}
	//写入新的组织功能组
	for _, v := range configData.FuncList {
		isFind := false
		for _, v2 := range orgData.OpenFunc {
			if v == v2 {
				isFind = true
				break
			}
		}
		if !isFind {
			orgData.OpenFunc = append(orgData.OpenFunc, v)
		}
	}
	if err = OrgCoreCore.UpdateOrgFunc(&OrgCoreCore.ArgsUpdateOrgFunc{
		ID:       orgData.ID,
		OpenFunc: orgData.OpenFunc,
	}); err != nil {
		CoreLog.Error("org sub set sub, update org func, org id: ", orgData.ID, ", new func: ", orgData.OpenFunc, ", err: ", err)
		err = nil
	}
	//通知过期时间
	err = BaseExpireTip.AppendTip(&BaseExpireTip.ArgsAppendTip{
		OrgID:      subData.OrgID,
		UserID:     0,
		SystemMark: "org_sub",
		BindID:     subData.ID,
		Hash:       getSubHash(&subData),
		ExpireAt:   subData.ExpireAt,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("append tip, ", err))
		return
	}
	//反馈
	return
}

// argsSetSubAdd 向后续约指定时间参数
type argsSetSubAdd struct {
	ConfigID int64 `json:"configID"`
	OrgID    int64 `json:"orgID"`
	Unit     int   `json:"unit"`
	OrderID  int64 `json:"orderID"`
}

// setSubAdd 向后续约指定时间
func setSubAdd(args *argsSetSubAdd) (err error) {
	//标记会员时间
	if _, err = SetSub(&ArgsSetSub{
		ConfigUnit: args.Unit,
		OrgID:      args.OrgID,
		ConfigID:   args.ConfigID,
		Params:     []CoreSQLConfig.FieldsConfigType{},
	}); err != nil {
		err = errors.New(fmt.Sprint("set sub, ", err))
		return
	}
	//完成订单
	ServiceOrderMod.UpdateFinish(args.OrderID, "用户订阅自动完成订单")
	//反馈
	return
}

// 过期收尾工作
func expireSubLast(configData *FieldsConfig, orgData *OrgCoreCore.FieldsOrg) {
	//查询组织的当前服务，剔除需要删除的服务项目
	var newFuncList []string
	for _, v := range orgData.OpenFunc {
		isFind := false
		for _, v2 := range configData.FuncList {
			if v == v2 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		newFuncList = append(newFuncList, v)
	}
	//更新组织服务内容
	if err := OrgCoreCore.UpdateOrgFunc(&OrgCoreCore.ArgsUpdateOrgFunc{
		ID:       orgData.ID,
		OpenFunc: newFuncList,
	}); err != nil {
		CoreLog.Error("org sub set sub, update org func, org id: ", orgData.ID, ", new func: ", orgData.OpenFunc, ", err: ", err)
		err = nil
	}
}
