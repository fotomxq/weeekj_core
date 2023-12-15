package RouterOrgCore

import (
	"errors"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/gin-gonic/gin"
)

// 获取CoreFieldsFrom的组织绑定来源
func GetDataFromByBind(bindData *OrgCore.FieldsBind) CoreSQLFrom.FieldsFrom {
	return CoreSQLFrom.FieldsFrom{
		System: "org",
		ID:     bindData.ID,
		Mark:   "bind",
		Name:   bindData.Name,
	}
}

// 确认绑定关系和来源用户是否一致
func CheckBindAndUser(c *gin.Context, bindID int64) bool {
	//获取用户数据
	userData := c.MustGet("UserData").(UserCore.DataUserDataType)
	//获取绑定关系
	bindData, err := OrgCore.GetBind(&OrgCore.ArgsGetBind{
		ID:     bindID,
		OrgID:  0,
		UserID: userData.Info.ID,
	})
	if err != nil {
		RouterReport.ErrorLog(c, "check bind and user, bind not exist, ", err, "no-organizational-bind", "no organizer bind permission")
		return false
	}
	if bindData.ID != bindID {
		RouterReport.ErrorLog(c, "check bind and user, bind not exist, ", err, "no-organizational-bind", "no organizer bind permission")
		return false
	}
	//反馈成功
	return true
}

// 通过用户获取在某个组织的绑定关系
func GetBindByUser(c *gin.Context, orgID int64) (bindData OrgCore.FieldsBind, err error) {
	//获取用户数据
	userData := c.MustGet("UserData").(UserCore.DataUserDataType)
	//获取绑定关系
	var bindList []OrgCore.FieldsBind
	bindList, _, err = OrgCore.GetBindList(&OrgCore.ArgsGetBindList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  1,
			Sort: "_id",
			Desc: false,
		},
		OrgID:    orgID,
		UserID:   userData.Info.ID,
		GroupID:  0,
		Manager:  "",
		IsRemove: false,
		Search:   "",
	})
	if err != nil {
		return
	}
	if len(bindList) < 1 {
		err = errors.New("no bind")
		return
	}
	bindData = bindList[0]
	return
}
