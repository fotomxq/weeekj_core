package Router2Mid

import (
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	"github.com/gin-gonic/gin"
)

// headerSelectedOrg 行政组织专用顶层头
// 用于组织和用户级别通用设计
// 必须建立在RouterMiddleware.HeaderLoggedUser(c)基础之上，即路由需设计多层函数关系构建
func headerSelectedOrg(c *gin.Context) {
	//获取数据
	tokenInfo := GetTokenInfo(c)
	//检查权限
	if !checkPermissionUser(c, tokenInfo.UserID, tokenInfo.OrgID, []string{"org"}) {
		return
	}
	//存储数据
	c.Set("selectOrgID", tokenInfo.OrgID)
	c.Set("selectOrgBindID", tokenInfo.OrgBindID)
	//获取权限列
	c.Set("selectOrgBindPermissions", OrgCore.GetPermissionByBindID(tokenInfo.OrgBindID))
	//继续执行
	c.Next()
}

func getOrgID(c *gin.Context) (orgID int64) {
	dataAny, b := c.Get("selectOrgID")
	if !b {
		data := c.MustGet("OrgData").(OrgCore.FieldsOrg)
		orgID = data.ID
	} else {
		orgID = dataAny.(int64)
	}
	return
}

func getOrgBindID(c *gin.Context) (bindID int64) {
	dataAny, b := c.Get("selectOrgBindID")
	if !b {
		data := c.MustGet("OrgBindData").(OrgCore.FieldsBind)
		bindID = data.ID
	} else {
		bindID = dataAny.(int64)
	}
	return
}

func getOrgBindPermissions(c *gin.Context) (data []string) {
	dataAny, b := c.Get("selectOrgBindPermissions")
	if !b {
		data = c.MustGet("OrgBindPermissions").([]string)
	} else {
		data = dataAny.([]string)
	}
	return
}
