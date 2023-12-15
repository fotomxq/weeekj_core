package UserAddress

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitDefault(t *testing.T) {
	ToolsTest.Init(t)
	TestCreate(t)
}

func TestSetDefault(t *testing.T) {
	err := SetDefault(&ArgsSetDefault{
		AddressID: newData.ID,
		UserID:    userID,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetDefaultAddress(t *testing.T) {
	data, err := GetDefaultAddress(&ArgsGetDefaultAddress{
		UserID: userID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestClear(t *testing.T) {
	TestDelete(t)
}
