package FinanceSafe

import (
	"fmt"
	BaseEarlyWarning "github.com/fotomxq/weeekj_core/v5/base/early_warning"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func runEW() {
	//遍历数据
	for {
		//获取安全事件
		var dataList []FieldsSafeType
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_info, ew_template_mark, message FROM finance_safe WHERE allow_open = true AND need_ew = true AND allow_ew = false LIMIT 10"); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		//发送预警数据
		for _, v := range dataList {
			contents := map[string]string{
				"Content": v.Message,
			}
			switch v.EWTemplateMark {
			case "FinanceSafeEWTemplateByLogHash":
				contents["LogID"] = fmt.Sprint(v.CreateInfo.ID)
			case "FinanceSafeEWTemplateByPayLost":
				contents["LogID"] = fmt.Sprint(v.CreateInfo.ID)
			case "FinanceSafeEWTemplateByPayPrice":
				contents["PayID"] = fmt.Sprint(v.CreateInfo.ID)
			case "FinanceSafeEWTemplateByPayLimit0":
				contents["PayID"] = fmt.Sprint(v.CreateInfo.ID)
			case "FinanceSafeEWTemplateByPayLimitMax":
				contents["PayID"] = fmt.Sprint(v.CreateInfo.ID)
			case "FinanceSafeEWTemplateByPayFrequencyOneFrom":
				contents["PayID"] = fmt.Sprint(v.CreateInfo.ID)
			case "FinanceSafeEWTemplateByPayFrequencyOneTo":
				contents["PayID"] = fmt.Sprint(v.CreateInfo.ID)
			case "FinanceSafeEWTemplateByPayFrequencyAll":
				contents["PayID"] = fmt.Sprint(v.CreateInfo.ID)
			default:
			}
			if err := BaseEarlyWarning.SendMod(&BaseEarlyWarning.ArgsSendMod{
				Mark:     v.EWTemplateMark,
				Contents: contents,
			}); err != nil {
				runLog("send ew mod, ", err)
			}
			if _, err := Router2SystemConfig.MainDB.Exec("UPDATE finance_safe SET allow_ew = true WHERE id = $1", v.ID); err != nil {
				CoreLog.Error("update finance safe id: ", v.ID, ", err: ", err)
				return
			}
		}
	}
}
