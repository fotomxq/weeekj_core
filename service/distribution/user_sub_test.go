package ServiceDistribution

import (
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	MarketGivingCore "github.com/fotomxq/weeekj_core/v5/market/giving_core"
	MarketGivingUserSub "github.com/fotomxq/weeekj_core/v5/market/giving_user_sub"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
	"testing"
)

var (
	newUserSubData       FieldsUserSub
	newConditionsData    MarketGivingUserSub.FieldsConditions
	newConfigData        MarketGivingCore.FieldsConfig
	newUserSubConfigData UserSubscription.FieldsConfig
)

func TestInitUserSub(t *testing.T) {
	TestInitDistribution(t)
	TestCreateDistribution(t)
	var err error
	//构建用户会员
	newUserSubConfigData, err = UserSubscription.CreateConfig(&UserSubscription.ArgsCreateConfig{
		OrgID:             newDistributionData.OrgID,
		Mark:              "",
		TimeType:          0,
		TimeN:             0,
		Currency:          86,
		Price:             99,
		PriceOld:          109,
		Title:             "测试会员",
		Des:               "",
		CoverFileID:       0,
		DesFiles:          []int64{},
		UserGroups:        []int64{},
		ExemptionPrice:    0,
		ExemptionDiscount: 0,
		ExemptionMinPrice: 0,
		Limits:            []UserSubscription.FieldsLimit{},
		ExemptionTime:     UserSubscription.FieldsExemptionTimes{},
		StyleID:           0,
		Params:            nil,
	})
	ToolsTest.ReportError(t, err)
}

func TestCreateConfig(t *testing.T) {
	var err error
	newConfigData, err = MarketGivingCore.CreateConfig(&MarketGivingCore.ArgsCreateConfig{
		OrgID:             TestOrg.OrgData.ID,
		Name:              "测试推荐处理",
		MarketConfigID:    0,
		LimitTimeType:     0,
		LimitCount:        0,
		UserIntegral:      0,
		UserSubs:          MarketGivingCore.FieldsConfigUserSubs{},
		UserTickets:       MarketGivingCore.FieldsConfigUserTickets{},
		DepositConfigMark: "",
		Price:             0,
		Count:             1,
		Params:            nil,
	})
	ToolsTest.ReportData(t, err, newConfigData)
}

func TestCreateConditions(t *testing.T) {
	var err error
	newConditionsData, err = MarketGivingUserSub.CreateConditions(&MarketGivingUserSub.ArgsCreateConditions{
		OrgID:       TestOrg.OrgData.ID,
		Name:        "测试条件",
		ConfigID:    newConfigData.ID,
		SubConfigID: newUserSubConfigData.ID,
		SubBuyCount: 1,
		Params:      nil,
	})
	ToolsTest.ReportData(t, err, newConditionsData)
}

func TestCreateUserSub(t *testing.T) {
	var err error
	//构建订阅关系
	err = CreateUserSub(&ArgsCreateUserSub{
		OrgID:             newDistributionData.OrgID,
		DistributionID:    newDistributionData.ID,
		SubConfigID:       newUserSubConfigData.ID,
		UnitPrice:         50,
		MarketGivingSubID: newConditionsData.ID,
		CoverFileID:       0,
		Des:               "测试内容",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetUserSubList(t *testing.T) {
	dataList, dataCount, err := GetUserSubList(&ArgsGetUserSubList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    -1,
		IsRemove: false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newUserSubData = dataList[0]
		t.Log("newUserSubData.UnitPrice: ", newUserSubData.UnitPrice)
	}
}

func TestUpdateUserSub(t *testing.T) {
	err := UpdateUserSub(&ArgsUpdateUserSub{
		ID:                newUserSubData.ID,
		OrgID:             newUserSubData.OrgID,
		DistributionID:    newUserSubData.DistributionID,
		SubConfigID:       newUserSubData.SubConfigID,
		UnitPrice:         newUserSubData.UnitPrice,
		MarketGivingSubID: newUserSubData.MarketGivingSubID,
		CoverFileID:       newUserSubData.CoverFileID,
		Des:               newUserSubData.Des,
	})
	ToolsTest.ReportError(t, err)
	TestGetUserSubList(t)
}

func TestInUserSub(t *testing.T) {
	data, err := InUserSub(&ArgsInUserSub{
		ID:    newUserSubData.ID,
		OrgID: newUserSubData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestBuyUserSub(t *testing.T) {
	data, errCode, err := BuyUserSub(&ArgsBuyUserSub{
		OrgID:                 TestOrg.OrgData.ID,
		UserID:                TestOrg.UserInfo.ID,
		CreateFrom:            0,
		AddressFrom:           CoreSQLAddress.FieldsAddress{},
		AddressTo:             CoreSQLAddress.FieldsAddress{},
		Des:                   "测试购物会员",
		SubConfigID:           newUserSubConfigData.ID,
		Unit:                  1,
		ReferrerNationCode:    "",
		ReferrerPhone:         "",
		Params:                nil,
		DistributionUserSubID: newUserSubData.ID,
	})
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		t.Error("errCode: ", errCode)
	} else {
		if data.Price != newUserSubData.UnitPrice {
			t.Error("price not: ", newUserSubData.UnitPrice, ", order price: ", data.Price)
		}
	}
}

func TestDeleteUserSub(t *testing.T) {
	err := DeleteUserSub(&ArgsDeleteUserSub{
		ID:    newUserSubData.ID,
		OrgID: newUserSubData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearUserSub(t *testing.T) {
	TestDeleteDistribution(t)
	TestClearDistribution(t)
}
