package RouterMidWeb

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"github.com/gin-gonic/gin"
)

// GetHtmlData 静态文件固定数据生成器
func GetHtmlData(c *gin.Context) gin.H {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("router base data, ", r)
		}
	}()
	data := gin.H{
		"orgID":     c.MustGet("orgID"),
		"themeData": c.MustGet("themeData"),
	}
	return data
}
