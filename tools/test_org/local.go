package TestOrg

import (
	"fmt"
	ClassConfig "gitee.com/weeekj/weeekj_core/v5/class/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	OrgTime "gitee.com/weeekj/weeekj_core/v5/org/time"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	"testing"
	"time"
)

//本地化操作

// LocalInit 本地化初始化
func LocalInit() {
	UserCore.Init()
	OrgCore.Init()
	OrgTime.Init()
}

// LocalCreateUser 创建用户
func LocalCreateUser(t *testing.T) {
	var errCode string
	var err error
	UserInfo, errCode, err = CreateUser(t)
	if err != nil {
		t.Error(errCode, err)
		return
	} else {
		t.Log("create new user, ", UserInfo)
		return
	}
}

func CreateUser(t *testing.T) (userInfo UserCore.FieldsUserType, errCode string, err error) {
	userInfo, errCode, err = UserCore.CreateUser(&UserCore.ArgsCreateUser{
		OrgID:                0,
		Name:                 fmt.Sprint("测试账户", CoreFilter.GetRandNumber(1, 1000)),
		Password:             "",
		NationCode:           "",
		Phone:                "",
		AllowSkipPhoneVerify: false,
		AllowSkipWaitEmail:   false,
		Email:                "",
		Username:             "",
		Avatar:               0,
		Status:               2,
		Parents:              nil,
		Groups:               nil,
		Infos:                nil,
		Logins:               nil,
		SortID:               0,
		Tags:                 nil,
	})
	if err != nil {
		t.Error(errCode, err)
		return
	} else {
		//t.Log("create new user, ", userInfo)
		return
	}
}

// LocalCreateOrg 创建新的组织
func LocalCreateOrg(t *testing.T) {
	LocalCreateUser(t)
	var errCode string
	key, err := CoreFilter.GetRandStr3(10)
	if err != nil {
		t.Error(err)
		return
	}
	OrgData, errCode, err = OrgCore.CreateOrg(&OrgCore.ArgsCreateOrg{
		UserID:     UserInfo.ID,
		Key:        key,
		Name:       fmt.Sprint("测试组织", key),
		Des:        "测试描述",
		ParentID:   0,
		ParentFunc: []string{"all"},
		OpenFunc:   []string{"all"},
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new org, ", OrgData)
	}
}

// LocalCreateWorkTime 创建分组的工作时间
func LocalCreateWorkTime(t *testing.T) {
	LocalCreateOrg(t)
	var err error
	WorkTimeData, err = OrgTime.Create(&OrgTime.ArgsCreate{
		ExpireAt: time.Time{},
		OrgID:    OrgData.ID,
		Groups:   []int64{},
		Binds:    []int64{},
		Configs:  OrgTime.FieldsConfigs{},
		Name:     "测试工作时间",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("create new work time, ", WorkTimeData)
	}
}

// LocalCreateGroup 创建分组
func LocalCreateGroup(t *testing.T) {
	LocalCreateWorkTime(t)
	var err error
	GroupData, err = OrgCore.CreateGroup(&OrgCore.ArgsCreateGroup{
		OrgID:   OrgData.ID,
		Manager: []string{"all"},
		Name:    "测试组织分组",
		Params: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "test_mark",
				Val:  "test_value",
			},
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("create new group, ", GroupData)
	}
}

// LocalCreateBind 创建绑定关系
func LocalCreateBind(t *testing.T) {
	LocalCreateGroup(t)
	var err error
	BindData, err = CreateBind(t, UserInfo)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("create new bind, ", BindData)
	}
}

func CreateBind(t *testing.T, userInfo UserCore.FieldsUserType) (bindData OrgCore.FieldsBind, err error) {
	bindData, err = OrgCore.SetBind(&OrgCore.ArgsSetBind{
		UserID:     UserInfo.ID,
		Avatar:     0,
		Name:       UserInfo.Name,
		OrgID:      OrgData.ID,
		GroupIDs:   []int64{GroupData.ID},
		Manager:    []string{"all"},
		NationCode: "",
		Phone:      "",
		Email:      "",
		SyncSystem: "",
		SyncID:     0,
		SyncHash:   "",
		Params: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "test_mark",
				Val:  "test_value",
			},
		},
	})
	if err != nil {
		t.Error("create new bind data, ", err)
	}
	return
}

// LocalOrgFinance 为商户新增储蓄的支持设置
func LocalOrgFinance(t *testing.T) {
	//建立储蓄关联处理
	if err := OrgCore.Config.SetConfig(&ClassConfig.ArgsSetConfig{
		BindID:    OrgData.ID,
		Mark:      "FinanceDepositDefaultMark",
		VisitType: "admin",
		Val:       "savings",
	}); err != nil {
		t.Error("create org finance default config, ", err)
		return
	}
}

// LocalClear 清理
// 如果存在组织或任意内容，则将清理
func LocalClear(t *testing.T) {
	var err error
	if BindData.Name != "" {
		err = OrgCore.DeleteBind(&OrgCore.ArgsDeleteBind{
			ID:    BindData.ID,
			OrgID: OrgData.ID,
		})
		if err != nil {
			t.Error("delete org bind: ", err)
		}
		BindData = OrgCore.FieldsBind{}
	}
	if GroupData.Name != "" {
		err = OrgCore.DeleteGroup(&OrgCore.ArgsDeleteGroup{
			ID: GroupData.ID,
		})
		if err != nil {
			t.Error("delete org group, ", err)
		}
		GroupData = OrgCore.FieldsGroup{}
	}
	if WorkTimeData.Name != "" {
		err = OrgTime.DeleteByID(&OrgTime.ArgsDeleteByID{
			ID:    WorkTimeData.ID,
			OrgID: OrgData.ID,
		})
		if err != nil {
			t.Error("delete org time, ", err)
		}
		WorkTimeData = OrgTime.FieldsWorkTime{}
	}
	if OrgData.Name != "" {
		err = OrgCore.DeleteOrg(&OrgCore.ArgsDeleteOrg{
			ID: OrgData.ID,
		})
		if err != nil {
			t.Error("delete org by id, ", err)
		}
		OrgData = OrgCore.FieldsOrg{}
	}
	if UserInfo.Name != "" {
		err = UserCore.DeleteUserByID(&UserCore.ArgsDeleteUserByID{
			ID: UserInfo.ID,
		})
		if err != nil {
			t.Error("delete user by id, ", err)
		}
		UserInfo = UserCore.FieldsUserType{}
	}
}
