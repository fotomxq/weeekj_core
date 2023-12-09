package RouterMidAPI

import (
	"errors"
	"fmt"
	BaseToken2 "gitee.com/weeekj/weeekj_core/v5/base/token2"
	"github.com/gin-gonic/gin"
)

// checkAPIAndToken 获取token并验证
func checkAPIAndToken(c *gin.Context, urlAction string) error {
	//从form获取数据
	action := c.GetHeader("action")
	timestamp := c.GetHeader("timestamp")
	nonce := c.GetHeader("nonce")
	secretID := c.GetHeader("secret_id")
	signatureKey := c.GetHeader("signature_key")
	signatureMethod := c.GetHeader("signature_method")
	//检查api
	tokenID, err := BaseToken2.Check(action, urlAction, timestamp, nonce, secretID, signatureKey, signatureMethod)
	if err != nil || tokenID < 1 {
		return errors.New(fmt.Sprint("check token api, ", err))
	}
	//保存token数据
	c.Set("tokenID", tokenID)
	//反馈
	return nil
}
