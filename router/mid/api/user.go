package RouterMidAPI

import (
	"errors"
	"fmt"
	BasePedometer "gitee.com/weeekj/weeekj_core/v5/base/pedometer"
	BaseSafe "gitee.com/weeekj/weeekj_core/v5/base/safe"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreRPCX "gitee.com/weeekj/weeekj_core/v5/core/rpcx"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	UserLogin2 "gitee.com/weeekj/weeekj_core/v5/user/login2"
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
func GetUserDataByToken(c *gin.Context) (UserCore.DataUserDataType, error) {
	//获取token数据
	tokenInfo := getTokenInfo(c)
	//通过token数据集合，找到用户
	if tokenInfo.UserID < 1 {
		return UserCore.DataUserDataType{}, errors.New("token from system is not user")
	}
	userData, err := UserCore.GetUserData(&UserCore.ArgsGetUserData{
		UserID: tokenInfo.UserID,
	})
	if err != nil {
		return UserCore.DataUserDataType{}, errors.New("cannot find user info, " + err.Error())
	}
	if userData.Info.ID < 1 {
		return UserCore.DataUserDataType{}, errors.New("user not exist")
	}
	//写入user数据
	c.Set("UserData", userData)
	//写入userID
	c.Set("loginUserID", userData.Info.ID)
	c.Set("loginUserOrgID", userData.Info.OrgID)
	c.Set("loginUserPermissions", UserLogin2.GetUserPermissions(userData.Info.ID))
	return userData, nil
}

// GetUserDataByC 获取用户数据
func GetUserDataByC(c *gin.Context) UserCore.DataUserDataType {
	return c.MustGet("UserData").(UserCore.DataUserDataType)
}

// GetUserIDByC 获取用户ID
func GetUserIDByC(c *gin.Context) (userID int64) {
	//尝试多个渠道获取用户ID
	userIDAny, b := c.Get("loginUserID")
	if b {
		userID = userIDAny.(int64)
	}
	if userID < 1 {
		//获取token数据
		tokenInfo := getTokenInfo(c)
		userID = tokenInfo.UserID
		if userID < 1 {
			return
		}
	}
	//反馈
	return
}

// TryGetUserDataByC 获取用户数据
func TryGetUserDataByC(c *gin.Context) (userData UserCore.DataUserDataType, b bool) {
	userDataI, b2 := c.Get("UserData")
	if !b2 {
		return
	}
	userData = userDataI.(UserCore.DataUserDataType)
	b = true
	return
}

// GetUserRoleDataByC 获取当前用户的角色数据包
func GetUserRoleDataByC(c *gin.Context, roleTypeMark string) (roleData UserRole.FieldsRole, b bool) {
	//获取用户数据
	userID := GetUserIDByC(c)
	//找到角色类型
	roleTypeData, err := UserRole.GetTypeMark(&UserRole.ArgsGetTypeMark{
		Mark: roleTypeMark,
	})
	if err != nil {
		RouterReport.BaseError(c, "no_permission", "没有该角色权限")
		return
	}
	//根据用户数据，找到角色数据
	roleData, err = UserRole.GetRoleUserID(&UserRole.ArgsGetRoleUserID{
		RoleType: roleTypeData.ID,
		UserID:   userID,
	})
	if err == nil && roleData.ID > 0 {
		b = true
		return
	}
	//反馈失败
	RouterReport.BaseError(c, "no_permission", "没有该角色权限")
	return
}

// TryGetUserIDDataByToken 后置获取用户ID方法
// 可以在登陆前API中使用，尝试找到token绑定用户，如果找不到则反馈失败
func TryGetUserIDDataByToken(c *gin.Context) (userID int64) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			//不记录错误，只是捕捉走
		}
	}()
	userData, err := GetUserDataByToken(c)
	if err != nil {
		return 0
	}
	if userData.Info.ID > 0 {
		return userData.Info.ID
	}
	return 0
}

// SetUserData 临时写入部分用户数据
// 用于刚登陆和特殊场景
func SetUserData(c *gin.Context, userInfo UserCore.FieldsUserType) {
	c.Set("UserData", UserCore.DataUserDataType{
		Info:        userInfo,
		Groups:      nil,
		Permissions: nil,
	})
	//写入userID
	c.Set("loginUserID", userInfo.ID)
	c.Set("loginUserOrgID", userInfo.OrgID)
	c.Set("loginUserPermissions", UserLogin2.GetUserPermissions(userInfo.ID))
}

// CheckUserPermission 检查用户权限模块
// 1、需提前设定好UserData上下文关系
// 2、用户可使用check查询用户是否具备对应权限
// 3、关系为一一对应，不能是多对一
func CheckUserPermission(c *gin.Context, permissionMark string) bool {
	//获取用户数据
	userData := GetUserDataByC(c)
	//检查权限是否存在？
	for _, v := range userData.Permissions {
		if v == permissionMark {
			return true
		}
	}
	//记录日志
	BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
		System: "user.user_permission",
		Level:  1,
		IP:     c.ClientIP(),
		UserID: userData.Info.ID,
		OrgID:  userData.Info.OrgID,
		Des:    fmt.Sprint("用户不具备权限[", permissionMark, "],但尝试访问API,URL:", c.Request.URL),
	})
	//安全事件
	if _, err := BasePedometer.NextData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "safe-user", ID: userData.Info.ID},
	}); err != nil {
		CoreLog.Error("cannot add user by safe, ", err)
	}
	//反馈
	RouterReport.BaseError(c, "no_permission", "")
	return false
}
