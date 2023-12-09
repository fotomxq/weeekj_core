package ERPWarehouse

import "testing"

var (
	newLocationData FieldsLocation
)

func TestLocationInit(t *testing.T) {
	TestAreaInit(t)
	TestCreateArea(t)
	TestGetAreaList(t)
}

//TODO: 等待货位管理模块完成后，进行单元测试

func TestLocationClear(t *testing.T) {
	TestDeleteArea(t)
	TestAreaClear(t)
}
