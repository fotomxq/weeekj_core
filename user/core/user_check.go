package UserCore

import CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"

// CheckUserHaveEmail 检查用户是否具备email
func CheckUserHaveEmail(userInfo *FieldsUserType) bool {
	return userInfo.Email != "" && CoreSQL.CheckTimeHaveData(userInfo.EmailVerify)
}

// CheckUserHavePhone 检查用户是否具备手机号
func CheckUserHavePhone(userInfo *FieldsUserType) bool {
	return userInfo.Phone != "" && CoreSQL.CheckTimeHaveData(userInfo.PhoneVerify)
}
