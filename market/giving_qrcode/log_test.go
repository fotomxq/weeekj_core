package MarketGivingQrcode

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
	errCode, err := CreateLog(&ArgsCreateLog{
		OrgID:          newConfigData.OrgID,
		QrcodeID:       newConditionsData.ID,
		UserID:         TestOrg.UserInfo.ID,
		ReferrerUserID: TestOrg.UserInfo.ID,
		ReferrerBindID: 0,
		PriceTotal:     500,
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
