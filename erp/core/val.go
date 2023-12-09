package ERPCore

import (
	"fmt"
	BaseFileSys2 "gitee.com/weeekj/weeekj_core/v5/base/filesys2"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLanguage "gitee.com/weeekj/weeekj_core/v5/core/language"
	"github.com/gin-gonic/gin"
)

// GetValOfType 获取组件的值，并解析为指定的类型
// 注意反馈为any泛型，可直接转为对应类型
// isShow 将转化为可直接显示的数据，例如bool为是否，而不是true/false
func GetValOfType(ctx *gin.Context, data *FieldsComponentVal, isShow bool) (result any) {
	switch data.ComponentType {
	case "number_int":
		//数字
		result = CoreFilter.GetInt64ByStringNoErr(data.Val)
	case "number_float":
		//浮点数
		result = CoreFilter.GetFloat64ByStringNoErr(data.Val)
	case "number_price":
		//价格数据
		if isShow {
			result = float64(CoreFilter.GetInt64ByStringNoErr(data.Val)) / 100
		} else {
			result = CoreFilter.GetInt64ByStringNoErr(data.Val)
		}
	case "number_p":
		//百分比数据
		if isShow {
			result = fmt.Sprint(float64(CoreFilter.GetInt64ByStringNoErr(data.Val))/100, "%")
		} else {
			result = CoreFilter.GetInt64ByStringNoErr(data.Val)
		}
	case "bool_open":
		//bool值
		if isShow {
			if CoreFilter.GetBoolByInterfaceNoErr(data.Val) {
				result = CoreLanguage.GetLanguageText(ctx, "bool_true")
			} else {
				result = CoreLanguage.GetLanguageText(ctx, "bool_false")
			}
		} else {
			result = CoreFilter.GetBoolByInterfaceNoErr(data.Val)
		}
	case "file_id":
		//文件URL
		fileID := CoreFilter.GetInt64ByStringNoErr(data.Val)
		if fileID > 0 {
			result = BaseFileSys2.GetPublicURLByClaimID(fileID)
		} else {
			result = ""
		}
	case "file_ids":
		//文件URL
		fileIDs := CoreFilter.GetIDsInString(data.Val, ",")
		if len(fileIDs) > 0 {
			for k := 0; k < len(fileIDs); k++ {
				result = BaseFileSys2.GetPublicURLByClaimID(fileIDs[k])
			}
		} else {
			result = ""
		}
	default:
		result = data.Val
	}
	return
}
