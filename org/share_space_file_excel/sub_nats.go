package OrgShareSpaceFileExcel

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//接收文件系统变更
	CoreNats.SubDataByteNoErr("org_share_space_file", "/org/share_space/file", subNatsFileUpdate)
}

// 接收文件系统变更
func subNatsFileUpdate(_ *nats.Msg, action string, fileID int64, system string, _ []byte) {
	if system != "excel" {
		return
	}
	switch action {
	case "delete":
		//删除文件
		if err := deleteDoc(fileID); err != nil {
			CoreLog.Error("org share space file excel, sub nats file update, delete file, ", err)
		}
	}
}
