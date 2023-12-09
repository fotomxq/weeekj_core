package RouterMidWeb

import (
	BaseStyle "gitee.com/weeekj/weeekj_core/v5/base/style"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	OrgDomain "gitee.com/weeekj/weeekj_core/v5/org/domain"
	"github.com/gin-gonic/gin"
)

// RouterMid 静态网页中间件
// 不需要任何中间件处理的头部
func RouterMid(c *gin.Context) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("router mid, ", r)
		}
	}()
	//获取路由地址对应的商户
	orgID, params, err := OrgDomain.GetDomainOrg(&OrgDomain.ArgsGetDomainOrg{
		Host: c.Request.Host,
	})
	c.Set("orgID", orgID)
	if err != nil {
		getOrgStyle(c, 0, "")
		CoreLog.Warn("router mid get host: ", c.Request.Host, ", get org id, err: ", err)
		return
	}
	//根据orgID，获取商户的theme
	themeStyleMark, _ := params.GetVal("theme")
	getOrgStyle(c, orgID, themeStyleMark)
}

// 获取样式
func getOrgStyle(c *gin.Context, orgID int64, themeMark string) {
	//初始化
	var themeData BaseStyle.DataGetOrgByStyleMark
	var err error
	//如果不存在主题，则给与默认的
	if themeMark == "" {
		themeMark = "_default"
	}
	//获取样式结构
	themeData, err = BaseStyle.GetOrgByStyleMark(&BaseStyle.ArgsGetOrgByStyleMark{
		Mark:  themeMark,
		OrgID: orgID,
	})
	//保存数据
	if err != nil {
		c.Set("themeData", themeData)
		c.Set("themeMark", themeMark)
		CoreLog.Warn("router mid get host: ", c.Request.Host, ", get org theme, ", themeMark, ", err: ", err)
		return
	}
	c.Set("themeData", themeData)
	c.Set("themeMark", themeData.Mark)
}
