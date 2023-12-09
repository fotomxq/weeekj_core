package OrgShareSpace

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetDirList 获取目录列表参数
type ArgsGetDirList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetDirList 获取目录列表
func GetDirList(args *ArgsGetDirList) (dataList []FieldsDir, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.OrgBindID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(org_bind_id = :org_bind_id)"
		maps["org_bind_id"] = args.OrgBindID
	} else {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(org_bind_id = 0)"
		maps["org_bind_id"] = args.OrgBindID
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
	tableName := "org_share_space_dir"
	var rawList []FieldsDir
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
		vData := getDirByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsCreateDir 创建新目录参数
type ArgsCreateDir struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//上级目录
	ParentID int64 `db:"parent_id" json:"parentID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
}

// CreateDir 创建新目录
func CreateDir(args *ArgsCreateDir) (data FieldsDir, err error) {
	//检查上级目录，归属权必须固定
	if args.ParentID > 0 {
		parentData := getDirByID(args.ParentID)
		if parentData.ID < 1 || !CoreFilter.EqID2(args.OrgID, parentData.OrgID) {
			err = errors.New("no parent data")
			return
		}
		if parentData.OrgBindID > 0 {
			if !CoreFilter.EqID2(args.OrgBindID, parentData.OrgBindID) {
				err = errors.New("no parent data")
				return
			}
		} else {
			if args.OrgBindID != parentData.OrgBindID {
				err = errors.New("no parent data")
				return
			}
		}
	}
	//构建目录
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_share_space_dir", "INSERT INTO org_share_space_dir (org_id, org_bind_id, parent_id, name) VALUES (:org_id,:org_bind_id,:parent_id,:name)", args, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateDir 修改目录参数
type ArgsUpdateDir struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//上级目录
	ParentID int64 `db:"parent_id" json:"parentID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
}

// UpdateDir 修改目录
func UpdateDir(args *ArgsUpdateDir) (err error) {
	//检查存在和权限
	data := getDirByID(args.ID)
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
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_share_space_dir SET parent_id = :parent_id, name = :name WHERE id = :id", map[string]interface{}{
		"id": data.ID,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteDirCache(data.ID)
	//反馈
	return
}

// ArgsDeleteDir 删除目录参数
type ArgsDeleteDir struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
}

// DeleteDir 删除目录
func DeleteDir(args *ArgsDeleteDir) (err error) {
	//检查存在和权限
	data := getDirByID(args.ID)
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
	//检查该目录是否存在文件？
	fileCount := GetFileCountByDir(data.ID)
	if fileCount > 0 {
		err = errors.New("dir have file")
		return
	}
	//执行删除
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "org_share_space_dir", "id", map[string]interface{}{
		"id": data.ID,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteDirCache(data.ID)
	//反馈
	return
}

func getDirByID(id int64) (data FieldsDir) {
	cacheMark := getDirCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, parent_id, name FROM org_share_space_dir WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getDirCacheMark(id int64) string {
	return fmt.Sprint("org:share:space:core:dir:id:", id)
}

func deleteDirCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getDirCacheMark(id))
}
