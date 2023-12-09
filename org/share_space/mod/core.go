package OrgShareSpaceMod

import CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"

// UpdateFileSize 请求更新文件的尺寸
func UpdateFileSize(system string, fileID int64, fileSize int64) {
	CoreNats.PushDataNoErr("/org/share_space/file", "size", fileID, system, map[string]interface{}{
		"size": fileSize,
	})
}
