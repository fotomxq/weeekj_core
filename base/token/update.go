package BaseToken

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsUpdateExpire 更新过期时间参数
type ArgsUpdateExpire struct {
	//ID
	ID int64
	//过期时间
	// 请使用RFC3339Nano结构时间，eg: 2020-11-03T08:31:13.314Z
	// JS中使用new Date().toISOString()
	ExpireAt string
}

// UpdateExpire 更新过期时间
func UpdateExpire(args *ArgsUpdateExpire) error {
	return UpdateExpireToValue(&ArgsUpdateExpireToValue{
		ID:       args.ID,
		ExpireAt: args.ExpireAt,
	})
}

// ArgsUpdateExpireByKey 根据Key更新过期时间参数
type ArgsUpdateExpireByKey struct {
	//Key
	Key string
	//过期时间
	// 请使用RFC3339Nano结构时间，eg: 2020-11-03T08:31:13.314Z
	// JS中使用new Date().toISOString()
	ExpireAt string
}

// UpdateExpireByKey 根据Key更新过期时间
func UpdateExpireByKey(args *ArgsUpdateExpireByKey) error {
	//获取数据
	data, err := GetByKey(&ArgsGetByKey{
		Key: args.Key,
	})
	if err != nil {
		return err
	}
	//生成过期时间
	var expireTime time.Time
	expireTime, err = time.Parse(time.RFC3339Nano, args.ExpireAt)
	if err != nil {
		err = errors.New("cannot get expire time, " + err.Error())
		return err
	}
	//如果低于生成的过期时间，则忽略
	if data.ExpireAt.Unix() >= expireTime.Unix() {
		return nil
	}
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_token SET expire_at=:expire_at WHERE key=:key AND expire_at >= NOW();", map[string]interface{}{
		"expire_at": expireTime,
		"key":       args.Key,
	})
	if err != nil {
		return err
	}
	//反馈
	return nil
}

// ArgsUpdateExpireToValue 修改token的过期时间参数
type ArgsUpdateExpireToValue struct {
	//ID
	ID int64
	//过期时间
	// 请使用RFC3339Nano结构时间，eg: 2020-11-03T08:31:13.314Z
	// JS中使用new Date().toISOString()
	ExpireAt string
}

// UpdateExpireToValue 修改token的过期时间
func UpdateExpireToValue(args *ArgsUpdateExpireToValue) error {
	//获取数据
	data, err := GetByID(&ArgsGetByID{
		ID: args.ID,
	})
	if err != nil {
		return err
	}
	//生成过期时间
	var expireTime time.Time
	expireTime, err = time.Parse(time.RFC3339Nano, args.ExpireAt)
	if err != nil {
		err = errors.New("cannot get expire time, " + err.Error())
		return err
	}
	//如果低于生成的过期时间，则忽略
	if data.ExpireAt.Unix() >= expireTime.Unix() {
		return nil
	}
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_token SET expire_at=:expire_at WHERE id=:id AND expire_at >= NOW();", map[string]interface{}{
		"expire_at": expireTime,
		"id":        args.ID,
	})
	if err != nil {
		return err
	}
	//反馈
	return nil
}
