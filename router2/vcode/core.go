package Router2VCode

import (
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	"github.com/gin-gonic/gin"
)

// 识别和获取头部上下文带数据
func getContextData(c any) (*gin.Context, Router2Mid.DataGetContextData) {
	return Router2Mid.GetContextData(c)
}
