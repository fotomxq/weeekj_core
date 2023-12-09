package BaseToken2

import (
	"errors"
	BaseExpireTip "gitee.com/weeekj/weeekj_core/v5/base/expire_tip"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsCreate 创建token参数
type ArgsCreate struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织绑定成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//登录渠道
	LoginFrom string `db:"login_from" json:"loginFrom"`
	//IP地址
	IP string `db:"ip" json:"ip"`
	//key
	// 钥匙，用于配对
	// 用户模式下，采用用户密码+SHA1计算后得出
	// 设备模式下，直接使用设备的key实现
	// 如果不提供，则自动生成并反馈数据
	Key string `db:"key" json:"key"`
	//是否记住我
	// 会延长过期时间
	IsRemember bool `db:"is_remember" json:"isRemember"`
}

// Create 创建token
func Create(args *ArgsCreate) (tokenID int64, errCode string, err error) {
	//获取相同渠道数据
	if args.UserID > 0 || args.DeviceID > 0 {
		data := GetByFrom(args.UserID, args.DeviceID, args.LoginFrom)
		if data.ID > 0 {
			DeleteToken(data.ID)
		}
	}
	//获取过期时间
	expireAt := getTokenNewExpire(args.IsRemember)
	//自动构建key和password
	if args.Key == "" {
		args.Key, err = getKeyNoReplace(20)
		if err != nil {
			errCode = "err_key"
			err = errors.New("cannot get rand key, " + err.Error())
			return
		}
	}
	//生成数据
	tokenID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_token2(expire_at, key, user_id, org_id, org_bind_id, device_id, login_from, ip, is_remember) VALUES(:expire_at, :key, :user_id, :org_id, :org_bind_id, :device_id, :login_from, :ip, :is_remember)", map[string]interface{}{
		"expire_at":   expireAt,
		"key":         args.Key,
		"user_id":     args.UserID,
		"org_id":      args.OrgID,
		"org_bind_id": args.OrgBindID,
		"device_id":   args.DeviceID,
		"login_from":  args.LoginFrom,
		"ip":          args.IP,
		"is_remember": args.IsRemember,
	})
	if err != nil || tokenID < 1 {
		err = errors.New("create token failed")
		errCode = "err_insert"
		return
	}
	//触发过期通知
	BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
		OrgID:      0,
		UserID:     0,
		SystemMark: "base_token2",
		BindID:     tokenID,
		Hash:       "",
		ExpireAt:   expireAt,
	})
	//反馈
	return
}

// 生成一个不重复的key
func getKeyNoReplace(limitKeyLen int) (string, error) {
	var key string
	var err error
	tryStep := 1
	maxTry := 10
	for {
		key, err = CoreFilter.GetRandStr3(limitKeyLen)
		if err != nil {
			return "", errors.New("cannot get new rand data, " + err.Error())
		}
		//检查是否重复
		_, err = getByKey(key)
		if err == nil {
			tryStep += 1
			if tryStep > maxTry {
				return "", errors.New("try is too many, pls delete some key or set limit more")
			}
			continue
		}
		return key, nil
	}
}

// 通过key获取token
func getByKey(key string) (data FieldsToken, err error) {
	//请求数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_token WHERE key=$1 AND expire_at >= NOW() LIMIT 1;", key)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
	}
	//反馈
	return
}
