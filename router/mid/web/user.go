package RouterMidWeb

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BasePedometer "github.com/fotomxq/weeekj_core/v5/base/pedometer"
	BaseSafe "github.com/fotomxq/weeekj_core/v5/base/safe"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	RouterMidAPI "github.com/fotomxq/weeekj_core/v5/router/mid/api"
	"github.com/gin-gonic/gin"
)

// 用户的标准化头部
// 反馈数据bool 代表是否需要主动处理跳转处理；否则被拦截后将自动跳转，外部调用方法不应该处理
func getUserData(c *gin.Context) (bool, error) {
	//获取用户结构
	userData, err := RouterMidAPI.GetUserDataByToken(c)
	if err != nil {
		//RouterReport.WarnLog(c, "token cannot get user data, ", err, "token_error", "无效用户")
		//RouterReport.BaseError(c, "token_error", "无效用户")
		return true, err
	}
	//对user进行安全事件检查
	SafetyUserON, err := BaseConfig.GetDataBool("SafetyUserON")
	if err != nil {
		SafetyUserON = true
	}
	if SafetyUserON && BasePedometer.CheckData(CoreSQLFrom.FieldsFrom{System: "safe_user", ID: userData.Info.ID}) {
		BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
			System: "api.token_ban",
			Level:  1,
			IP:     c.ClientIP(),
			UserID: 0,
			OrgID:  0,
			Des:    fmt.Sprint("用户[", userData.Info.ID, "]被禁用,但尝试访问URL:", c.Request.URL),
		})
		c.Redirect(200, "/ban")
		return false, errors.New("user have ban")
	}
	//反馈成功
	return true, nil
}
