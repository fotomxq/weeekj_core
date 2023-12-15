package Router2Mid

import (
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
	"net/http"
)

//gin中间件处理器

// HeaderBase 顶部设定
func HeaderBase(c *gin.Context) {
	//设定options请求处理
	if c.Request.Method == "OPTIONS" {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "content-type,action,nonce,secret_id,signature_key,signature_method,timestamp")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.AbortWithStatus(http.StatusOK)
	}
	//设置origin
	switch Router2SystemConfig.HeaderOrigin {
	case "*":
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	case "":
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	default:
		c.Writer.Header().Set("Access-Control-Allow-Origin", Router2SystemConfig.HeaderOrigin)
	}
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
}

// headerOptions 头基本设定
// 空包，用于options路由
//func headerOptions(c *gin.Context) {
//	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE")
//	c.Writer.Header().Set("Access-Control-Allow-Headers", "content-type,action,nonce,secret_id,signature_key,signature_method,timestamp")
//}
