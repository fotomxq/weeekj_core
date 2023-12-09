package TestAPI

import (
	"testing"

	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	if !isInit {
		//获取hash随机值
		newUpdateHash := CoreFilter.GetRandNumber(1, 10)
		newUpdateHash2, err := CoreFilter.GetRandStr3(10)
		if err != nil {
			t.Error(err)
		}
		SetToken(int64(newUpdateHash), newUpdateHash2)
		SetBaseURL("http://localhost:8000")
	}
	isInit = true
}
