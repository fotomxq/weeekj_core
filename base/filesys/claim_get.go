package BaseFileSys

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsGetFileClaimList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//指定文件ID
	FileID int64 `json:"fileID" check:"id" empty:"true"`
	//是否公开
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//搜索内容
	Search string `json:"search" check:"search" empty:"true"`
}

func GetFileClaimList(args *ArgsGetFileClaimList) (dataList []FieldsFileClaimType, dataCount int64, err error) {
	//获取缓冲
	cacheMark := fmt.Sprint(getClaimCacheMark(0), ":GetFileClaimList:", args.Pages.GetCacheMark(), ".", args.UserID, ".", args.OrgID, ".", args.FileID, ".", args.IsPublic, ".", args.Search)
	type cacheType struct {
		DataList  []FieldsFileClaimType
		DataCount int64
	}
	var cacheData cacheType
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &cacheData); err == nil && len(cacheData.DataList) > 0 {
		dataList = cacheData.DataList
		dataCount = cacheData.DataCount
		return
	}
	//获取数据
	where := "(des ILIKE '%' || :search || '%') AND is_public = :is_public"
	maps := map[string]interface{}{
		"search":    args.Search,
		"is_public": args.IsPublic,
	}
	if args.IsPublic {
		if args.UserID > 0 {
			where = where + " AND user_id = :user_id"
			maps["user_id"] = args.UserID
		}
		if args.OrgID > 0 {
			where = where + " AND org_id = :org_id"
			maps["org_id"] = args.OrgID
		}
	} else {
		where = where + " AND user_id = :user_id AND org_id = :org_id"
		maps["user_id"] = args.UserID
		maps["org_id"] = args.OrgID
	}
	if args.FileID > 0 {
		where = where + " AND file_id = :file_id"
		maps["file_id"] = args.FileID
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_file_claim",
		"id",
		"SELECT id, create_at, update_at, update_hash, user_id, org_id, is_public, file_id, expire_at, visit_last_at, visit_count, des, infos FROM core_file_claim WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "expire_at", "visit_last_at", "visit_count"},
	)
	if err != nil {
		return
	}
	//写入缓冲
	cacheData = cacheType{
		DataList:  dataList,
		DataCount: dataCount,
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, cacheData, cacheTime)
	//反馈
	return
}

// ArgsGetFileClaimByID 获取认领信息参数
type ArgsGetFileClaimByID struct {
	//引用文件ID
	ClaimID int64
	//用户ID
	// 可选，用于检测
	UserID int64 `json:"userID"`
	//组织ID
	// 可选，用于检测
	OrgID int64 `json:"orgID"`
}

// GetFileClaimByID 获取认领信息
func GetFileClaimByID(args *ArgsGetFileClaimByID) (data FieldsFileClaimType, err error) {
	//获取缓冲
	cacheMark := fmt.Sprint(getClaimCacheMark(args.ClaimID), ":GetFileClaimByID:", args.UserID, args.OrgID)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	//获取数据
	where := "id = :id"
	maps := map[string]interface{}{
		"id": args.ClaimID,
	}
	if args.UserID > 0 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	err = CoreSQL.GetOne(
		Router2SystemConfig.MainDB.DB,
		&data,
		"SELECT id, create_at, update_at, update_hash, user_id, org_id, is_public, file_id, expire_at, visit_last_at, visit_count, des, infos FROM core_file_claim WHERE "+where,
		maps,
	)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("data is empty")
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTime)
	//反馈
	return
}

// ArgsGetFileClaimByIDList 获取一组认领文件的数据参数
type ArgsGetFileClaimByIDList struct {
	//引用ID列
	ClaimIDList []int64 `json:"claimIDList"`
	//用户ID
	// 可选，用于检测
	UserID int64 `json:"userID"`
	//组织ID
	// 可选，用于检测
	OrgID int64 `json:"orgID"`
	//是否仅公开数据
	IsPublic bool `json:"isPublic"`
}

