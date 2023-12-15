package RouterGinSet

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	//Router 路由基础
	Router *gin.Engine
)

// Init 初始化路由设定
// 全局通用设计
func Init() {
	//debug
	if Router2SystemConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	//初始化gin
	Router = gin.New()
	//设置路由基础
	Router.Use(CoreLog.GinLogger())
	Router.Use(gin.Recovery())
	Router.Use(Router2Mid.HeaderBase)
	//设置404页面
	Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "")
	})
	Router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "")
	})
}

// PageDefaultFavicon 设置/favicon.ico
func PageDefaultFavicon() {
	Router.GET("/favicon.ico", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
}

// PageDefaultRoot 设置根服务
func PageDefaultRoot() {
	Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
}

// RunServer 启动服务
func RunServer() bool {
	if err := Router.Run(Router2SystemConfig.RouterHost); err != nil {
		CoreLog.Error("cannot run server, " + err.Error())
	}
	return false
}
