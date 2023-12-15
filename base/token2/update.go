package BaseToken2

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// SelectOrgAndBindID 选择组织和组织成员
func SelectOrgAndBindID(id int64, orgID int64, orgBindID int64) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_token2 SET org_id = :org_id, org_bind_id = :org_bind_id WHERE id = :id", map[string]interface{}{
		"id":          id,
		"org_id":      orgID,
		"org_bind_id": orgBindID,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteTokenCache(id)
	//反馈
	return
}

// 更新过期时间
func updateExpire(data *FieldsToken) {
	expireAt := getTokenNewExpire(data.IsRemember)
	_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_token2 SET expire_at = :expire_at WHERE id = :id", map[string]interface{}{
		"id":        data.ID,
		"expire_at": expireAt,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteTokenCache(data.ID)
	//触发过期通知
	BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
		OrgID:      0,
		UserID:     0,
		SystemMark: "base_token2",
		BindID:     data.ID,
		Hash:       "",
		ExpireAt:   expireAt,
	})
}

// 获取token过期时间
func getTokenNewExpire(isRemember bool) time.Time {
	var tokenDefaultExpire string
	tokenDefaultExpire, _ = BaseConfig.GetDataString("TokenDefaultExpire")
	if tokenDefaultExpire == "" {
		tokenDefaultExpire = "1h"
	}
	if isRemember {
		tokenDefaultExpire, _ = BaseConfig.GetDataString("TokenRememberExpire")
		if tokenDefaultExpire == "" {
			tokenDefaultExpire = "72h"
		}
	}
	nowAt := CoreFilter.GetNowTimeCarbon()
	expireAt, err := CoreFilter.GetTimeByAdd(tokenDefaultExpire)
	if err != nil {
		expireAt = nowAt.AddMinutes(30).Time
	}
	if expireAt.Unix()-nowAt.Time.Unix() < 1800 {
		expireAt = nowAt.AddMinutes(15).Time
	}
	return expireAt
}
