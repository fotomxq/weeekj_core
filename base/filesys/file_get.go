package BaseFileSys

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetFileList 获取文件列表参数
type ArgsGetFileList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom `json:"createInfo"`
	//渠道来源
	FromInfo CoreSQLFrom.FieldsFrom `json:"fromInfo"`
	//文件类型
	FileType string `json:"fileType" check:"mark" empty:"true"`
	//文件SHA1
	FileShaSearch string `json:"fileShaSearch" check:"mark" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetFileList 获取文件列表
func GetFileList(args *ArgsGetFileList) (dataList []FieldsFileType, dataCount int64, err error) {
	//获取缓冲
	cacheMark := fmt.Sprint(getFileCacheMark(0), ":GetFileList:", args.Pages.GetCacheMark(), ".", args.CreateInfo.GetString(), ".", args.FromInfo.GetString(), ".", args.FileType, ".", args.FileShaSearch, ".", args.Search)
	type cacheType struct {
		DataList  []FieldsFileType
		DataCount int64
	}
	var cacheData cacheType
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &cacheData); err == nil && len(cacheData.DataList) > 0 {
		dataList = cacheData.DataList
		dataCount = cacheData.DataCount
		return
	}
	//获取数据
	where := "(create_ip ILIKE '%' || :search || '%')"
	maps := map[string]interface{}{
		"search": args.Search,
	}
	where, maps, err = args.CreateInfo.GetListAnd("create_info", "create_info", where, maps)
	where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
	if args.FileType != "" {
		where = where + " AND file_type = :file_type"
		maps["file_type"] = args.FileType
	}
	if args.FileShaSearch != "" {
		where = where + " AND file_hash = :file_hash"
		maps["file_hash"] = args.FileShaSearch
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_file",
		"id",
		"SELECT id, create_at, update_at, update_hash, create_ip, create_info, file_size, file_type, file_hash, file_src, from_info, infos FROM core_file WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "create_ip", "file_size", "file_type"},
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

// ArgsGetFileByID 获取文件信息参数
type ArgsGetFileByID struct {
	//文件ID
	ID int64 `json:"id" check:"id"`
	//创建来源
	// 可选，用于验证
	CreateInfo CoreSQLFrom.FieldsFrom `json:"createInfo"`
}

// Deprecated: 建议采用BaseFileSys2
// GetFileByID 获取文件信息
func GetFileByID(args *ArgsGetFileByID) (data FieldsFileType, err error) {
	//获取缓冲
	cacheMark := fmt.Sprint(getFileCacheMark(args.ID), ":GetFileByID:", args.CreateInfo.GetString())
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	where := "id = :id"
	maps := map[string]interface{}{
		"id": args.ID,
	}
	if args.CreateInfo.System != "" {
		where = where + " AND create_info @> :create_info"
		maps, err = args.CreateInfo.GetMaps("create_info", maps)
		if err != nil {
			return
		}
	}
	err = CoreSQL.GetOne(
		Router2SystemConfig.MainDB.DB,
		&data,
		"SELECT id, create_at, update_at, update_hash, create_ip, create_info, file_size, file_type, file_hash, file_src, from_info, infos FROM core_file WHERE "+where,
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

// ArgsGetFileByIDs 获取一组文件信息参数
type ArgsGetFileByIDs struct {
	//一组ID
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//来源，用于验证
	CreateInfo CoreSQLFrom.FieldsFrom `json:"createInfo"`
}

// Deprecated: 建议采用BaseFileSys2
// GetFileByIDs 获取一组文件信息
func GetFileByIDs(args *ArgsGetFileByIDs) (dataList []FieldsFileType, dataCount int64, err error) {
	//获取缓冲
	cacheMark := fmt.Sprint(getFileCacheMark(0), ":GetFileByIDs:", args.IDs, ".", args.CreateInfo.GetString())
	type cacheType struct {
		DataList  []FieldsFileType
		DataCount int64
	}
	var cacheData cacheType
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &cacheData); err == nil && len(cacheData.DataList) > 0 {
		dataList = cacheData.DataList
		dataCount = cacheData.DataCount
		return
	}
	//获取数据
	where := "id = ANY(:ids)"
	maps := map[string]interface{}{
		"ids": args.IDs,
	}
	if args.CreateInfo.System != "" {
		where = where + " AND create_info @> :create_info"
		maps, err = args.CreateInfo.GetMaps("create_info", maps)
		if err != nil {
			return
		}
	}
	dataCount, err = CoreSQL.GetListAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_file",
		"id",
		"SELECT id, create_at, update_at, update_hash, create_ip, create_info, file_size, file_type, file_hash, file_src, from_info, infos FROM core_file WHERE "+where,
		where,
		maps,
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

// ArgsGetFileByIDsAndClaim 获取一组文件信息并检查权限
type ArgsGetFileByIDsAndClaim struct {
	//一组ID
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//来源，用于验证
	CreateInfo CoreSQLFrom.FieldsFrom `json:"createInfo"`
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
// GetFileByIDsAndClaim 获取一组文件信息并检查权限
func GetFileByIDsAndClaim(args *ArgsGetFileByIDsAndClaim) (dataList []FieldsFileType, dataCount int64, err error) {
	//获取缓冲
	cacheMark := fmt.Sprint(getFileCacheMark(0), ":GetFileByIDsAndClaim:", args.IDs, ".", args.CreateInfo.GetString(), ".", args.UserID, ".", args.OrgID, ".", args.IsPublic)
	type cacheType struct {
		DataList  []FieldsFileType
		DataCount int64
	}
	var cacheData cacheType
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &cacheData); err == nil && len(cacheData.DataList) > 0 {
		dataList = cacheData.DataList
		dataCount = cacheData.DataCount
		return
	}
	//获取文件底层数据包
	dataList, dataCount, err = GetFileByIDs(&ArgsGetFileByIDs{
		IDs:        args.IDs,
		CreateInfo: args.CreateInfo,
	})
	if err != nil || dataCount < 1 {
		err = errors.New(fmt.Sprint("get file by ids, ", err))
		return
	}
	//检查所属权
	if args.UserID > -1 || args.OrgID > -1 || args.IsPublic {
		var claimList []FieldsFileClaimType
		err = Router2SystemConfig.MainDB.Select(&claimList, "SELECT id, file_id FROM core_file_claim WHERE (expire_at >= NOW() OR expire_at < to_timestamp(1000000)) AND ($1 < 0 OR user_id = $1) AND ($2 < 0 OR org_id = $2) AND ($3 = false OR is_public = $3 AND file_id = ANY($4))", args.UserID, args.OrgID, args.IsPublic, args.IDs)
		if err != nil {
			err = errors.New(fmt.Sprint("cannot get claim data, ", err))
			return
		}
		if len(claimList) < 1 {
			err = errors.New("cannot get claim data")
			return
		}
		for _, v := range dataList {
			isFind := false
			for _, v2 := range claimList {
				if v.ID == v2.FileID {
					isFind = true
					break
				}
			}
			if !isFind {
				err = errors.New(fmt.Sprint("have file not this user or org, file id: ", v.ID, ", user id: ", args.UserID, ", org id: ", args.OrgID, ", is public: ", args.IsPublic))
				return
			}
		}
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
