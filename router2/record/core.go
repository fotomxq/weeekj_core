package Router2Record

import (
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	UserRecord2Mod "gitee.com/weeekj/weeekj_core/v5/user/record2/mod"
	"github.com/gin-gonic/gin"
	"strings"
)

// 识别和获取头部上下文带数据
func getContextData(c any) (*gin.Context, Router2Mid.DataGetContextData) {
	return Router2Mid.GetContextData(c)
}

// AddRecord 结合C写入数据包
func AddRecord(context any, system string, modID int64, des ...interface{}) {
	c, cData := getContextData(context)
	//生成mark
	mark := strings.Replace(c.Request.URL.String(), "/", "_", -1)
	mark = strings.Replace(mark, ":", "", -1)
	//写入日志
	UserRecord2Mod.AppendData(cData.OrgID, cData.OrgBindID, cData.UserID, system, modID, mark, des)
}
