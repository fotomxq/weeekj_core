package UserLogin

import (
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseFileSys2 "github.com/fotomxq/weeekj_core/v5/base/filesys2"
	BasePedometer "github.com/fotomxq/weeekj_core/v5/base/pedometer"
	BaseQiniu "github.com/fotomxq/weeekj_core/v5/base/qiniu"
	BaseToken2 "github.com/fotomxq/weeekj_core/v5/base/token2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreRPCX "github.com/fotomxq/weeekj_core/v5/core/rpcx"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/gin-gonic/gin"
	"time"
)

// LoginAfter 用户登陆成功后标准步骤
func LoginAfter(c *gin.Context, userInfo *UserCore.FieldsUserType, loginMark string, rememberMe bool, loginFrom string) {
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
	userData, orgBindList, err := LoginData(c, userInfo.ID)
	//检查用户是否具备组织
	if err == nil {
		if len(orgBindList) > 0 {
			//反馈的数据结构
			type ReportDataToken struct {
				Token    int64     `json:"token"`
				Key      string    `json:"key"`
				ExpireAt time.Time `json:"expireAt"`
			}
			type ReportData struct {
				//会话
				Token ReportDataToken `json:"token"`
				//用户脱敏数据
				UserData UserCore.DataUserDataType `json:"userData"`
				//绑定关系和组织关系
				OrgBindData []OrgCoreCore.DataGetBindByUserMarge `json:"orgBindData"`
				//文件集
				// id => url
				FileList map[int64]string `json:"fileList"`
			}
			var waitFiles []int64
			for _, v := range orgBindList {
				if v.OrgCoverFileID > 0 {
					waitFiles = append(waitFiles, v.OrgCoverFileID)
				}
			}
			if userData.Info.Avatar > 0 {
				waitFiles = append(waitFiles, userData.Info.Avatar)
			}
			fileList := map[int64]string{}
			if len(waitFiles) > 0 {
				fileList, err = BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
					ClaimIDList: waitFiles,
					UserID:      0,
					OrgID:       0,
					IsPublic:    true,
				})
				if err != nil {
					//找不到文件，跳过
					err = nil
				}
			}
			RouterReport.BaseData(c, ReportData{
				Token: ReportDataToken{
					Token:    newTokenInfo.ID,
					Key:      newTokenInfo.Key,
					ExpireAt: newTokenInfo.ExpireAt,
				},
				UserData:    userData,
				OrgBindData: orgBindList,
				FileList:    fileList,
			})
		} else {
			//反馈的数据结构
			type ReportDataToken struct {
				Token    int64     `json:"token"`
				Key      string    `json:"key"`
				ExpireAt time.Time `json:"expireAt"`
			}
			type ReportData struct {
				//会话
				Token ReportDataToken `json:"token"`
				//用户脱敏数据
				UserData UserCore.DataUserDataType `json:"userData"`
				//文件集
				// id => url
				FileList map[int64]string `json:"fileList"`
			}
			var waitFiles []int64
			if userData.Info.Avatar > 0 {
				waitFiles = append(waitFiles, userData.Info.Avatar)
			}
			fileList := map[int64]string{}
			if len(waitFiles) > 0 {
				fileList, err = BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
					ClaimIDList: waitFiles,
					UserID:      0,
					OrgID:       0,
					IsPublic:    true,
				})
				if err != nil {
					//找不到文件，跳过
					err = nil
				}
			}
			RouterReport.BaseData(c, ReportData{
				Token: ReportDataToken{
					Token:    newTokenInfo.ID,
					Key:      newTokenInfo.Key,
					ExpireAt: newTokenInfo.ExpireAt,
				},
				UserData: userData,
				FileList: fileList,
			})
		}
	} else {
		//用户数据包获取失败，需反馈失败
		return
	}
	//强制更新用户数据包
	if userInfo.OrgID > 0 {
		OrgUserMod.PushUpdateUserData(userInfo.OrgID, userInfo.ID)
	}
}

