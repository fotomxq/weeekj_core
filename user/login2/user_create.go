package UserLogin2

import (
	AnalysisUserVisit "github.com/fotomxq/weeekj_core/v5/analysis/user_visit"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

// ArgsCreateUser 创建用户并完成邀请处理等机制的汇总处理参数
type ArgsCreateUser struct {
	//注册方式
	// phone 手机号注册；admin 后台强制注册; email 邮箱注册
	// weixin_wxx 微信小程序授权注册; weixin_wxx_phone 微信小程序授权手机号注册; weixin_app 微信
	RegFrom string `json:"regFrom"`
	//推荐人手机号
	ReferrerNationCode string `json:"referrerNationCode" check:"nationCode" empty:"true"`
	ReferrerPhone      string `json:"referrerPhone" check:"phone" empty:"true"`
}

// CreateUser 创建用户并完成邀请处理等机制的汇总处理
func CreateUser(createArgs *UserCore.ArgsCreateUser, args *ArgsCreateUser) (userInfo UserCore.FieldsUserType, errCode string, err error) {
	//尝试获取推荐人用户数据
	var referrerUserID, referrerBindID int64
	if args.ReferrerPhone != "" {
		//如果referrerNationCode为空，则默认采用86
		if args.ReferrerNationCode == "" {
			args.ReferrerNationCode = "86"
		}
		//获取推荐人用户信息
		refUserData, _ := UserCore.GetUserByPhone(&UserCore.ArgsGetUserByPhone{
			OrgID:      createArgs.OrgID,
			NationCode: args.ReferrerNationCode,
			Phone:      args.ReferrerPhone,
		})
		//如果没有找到用户，则在全局查询
		if refUserData.ID < 1 {
			refUserData, _ = UserCore.GetUserByPhone(&UserCore.ArgsGetUserByPhone{
				OrgID:      -1,
				NationCode: args.ReferrerNationCode,
				Phone:      args.ReferrerPhone,
			})
		}
		if refUserData.ID > 0 {
			referrerUserID = refUserData.ID
			bindData, _ := OrgCore.GetBindByUserAndOrg(&OrgCore.ArgsGetBindByUserAndOrg{
				UserID: refUserData.ID,
				OrgID:  refUserData.OrgID,
			})
			if bindData.ID > 0 {
				referrerBindID = bindData.ID
			}
		}
	}
	//将推荐人数据覆盖到创建参数中
	if referrerUserID > 0 || referrerBindID > 0 {
		createArgs.Infos = CoreSQLConfig.Set(createArgs.Infos, "referrerUserID", referrerUserID)
		createArgs.Infos = CoreSQLConfig.Set(createArgs.Infos, "referrerBindID", referrerBindID)
		createArgs.Infos = CoreSQLConfig.Set(createArgs.Infos, "referrerNationCode", args.ReferrerNationCode)
		createArgs.Infos = CoreSQLConfig.Set(createArgs.Infos, "referrerPhone", args.ReferrerPhone)
	}
	//创建用户
	userInfo, errCode, err = UserCore.CreateUser(createArgs)
	if err != nil {
		return
	}
	//收尾处理
	CreateAndFinal(userInfo.OrgID, userInfo.ID, args)
	//反馈
	return
}

// CreateAndFinal 外部注册使用内部方法
func CreateAndFinal(orgID int64, userID int64, args *ArgsCreateUser) {
	//处理邀请人机制
	if args.ReferrerPhone != "" {
		//TODO: 该设计可以撤销，后续将根据订阅/user/core/create_user来获取数据，内部的infos包含了邀请人信息
		CoreNats.PushDataNoErr("user_login2_new", "/user/login2/new", "new", userID, "", map[string]interface{}{
			"orgID":              orgID,
			"referrerNationCode": args.ReferrerNationCode,
			"referrerPhone":      args.ReferrerPhone,
		})
	}
	//根据注册渠道，处理统计
	switch args.RegFrom {
	case "phone":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: orgID,
			Mark:  5,
			Count: 1,
		})
	case "email":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: orgID,
			Mark:  8,
			Count: 1,
		})
	case "admin":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: orgID,
			Mark:  7,
			Count: 1,
		})
	case "weixin_wxx":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: orgID,
			Mark:  6,
			Count: 1,
		})
	case "weixin_wxx_phone":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: orgID,
			Mark:  6,
			Count: 1,
		})
	case "weixin_app":
		_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
			OrgID: orgID,
			Mark:  9,
			Count: 1,
		})
	}
	_ = AnalysisUserVisit.CreateCount(&AnalysisUserVisit.ArgsCreateCount{
		OrgID: orgID,
		Mark:  0,
		Count: 1,
	})
}
