package Router2Mid

import (
	"github.com/gin-gonic/gin"
)

// RouterURLOrgC 普通级别头部
type RouterURLOrgC struct {
	//上下文
	Context *gin.Context
	//日志头部
	LogAppend string
	//用户ID
	UserID int64
	//组织ID
	OrgID int64
	//组织成员ID
	OrgBindID int64
}

func (t *RouterURLOrg) GET(urlPath string, handle func(*RouterURLOrgC)) {
	//映射URL
	t.BaseData.Routers.GET(urlPath, func(c *gin.Context) {
		getURLOrg(c, handle)
	})
}

func (t *RouterURLOrg) POST(urlPath string, handle func(*RouterURLOrgC)) {
	//映射URL
	t.BaseData.Routers.POST(urlPath, func(c *gin.Context) {
		getURLOrg(c, handle)
	})
}

func (t *RouterURLOrg) PUT(urlPath string, handle func(*RouterURLOrgC)) {
	//映射URL
	t.BaseData.Routers.PUT(urlPath, func(c *gin.Context) {
		getURLOrg(c, handle)
	})
}

func (t *RouterURLOrg) DELETE(urlPath string, handle func(*RouterURLOrgC)) {
	//映射URL
	t.BaseData.Routers.DELETE(urlPath, func(c *gin.Context) {
		getURLOrg(c, handle)
	})
}

// getURLPublic 方法集合处理封装
func getURLOrg(c *gin.Context, handle func(*RouterURLOrgC)) {
	//获取组织基本信息
	userID := getUserID(c)
	orgID := getOrgID(c)
	bindID := getOrgBindID(c)
	//映射接口
	handle(&RouterURLOrgC{
		Context:   c,
		LogAppend: "",
		UserID:    userID,
		OrgID:     orgID,
		OrgBindID: bindID,
	})
}
