package BaseToken

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// Deprecated: 建议采用Create2
// ArgsCreate 创建新的参数
type ArgsCreate struct {
	//创建来源
	FromInfo CoreSQLFrom.FieldsFrom
	//登陆渠道
	LoginInfo CoreSQLFrom.FieldsFrom
	//登陆key长度限制
	LimitKeyLen int
	//IP地址
	IP string
	//过期时间
	// 请使用RFC3339Nano结构时间，eg: 2020-11-03T08:31:13.314Z
	// JS中使用new Date().toISOString()
	ExpireAt string
	//是否记住我
	IsRemember bool `db:"is_remember" json:"isRemember"`
}

// Deprecated: 建议采用Create2
// Create 创建新的
func Create(args *ArgsCreate) (data FieldsTokenType, errCode string, err error) {
	//尝试获取数据
	data, err = GetByFromAndLogin(&ArgsGetByFromAndLogin{
		FromInfo:  args.FromInfo,
		LoginInfo: args.LoginInfo,
	})
	if err == nil {
		//反馈之前的登录头
		return
	}
	//自动构建key和password
	var key string
	key, err = getKeyNoReplace(args.LimitKeyLen)
	if err != nil {
		errCode = "get_new_key"
		err = errors.New("cannot get rand key, " + err.Error())
		return
	}
	var password string
	password, err = CoreFilter.GetRandStr3(30)
	if err != nil {
		errCode = "get_new_password"
		err = errors.New("cannot get rand password, " + err.Error())
		return
	}
	//生成过期时间
	var expireTime time.Time
	expireTime, err = time.Parse(time.RFC3339Nano, args.ExpireAt)
	if err != nil {
		errCode = "get_expire_time"
		err = errors.New("cannot get expire time, " + err.Error())
		return
	}
	//生成数据
	var lastID int64
	sets := map[string]interface{}{
		"expire_at":   expireTime,
		"key":         key,
		"password":    password,
		"from_info":   args.FromInfo,
		"login_info":  args.LoginInfo,
		"ip":          args.IP,
		"is_remember": args.IsRemember,
	}
	lastID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_token(expire_at, key, password, from_info, login_info, ip, is_remember) VALUES(:expire_at, :key, :password, :from_info, :login_info, :ip, :is_remember)", sets)
	if err == nil {
		data, err = GetByID(&ArgsGetByID{
			ID: lastID,
		})
		if err != nil {
			errCode = "insert_data"
			return
		}
	} else {
		errCode = "insert_data"
	}
	return
}

// Create2 创建新的
func Create2(args *ArgsCreate) (data FieldsTokenType, errCode string, err error) {
	//尝试获取数据
	data, err = GetByFromAndLogin(&ArgsGetByFromAndLogin{
		FromInfo:  args.FromInfo,
		LoginInfo: args.LoginInfo,
	})
	if err == nil {
		//反馈之前的登录头
		return
	}
	//自动构建key和password
	var key string
	key, err = getKeyNoReplace(args.LimitKeyLen)
	if err != nil {
		errCode = "err_key"
		err = errors.New("cannot get rand key, " + err.Error())
		return
	}
	var password string
	password, err = CoreFilter.GetRandStr3(30)
	if err != nil {
		errCode = "err_key"
		err = errors.New("cannot get rand password, " + err.Error())
		return
	}
	//生成过期时间
	var expireTime time.Time
	expireTime, err = time.Parse(time.RFC3339Nano, args.ExpireAt)
	if err != nil {
		errCode = "err_time"
		err = errors.New("cannot get expire time, " + err.Error())
		return
	}
	//生成数据
	var lastID int64
	sets := map[string]interface{}{
		"expire_at":   expireTime,
		"key":         key,
		"password":    password,
		"from_info":   args.FromInfo,
		"login_info":  args.LoginInfo,
		"ip":          args.IP,
		"is_remember": args.IsRemember,
	}
	lastID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_token(expire_at, key, password, from_info, login_info, ip, is_remember) VALUES(:expire_at, :key, :password, :from_info, :login_info, :ip, :is_remember)", sets)
	if err == nil {
		data, err = GetByID(&ArgsGetByID{
			ID: lastID,
		})
		if err != nil {
			errCode = "err_insert"
			return
		}
	} else {
		errCode = "err_insert"
	}
	return
}
