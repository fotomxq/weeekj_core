package ServiceUserInfo

import (
	"fmt"
	OrgWorkTipMod "github.com/fotomxq/weeekj_core/v5/org/work_tip/mod"
)

// CheckInfoArg 检查档案完整性，并推送工作消息
func CheckInfoArg(infoID int64, orgBindID int64) {
	data := getInfoID(infoID)
	if data.ID < 1 {
		return
	}
	if data.Address == "" || data.Name == "" || data.Phone == "" || data.DateOfBirth.Unix() < 1000000 || data.CoverFileID < 1 || data.EmergencyContactPhone == "" || data.EmergencyContact == "" || data.IDCard == "" {
		OrgWorkTipMod.AppendTip(&OrgWorkTipMod.ArgsAppendTip{
			OrgID:     data.OrgID,
			OrgBindID: orgBindID,
			Msg:       fmt.Sprint("信息档案存在多个关键信息缺失，请尽快核对并完善。"),
			System:    "service_user_info",
			BindID:    data.ID,
		})
	}
}
