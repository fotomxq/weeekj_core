package Router2Mid

import "github.com/gin-gonic/gin"

// RouterURLPublicC 普通级别头部
type RouterURLPublicC struct {
	//上下文
	Context *gin.Context
	//日志头部
	LogAppend string
}

func (t *RouterURLPublic) GET(urlPath string, handle func(*RouterURLPublicC)) {
	//映射URL
	t.BaseData.Routers.GET(urlPath, func(c *gin.Context) {
		getURLPublic(c, handle)
	})
}

func (t *RouterURLPublic) POST(urlPath string, handle func(*RouterURLPublicC)) {
	//映射URL
	t.BaseData.Routers.POST(urlPath, func(c *gin.Context) {
		getURLPublic(c, handle)
	})
}

func (t *RouterURLPublic) PUT(urlPath string, handle func(*RouterURLPublicC)) {
	//映射URL
	t.BaseData.Routers.PUT(urlPath, func(c *gin.Context) {
		getURLPublic(c, handle)
	})
}

func (t *RouterURLPublic) DELETE(urlPath string, handle func(*RouterURLPublicC)) {
	//映射URL
	t.BaseData.Routers.DELETE(urlPath, func(c *gin.Context) {
		getURLPublic(c, handle)
	})
}

// getURLPublic 方法集合处理封装
func getURLPublic(c *gin.Context, handle func(*RouterURLPublicC)) {
	handle(&RouterURLPublicC{
		Context:   c,
		LogAppend: "",
	})
}
