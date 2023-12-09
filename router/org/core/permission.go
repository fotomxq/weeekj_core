package RouterOrgCore

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	RouterMidAPI "gitee.com/weeekj/weeekj_core/v5/router/mid/api"
	RouterMidOrg "gitee.com/weeekj/weeekj_core/v5/router/mid/org"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
)

// CheckPermissionByUser 检查用户是否具备对应组织权限，同时附带权限拦截
func CheckPermissionByUser(c *gin.Context, permissions []string) bool {
	//检查个人的组织操作权限
	if !RouterMidAPI.CheckUserPermission(c, "org") {
		//方法内包含了路由设置，此处反馈即可
		return false
	}
	//获取用户绑定关系，检查在组织内的附加权限是否具备
	//获取用户数据
	userData := RouterMidAPI.GetUserDataByC(c)
	//获取绑定关系
	bindData := RouterMidOrg.GetOrgBindData(c)
	//获取权限列
	bindPermissions := RouterMidOrg.GetOrgBindPermissions(c)
	//检查权限
	isOK := false
	for _, v := range bindPermissions {
		for _, v2 := range permissions {
			if v == v2 {
				isOK = true
				break
			}
		}
		if isOK {
			break
		}
	}
	if !isOK {
		//反馈失败
		CoreLog.Warn("user(", userData.Info.ID, ") try visit page url: ", c.Request.RequestURI, ", but not have work bind permission, need bind permission by org id: ", bindData.OrgID, ", bind id: ", bindData.ID, ", need org permission: ", permissions)
		RouterReport.BaseError(c, "no-work-permission", "组织权限不足")
		return false
	}
	//反馈成功
	return true
}
