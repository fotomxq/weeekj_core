package RouterParams

import (
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
)

// 常见action类动作请求参数头
type ActionIDType struct {
	ID int64 `json:"id"`
}

type ActionIDAndUserType struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userID"`
}

// 获取参数头
func GetActionID(c *gin.Context) (ActionIDType, bool) {
	params := ActionIDType{}
	if err := c.BindJSON(&params); err != nil {
		RouterReport.BaseError(c, "params_lost", "参数丢失")
		return ActionIDType{}, false
	}
	//过滤参数
	if params.ID < 1 {
		RouterReport.BaseError(c, "params_error", "ID无效")
		return ActionIDType{}, false
	}
	//反馈
	return params, true
}

func GetActionIDAndUser(c *gin.Context) (ActionIDAndUserType, bool) {
	params := ActionIDAndUserType{}
	if err := c.BindJSON(&params); err != nil {
		RouterReport.BaseError(c, "params_lost", "参数丢失")
		return ActionIDAndUserType{}, false
	}
	//过滤参数
	if params.ID < 1 || params.UserID < 1 {
		RouterReport.BaseError(c, "params_error", "ID或用户ID无效")
		return ActionIDAndUserType{}, false
	}
	//反馈
	return params, true
}
