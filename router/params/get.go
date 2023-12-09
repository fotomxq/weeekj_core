package RouterParams

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
)

// GetID 获取参数带有ID的头
func GetID(c *gin.Context) (int64, bool) {
	return GetIDByName(c, "id")
}

func GetMark(c *gin.Context) (string, bool) {
	mark := c.Param("mark")
	if !CoreFilter.CheckMark(mark) {
		RouterReport.BaseError(c, "params_lost", "无效的标识码")
		return "", false
	}
	return mark, true
}

// GetIDByName 获取指定的ID数据
func GetIDByName(c *gin.Context, name string) (int64, bool) {
	id := c.Param(name)
	if !CoreFilter.CheckID(id) {
		RouterReport.BaseError(c, "params_lost", "无效的ID")
		return 0, false
	}
	idInt64, err := CoreFilter.GetInt64ByString(id)
	if err != nil {
		RouterReport.BaseError(c, "params_lost", "无效的ID")
		return 0, false
	}
	return idInt64, true
}

func GetIDByNameNoErr(c *gin.Context, name string) (int64, bool) {
	id := c.Param(name)
	if !CoreFilter.CheckID(id) {
		return 0, false
	}
	idInt64, err := CoreFilter.GetInt64ByString(id)
	if err != nil {
		return 0, false
	}
	return idInt64, true
}
