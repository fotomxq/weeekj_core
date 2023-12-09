package BaseFileSys

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteClaim 删除文件认领参数
type ArgsDeleteClaim struct {
	//Hash
	// 可选，用于交叉验证
	UpdateHash string
	//引用文件ID
	ClaimID int64
	//用户ID
	// 可选，用于检测
	UserID int64 `db:"user_id" json:"userID"`
	//组织ID
	// 可选，用于检测
	OrgID int64 `db:"org_id" json:"orgID"`
}

// DeleteClaim 删除文件认领
// 如果系统启动了特殊设定，则自动检查是否存在其他引用，否则将删除文件底层数据
func DeleteClaim(args *ArgsDeleteClaim) (err error) {
	//获取领用数据
	var data FieldsFileClaimType
	data, err = GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: args.ClaimID,
		UserID:  args.UserID,
		OrgID:   args.OrgID,
	})
	if err != nil {
		return
	}
	//删除文件认领
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_file_claim", "id = :id AND (update_hash = :update_hash OR :update_hash = '')", map[string]interface{}{
		"update_hash": args.UpdateHash,
		"id":          args.ClaimID,
	})
	if err != nil {
		return
	}
	//删除文件实体
	count := GetFileClaimCount(&ArgsGetFileClaimCount{
		FileID: data.FileID,
	})
	if count < 1 {
		if err = DeleteFile(&ArgsDeleteFile{
			UpdateHash: "",
			FileID:     data.FileID,
			CreateInfo: CoreSQLFrom.FieldsFrom{},
		}); err != nil {
			return err
		}
	}
	//删除缓冲
	clearClaimCache(args.ClaimID)
	//反馈
	return
}

// ArgsDeleteClaimByFileID 删除文件所有引用关系参数
type ArgsDeleteClaimByFileID struct {
	//原始文件ID
	FileID int64 `db:"file_id"`
}

// DeleteClaimByFileID 删除文件所有引用关系
func DeleteClaimByFileID(args *ArgsDeleteClaimByFileID) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_file_claim", "file_id = :file_id", args)
	if err != nil {
		return
	}
	if err = DeleteFile(&ArgsDeleteFile{
		UpdateHash: "",
		FileID:     args.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	}); err != nil {
		return err
	}
	//删除缓冲
	clearClaimCache(0)
	clearFileCache(args.FileID)
	//反馈
	return
}
