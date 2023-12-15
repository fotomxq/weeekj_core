package BaseFileSys

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsUpdateClaimInfo 修改认领信息参数
type ArgsUpdateClaimInfo struct {
	//更新Hash
	UpdateHash string `db:"update_hash" json:"updateHash" check:"mark"`
	//引用ID
	ClaimID int64 `db:"id" json:"claimID" check:"id"`
	//用户ID
	// 可选，用于检测
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//组织ID
	// 可选，用于检测
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//扩展信息
	ClaimInfos CoreSQLConfig.FieldsConfigsType `db:"infos" json:"claimInfos"`
	//是否公开
	IsPublic bool `db:"is_public" json:"isPublic" check:"bool" empty:"true"`
	//描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
}

// UpdateClaimInfo 修改认领信息
func UpdateClaimInfo(args *ArgsUpdateClaimInfo) (err error) {
	where := "id = :id"
	var newUpdateHash string
	newUpdateHash, err = CoreFilter.GetRandStr3(10)
	if err != nil {
		return
	}
	maps := map[string]interface{}{
		"id":              args.ClaimID,
		"new_update_hash": newUpdateHash,
		"is_public":       args.IsPublic,
		"expire_at":       args.ExpireAt,
		"des":             args.Des,
		"infos":           args.ClaimInfos,
	}
	if args.UpdateHash != "" {
		maps["old_update_hash"] = args.UpdateHash
		where = where + " AND update_hash = :old_update_hash"
	}
	if args.UserID > 0 {
		maps["user_id"] = args.UserID
		where = where + " AND user_id = :user_id"
	}
	if args.OrgID > 0 {
		maps["org_id"] = args.OrgID
		where = where + " AND org_id = :org_id"
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_file_claim SET  update_at = NOW(), update_hash = :new_update_hash, is_public = :is_public, expire_at = :expire_at, des = :des, infos = :infos WHERE "+where, maps)
	if err != nil {
		return
	}
	//清理缓冲
	clearClaimCache(args.ClaimID)
	//反馈
	return
}

// ArgsAddVisit 增加访问次数参数
type ArgsAddVisit struct {
	//文件引用ID
	ClaimID int64 `db:"id" json:"id"`
	//访问用户ID
	// 可以留空
	UserID int64 `db:"user_id" check:"id" empty:"true"`
	//访问IP
	IP string `db:"ip" check:"ip" empty:"true"`
}

// AddVisit 增加访问次数
func AddVisit(args *ArgsAddVisit) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_file_claim SET visit_last_at = NOW(), visit_count = coalesce(visit_count, 0) + 1 WHERE id = :id", args)
	if err != nil {
		return
	}
	var claimData FieldsFileClaimType
	err = Router2SystemConfig.MainDB.Get(&claimData, "SELECT id, file_id FROM core_file_claim WHERE id = $1", args.ClaimID)
	if err != nil {
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_file_claim_visit (claim_id, file_id, user_id, create_ip) VALUES (:claim_id,:file_id,:user_id,:create_ip)", map[string]interface{}{
		"claim_id":  args.ClaimID,
		"file_id":   claimData.FileID,
		"user_id":   args.UserID,
		"create_ip": args.IP,
	})
	if err != nil {
		return
	}
	//清理缓冲
	clearClaimCache(args.ClaimID)
	//反馈
	return
}

// ArgsClaimFile 认领文件参数
type ArgsClaimFile struct {
	//原始文件ID
	FileID int64 `json:"fileID" check:"id"`
	//创建用户
	UserID int64 `json:"userID" check:"id"`
	//创建组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//是否为公开的文件
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//过期时间
	// 留空则不过期
	ExpireAt time.Time `json:"expireAt" check:"time" empty:"true"`
	//扩展信息
	ClaimInfos CoreSQLConfig.FieldsConfigsType `json:"claimInfos"`
	//描述
	Des string `json:"des" check:"des" min:"1" max:"1000" empty:"true"`
}

// ClaimFile 认领文件
func ClaimFile(args *ArgsClaimFile) (data FieldsFileClaimType, errCode string, err error) {
	var fileInfo FieldsFileType
	fileInfo, err = GetFileByID(&ArgsGetFileByID{
		ID:         args.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	})
	if err != nil {
		errCode = "file_not_exist"
		err = errors.New("file is not exist, " + err.Error())
		return
	}
	var newUpdateHash string
	newUpdateHash, err = CoreFilter.GetRandStr3(10)
	if err != nil {
		errCode = "update_hash"
		err = errors.New("rand hash, " + err.Error())
		return
	}
	maps := map[string]interface{}{
		"update_hash": newUpdateHash,
		"user_id":     args.UserID,
		"org_id":      args.OrgID,
		"is_public":   args.IsPublic,
		"infos":       args.ClaimInfos,
		"des":         args.Des,
		"file_id":     fileInfo.ID,
	}
	if args.ExpireAt.Unix() > 0 {
		maps["expire_at"] = args.ExpireAt
	} else {
		maps["expire_at"] = time.Time{}
	}
	if err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_file_claim", "INSERT INTO core_file_claim (update_hash, user_id, org_id, is_public, file_id, expire_at, visit_last_at, visit_count, des, infos) VALUES (:update_hash, :user_id, :org_id, :is_public, :file_id, :expire_at, NOW(), 0, :des, :infos)", maps, &data); err != nil {
		errCode = "insert_claim"
		return
	}
	//清理缓冲
	clearClaimCache(0)
	clearFileCache(args.FileID)
	//反馈
	return
}
