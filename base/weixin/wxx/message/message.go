package BaseWeixinWXXMessage

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

var (
	//默认过期时间
	defaultExpireTime = "72h"
)

// ArgsCreate 写入新的数据参数
type ArgsCreate struct {
	//组织ID
	OrgID int64
	//用户ID
	UserID int64
	//用户OpenID
	OpenID string
	//表单ID
	FormID string
}

// Create 写入新的数据
// 系统将自动去重处理，如果重复则拒绝
func Create(args *ArgsCreate) (err error) {
	if _, err = getByOpenIDAndFormID(args.OrgID, args.UserID, args.OpenID, args.FormID); err != nil {
		return errors.New("data is exist")
	}
	var expireAt time.Time
	expireAt, err = CoreFilter.GetTimeByAdd(defaultExpireTime)
	if err != nil {
		return errors.New("expire time is error, " + err.Error())
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_weixin_wxx_message (expire_at, org_id, user_id, open_id, from_id) VALUES (:expire_at, :org_id, :user_id, :open_id, :from_id)", map[string]interface{}{
		"expire_at": expireAt,
		"org_id":    args.OrgID,
		"user_id":   args.UserID,
		"open_id":   args.OpenID,
		"form_id":   args.FormID,
	})

	return
}

// getByOpenID 获取任意一个formID
// 不需要对外广播，该方法只会在微信接口处理内部使用
func getByOpenID(orgID, userID int64, openID string) (formID string, err error) {
	var data FieldsWeixinMessageType
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, org_id, user_id, open_id, from_id FROM core_weixin_wxx_message WHERE org_id = $1 AND user_id = $2 AND open_id = $3 LIMIT 1", orgID, userID, openID)
	if err == nil {
		_ = deleteByOpenIDAndFormID(data.ID)
	}
	return data.FormID, err
}

// getByOpenIDAndFormID 通过openid和formID获取数据
func getByOpenIDAndFormID(orgID, userID int64, openID, formID string) (data FieldsWeixinMessageType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, org_id, user_id, open_id, from_id FROM core_weixin_wxx_message WHERE org_id = $1 AND user_id = $2 AND open_id = $3 AND from_id = $4", orgID, userID, openID, formID)
	return
}

// deleteByOpenIDAndFormID 删除数据
func deleteByOpenIDAndFormID(id int64) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_weixin_wxx_message", "id", map[string]interface{}{
		"id": id,
	})
	return err
}
