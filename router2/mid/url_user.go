package Router2Mid

import "github.com/gin-gonic/gin"

// RouterURLUserC 普通级别头部
type RouterURLUserC struct {
	//上下文
	Context *gin.Context
	//日志头部
	LogAppend string
	//用户ID
	UserID int64
}

func (t *RouterURLUser) GET(urlPath string, handle func(*RouterURLUserC)) {
	//映射URL
	t.BaseData.Routers.GET(urlPath, func(c *gin.Context) {
		getURLUser(c, handle)
	})
}

func (t *RouterURLUser) POST(urlPath string, handle func(*RouterURLUserC)) {
	//映射URL
	t.BaseData.Routers.POST(urlPath, func(c *gin.Context) {
		getURLUser(c, handle)
	})
}

func (t *RouterURLUser) PUT(urlPath string, handle func(*RouterURLUserC)) {
	//映射URL
	t.BaseData.Routers.PUT(urlPath, func(c *gin.Context) {
		getURLUser(c, handle)
	})
}

func (t *RouterURLUser) DELETE(urlPath string, handle func(*RouterURLUserC)) {
	//映射URL
	t.BaseData.Routers.DELETE(urlPath, func(c *gin.Context) {
		getURLUser(c, handle)
	})
}

// getURLUser 方法集合处理封装
func getURLUser(c *gin.Context, handle func(*RouterURLUserC)) {
	//获取组织基本信息
	userID := getUserID(c)
	//映射接口
	handle(&RouterURLUserC{
		Context:   c,
		LogAppend: "",
		UserID:    userID,
	})
}
