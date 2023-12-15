package UserCore

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsReSendUserEmail 重新发送验证码参数
type ArgsReSendUserEmail struct {
	//用户ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// ReSendUserEmail 重新发送验证码
func ReSendUserEmail(args *ArgsReSendUserEmail) (err error) {
	var data FieldsUserType
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id FROM user_core WHERE delete_at < to_timestamp(1000000) AND id = $1 AND ($2 < 1 OR org_id = $2) AND email != '' LIMIT 1", args.ID, args.OrgID)
	if err != nil {
		return
	}
	pushNatsUserEmailWait(args.ID)
	return
}

// ArgsCheckUserEmailVCode 验证用户邮件验证码参数
type ArgsCheckUserEmailVCode struct {
	//用户ID
	ID int64 `db:"id" json:"id" check:"id"`
	//验证码
	VCode string `db:"vcode" json:"vcode" check:"mark"`
}

// CheckUserEmailVCode 验证用户邮件验证码
func CheckUserEmailVCode(args *ArgsCheckUserEmailVCode) (errCode string, err error) {
	//获取用户数据
	userData := getUserByID(args.ID)
	if userData.ID < 1 || CoreSQL.CheckTimeThanNow(userData.DeleteAt) {
		errCode = "no_user"
		err = errors.New("no data")
		return
	}
	//识别验证码
	var emailData FieldsRegWaitEmail
	err = Router2SystemConfig.MainDB.Get(&emailData, "SELECT id, expire_at, vcode FROM user_reg_wait_email WHERE user_id = $1 AND delete_at < to_timestamp(1000000) ORDER BY id DESC LIMIT 1", args.ID)
	if emailData.ExpireAt.Unix() < CoreFilter.GetNowTime().Unix() {
		errCode = "expire"
		err = errors.New("is expire")
		return
	}
	if emailData.VCode != args.VCode {
		errCode = "vcode"
		err = errors.New("vcode error")
		return
	}
	//修改用户状态
	if userData.Status != 2 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET status = 2, email_verify = NOW() WHERE id = :id", map[string]interface{}{
			"id": userData.ID,
		})
		if err != nil {
			errCode = "update_status"
			err = errors.New(fmt.Sprint("update user status, ", err))
			return
		}
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET email_verify = NOW() WHERE id = :id", map[string]interface{}{
			"id": userData.ID,
		})
		if err != nil {
			errCode = "update_status"
			err = errors.New(fmt.Sprint("update user status, ", err))
			return
		}
	}
	//删除缓冲
	deleteUserCache(args.ID)
	//反馈成功
	return
}
