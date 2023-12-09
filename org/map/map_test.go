package OrgMap

import (
	"errors"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	isInit     = false
	newMapData FieldsMap
)

func TestTimeInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
		isInit = true
	}
	TestOrg.LocalCreateOrg(t)
}

func TestSetMap(t *testing.T) {
	data, err := SetMap(&ArgsSetMap{
		OrgID:       TestOrg.OrgData.ID,
		ParentID:    0,
		CoverFileID: 0,
		Name:        "测试地址",
		Des:         "测试描述",
		Country:     86,
		Province:    140000,
		City:        140100,
		Address:     "测试地址描述",
		MapType:     1,
		Longitude:   112.599725,
		Latitude:    37.798266,
		Params:      nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newMapData = data
	}
}

func TestGetMapList(t *testing.T) {
	dataList, dataCount, err := GetMapList(&ArgsGetMapList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:       -1,
		ParentID:    -1,
		NeedIsAudit: true,
		IsAudit:     false,
		IsRemove:    false,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		if dataList[0].ID != newMapData.ID {
			err = errors.New("no data by dataList key 0")
		}
	}
}

func TestAuditMap(t *testing.T) {
	err := AuditMap(&ArgsAuditMap{
		ID: newMapData.ID,
	})
	ToolsTest.ReportError(t, err)
	if err == nil {
		dataList, dataCount, err := GetMapList(&ArgsGetMapList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: 1,
				Max:  10,
				Sort: "id",
				Desc: true,
			},
			OrgID:       -1,
			ParentID:    -1,
			NeedIsAudit: true,
			IsAudit:     true,
			IsRemove:    false,
			Search:      "",
		})
		ToolsTest.ReportDataList(t, err, dataList, dataCount)
		if err == nil && len(dataList) > 0 {
			if dataList[0].ID != newMapData.ID {
				err = errors.New("no data by dataList key 0")
			}
		}
	}
}

func TestFindMapByArea(t *testing.T) {
	data, err := FindMapByArea(&ArgsFindMapByArea{
		Country:   86,
		Province:  140000,
		City:      140100,
		MapType:   1,
		Longitude: 112.602702,
		Latitude:  37.797206,
		Radius:    0.1,
	})
	ToolsTest.ReportData(t, err, data)
	data, err = FindMapByArea(&ArgsFindMapByArea{
		Country:   newMapData.Country,
		Province:  newMapData.Province,
		City:      newMapData.City,
		MapType:   newMapData.MapType,
		Longitude: 15.0,
		Latitude:  15.1,
		Radius:    3.0,
	})
	if err == nil && len(data) > 0 {
		t.Error("have data")
	}
}

func TestDeleteMap(t *testing.T) {
	err := DeleteMap(&ArgsDeleteMap{
		ID:    newMapData.ID,
		OrgID: newMapData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClear(t *testing.T) {
	TestOrg.LocalClear(t)
}
