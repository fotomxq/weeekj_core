package ServiceAD

import (
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	TestOrgArea "github.com/fotomxq/weeekj_core/v5/tools/test_org_area"
	"testing"
)

func TestInit(t *testing.T) {
	TestOrgArea.CreateChildArea(t, "ad")
}

// 普通绑定测试
func TestPutAD(t *testing.T) {
	TestCreateAD(t)
	TestSetBind(t)
	data, errCode, err := PutAD(&ArgsPutAD{
		OrgID:  TestOrg.OrgData.ID,
		AreaID: TestOrgArea.AreaData.ID,
		Mark:   "test",
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("find put ad: ", data)
	}
}

func TestPutADByGPS(t *testing.T) {
	data, errCode, err := PutADByGPS(&ArgsPutADByGPS{
		MapType: 0,
		Point: CoreSQLGPS.FieldsPoint{
			Longitude: 5,
			Latitude:  6,
		},
		OrgID: TestOrg.OrgData.ID,
		Mark:  "test",
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("find put ad: ", data)
	}
}

func TestClearPut(t *testing.T) {
	TestDeleteBind(t)
	TestDeleteAD(t)
}

func TestClear(t *testing.T) {
	TestOrgArea.Clear(t)
}
