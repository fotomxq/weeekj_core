package OrgShareSpaceMod

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// UpdateFileSize 请求更新文件的尺寸
func UpdateFileSize(system string, fileID int64, fileSize int64) {
	CoreNats.PushDataNoErr("org_share_space_file", "/org/share_space/file", "size", fileID, system, map[string]interface{}{
		"size": fileSize,
	})
}
