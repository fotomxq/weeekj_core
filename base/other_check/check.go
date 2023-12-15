package BaseOtherCheck

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"strings"
	"time"
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsCheck, dataCount int64, err error) {
	where := "true"
	maps := map[string]interface{}{}
	tableName := "core_other_check"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, expire_at, url, data FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "expire_at"},
	)
	return
}

// ArgsCreate 申请新的数据参数
type ArgsCreate struct {
	//路由地址
	URL string `db:"url" json:"url"`
	//数据
	Data string `db:"data" json:"data"`
}

// Create 申请新的数据
func Create(args *ArgsCreate) (data FieldsCheck, err error) {
	if args.URL == "" {
		err = errors.New("url is empty")
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_other_check WHERE url = $1 AND expire_at >= NOW()", args.URL)
	if err == nil && data.ID > 0 {
		err = errors.New("data is exist")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_other_check", "INSERT INTO core_other_check (expire_at, url, data) VALUES (:expire_at,:url,:data)", map[string]interface{}{
		"expire_at": CoreFilter.GetNowTimeCarbon().AddHours(3).Time,
		"url":       args.URL,
		"data":      args.Data,
	}, &data)
	return
}

// ArgsUpdateExpire 延长配置参数
type ArgsUpdateExpire struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//过期时间
	ExpireAt string `db:"expire_at" json:"expireAt" check:"isoTime"`
}

// UpdateExpire 延长配置
func UpdateExpire(args *ArgsUpdateExpire) (err error) {
	var expireAt time.Time
	expireAt, err = CoreFilter.GetTimeByISO(args.ExpireAt)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_other_check SET expire_at = :expire_at WHERE id = :id", map[string]interface{}{
		"id":        args.ID,
		"expire_at": expireAt,
	})
	return
}

// ArgsGetURL 检查指定路由设置参数
type ArgsGetURL struct {
	//路由地址
	URL string `db:"url" json:"url"`
}

// GetURL 检查指定路由设置
func GetURL(args *ArgsGetURL) (result string, err error) {
	//路由地址
	if args.URL == "" {
		err = errors.New("url error")
		return
	}
	urlPath := strings.Split(args.URL, "/")
	if len(urlPath) < 1 {
		err = errors.New("url error")
		return
	}
	//验证路由
	var data FieldsCheck
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, expire_at, data FROM core_other_check WHERE url = $1 AND expire_at > NOW()", urlPath[len(urlPath)-1])
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	result = data.Data
	return
}
