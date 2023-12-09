package BaseFileSys

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateFileInfo 修改文件信息参数
type ArgsUpdateFileInfo struct {
	//Hash
	UpdateHash string `db:"update_hash" json:"updateHash" check:"mark"`
	//文件ID
	FileID int64 `db:"id" json:"fileID" check:"id"`
	//扩展信息
	Infos CoreSQLConfig.FieldsConfigsType `db:"infos" json:"infos"`
}

// Deprecated: 建议采用BaseFileSys2
// UpdateFileInfo 修改文件信息
func UpdateFileInfo(args *ArgsUpdateFileInfo) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint("recover, ", e))
			return
		}
	}()
	var newUpdateHash string
	newUpdateHash, err = CoreFilter.GetRandStr3(10)
	if err != nil {
		return errors.New("rand hash, " + err.Error())
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_file SET  update_at = NOW(), update_hash = :update_hash, infos = :infos WHERE id = :id AND update_hash = :old_update_hash", map[string]interface{}{
		"old_update_hash": args.UpdateHash,
		"update_hash":     newUpdateHash,
		"id":              args.FileID,
		"infos":           args.Infos,
	})
	if err != nil {
		return
	}
	//清理缓冲
	clearFileCache(args.FileID)
	//反馈
	return
}
