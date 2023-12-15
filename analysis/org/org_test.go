package AnalysisOrg

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	TMSTransport "github.com/fotomxq/weeekj_core/v5/tms/transport"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"testing"
)

func TestInitMarge(t *testing.T) {
	TestInit(t)
	_, _, _ = UserCore.CreateUser(&UserCore.ArgsCreateUser{
		OrgID:              TestOrg.OrgData.ID,
		Name:               "测试账户",
		Password:           "",
		NationCode:         "",
		Phone:              "",
		AllowSkipWaitEmail: false,
		Email:              "",
		Username:           "",
		Avatar:             0,
		Status:             2,
		Parents:            UserCore.FieldsUserParents{},
		Groups:             UserCore.FieldsUserGroupsType{},
		Infos:              CoreSQLConfig.FieldsConfigsType{},
		Logins:             UserCore.FieldsUserLoginsType{},
	})
	for i := 0; i < 11; i++ {
		tmsData, _, err := TMSTransport.CreateTransport(&TMSTransport.ArgsCreateTransport{
			OrgID:       TestOrg.OrgData.ID,
			BindID:      TestOrg.BindData.ID,
			InfoID:      0,
			UserID:      TestOrg.UserInfo.ID,
			FromAddress: CoreSQLAddress.FieldsAddress{},
			ToAddress:   CoreSQLAddress.FieldsAddress{},
			OrderID:     0,
			Goods: []TMSTransport.FieldsTransportGood{
				{
					System: "mall",
					ID:     1,
					Mark:   "",
					Name:   "测试商品",
					Count:  1,
				},
			},
			Weight:    0,
			Length:    0,
			Width:     0,
			Currency:  0,
			Price:     0,
			PayFinish: true,
			TaskAt:    "",
			Params:    nil,
		})
		if err == nil && i < 9 {
			if err := TMSTransport.UpdateTransportPick(&TMSTransport.ArgsUpdateTransportPick{
				ID:     tmsData.ID,
				OrgID:  -1,
				BindID: -1,
			}); err != nil {
				t.Error(err)
			}
			if err := TMSTransport.UpdateTransportSend(&TMSTransport.ArgsUpdateTransportSend{
				ID:     tmsData.ID,
				OrgID:  -1,
				BindID: -1,
			}); err != nil {
				t.Error(err)
			}
			if err := TMSTransport.UpdateTransportFinish(&TMSTransport.ArgsUpdateTransportFinish{
				ID:            tmsData.ID,
				OrgID:         -1,
				BindID:        -1,
				OperateBindID: tmsData.BindID,
				IsOrderRefund: false,
			}); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestGetMarge(t *testing.T) {
	data, err := GetMarge(&ArgsMarge{
		Marks: []ArgsMargeMark{
			{
				Mark: "user_new",
				TimeBetween: CoreSQLTime.DataCoreTime{
					MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubMonth().Time),
					MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().AddHour().Time),
				},
				Limit: 1,
			},
			{
				Mark: "tms_wait_count",
				TimeBetween: CoreSQLTime.DataCoreTime{
					MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubMonth().Time),
					MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().AddHour().Time),
				},
				Limit: 1,
			},
		},
		OrgID: TestOrg.OrgData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
	data, err = GetMarge(&ArgsMarge{
		Marks: []ArgsMargeMark{
			{
				Mark: "tms_time_count_month",
				TimeBetween: CoreSQLTime.DataCoreTime{
					MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubMonth().Time),
					MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().AddHour().Time),
				},
				Limit: 12,
			},
		},
		OrgID: TestOrg.OrgData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestClearMarge(t *testing.T) {
	TestClear(t)
}
