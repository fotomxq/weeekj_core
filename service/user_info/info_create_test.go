package ServiceUserInfo

import (
	"errors"
	"fmt"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	newInfoData FieldsInfo
)

func TestInitInfo(t *testing.T) {
	TestInit(t)
}

func TestCreateInfo(t *testing.T) {
	TestCreateTemplate(t)
	data, err := CreateInfo(&ArgsCreateInfo{
		OrgID:                 1,
		UserID:                0,
		BindID:                0,
		BindType:              0,
		Name:                  "测试文件",
		Country:               0,
		Gender:                0,
		IDCard:                "",
		Phone:                 "",
		CoverFileID:           0,
		DesFiles:              []int64{},
		Address:               "",
		DateOfBirth:           time.Time{},
		MaritalStatus:         false,
		EducationStatus:       0,
		Profession:            "",
		Level:                 0,
		EmergencyContact:      "",
		EmergencyContactPhone: "",
		SortID:                0,
		Tags:                  []int64{},
		DocID:                 0,
		Des:                   "",
		Director1:             0,
		Director2:             0,
		Params:                nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newInfoData = data
	} else {
		err = errors.New(fmt.Sprint("create info, ", err))
	}
}

func TestGetInfoID(t *testing.T) {
	data, err := GetInfoID(&ArgsGetInfoID{
		ID:    newInfoData.ID,
		OrgID: newInfoData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetInfoMore(t *testing.T) {
	data, err := GetInfoMore(&ArgsGetInfoMore{
		IDs:        []int64{newInfoData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetInfoMoreNames(t *testing.T) {
	data, err := GetInfoMoreNames(&ArgsGetInfoMore{
		IDs:        []int64{newInfoData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetOrgInfoMore(t *testing.T) {
	data, err := GetOrgInfoMore(&ArgsGetOrgInfoMore{
		IDs:        []int64{newInfoData.ID},
		HaveRemove: false,
		OrgID:      newInfoData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetOrgInfoMoreNames(t *testing.T) {
	data, err := GetOrgInfoMoreNames(&ArgsGetOrgInfoMore{
		IDs:        []int64{newInfoData.ID},
		HaveRemove: false,
		OrgID:      newInfoData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestCheckInfo(t *testing.T) {
	err := CheckInfo(&ArgsCheckInfo{
		ID:    newInfoData.ID,
		OrgID: newInfoData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetInfoList(t *testing.T) {
	dataList, dataCount, err := GetInfoList(&ArgsGetInfoList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		UserID:   -1,
		BindID:   -1,
		Country:  -1,
		SortID:   -1,
		Tags:     []int64{},
		Director: -1,
		IsDie:    false,
		IsOut:    false,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateInfo(t *testing.T) {
	err := UpdateInfo(&ArgsUpdateInfo{
		ID:                    newInfoData.ID,
		OrgID:                 newInfoData.OrgID,
		UserID:                newInfoData.UserID,
		BindID:                newInfoData.BindID,
		BindType:              newInfoData.BindType,
		Name:                  newInfoData.Name,
		Country:               newInfoData.Country,
		Gender:                newInfoData.Gender,
		IDCard:                newInfoData.IDCard,
		Phone:                 newInfoData.Phone,
		CoverFileID:           newInfoData.CoverFileID,
		DesFiles:              newInfoData.DesFiles,
		Address:               newInfoData.Address,
		DateOfBirth:           newInfoData.DateOfBirth,
		MaritalStatus:         newInfoData.MaritalStatus,
		EducationStatus:       newInfoData.EducationStatus,
		Profession:            newInfoData.Profession,
		Level:                 newInfoData.Level,
		EmergencyContact:      newInfoData.EmergencyContact,
		EmergencyContactPhone: newInfoData.EmergencyContactPhone,
		SortID:                newInfoData.SortID,
		Tags:                  newInfoData.Tags,
		DocID:                 newInfoData.DocID,
		Des:                   newInfoData.Des,
		Director1:             newInfoData.Director1,
		Director2:             newInfoData.Director2,
		Params:                newInfoData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteInfo(t *testing.T) {
	err := DeleteInfo(&ArgsDeleteInfo{
		ID:    newInfoData.ID,
		OrgID: newInfoData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteTemplate(t)
}
