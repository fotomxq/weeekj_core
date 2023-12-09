package Router2Mid

import (
	"fmt"
	BasePedometer "gitee.com/weeekj/weeekj_core/v5/base/pedometer"
	BaseSafe "gitee.com/weeekj/weeekj_core/v5/base/safe"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreRPCX "gitee.com/weeekj/weeekj_core/v5/core/rpcx"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
)

// CheckPermission 检查权限
// 自动识别组织或用户
func CheckPermission(context any, permissionMarks []string) bool {
	if len(permissionMarks) < 1 {
		return true
	}
	var userID, orgID, orgBindID int64
	var c *gin.Context
	mode := ""
	userC, ok := context.(*RouterURLUserC)
	if ok {
		c = userC.Context
		mode = "user"
		userID = userC.UserID
		orgID = getUserOrgID(c)
	} else {
		orgC, ok := context.(*RouterURLOrgC)
		if ok {
			c = orgC.Context
			mode = "org"
			userID = orgC.UserID
			orgID = orgC.OrgID
			orgBindID = orgC.OrgBindID
		} else {
			roleC, ok := context.(*RouterURLRoleC)
			if ok {
				c = roleC.Context
				mode = "role"
				userID = roleC.UserID
			} else {
				//抛出异常，在非登录模式下检查权限
				panic(fmt.Sprint("check permission failed"))
				return false
			}
		}
	}
	switch mode {
	case "user":
		return checkPermissionUser(c, userID, orgID, permissionMarks)
	case "org":
		return checkPermissionOrg(c, userID, orgID, orgBindID, permissionMarks)
	case "role":
		return checkPermissionUser(c, userID, orgID, permissionMarks)
	default:
		//抛出异常，在非登录模式下检查权限
		panic(fmt.Sprint("check permission failed"))
		return false
	}
}

func checkPermissionUser(c *gin.Context, userID, orgID int64, permissionMarks []string) bool {
	//获取用户数据
	permissions := getUserPermissions(userID)
	//检查权限是否存在？
	haveAll := true
	for _, v := range permissionMarks {
		isFind := false
		for _, v2 := range permissions {
			if v == v2 {
				isFind = true
				break
			}
		}
		if !isFind {
			//记录日志
			BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
				System: "user.user_permission",
				Level:  1,
				IP:     c.ClientIP(),
				UserID: userID,
				OrgID:  orgID,
				Des:    fmt.Sprint("用户不具备权限[", v, "],但尝试访问API,URL:", c.Request.URL),
			})
			haveAll = false
			//记录日志
			if Router2SystemConfig.GlobConfig.Router.NeedTokenLog {
				CoreLog.Warn("router mid check user permissions failed, user id: ", userID, ", user have: ", permissions, ", need: ", permissionMarks)
			}
			//反馈
			break
		}
	}
	if haveAll {
		return true
	}
	//安全事件
	if _, err := BasePedometer.NextData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "safe-user", ID: userID},
	}); err != nil {
		reportGin(c, true, 0, err, "add user by safe", false, "err_permission", 0, nil)
		return false
	}
	//反馈
	reportGin(c, false, 0, nil, "add user by safe", false, "err_permission", 0, nil)
	//反馈
	return false
}

func checkPermissionOrg(c *gin.Context, userID, orgID, orgBindID int64, permissionMarks []string) bool {
	//检查个人的组织操作权限
	if !checkPermissionUser(c, userID, orgID, []string{"org"}) {
		//记录日志
		if Router2SystemConfig.GlobConfig.Router.NeedTokenLog {
			CoreLog.Warn("router mid check org bind permissions failed, user id: ", userID, ", no org base permission.")
		}
		//反馈
		//方法内包含了路由设置，此处反馈即可
		return false
	}
	//检查组织成员权限
	haveOrgPermission := OrgCore.CheckPermissionByBindID(orgBindID, permissionMarks)
	if !haveOrgPermission {
		//反馈失败
		reportGin(c, true, 0, nil, fmt.Sprint("need org permissions:", permissionMarks), false, "err_permission_org", 0, nil)
		//记录日志
		if Router2SystemConfig.GlobConfig.Router.NeedTokenLog {
			CoreLog.Warn("router mid check org bind permissions failed, user id: ", userID, ", org bind id: ", orgBindID, ", org id: ", orgID, ", need: ", permissionMarks)
		}
		//反馈
		return false
	}
	//反馈成功
	return true
}
