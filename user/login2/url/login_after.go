package UserLogin2URL

import (
	"encoding/json"
	"fmt"
	AnalysisUserVisit "github.com/fotomxq/weeekj_core/v5/analysis/user_visit"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseFileSys2 "github.com/fotomxq/weeekj_core/v5/base/filesys2"
	BasePedometer "github.com/fotomxq/weeekj_core/v5/base/pedometer"
	BaseToken2 "github.com/fotomxq/weeekj_core/v5/base/token2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreRPCX "github.com/fotomxq/weeekj_core/v5/core/rpcx"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

// LoginAfter 用户登陆成功后标准步骤
func LoginAfter(c *Router2Mid.RouterURLHeaderC, userInfo *UserCore.FieldsUserType, loginMark string, rememberMe bool, loginFrom string) {
	//用户是否被禁止登陆？
	if !checkUserNotBan(c, userInfo) {
		return
	}
	//重建token
	newTokenInfo, b := clearAndCreateToken(c, userInfo, fmt.Sprint(loginFrom, "-", loginMark), rememberMe)
	if !b {
		return
	}
	//获取用户和组织数据集合
	userData, orgBindList, err := getLoginData(c, userInfo.ID)
	//如果包获取失败，则反馈失败
	if err != nil {
		return
	}
	//反馈的数据结构
	type reportData struct {
		//会话
		TokenID       int64  `json:"tokenID"`
		TokenKey      string `json:"tokenKey"`
		TokenExpireAt string `json:"tokenExpireAt"`
		//用户信息
		UserID         int64    `json:"userID"`
		NiceName       string   `json:"niceName"`
		UserAvatar     string   `json:"userAvatar"`
		UserPermission []string `json:"userPermission"`
		//用户所属组织
		OrgID int64 `json:"orgID"`
		//是否存在email
		HaveEmail bool `json:"haveEmail"`
		//是否绑定手机号
		HavePhone bool `json:"havePhone"`
		//用户登录名
		UserName string `json:"userName"`
	}
	//组装用户基本结构体
	newData := reportData{
		TokenID:        newTokenInfo.ID,
		TokenKey:       newTokenInfo.Key,
		TokenExpireAt:  CoreFilter.GetTimeToDefaultTime(newTokenInfo.ExpireAt),
		UserID:         userData.ID,
		NiceName:       userData.Name,
		UserAvatar:     BaseFileSys2.GetPublicURLByClaimID(userData.Avatar),
		UserPermission: getUserPermissionsByUserInfo(&userData),
		OrgID:          userInfo.OrgID,
		HaveEmail:      UserCore.CheckUserHaveEmail(userInfo),
		HavePhone:      UserCore.CheckUserHavePhone(userInfo),
		UserName:       userData.Username,
	}
	//将newData解析为interface结构
	newDataByte, err := json.Marshal(newData)
	if err != nil {
		Router2Mid.ReportWarnLog(c, "get json data", err, "err_json")
		return
	}
	var newDataAny map[string]interface{}
	if err = json.Unmarshal(newDataByte, &newDataAny); err != nil {
		Router2Mid.ReportWarnLog(c, "get json data", err, "err_json")
		return
	}
	//填入用户是否可见加密信息
	if Router2SystemConfig.GlobConfig.User.LoginViewPhone {
		newDataAny["nationCode"] = userInfo.NationCode
		newDataAny["phone"] = userInfo.Phone
	}
	if Router2SystemConfig.GlobConfig.User.LoginViewEmail {
		newDataAny["email"] = userInfo.Email
	}
	//检查用户是否具备组织
	if len(orgBindList) > 0 {
		//反馈的数据结构
		type reportOrgData struct {
			ID       int64    `json:"id"`
			Key      string   `json:"key"`
			Name     string   `json:"name"`
			Avatar   string   `json:"avatar"`
			OpenFunc []string `json:"openFunc"`
			Manager  []string `json:"manager"`
			BindID   int64    `json:"bindID"`
		}
		var newDataOrg []reportOrgData
		for _, v := range orgBindList {
			newDataOrg = append(newDataOrg, reportOrgData{
				ID:       v.OrgID,
				Key:      v.OrgKey,
				Name:     v.OrgName,
				Avatar:   BaseFileSys2.GetPublicURLByClaimID(v.OrgCoverFileID),
				OpenFunc: v.OrgOpenFunc,
				Manager:  v.Manager,
				BindID:   v.BindID,
			})
		}
		//填充数据
		newDataAny["orgBindData"] = newDataOrg
	}
	//强制更新用户数据包
	if userInfo.OrgID > 0 {
		OrgUserMod.PushUpdateUserData(userInfo.OrgID, userInfo.ID)
	}
	//统计行为
	//根据注册渠道，处理统计
	switch loginFrom {
	//case "password":
	//	_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
	//		OrgID: userInfo.OrgID,
	//		Mark:  100,
	//		Count: 1,
	//	})
	//case "phone":
	//	_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
	//		OrgID: userInfo.OrgID,
	//		Mark:  101,
	//		Count: 1,
	//	})
	//case "email":
	//	_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
	//		OrgID: userInfo.OrgID,
	//		Mark:  104,
	//		Count: 1,
	//	})
	case "marge_default":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: userInfo.OrgID,
			Mark:  100,
			Count: 1,
		})
	case "weixin_wxx":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: userInfo.OrgID,
			Mark:  102,
			Count: 1,
		})
	case "weixin_wxx_phone":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: userInfo.OrgID,
			Mark:  102,
			Count: 1,
		})
	case "weixin_app":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: userInfo.OrgID,
			Mark:  102,
			Count: 1,
		})
	}
	_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
		OrgID: userInfo.OrgID,
		Mark:  1,
		Count: 1,
	})
	//反馈结构体
	Router2Mid.BaseData(c, newDataAny)
}

