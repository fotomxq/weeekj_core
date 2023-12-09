package UserLogin

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

/**
流程设计：
1、网页或其他端生成二维码，提供给其他渠道扫码匹配
2、用户持有手机端APP扫码，该扫码行为会同时提供用户ID信息
3、扫码完成后，用户ID会被写入表内
4、网页端可以通过被动接口，验证是否完成扫码，验证后将自动删除二维码
*/

// ArgsMakeQrcode 构建扫码的二维码参数
type ArgsMakeQrcode struct {
	//会话ID
	TokenID int64 `json:"tokenID"`
	//系统类型
	SystemMark string `json:"systemMark"`
}

// MakeQrcode 构建扫码的二维码
// 会话新增新的二维码后，将自动删除不是同一个系统下旧的二维码
func MakeQrcode(args *ArgsMakeQrcode) (dataQrcode DataQrcode, err error) {
	var data FieldsQrcode
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, user_id, key FROM user_login_qrcode WHERE token_id = $1 AND system_mark = $2 AND expire_at >= NOW()", args.TokenID, args.SystemMark)
	if err == nil && data.UserID < 1 {
		dataQrcode = DataQrcode{
			Mark: globAppName,
			ID:   data.ID,
			Key:  data.Key,
		}
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "user_login_qrcode", "id", map[string]interface{}{
		"id": data.ID,
	})
	//生成新的key
	var newKey string
	newKey, err = CoreFilter.GetRandStr3(20)
	if err != nil || newKey == "" {
		err = errors.New(fmt.Sprint("key is empty, ", err))
		return
	}
	//创建新的数据
	var lastID int64
	lastID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO user_login_qrcode (expire_at, token_id, system_mark, key, user_id) VALUES (:expire_at, :token_id, :system_mark, :key, 0)", map[string]interface{}{
		"expire_at":   CoreFilter.GetNowTimeCarbon().AddMinutes(3).Time,
		"token_id":    args.TokenID,
		"system_mark": args.SystemMark,
		"key":         newKey,
	})
	if err != nil {
		return
	}
	dataQrcode = DataQrcode{
		Mark: globAppName,
		ID:   lastID,
		Key:  newKey,
	}
	return
}

// 完成扫码操作
type ArgsFinishQrcode struct {
	//ID
	ID int64 `json:"id"`
	//Key key
	Key string `json:"key"`
	//UserID 用户ID
	UserID int64 `json:"userID"`
}

func FinishQrcode(args *ArgsFinishQrcode) (err error) {
	var data FieldsQrcode
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, user_id FROM user_login_qrcode WHERE id = $1 AND key = $2 AND expire_at >= NOW()", args.ID, args.Key)
	if err != nil {
		return
	}
	if data.UserID > 0 {
		err = errors.New("qrcode is used")
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_login_qrcode SET user_id = :user_id WHERE id = :id", map[string]interface{}{
		"id":      args.ID,
		"user_id": args.UserID,
	})
	return
}

// ArgsCheckQrcode 验证扫码行为参数
type ArgsCheckQrcode struct {
	//ID
	ID int64 `json:"id"`
	//key
	Key string `json:"key"`
}

// CheckQrcode 验证扫码行为
func CheckQrcode(args *ArgsCheckQrcode) (userID int64, err error) {
	var data FieldsQrcode
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, user_id FROM user_login_qrcode WHERE id = $1 AND key = $2 AND expire_at >= NOW()", args.ID, args.Key)
	if err != nil {
		return
	}
	if data.UserID < 1 {
		err = errors.New("wait qrcode")
		return
	}
	userID = data.UserID
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "user_login_qrcode", "id", map[string]interface{}{
		"id": args.ID,
	})
	return
}
