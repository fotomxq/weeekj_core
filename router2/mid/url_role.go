package Router2Mid

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RouterURLRoleC 用户角色头部
type RouterURLRoleC struct {
	//上下文
	Context *gin.Context
	//日志头部
	LogAppend string
	//用户ID
	UserID int64
	//RoleID
	RoleID int64
	//角色类型
	RoleType string
}

func (t *RouterURLRole) GET(urlPath string, handle func(*RouterURLRoleC)) {
	//映射URL
	t.BaseData.Routers.GET(urlPath, func(c *gin.Context) {
		getURLRole(t, c, handle)
	})
}

func (t *RouterURLRole) POST(urlPath string, handle func(*RouterURLRoleC)) {
	//映射URL
	t.BaseData.Routers.POST(urlPath, func(c *gin.Context) {
		getURLRole(t, c, handle)
	})
}

func (t *RouterURLRole) PUT(urlPath string, handle func(*RouterURLRoleC)) {
	//映射URL
	t.BaseData.Routers.PUT(urlPath, func(c *gin.Context) {
		getURLRole(t, c, handle)
	})
}

func (t *RouterURLRole) DELETE(urlPath string, handle func(*RouterURLRoleC)) {
	//映射URL
	t.BaseData.Routers.DELETE(urlPath, func(c *gin.Context) {
		getURLRole(t, c, handle)
	})
}

// getURLRole 方法集合处理封装
func getURLRole(t *RouterURLRole, c *gin.Context, handle func(*RouterURLRoleC)) {
	//获取组织基本信息
	userID := getUserID(c)
	//获取用户角色
	roleData, b := getUserRoleDataByC(c, t.RoleType)
	if !b {
		c.Redirect(http.StatusMovedPermanently, userRoleBanURL)
		c.Abort()
		return
	}
	//映射接口
	handle(&RouterURLRoleC{
		Context:   c,
		LogAppend: "",
		UserID:    userID,
		RoleID:    roleData.ID,
		RoleType:  t.RoleType,
	})
}
