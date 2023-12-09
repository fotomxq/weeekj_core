package MapAddress

import "testing"

func TestGetAddressByStr(t *testing.T) {
	data := GetAddressByStr("")
	t.Log("Province:", data.Province)
	t.Log("City", data.City)
	t.Log("Street:", data.Street)
	t.Log("Address:", data.Address)
}
