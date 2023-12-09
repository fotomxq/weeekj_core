package BaseToken

import (
	"errors"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetList 获取token列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//来源
	FromInfo CoreSQLFrom.FieldsFrom
	//搜索
	Search string
}

// GetList 获取token列表
func GetList(args *ArgsGetList) (dataList []FieldsTokenType, dataCount int64, err error) {
	where, maps, err := args.FromInfo.GetList("from_info", "from_info", nil)
	if err != nil {
		return
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_token",
		"id",
		"SELECT id, create_at, update_at, expire_at, key, password, from_info, login_info, ip, is_remember FROM core_token WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "expire_at", "ip"},
	)
	return
}

// ArgsGetByID 获取token参数
type ArgsGetByID struct {
	//ID
	ID int64
}

// GetByID 获取token
func GetByID(args *ArgsGetByID) (data FieldsTokenType, err error) {
	//请求数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, expire_at, key, password, from_info, login_info, ip, is_remember FROM core_token WHERE id=$1 AND expire_at >= NOW() LIMIT 1;", args.ID)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
	}
	//反馈
	return
}

type ArgsGetByKey struct {
	//Key
	Key string
}

func GetByKey(args *ArgsGetByKey) (data FieldsTokenType, err error) {
	//请求数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, expire_at, key, password, from_info, login_info, ip, is_remember FROM core_token WHERE key=$1 AND expire_at >= NOW() LIMIT 1;", args.Key)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
	}
	//反馈
	return
}

//通过from/fromID/渠道获取数据

type ArgsGetByFromAndLogin struct {
	//来源
	FromInfo CoreSQLFrom.FieldsFrom
	//登陆渠道
	LoginInfo CoreSQLFrom.FieldsFrom
}

func GetByFromAndLogin(args *ArgsGetByFromAndLogin) (data FieldsTokenType, err error) {
	//构建请求头
	fromInfo, err := args.FromInfo.GetRaw()
	if err != nil {
		return
	}
	loginInfo, err := args.LoginInfo.GetRaw()
	if err != nil {
		return
	}
	//请求数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, expire_at, key, password, from_info, login_info, ip, is_remember FROM core_token WHERE expire_at >= NOW() AND from_info @> $1 AND login_info @> $2 ORDER BY id DESC LIMIT 1;", fromInfo, loginInfo)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
	}
	//反馈
	return
}

// 获取默认过期时间差（秒）
func getTokenExpireConfig() (expireSec int64) {
	//获取缓存设置
	tokenDefaultExpire, err := BaseConfig.GetDataString("TokenDefaultExpire")
	if err != nil {
		expireSec = 600
		return
	}
	//获取和转化配置项
	expireTime, err := CoreFilter.GetTimeByAdd(tokenDefaultExpire)
	if err != nil {
		expireSec = 600
		return
	}
	expireSec = expireTime.Unix() - CoreFilter.GetNowTime().Unix()
	if expireSec < 60 {
		expireSec = 600
	}
	return
}
