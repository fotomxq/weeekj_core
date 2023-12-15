package ERPAudit

import (
	ERPCore "github.com/fotomxq/weeekj_core/v5/erp/core"
)

// GetStepChildAllVal 获取节点内容
func GetStepChildAllVal(stepChildID int64) (dataList []ERPCore.FieldsComponentVal) {
	return componentValObj.GetAllVal(stepChildID)
}
