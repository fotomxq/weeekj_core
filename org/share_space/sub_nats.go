package OrgShareSpace

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
	"time"
)

func subNats() {
	//变更文件大小
	CoreNats.SubDataByteNoErr("/org/share_space/file", subNatsFileUpdate)
}

func subNatsFileUpdate(_ *nats.Msg, action string, fileID int64, system string, data []byte) {
	logAppend := "org share space sub nats file update, "
	switch action {
	case "size":
		//等待5秒再修改，因为创建文档在核心内容构建之后发生，可避免文件不存在去修改的问题发生
		time.Sleep(time.Second * 5)
		//获取参数
		fileSize := gjson.GetBytes(data, "size").Int()
		//修改文件尺寸
		if err := updateFileSize(system, fileID, fileSize); err != nil {
			CoreLog.Error(logAppend, "update file size, ", err)
		}
	}
}
