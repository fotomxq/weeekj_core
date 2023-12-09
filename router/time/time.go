package RouterTime

import (
	"fmt"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
)

// GetBetweenTime 获取和分解统计类时间结构，并解决路由问题
func GetBetweenTime(c *gin.Context, args CoreSQLTime.DataCoreTime) (newTime CoreSQLTime.FieldsCoreTime, b bool) {
	var err error
	newTime, err = CoreSQLTime.GetBetweenByISO(args)
	if err != nil {
		RouterReport.ErrorLog(c, fmt.Sprint("get between time by url: ", c.Request.URL, ", err: "), err, "time", "时间格式错误")
		return
	}
	b = true
	return
}
