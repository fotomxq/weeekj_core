package MarketGivingNewUser

import (
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

func TestInitLog(t *testing.T) {
	TestInitConditions(t)
	TestCreateConfig(t)
	TestCreateConditions(t)
}

func TestCreateLog(t *testing.T) {
	errCode, err := createLog(&argsCreateLog{
		OrgID:          newConfigData.OrgID,
		UserID:         TestOrg.UserInfo.ID,
		ReferrerUserID: TestOrg.UserInfo.ID,
		ReferrerBindID: 0,
		PriceTotal:     500,
		IsOrder:        false,
	})
	if err != nil {
		t.Error(errCode, ", err: ", err)
	}
}

func TestClearLog(t *testing.T) {
	TestDeleteConditions(t)
	TestDeleteConfig(t)
	TestClearConditions(t)
}
