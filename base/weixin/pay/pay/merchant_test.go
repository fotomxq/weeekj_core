package BaseWeixinPayPay

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"testing"
)

func TestInit2(t *testing.T) {
	TestInit(t)
}

// 验证注意
// 以用户 openID: o3LjF5AW-BgZVL8vjy5VUzQ6EhSA 为案例
func TestMerchantChange(t *testing.T) {
	orderID, err := CoreFilter.GetRandStr3(10)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	openID := "o3LjF5AW-BgZVL8vjy5VUzQ6EhSA"
	//发起请求
	t.Log("send: ", orderID, openID, "刘子路", "测试交易", 100)
	jsonByte, err := MerchantChange(&ArgsMerchantChange{
		OrgID:      0,
		PayKey:     orderID,
		UserOpenID: openID,
		UserName:   "刘子路",
		PayDes:     "测试交易",
		Price:      100,
	})
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	t.Log(string(jsonByte))
}
