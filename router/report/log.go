package RouterReport

import (
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	"github.com/gin-gonic/gin"
)

// 日志内部附加处理
func getLogMsg(c *gin.Context, msg string) string {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("log failed, ", r, ", url: ", c.Request.RequestURI, ", msg: ", msg)
		}
	}()
	var appendMsg string
	userDataIn, b := c.Get("UserData")
	if b {
		userData := userDataIn.(UserCore.DataUserDataType)
		appendMsg = fmt.Sprint("userID:", userData.Info.ID, ", userPhone:", userData.Info.Phone, ", ", appendMsg)
	}
	orgDataIn, b := c.Get("OrgData")
	if b {
		orgData := orgDataIn.(OrgCore.FieldsOrg)
		appendMsg = fmt.Sprint("orgID:", orgData.ID, ", ", appendMsg)
	}
	orgBindDataIn, b := c.Get("OrgBindData")
	if b {
		bindData := orgBindDataIn.(OrgCore.FieldsBind)
		appendMsg = fmt.Sprint("orgBindID:", bindData.ID, ", ", appendMsg)
	}
	if appendMsg != "" {
		msg = appendMsg + " -> " + msg
	}
	return msg
}
