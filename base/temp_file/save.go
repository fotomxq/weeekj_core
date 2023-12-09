package BaseTempFile

import (
	"errors"
	"fmt"
	BaseExpireTip "gitee.com/weeekj/weeekj_core/v5/base/expire_tip"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// CreateTempFileSrc 创建临时文件
// Deprecated
func CreateTempFileSrc(dir string, fileName string, fileType string) (src string, name string, err error) {
	tempDir := Router2SystemConfig.RootDir + CoreFile.Sep + "temp" + CoreFile.Sep + dir
	if err = CoreFile.CreateFolder(tempDir); err != nil {
		err = errors.New(fmt.Sprint("create temp dir, ", err))
		return
	}
	if fileName == "" {
		fileName = CoreFilter.GetRandStr4(50)
	}
	name = CoreFilter.GetNowTimeCarbon().Format("2006010215") + "_" + fileName + fileType
	src = tempDir + CoreFile.Sep + name
	return
}

// SaveFileBefore 检查预加载文件
func SaveFileBefore(fileParams string) (newID int64, hash string, b bool) {
	var data FieldsFile
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, file_params, file_src, file_sha1, name, file_type FROM core_temp_file WHERE file_params = $1", fileParams)
	if data.ID < 1 {
		return
	}
	newID = data.ID
	hash = data.FileSHA1
	b = true
	return
}

// SaveFile 保存文件
func SaveFile(expireSec int, fileParams string, fileName string, fileSha1 string, fileType string) (fileSrc string, newID int64, hash string, err error) {
	//生成临时目录
	tempDir := Router2SystemConfig.RootDir + CoreFile.Sep + "temp"
	if !CoreFile.IsFolder(tempDir) {
		if err = CoreFile.CreateFolder(tempDir); err != nil {
			err = errors.New(fmt.Sprint("create temp dir, ", err))
			return
		}
	}
	if fileSha1 == "" {
		fileSha1 = CoreFilter.GetRandStr4(50)
	}
	//构建文件路径
	fileSrc = tempDir + CoreFile.Sep + CoreFilter.GetNowTimeCarbon().Format("2006010215") + "_" + fileSha1 + "." + fileType
	//过期时间
	if expireSec < 10 {
		expireSec = 120
	}
	expireAt := CoreFilter.GetNowTimeCarbon().AddSeconds(expireSec)
	//构建记录
	newID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_temp_file (expire_at, file_params, file_src, file_sha1, name, file_type) VALUES (:expire_at,:file_params,:file_src,:file_sha1,:name,:file_type)", map[string]interface{}{
		"expire_at":   expireAt,
		"file_params": fileParams,
		"file_src":    fileSrc,
		"file_sha1":   fileSha1,
		"name":        fileName,
		"file_type":   fileType,
	})
	if err != nil {
		return
	}
	hash = fileSha1
	BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
		OrgID:      0,
		UserID:     0,
		SystemMark: "core_temp_file",
		BindID:     newID,
		Hash:       "",
		ExpireAt:   expireAt.Time,
	})
	return
}

func getFileID(id int64) (data FieldsFile) {
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, file_params, file_src, file_sha1, name, file_type FROM core_temp_file WHERE id = $1", id)
	return
}
