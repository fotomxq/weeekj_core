package RouterAPISecret

import (
	"testing"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

func TestMakeSignatureKey(t *testing.T) {
	t.Log(makeSignatureKey("/v2/login/login/header/qrcode/make", "1657178584", "32917", "236336", "d168b7d3a34dc99f5e88", "sha256"))
	str := "/v2/login/login/header/qrcode/make16571814175955221836481d583487460e5c4"
	strSha256, _ := CoreFilter.GetSha256([]byte(str))
	t.Log("strSha256: ", string(strSha256))
}
