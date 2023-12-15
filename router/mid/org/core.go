package RouterMidOrg

import (
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	RouterMidAPI "github.com/fotomxq/weeekj_core/v5/router/mid/api"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
)

// HeaderSelectedOrg 行政组织专用顶层头
// 用于组织和用户级别通用设计
// 必须建立在RouterMiddleware.HeaderLoggedUser(c)基础之上，即路由需设计多层函数关系构建
func HeaderSelectedOrg(c *gin.Context) {
	//权限检查
	if !RouterMidAPI.CheckUserPermission(c, "org") {
		return
	}
	//获取用户结构体
	//userData := c.MustGet("UserData").(UserCore.DataUserDataType)
	userID := c.MustGet("loginUserID").(int64)
	//获取组织和绑定关系数据包
	orgData, bindData, permissions, err := OrgCore.GetSelectAndData(&OrgCore.ArgsGetSelectAndData{
		UserID: userID,
	})
	if err != nil {
		RouterReport.WarnLog(c, "get select org by user id, ", err, "org_not_select", "尚未选择组织")
		return
	}
	//存储数据
	c.Set("OrgData", orgData)
	c.Set("OrgBindData", bindData)
	c.Set("OrgBindPermissions", permissions)
	//存储数据
	c.Set("selectOrgID", orgData.ID)
	c.Set("selectOrgBindID", bindData.ID)
	//获取权限列
	c.Set("selectOrgBindPermissions", OrgCore.GetPermissionByBindID(bindData.ID))
	//继续执行
	c.Next()
}

// GetOrg 通过头数据获取组织相关信息
func GetOrg(c *gin.Context) (data OrgCore.FieldsOrg) {
	data = c.MustGet("OrgData").(OrgCore.FieldsOrg)
	return
}

func GetOrgBindData(c *gin.Context) (data OrgCore.FieldsBind) {
	data = c.MustGet("OrgBindData").(OrgCore.FieldsBind)
	return
}

func TryGetOrgBindData(c *gin.Context) (data OrgCore.FieldsBind, b bool) {
	dataI, b2 := c.Get("OrgBindData")
	if !b2 {
		return
	}
	data = dataI.(OrgCore.FieldsBind)
	b = true
	return
}

func GetOrgBindPermissions(c *gin.Context) (data []string) {
	data = c.MustGet("OrgBindPermissions").([]string)
	return
}
