package Router2Mid

import "github.com/gin-gonic/gin"

// RouterURLHeaderC 普通级别头部
type RouterURLHeaderC struct {
	//上下文
	Context *gin.Context
	//日志头部
	LogAppend string
	//会话ID
	TokenID int64
}

func (t *RouterURLHeader) GET(urlPath string, handle func(*RouterURLHeaderC)) {
	//映射URL
	t.BaseData.Routers.GET(urlPath, func(c *gin.Context) {
		getURLHeader(c, handle)
	})
}

func (t *RouterURLHeader) POST(urlPath string, handle func(*RouterURLHeaderC)) {
	//映射URL
	t.BaseData.Routers.POST(urlPath, func(c *gin.Context) {
		getURLHeader(c, handle)
	})
}

func (t *RouterURLHeader) PUT(urlPath string, handle func(*RouterURLHeaderC)) {
	//映射URL
	t.BaseData.Routers.PUT(urlPath, func(c *gin.Context) {
		getURLHeader(c, handle)
	})
}

func (t *RouterURLHeader) DELETE(urlPath string, handle func(*RouterURLHeaderC)) {
	//映射URL
	t.BaseData.Routers.DELETE(urlPath, func(c *gin.Context) {
		getURLHeader(c, handle)
	})
}

// getURLHeader 方法集合处理封装
func getURLHeader(c *gin.Context, handle func(*RouterURLHeaderC)) {
	//获取当前会话
	tokenID := GetTokenID(c)
	//映射接口
	handle(&RouterURLHeaderC{
		Context:   c,
		LogAppend: "",
		TokenID:   tokenID,
	})
}
