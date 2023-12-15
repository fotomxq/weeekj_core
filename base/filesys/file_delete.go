package BaseFileSys

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteFile 删除文件参数
type ArgsDeleteFile struct {
	//Hash
	// 可选，用于交叉验证
	UpdateHash string `db:"update_hash"`
	//文件ID
	FileID int64 `db:"id"`
	//创建来源
	// 可选，用于验证
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info"`
}

// Deprecated: 建议采用BaseFileSys2
// DeleteFile 删除文件
// 可用于强制删除违法文件的方案
// 会直接删除文件数据，但不会删除文件引用，方便查询关联的用户等信息，以及做提示用
func DeleteFile(args *ArgsDeleteFile) (err error) {
	where := "id = :id"
	maps := map[string]interface{}{
		"id": args.FileID,
	}
	if args.UpdateHash != "" {
		maps["update_hash"] = args.UpdateHash
		where = where + " AND update_hash = :update_hash"
	}
	if args.CreateInfo.System != "" {
		where = where + " AND create_info @> :create_info"
		maps, err = args.CreateInfo.GetMaps("create_info", maps)
		if err != nil {
			return
		}
	}
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_file", where, maps)
	if err != nil {
		return
	}
	//写入缓冲
	clearFileCache(args.FileID)
	clearClaimCache(0)
	//反馈
	return
}
