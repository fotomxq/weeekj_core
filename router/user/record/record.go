package RouterUserRecord

import (
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	RouterMidAPI "gitee.com/weeekj/weeekj_core/v5/router/mid/api"
	UserRecord "gitee.com/weeekj/weeekj_core/v5/user/record"
	"github.com/gin-gonic/gin"
	"strings"
)

// CreateByC 结合C写入数据包
func CreateByC(c *gin.Context, content ...interface{}) {
	userData := RouterMidAPI.GetUserDataByC(c)
	mark := strings.Replace(c.Request.URL.String(), "/", "_", -1)
	mark = strings.Replace(mark, ":", "", -1)
	//重组消息
	var contentMsg string
	for _, v := range content {
		contentMsg = contentMsg + fmt.Sprint(v)
	}
	if err := UserRecord.Create(&UserRecord.ArgsCreate{
		OrgID:       userData.Info.OrgID,
		UserID:      userData.Info.ID,
		UserName:    userData.Info.Name,
		ContentMark: mark,
		Content:     contentMsg,
	}); err != nil {
		CoreLog.Error("create user record, ", err)
	}
}
