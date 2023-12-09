package RouterMidCore

import (
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
)

//gin中间件处理器

// SetOriginConfig 设置origin头部
func SetOriginConfig(c *gin.Context) {
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
