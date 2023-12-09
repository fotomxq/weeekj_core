package CoreFilter

import "testing"

func TestGetFileIDListByContent(t *testing.T) {
	str := "abc{file:abc01}abc{file:abc02}abcDef{file:abc03}"
	finds, err := GetFileIDListByContent(str)
	if err != nil{
		t.Error(err)
	}else{
		t.Log(finds)
	}
}