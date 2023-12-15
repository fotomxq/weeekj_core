package RouterFinance

import (
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ReportSuccessPage 返回的页面URL地址，并跳转到该页面地址
func ReportSuccessPage(c *gin.Context, orgID int64, payID int64) {
	//获取配置
	appURL, err := BaseConfig.GetDataString("AppURL")
	if err != nil {
		return
	}
	appendURL, err := BaseConfig.GetDataString("FinancePayResultSuccessURL")
	if err != nil {
		c.Redirect(202, appURL)
		return
	}
	//跳转页面你
	c.Redirect(http.StatusMovedPermanently, fmt.Sprint(appURL, appendURL, "orgid=", orgID, "&payid=", payID))
}

// ReportFailedPage 支付失败后跳转URL
func ReportFailedPage(c *gin.Context, orgID int64, payID int64) {
	//获取配置
	appURL, err := BaseConfig.GetDataString("AppURL")
	if err != nil {
		return
	}
	appendURL, err := BaseConfig.GetDataString("FinancePayResultFailedURL")
	if err != nil {
		c.Redirect(202, appURL)
		return
	}
	//跳转页面你
	c.Redirect(http.StatusMovedPermanently, fmt.Sprint(appURL, appendURL, "orgid=", orgID, "&payid=", payID))
}
