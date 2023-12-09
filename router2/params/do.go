package Router2Params

import (
	Router2Report "gitee.com/weeekj/weeekj_core/v5/router2/mid"
)

// ActionIDType 常见action类动作请求参数头
type ActionIDType struct {
	ID int64 `json:"id"`
}

type ActionIDAndUserType struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userID"`
}

// GetActionID 获取参数头
func GetActionID(context any) (ActionIDType, bool) {
	c := getContext(context)
	params := ActionIDType{}
	if err := c.BindJSON(&params); err != nil {
		Router2Report.ReportBaseError(c, "report_params_lost")
		return ActionIDType{}, false
	}
	//过滤参数
	if params.ID < 1 {
		Router2Report.ReportBaseError(c, "err_params_id")
		return ActionIDType{}, false
	}
	//反馈
	return params, true
}

func GetActionIDAndUser(context any) (ActionIDAndUserType, bool) {
	c := getContext(context)
	params := ActionIDAndUserType{}
	if err := c.BindJSON(&params); err != nil {
		Router2Report.ReportBaseError(c, "report_params_lost")
		return ActionIDAndUserType{}, false
	}
	//过滤参数
	if params.ID < 1 || params.UserID < 1 {
		Router2Report.ReportBaseError(c, "err_params_user")
		return ActionIDAndUserType{}, false
	}
	//反馈
	return params, true
}