// LoginAfter2 用户登陆成功后标准步骤
// TODO: 未来需要对内部函数做一些整理和优化
func LoginAfter2(c *gin.Context, userInfo *UserCore.FieldsUserType, loginMark string, rememberMe bool, loginFrom string) {
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
	userData, orgBindList, err := LoginData(c, userInfo.ID)
	//检查用户是否具备组织
	if err == nil {
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
			type reportData struct {
				//会话
				TokenID       int64  `json:"tokenID"`
				TokenKey      string `json:"tokenKey"`
				TokenExpireAt string `json:"tokenExpireAt"`
				//用户信息
				UserID         int64    `json:"userID"`
				UserName       string   `json:"userName"`
				UserAvatar     string   `json:"userAvatar"`
				UserPermission []string `json:"userPermission"`
				//绑定关系和组织关系
				OrgBindData []reportOrgData `json:"orgBindData"`
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
			newData := reportData{
				TokenID:        newTokenInfo.ID,
				TokenKey:       newTokenInfo.Key,
				TokenExpireAt:  CoreFilter.GetTimeToDefaultTime(newTokenInfo.ExpireAt),
				UserID:         userData.Info.ID,
				UserName:       userData.Info.Name,
				UserAvatar:     BaseFileSys2.GetPublicURLByClaimID(userData.Info.Avatar),
				UserPermission: userData.Permissions,
				OrgBindData:    newDataOrg,
			}
			RouterReport.BaseData(c, newData)
		} else {
			//反馈的数据结构
			type reportData struct {
				//会话
				TokenID       int64  `json:"tokenID"`
				TokenKey      string `json:"tokenKey"`
				TokenExpireAt string `json:"tokenExpireAt"`
				//用户信息
				UserID         int64    `json:"userID"`
				UserName       string   `json:"userName"`
				UserAvatar     string   `json:"userAvatar"`
				UserPermission []string `json:"userPermission"`
			}
			newData := reportData{
				TokenID:        newTokenInfo.ID,
				TokenKey:       newTokenInfo.Key,
				TokenExpireAt:  CoreFilter.GetTimeToDefaultTime(newTokenInfo.ExpireAt),
				UserID:         userData.Info.ID,
				UserName:       userData.Info.Name,
				UserAvatar:     BaseFileSys2.GetPublicURLByClaimID(userData.Info.Avatar),
				UserPermission: userData.Permissions,
			}
			RouterReport.BaseData(c, newData)
		}
	} else {
		//用户数据包获取失败，需反馈失败
		return
	}
	//强制更新用户数据包
	if userInfo.OrgID > 0 {
		OrgUserMod.PushUpdateUserData(userInfo.OrgID, userInfo.ID)
	}
}

// LoginAfterSave 用户登陆成功后标准步骤
// 该设计会将数据，存储到列队等待处理
func LoginAfterSave(c *gin.Context, userInfo *UserCore.FieldsUserType, loginMark string, rememberMe bool, loginFrom string) {
	//用户是否被禁止登陆？
	if !checkUserNotBan(c, userInfo) {
		return
	}
	//重建token
	newTokenInfo, b := clearAndCreateToken(c, userInfo, loginMark, rememberMe)
	if !b {
		return
	}
	//获取用户和组织数据集合
	userData, orgBindList, err := LoginData(c, userInfo.ID)
	//检查用户是否具备组织
	if err == nil {
		if len(orgBindList) > 0 {
			var waitFiles []int64
			for _, v := range orgBindList {
				if v.OrgCoverFileID > 0 {
					waitFiles = append(waitFiles, v.OrgCoverFileID)
				}
			}
			if userData.Info.Avatar > 0 {
				waitFiles = append(waitFiles, userData.Info.Avatar)
			}
			fileList := map[int64]string{}
			if len(waitFiles) > 0 {
				fileList, err = BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
					ClaimIDList: waitFiles,
					UserID:      0,
					OrgID:       0,
					IsPublic:    true,
				})
				if err != nil {
					//找不到文件，跳过
					err = nil
				}
			}
			loginData := FieldsSaveReportData{
				Token: FieldsSaveReportDataToken{
					Token: newTokenInfo.ID,
					Key:   newTokenInfo.Key,
				},
				UserData:    userData,
				OrgBindData: orgBindList,
				FileList:    fileList,
			}
			var newKey string
			newKey, err = appendSave(loginData)
			RouterReport.Data(c, "login user save, json error, ", "登陆异常", err, newKey)
		} else {
			var waitFiles []int64
			if userData.Info.Avatar > 0 {
				waitFiles = append(waitFiles, userData.Info.Avatar)
			}
			fileList := map[int64]string{}
			if len(waitFiles) > 0 {
				fileList, err = BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
					ClaimIDList: waitFiles,
					UserID:      0,
					OrgID:       0,
					IsPublic:    true,
				})
				if err != nil {
					//找不到文件，跳过
					err = nil
				}
			}
			loginData := FieldsSaveReportData{
				Token: FieldsSaveReportDataToken{
					Token: newTokenInfo.ID,
					Key:   newTokenInfo.Key,
				},
				UserData: userData,
				FileList: fileList,
			}
			var newKey string
			newKey, err = appendSave(loginData)
			RouterReport.Data(c, "login user save, json error, ", "登陆异常", err, newKey)
		}
	} else {
		//用户数据包获取失败，需反馈失败
		return
	}
	//强制更新用户数据包
	if userInfo.OrgID > 0 {
		OrgUserMod.PushUpdateUserData(userInfo.OrgID, userInfo.ID)
	}
}

