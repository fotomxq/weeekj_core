package CoreLanguage

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	"github.com/gin-gonic/gin"
)

// GetLanguageText 获取指定语言数据
func GetLanguageText(c *gin.Context, mark string) string {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("get language text: ", mark, ", err: ", r)
			return
		}
	}()
	lang := getLanguageObj(c)
	str := loadVal(&lang, mark)
	if str == "" {
		CoreLog.Warn("get language failed, text is unknow, mark: ", mark)
		return "unknow"
	}
	return str
}

// getLanguageObj 通过浏览器头获取语言包类型
// 语言包参考：https://github.com/gohouse/i18n
func getLanguageObj(c *gin.Context) dataLanguage {
	//获取语言包类型
	language := c.Request.Header.Get("language")
	if language == "" {
		language = "zh_cn"
	} else {
		switch language {
		case "zh_cn":
		default:
			language = "zh_cn"
		}
	}
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("get language data: ", language, ", err: ", r)
			return
		}
	}()
	//获取语言包
	return loadLanguage(language)
}
