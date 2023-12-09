package Router2Mid

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	UserRole "gitee.com/weeekj/weeekj_core/v5/user/role"
	"github.com/gin-gonic/gin"
)

//GetUserDataByToken 通过token获取user数据
/**
action					URL动作绑定
timestamp				时间戳
nonce 					随机码
secretID 				用户唯一info.value值
signatureKey 			用户密码加密后的数据
signatureMethod 	加密方式sha256
key						用户密码
*/
/**
cookie方法先使用一般路由，构建core.session.create，创建后重新建立token.create即可。
本中间件将使用token.from id判断用户ID是否匹配、token.from 判断来源是否为from、其他内容可根据不同路由的特殊性自行调整判断。
*/
func getUserDataByToken(c *gin.Context) (userID int64, err error) {
	//获取token数据
	tokenInfo := GetTokenInfo(c)
	//通过token数据集合，找到用户
	if tokenInfo.UserID < 1 {
		err = errors.New("token from system is not user")
		return
	}
	//TODO: 兼容v2 写入user数据，兼容v2路由处理
	userData, err := UserCore.GetUserData(&UserCore.ArgsGetUserData{
		UserID: tokenInfo.UserID,
	})
	if err != nil {
		err = errors.New("cannot find user info, " + err.Error())
		return
	}
	if userData.Info.ID < 1 {
		err = errors.New("user not exist")
		return
	}
	c.Set("UserData", userData)
	//写入userID
	c.Set("loginUserID", userData.Info.ID)
	c.Set("loginUserOrgID", userData.Info.OrgID)
	c.Set("loginUserPermissions", getUserPermissions(userData.Info.ID))
	//反馈
	userID = userData.Info.ID
	return
}

// 获取用户ID
func getUserID(c *gin.Context) (userID int64) {
	return c.MustGet("loginUserID").(int64)
}

// TryGetUserID 尝试获取用户ID
func TryGetUserID(c *gin.Context) (userID int64, b bool) {
	////尝试通过userID获取用户ID
	//var userIDStr interface{}
	//userIDStr, b = c.Get("loginUserID")
	//if !b {
	//	return
	//}
	//userID, _ = CoreFilter.GetInt64ByInterface(userIDStr)
	//if userID < 1 {
	//	b = false
	//	return
	//}
	//通过会话获取用户ID
	var err error
	userID, err = getUserDataByToken(c)
	if err != nil {
		return
	}
	//反馈
	return
}

// getUserRoleDataByC 获取当前用户的角色数据包
func getUserRoleDataByC(c *gin.Context, roleTypeMark string) (roleData UserRole.FieldsRole, b bool) {
	//获取用户数据
	userID := getUserID(c)
	//找到角色类型
	roleTypeData := UserRole.GetTypeMarkNoErr(roleTypeMark)
	if roleTypeData.ID < 1 {
		return
	}
	//根据用户数据，找到角色数据
	var err error
	roleData, err = UserRole.GetRoleUserID(&UserRole.ArgsGetRoleUserID{
		RoleType: roleTypeData.ID,
		UserID:   userID,
	})
	if err == nil && roleData.ID > 0 {
		b = true
		return
	}
	return
}

// UpdateUserLogin 用户登陆操作处理
func UpdateUserLogin(c *RouterURLHeaderC, userInfo *UserCore.FieldsUserType) {
	//TODO: v2过度处理，后续改进所有路由后删除
	c.Context.Set("UserData", UserCore.DataUserDataType{
		Info: UserCore.FieldsUserType{
			ID:          userInfo.ID,
			CreateAt:    userInfo.CreateAt,
			UpdateAt:    userInfo.UpdateAt,
			DeleteAt:    userInfo.DeleteAt,
			Status:      userInfo.Status,
			OrgID:       userInfo.OrgID,
			Name:        userInfo.Name,
			Password:    userInfo.Password,
			NationCode:  userInfo.NationCode,
			Phone:       userInfo.Phone,
			PhoneVerify: userInfo.PhoneVerify,
			Email:       userInfo.Email,
			EmailVerify: userInfo.EmailVerify,
			Username:    userInfo.Username,
			Avatar:      userInfo.Avatar,
			Parents:     userInfo.Parents,
			Groups:      userInfo.Groups,
			Infos:       userInfo.Infos,
			Logins:      userInfo.Logins,
			SortID:      userInfo.SortID,
			Tags:        userInfo.Tags,
		},
		Groups:      nil,
		Permissions: nil,
	})
	//保存用户基本信息
	c.Context.Set("loginUserID", userInfo.ID)
	c.Context.Set("loginUserOrgID", userInfo.OrgID)
	c.Context.Set("loginUserPermissions", getUserPermissions(userInfo.ID))
}

// 获取用户依存的组织ID
func getUserOrgID(c *gin.Context) (orgID int64) {
	return c.MustGet("loginUserOrgID").(int64)
}

// getUserPermissions 获取用户登录常用的权限组
func getUserPermissions(userID int64) []string {
	//获取用户数据
	userInfo, _ := UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    userID,
		OrgID: -1,
	})
	if userInfo.ID < 1 {
		return []string{}
	}
	//反馈权限
	return getUserPermissionsByUserData(&userInfo)
}

func getUserPermissionsByUserData(userInfo *UserCore.FieldsUserType) []string {
	var permissions []string
	var groupIDs []int64
	for _, v := range userInfo.Groups {
		if v.ExpireAt.Unix() < 1000000 || v.ExpireAt.Unix() >= CoreFilter.GetNowTime().Unix() {
			groupIDs = append(groupIDs, v.GroupID)
		}
	}
	if len(groupIDs) > 0 {
		permissions = UserCore.GetGroupPermissionList(groupIDs)
	}
	return permissions
}
