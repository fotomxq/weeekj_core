package BaseSMS

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newConfigData FieldsConfig
)

func TestInitConfig(t *testing.T) {
	TestInit(t)
}

func TestCreateConfig(t *testing.T) {
	var err error
	newConfigData, err = CreateConfig(&ArgsCreateConfig{
		OrgID:         0,
		System:        "aliyun",
		Name:          "阿里云短信",
		AppID:         "LTAI4FpWxT9qsYh6XX87h8mv",
		AppKey:        "XA6JjdktsQP18TSCcKrrTIkkc6pJrZ",
		DefaultExpire: "300s",
		TimeSpacing:   60,
		TemplateID:    "SMS_181852664",
		TemplateSign:  "海多智链",
		TemplateParams: CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "val",
				Val:  "code",
			},
			{
				Mark: "time",
				Val:  "__skip",
			},
		},
		Params: CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "aliyunRegionID",
				Val:  "cn-hangzhou",
			},
			{
				Mark: "check",
				Val:  "true",
			},
		},
	})
	ToolsTest.ReportData(t, err, newConfigData)
}

func TestGetConfigByID(t *testing.T) {
	var err error
	newConfigData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    newConfigData.ID,
		OrgID: newConfigData.OrgID,
	})
	ToolsTest.ReportData(t, err, newConfigData)
}

func TestGetConfigList(t *testing.T) {
	dataList, dataCount, err := GetConfigList(&ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil {
		newConfigData = dataList[0]
	}
}

func TestUpdateConfig(t *testing.T) {
	err := UpdateConfig(&ArgsUpdateConfig{
		ID:             newConfigData.ID,
		OrgID:          newConfigData.OrgID,
		System:         newConfigData.System,
		Name:           newConfigData.Name,
		AppID:          newConfigData.AppID,
		AppKey:         newConfigData.AppKey,
		DefaultExpire:  newConfigData.DefaultExpire,
		TimeSpacing:    newConfigData.TimeSpacing,
		TemplateID:     newConfigData.TemplateID,
		TemplateSign:   newConfigData.TemplateSign,
		TemplateParams: newConfigData.TemplateParams,
		Params:         newConfigData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteConfig(t *testing.T) {
	err := DeleteConfig(&ArgsDeleteConfig{
		ID:    newConfigData.ID,
		OrgID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearConfig(t *testing.T) {
	TestClear(t)
}
