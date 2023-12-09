package OrgShareSpace

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetFileList 获取文件列表参数
type ArgsGetFileList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//目录ID
	DirID int64 `db:"dir_id" json:"dirID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetFileList 获取文件列表
func GetFileList(args *ArgsGetFileList) (dataList []FieldsFile, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.OrgBindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(org_bind_id = :org_bind_id OR org_bind_id = 0)"
		maps["org_bind_id"] = args.OrgBindID
	}
	if args.DirID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "dir_id = :dir_id"
		maps["dir_id"] = args.DirID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "org_share_space_file"
	var rawList []FieldsFile
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getFileByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetFileByID 查看文件信息参数
type ArgsGetFileByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
}

// GetFileByID 查看文件信息
func GetFileByID(args *ArgsGetFileByID) (data FieldsFile, isEdit bool, err error) {
	//检查存在和权限
	data = getFileByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	if data.OrgBindID > 0 {
		if !CoreFilter.EqID2(args.OrgBindID, data.OrgBindID) {
			isFind := false
			for _, v := range data.ShareOrgBindIDs {
				if v.OrgBindID == args.OrgBindID {
					if v.Mode == 1 {
						isEdit = true
					}
					isFind = true
					break
				}
			}
			if !isFind {
				err = errors.New("no data")
				return
			}
		} else {
			isEdit = true
		}
	}
	//反馈
	return
}

// GetFileCountByDir 获取目录下有多少文件
func GetFileCountByDir(dirID int64) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM org_share_space_file WHERE dir_id = $1", dirID)
	return
}

// ArgsCreateFile 创建文件参数
type ArgsCreateFile struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//协同人列
	ShareOrgBindIDs FieldsFileShareOrgBindList `db:"share_org_bind_ids" json:"shareOrgBindIDs"`
	//目录ID
	DirID int64 `db:"dir_id" json:"dirID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//文件系统
	System string `db:"system" json:"system" check:"mark"`
	//文件ID
	FileID int64 `db:"file_id" json:"fileID" check:"id"`
	//文件尺寸
	FileSize int64 `db:"file_size" json:"fileSize" check:"int64Than0"`
}

// CreateFile 创建文件
func CreateFile(args *ArgsCreateFile) (data FieldsFile, err error) {
	//修正参数
	if len(args.ShareOrgBindIDs) < 1 {
		args.ShareOrgBindIDs = FieldsFileShareOrgBindList{}
	}
	//检查目录权限
	if args.DirID > 0 {
		dirData := getDirByID(args.DirID)
		if args.OrgID != dirData.OrgID || (args.OrgBindID != dirData.OrgBindID || dirData.OrgBindID == 0) {
			err = errors.New("no data")
			return
		}
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_share_space_file", "INSERT INTO org_share_space_file (org_id, org_bind_id, share_org_bind_ids, dir_id, name, system, file_id, file_size) VALUES (:org_id,:org_bind_id,:share_org_bind_ids,:dir_id,:name,:system,:file_id,:file_size)", args, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateFile 修改文件信息参数
type ArgsUpdateFile struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//目录ID
	DirID int64 `db:"dir_id" json:"dirID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
}

// UpdateFile 修改文件信息
func UpdateFile(args *ArgsUpdateFile) (err error) {
	//检查存在和权限
	data := getFileByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	if data.OrgBindID > 0 {
		if !CoreFilter.EqID2(args.OrgBindID, data.OrgBindID) {
			err = errors.New("no data")
			return
		}
	}
	//检查目录权限
	if args.DirID > 0 {
		dirData := getDirByID(args.DirID)
		if args.OrgID != dirData.OrgID || (args.OrgBindID != dirData.OrgBindID || dirData.OrgBindID == 0) {
			err = errors.New("no data")
			return
		}
	}
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_share_space_file SET dir_id = :dir_id, name = :name WHERE id = :id", map[string]interface{}{
		"id":     data.ID,
		"dir_id": args.DirID,
		"name":   args.Name,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteFileCache(data.ID)
	//反馈
	return
}

// ArgsMoveFile 批量转移文件参数
type ArgsMoveFile struct {
	//ID
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//目录ID
	DirID int64 `db:"dir_id" json:"dirID" check:"id"`
}

// MoveFile 批量转移文件
func MoveFile(args *ArgsMoveFile) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_share_space_file SET dir_id = :dir_id WHERE id = ANY(:ids) AND org_id = :org_id AND org_bind_id = :org_bind_id", map[string]interface{}{
		"ids":         args.IDs,
		"dir_id":      args.DirID,
		"org_id":      args.OrgID,
		"org_bind_id": args.OrgBindID,
	})
	if err != nil {
		return
	}
	//删除缓冲
	for _, v := range args.IDs {
		deleteFileCache(v)
	}
	//反馈
	return
}

// ArgsDeleteFile 删除文件参数
type ArgsDeleteFile struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
}

// DeleteFile 删除文件
func DeleteFile(args *ArgsDeleteFile) (err error) {
	//检查存在和权限
	data := getFileByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	if data.OrgBindID > 0 {
		if !CoreFilter.EqID2(args.OrgBindID, data.OrgBindID) {
			err = errors.New("no data")
			return
		}
	}
	//执行删除
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "org_share_space_file", "id", map[string]interface{}{
		"id": data.ID,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteFileCache(data.ID)
	//发出通知
	CoreNats.PushDataNoErr("/org/share_space/file", "delete", data.FileID, data.System, nil)
	//反馈
	return
}

// ArgsDeleteFiles 批量删除文件参数
type ArgsDeleteFiles struct {
	//ID
	IDs []int64 `db:"ids" json:"ids" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
}

// DeleteFiles 批量删除文件
func DeleteFiles(args *ArgsDeleteFiles) (err error) {
	for _, v := range args.IDs {
		err = DeleteFile(&ArgsDeleteFile{
			ID:        v,
			OrgID:     args.OrgID,
			OrgBindID: args.OrgID,
		})
		if err != nil {
			return
		}
	}
	return
}

// 变更文件大小
func updateFileSize(system string, fileID int64, fileSize int64) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_share_space_file SET file_size = :file_size WHERE system = :system AND file_id = :file_id", map[string]interface{}{
		"file_size": fileSize,
		"system":    system,
		"file_id":   fileID,
	})
	if err != nil {
		return
	}
	//获取文件ID
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_share_space_file WHERE system = $1 AND file_id = $2 LIMIT 1", system, fileID)
	if err != nil {
		return
	}
	//删除缓冲
	deleteFileCache(id)
	//反馈
	return
}

// 获取文件数据
func getFileByID(id int64) (data FieldsFile) {
	cacheMark := getFileCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, share_org_bind_ids, dir_id, name, system, file_id, file_size FROM org_share_space_file WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getFileCacheMark(id int64) string {
	return fmt.Sprint("org:share:space:core:file:id:", id)
}

func deleteFileCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getFileCacheMark(id))
}
