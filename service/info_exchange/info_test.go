package ServiceInfoExchange

import (
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newInfoData FieldsInfo
)

func TestInitInfo(t *testing.T) {
	TestInit(t)
}

func TestCreateInfo(t *testing.T) {
	data, err := CreateInfo(&ArgsCreateInfo{
		InfoType:     "none",
		ExpireAt:     "",
		OrgID:        TestOrg.OrgData.ID,
		UserID:       TestOrg.UserInfo.ID,
		SortID:       0,
		Tags:         []int64{},
		Title:        "测试信息",
		TitleDes:     "测试信息标题des",
		Des:          "测试信息des",
		CoverFileIDs: []int64{},
		Currency:     86,
		Price:        10,
		LimitCount:   10,
		Address: CoreSQLAddress.FieldsAddress{
			Country:    86,
			Province:   0,
			City:       10010,
			Address:    "测试地址",
			MapType:    0,
			Longitude:  1,
			Latitude:   1,
			Name:       "测试姓名",
			NationCode: "86",
			Phone:      "17777777777",
		},
		Params: CoreSQLConfig.FieldsConfigsType{},
	})
	if err != nil {
		t.Error(err)
	} else {
		newInfoData = data
		t.Log(newInfoData)
	}
}

func TestGetInfoList(t *testing.T) {
	dataList, dataCount, err := GetInfoList(&ArgsGetInfoList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:         -1,
		UserID:        -1,
		InfoType:      "",
		SortID:        -1,
		Tags:          []int64{},
		NeedIsAudit:   false,
		IsAudit:       false,
		NeedIsPublish: false,
		IsPublish:     false,
		PriceMin:      -1,
		PriceMax:      -1,
		NeedIsExpire:  false,
		IsExpire:      false,
		IsRemove:      false,
		Search:        "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetInfoID(t *testing.T) {
	data, err := GetInfoID(&ArgsGetInfoID{
		ID:     newInfoData.ID,
		OrgID:  newInfoData.OrgID,
		UserID: newInfoData.UserID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestPublishInfo(t *testing.T) {
	err := PublishInfo(&ArgsPublishInfo{
		ID:     newInfoData.ID,
		OrgID:  newInfoData.OrgID,
		UserID: newInfoData.UserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestAuditInfo(t *testing.T) {
	err := AuditInfo(&ArgsAuditInfo{
		ID:       newInfoData.ID,
		OrgID:    newInfoData.OrgID,
		IsAudit:  true,
		AuditDes: "",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetInfoPublishID(t *testing.T) {
	data, err := GetInfoPublishID(&ArgsGetInfoID{
		ID:     newInfoData.ID,
		OrgID:  newInfoData.OrgID,
		UserID: newInfoData.UserID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateInfo(t *testing.T) {
	err := UpdateInfo(&ArgsUpdateInfo{
		ID:           newInfoData.ID,
		OrgID:        newInfoData.OrgID,
		UserID:       newInfoData.UserID,
		ExpireAt:     "",
		SortID:       newInfoData.SortID,
		Tags:         newInfoData.Tags,
		Title:        newInfoData.Title,
		TitleDes:     newInfoData.TitleDes,
		Des:          newInfoData.Des,
		CoverFileIDs: newInfoData.CoverFileIDs,
		Currency:     newInfoData.Currency,
		Price:        newInfoData.Price,
		LimitCount:   newInfoData.LimitCount,
		Address:      newInfoData.Address,
		Params:       newInfoData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteInfo(t *testing.T) {
	err := DeleteInfo(&ArgsDeleteInfo{
		ID:     newInfoData.ID,
		OrgID:  newInfoData.OrgID,
		UserID: newInfoData.UserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearInfo(t *testing.T) {
	TestClear(t)
}
