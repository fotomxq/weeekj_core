package BaseToken2

import (
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"testing"
)

func TestCheck(t *testing.T) {
	tokenID, err := Check("/v4/finance/deposit/om/org_self", "/v4/finance/deposit/om/org_self", "1679825386", "16798", "14170", "e0dff2e9b005c2e2e629", "sha256")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(tokenID)
	}
}

func TestCheckByKey(t *testing.T) {
	err := CheckByKey("/v2/finance/pay/org/manager", "/v2/finance/pay/org/manager", "1679826199", "29465", "14170", "e007deb30ad0f89e291a7cc43045c13e3610cc09e498958ee2ab18dda6c44a24", "sha256", "e0dff2e9b005c2e2e629")
	if err != nil {
		t.Log("ok")
	} else {
		t.Log("no")
	}
	strSHA256Str := "/v2/finance/pay/org/manager29465167982619914170e0dff2e9b005c2e2e629"
	strSHA256, _ := CoreFilter.GetSha256Str(strSHA256Str)
	t.Log("SHA256: ", strSHA256)
	str := fmt.Sprint("/v2/finance/pay/org/manager", 29465, 1679826199, 14170, "e0dff2e9b005c2e2e629")
	t.Log(CoreFilter.GetSha1Str(str))
	str2, err := CoreFilter.GetSha1([]byte(str))
	if err != nil {
		t.Error(err)
	} else {
		t.Log("SHA1 byte: ", string(str2))
	}
}
