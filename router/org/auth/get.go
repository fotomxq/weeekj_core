package RouterOrgAuth

import (
	"fmt"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	RouterMidOrg "github.com/fotomxq/weeekj_core/v5/router/mid/org"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context, orgID int64, childOrgID int64, configMark string) (int64, int64, bool) {
	//获取组织
	orgData := RouterMidOrg.GetOrg(c)
	//如果相同，则跳出
	if orgID == orgData.ID {
		return orgID, childOrgID, true
	}
	//检查子商户
	nowOrgData, err := OrgCoreCore.GetOrg(&OrgCoreCore.ArgsGetOrg{
		ID: orgID,
	})
	if err != nil {
		RouterReport.WarnLog(c, fmt.Sprint("org manager get org auth data failed, now org parent id: ", orgData.ParentID, ", now org id: ", orgData.ID, ", need org id: ", orgID, ", need child org id: ", childOrgID, ", need auth org: ", orgData.ID), err, "no_auth", "无授权")
	} else {
		//上级商户查看子商户数据
		if nowOrgData.ParentID > 0 && orgData.ID == nowOrgData.ParentID {
			return orgID, childOrgID, true
		}
	}
	//不相同，则检查是否为子商户在访问数据
	if orgData.ParentID > 0 {
		if orgData.ParentID == orgID && orgData.ID == childOrgID {
			return orgData.ParentID, orgData.ID, true
		} else {
			RouterReport.WarnLog(c, fmt.Sprint("org manager get org auth data failed, now org parent id: ", orgData.ParentID, ", now org id: ", orgData.ID, ", need org id: ", orgID, ", need child org id: ", childOrgID, ", need auth org: ", orgData.ID), nil, "no_auth", "无授权")
		}
	} else {
		//和上级还是不同，则检查平级授权
		if CheckAuth(c, orgID, configMark) {
			return orgID, childOrgID, true
		}
	}
	//非法条件处理
	return -1, -1, false
}
