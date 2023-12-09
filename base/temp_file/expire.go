package BaseTempFile

import (
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
)

func subNatsExpireID(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	if action != "core_temp_file" {
		return
	}
	data := getFileID(id)
	if data.ID < 1 {
		return
	}
	err := CoreFile.DeleteF(data.FileSrc)
	if err != nil {
		CoreLog.Error("core temp file sub nats expire id, delete file: ", err)
	}
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_temp_file", "id", map[string]interface{}{
		"id": data.ID,
	})
	if err != nil {
		CoreLog.Error("core temp file sub nats expire id, delete file data: ", err)
	}
}
