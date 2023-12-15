package CoreLog

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/gin-gonic/gin"
)

//设置gin路由的日志部分
// 采用logrus实现，go-file-rotatelogs插件分割文件
// 参考: https://hacpai.com/article/1531819038419
// 修改了日志生成结构部分

// GinLogger 给gin的日志处理器
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := CoreFilter.GetNowTime()
		// 处理请求
		c.Next()
		// 结束时间
		end := CoreFilter.GetNowTime()
		// 执行时间
		latency := end.Sub(start)
		//获取基本数据
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		//输出信息
		ginLog.LogHandle.Infof("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}
