package Router2Params

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	Router2Report "gitee.com/weeekj/weeekj_core/v5/router2/mid"
)

// GetID 获取参数带有ID的头
func GetID(context any) (int64, bool) {
	return GetIDByName(context, "id")
}

func GetMark(context any) (string, bool) {
	return GetMarkByName(context, "mark")
}

func GetMarkByName(context any, name string) (string, bool) {
	c := getContext(context)
	mark := c.Param(name)
	if !CoreFilter.CheckMark(mark) {
		Router2Report.ReportBaseError(context, "err_params_mark")
		return "", false
	}
	return mark, true
}

// GetInt64ByName 获取指定的Int64数据
func GetInt64ByName(context any, name string) (int64, bool) {
	c := getContext(context)
	id := c.Param(name)
	idInt64, err := CoreFilter.GetInt64ByString(id)
	if err != nil {
		Router2Report.ReportBaseError(context, "err_params_id")
		return 0, false
	}
	return idInt64, true
}

// GetIDByName 获取指定的ID数据
func GetIDByName(context any, name string) (int64, bool) {
	c := getContext(context)
	id := c.Param(name)
	if !CoreFilter.CheckID(id) {
		Router2Report.ReportBaseError(context, "err_params_id")
		return 0, false
	}
	idInt64, err := CoreFilter.GetInt64ByString(id)
	if err != nil {
		Router2Report.ReportBaseError(context, "err_params_id")
		return 0, false
	}
	return idInt64, true
}

func GetIDByNameNoErr(context any, name string) (int64, bool) {
	c := getContext(context)
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
