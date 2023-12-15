package RouterMidAPI

import (
	"errors"
	"fmt"
	BaseToken2 "github.com/fotomxq/weeekj_core/v5/base/token2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/gin-gonic/gin"
)

// TokenGetCookie 通过cookie获取token
func TokenGetCookie(c *gin.Context) error {
	//获取token数据
	tokenIDStr, err := c.Cookie("tokenID")
	if err != nil {
		return err
	}
	tokenKey, err := c.Cookie("tokenKey")
	if err != nil {
		return err
	}
	//解析并获取会话数据
	tokenID, err := CoreFilter.GetInt64ByString(tokenIDStr)
	if err != nil {
		return err
	}
	//获取token
	tokenData := BaseToken2.GetByID(tokenID)
	if tokenData.ID < 1 {
		return errors.New("no token")
	}
	//检查密钥是否匹配
	if tokenData.Key != tokenKey {
		return errors.New("token key error")
	}
	//反馈成功
	return nil
}

// SetTokenCookie 设置头部到cookie
func SetTokenCookie(c *gin.Context, tokenData BaseToken2.FieldsToken) {
	//设置token数据
	expireUnix := tokenData.ExpireAt.Unix() - CoreFilter.GetNowTime().Unix()
	if expireUnix < 1 {
		expireUnix = 3600
	}
	c.SetCookie("tokenID", fmt.Sprint(tokenData.ID), int(expireUnix), "/", "", true, true)
	c.SetCookie("tokenKey", fmt.Sprint(tokenData.Key), int(expireUnix), "/", "", true, true)
}

// UpdateTokenCookie 更新cookie会话
func UpdateTokenCookie(c *gin.Context) error {
	//获取token数据
	tokenInfo := getTokenInfo(c)
	if tokenInfo.ID < 1 {
		return errors.New("no token")
	}
	//更新cookie
	SetTokenCookie(c, tokenInfo)
	//反馈
	return nil
}

// ClearTokenCookie 清理token
func ClearTokenCookie(c *gin.Context) {
	c.SetCookie("tokenID", "", -1, "/", "", true, true)
	c.SetCookie("tokenKey", "", -1, "/", "", true, true)
}