// 通用检查用户是否可用？
func checkUserNotBan(c *Router2Mid.RouterURLHeaderC, userInfo *UserCore.FieldsUserType) bool {
	if userInfo.Status == 1 {
		Router2Mid.ReportWarnLog(c, "login user is ban", nil, "err_user_audit")
		return false
	}
	if userInfo.Status != 2 {
		Router2Mid.ReportWarnLog(c, "login user is ban", nil, "err_user_ban")
		return false
	}
	SafetyUserON, err := BaseConfig.GetDataBool("SafetyUserON")
	if err != nil {
		SafetyUserON = true
	}
	if SafetyUserON && BasePedometer.CheckData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "safe-user", ID: userInfo.ID},
	}) {
		Router2Mid.ReportWarnLog(c, "login user is ban", nil, "err_user_ban")
		return false
	}
	return true
}

// 清理和重建token
func clearAndCreateToken(c *Router2Mid.RouterURLHeaderC, userInfo *UserCore.FieldsUserType, loginMark string, rememberMe bool) (data BaseToken2.FieldsToken, b bool) {
	var err error
	//获取token数据
	oldTokenID := Router2Mid.GetTokenID(c.Context)
	//建立新的token
	BaseToken2.DeleteToken(oldTokenID)
	tokenID, errCode, err := BaseToken2.Create(&BaseToken2.ArgsCreate{
		UserID:     userInfo.ID,
		OrgID:      userInfo.OrgID,
		OrgBindID:  0,
		DeviceID:   0,
		LoginFrom:  loginMark,
		IP:         c.Context.ClientIP(),
		Key:        "",
		IsRemember: rememberMe,
	})
	if err != nil {
		Router2Mid.ReportWarnLog(c, fmt.Sprint("delete token, old token id: ", oldTokenID, ", new token id: ", tokenID), err, errCode)
		return
	}
	data = BaseToken2.GetByID(tokenID)
	//反馈成功
	b = true
	return
}

// getLoginData 登陆后数据包
func getLoginData(c *Router2Mid.RouterURLHeaderC, userID int64) (userData UserCore.FieldsUserType, orgData []OrgCoreCore.DataGetBindByUserMarge, err error) {
	//获取用户基础数据
	userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    userID,
		OrgID: -1,
	})
	if err != nil {
		Router2Mid.ReportWarnLog(c, fmt.Sprint("not get user data, user id: ", userID), err, "err_user")
		return
	}
	//获取该用户的所有绑定关系
	orgData, err = OrgCoreCore.GetBindByUserMarge(&OrgCoreCore.ArgsGetBindByUser{
		UserID: userData.ID,
	})
	if err != nil {
		err = nil
	}
	return
}
