package Router2APISecret

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"testing"
)

func TestCheckAPI(t *testing.T) {
	str := "/v4/base/channel/public/one16798248311080614170e0dff2e9b005c2e2e629"
	strSha256, _ := CoreFilter.GetSha256([]byte(str))
	t.Log("key: ", string(strSha256))
}