// LoginData 登陆后数据包
func LoginData(c *gin.Context, userID int64) (userData UserCore.DataUserDataType, orgData []OrgCoreCore.DataGetBindByUserMarge, err error) {
	//获取用户基础数据
	userData, err = UserCore.GetUserData(&UserCore.ArgsGetUserData{
		UserID: userID,
	})
	if err != nil {
		RouterReport.ErrorLog(c, "not get user data, ", err, "user_not_exist", "用户不存在")
		return
	}
	//获取该用户的所有绑定关系
	orgData, err = OrgCoreCore.GetBindByUserMarge(&OrgCoreCore.ArgsGetBindByUser{
		UserID: userData.Info.ID,
	})
	if err != nil {
		err = nil
	}
	return
}

// LoginData2 登陆后数据包
func LoginData2(userID int64) (userData UserCore.DataUserDataType, orgData []OrgCoreCore.DataGetBindByUserMarge, err error) {
	//获取用户基础数据
	userData, err = UserCore.GetUserData(&UserCore.ArgsGetUserData{
		UserID: userID,
	})
	if err != nil {
		return
	}
	//获取该用户的所有绑定关系
	orgData, err = OrgCoreCore.GetBindByUserMarge(&OrgCoreCore.ArgsGetBindByUser{
		UserID: userData.Info.ID,
	})
	if err != nil {
		err = nil
	}
	return
}

// 通用检查用户是否可用？
func checkUserNotBan(c *gin.Context, userInfo *UserCore.FieldsUserType) bool {
	if userInfo.Status == 1 {
		RouterReport.ErrorLog(c, "login user is ban", nil, "user_audit", "用户正在审核或需要激活，请检查您的短信或邮箱激活。")
		return false
	}
	if userInfo.Status != 2 {
		RouterReport.ErrorLog(c, "login user is ban", nil, "user_ban", "用户已被禁用")
		return false
	}
	SafetyUserON, err := BaseConfig.GetDataBool("SafetyUserON")
	if err != nil {
		SafetyUserON = true
	}
	if SafetyUserON && BasePedometer.CheckData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "safe-user", ID: userInfo.ID},
	}) {
		RouterReport.ErrorLog(c, "login user is ban", nil, "user_ban", "用户登陆异常，请稍后再试")
		return false
	}
	return true
}

// 清理和重建token
func clearAndCreateToken(c *gin.Context, userInfo *UserCore.FieldsUserType, loginMark string, rememberMe bool) (data BaseToken2.FieldsToken, b bool) {
	var err error
	//获取token数据
	tokenInfo := Router2Mid.GetTokenInfo(c)
	//建立新的token
	BaseToken2.DeleteToken(tokenInfo.ID)
	tokenID, errCode, err := BaseToken2.Create(&BaseToken2.ArgsCreate{
		UserID:     userInfo.ID,
		OrgID:      userInfo.OrgID,
		OrgBindID:  0,
		DeviceID:   0,
		LoginFrom:  loginMark,
		IP:         c.ClientIP(),
		Key:        "",
		IsRemember: rememberMe,
	})
	if err != nil {
		RouterReport.ErrorLog(c, fmt.Sprint("cannot delete token, token id: ", tokenInfo.ID), err, "token_error", errCode)
		return
	}
	data = BaseToken2.GetByID(tokenID)
	//反馈成功
	b = true
	return
}