// Deprecated: 建议采用BaseFileSys2
// GetFileClaimByIDList 获取一组认领文件的数据
func GetFileClaimByIDList(args *ArgsGetFileClaimByIDList) (dataList []FieldsFileClaimType, err error) {
	for _, v := range args.ClaimIDList {
		var data FieldsFileClaimType
		if args.IsPublic {
			data, err = GetFileClaimByID(&ArgsGetFileClaimByID{
				ClaimID: v,
				UserID:  0,
				OrgID:   0,
			})
		} else {
			if args.UserID < 1 {
				err = errors.New("data not exist")
				return
			}
			data, err = GetFileClaimByID(&ArgsGetFileClaimByID{
				ClaimID: v,
				UserID:  args.UserID,
				OrgID:   args.OrgID,
			})
		}
		if err != nil {
			continue
		}
		dataList = append(dataList, data)
	}
	if len(dataList) < 1 {
		err = errors.New("data less 1")
	} else {
		err = nil
	}
	return
}

// ArgsGetFileByClaimIDs 通过认领ID组，获取文件数据参数
type ArgsGetFileByClaimIDs struct {
	//引用ID列
	ClaimIDList []int64 `json:"claimIDList"`
	//用户ID
	// 可选，用于检测
	UserID int64 `json:"userID"`
	//组织ID
	// 可选，用于检测
	OrgID int64 `json:"orgID"`
	//是否仅公开数据
	IsPublic bool `json:"isPublic"`
}

// Deprecated: 建议采用BaseFileSys2
// GetFileByClaimIDs 通过认领ID组，获取文件数据
func GetFileByClaimIDs(args *ArgsGetFileByClaimIDs) (fileClaimList []FieldsFileClaimType, fileList []FieldsFileType, dataCount int64, err error) {
	var finishList []int64
	for k := 0; k < len(args.ClaimIDList); k++ {
		v := args.ClaimIDList[k]
		//跳过重复数据
		isFind := false
		for k2 := 0; k2 < len(finishList); k2++ {
			v2 := finishList[k2]
			if v == v2 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		finishList = append(finishList, v)
		//获取引用数据
		var data FieldsFileClaimType
		if args.IsPublic {
			data, err = GetFileClaimByID(&ArgsGetFileClaimByID{
				ClaimID: v,
				UserID:  0,
				OrgID:   0,
			})
		} else {
			data, err = GetFileClaimByID(&ArgsGetFileClaimByID{
				ClaimID: v,
				UserID:  args.UserID,
				OrgID:   args.OrgID,
			})
			if !data.IsPublic {
				continue
			}
		}
		if err != nil {
			continue
		}
		//获取文件数据
		var vFile FieldsFileType
		vFile, err = GetFileByID(&ArgsGetFileByID{
			ID:         data.FileID,
			CreateInfo: CoreSQLFrom.FieldsFrom{},
		})
		if err != nil {
			continue
		}
		fileClaimList = append(fileClaimList, data)
		fileList = append(fileList, vFile)
	}
	if len(fileList) < 1 {
		err = errors.New("data less 1")
	} else {
		err = nil
		dataCount = int64(len(fileList))
	}
	return
}

// ArgsGetFileClaimCount 获取文件引用数量参数
type ArgsGetFileClaimCount struct {
	//文件ID
	FileID int64 `db:"file_id" check:"id"`
}

// GetFileClaimCount 获取文件引用数量
func GetFileClaimCount(args *ArgsGetFileClaimCount) (count int64) {
	var err error
	//获取缓冲
	cacheMark := fmt.Sprint(getClaimCacheMark(0), ":GetFileClaimCount:", args.FileID)
	count, err = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if err == nil && count > 0 {
		return
	}
	//获取数据
	count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "core_file_claim", "id", "file_id = $1", args.FileID)
	if err != nil {
		count = 0
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetInt64(cacheMark, count, cacheTime)
	//反馈
	return
}
