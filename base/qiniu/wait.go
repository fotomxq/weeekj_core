package BaseQiniu

import (
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsCreateWait 创建新的token参数
type ArgsCreateWait struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//是否公开
	IsPublic bool `db:"is_public" json:"isPublic"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//扩展参数
	ClaimInfos CoreSQLConfig.FieldsConfigsType `db:"claim_infos" json:"claimInfos"`
	//描述
	Des string `db:"des" json:"des"`
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//IP
	IP string `db:"ip" json:"ip"`
}

// CreateWait 创建新的token
func CreateWait(args *ArgsCreateWait) (waitID int64, err error) {
	//获取配置
	var fileQiNiuTokenMaxCount int
	fileQiNiuTokenMaxCount, err = BaseConfig.GetDataInt("FileQiNiuTokenMaxCount")
	if err != nil {
		fileQiNiuTokenMaxCount = 5
		err = nil
	}
	if fileQiNiuTokenMaxCount > 0 {
		type dataType struct {
			Count int `db:"count" json:"count"`
		}
		var countData dataType
		err = Router2SystemConfig.MainDB.Get(&countData, "SELECT count(id) as count FROM core_file_qiniu_wait WHERE user_id = $1", args.UserID)
		if err == nil && countData.Count >= fileQiNiuTokenMaxCount {
			//将最早的数据删除
			type dataIDType struct {
				ID int64 `db:"id" json:"id"`
			}
			var dataID dataIDType
			err = Router2SystemConfig.MainDB.Get(&dataID, "SELECT id FROM core_file_qiniu_wait WHERE user_id = $1 ORDER BY id LIMIT 1", args.UserID)
			if err == nil {
				_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_file_qiniu_wait", "id", dataID)
				if err != nil {
					return
				}
			}
		}
	}
	//创建新的数据
	waitID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_file_qiniu_wait (user_id, org_id, is_public, expire_at, claim_infos, des, create_info, ip, file_claim_id) VALUES (:user_id, :org_id, :is_public, :expire_at, :claim_infos, :des, :create_info, :ip, 0)", args)
	return
}

// GetWaitByID 获取指定的token
func GetWaitByID(id int64) (data FieldsWait, err error) {
	//获取数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, user_id, org_id, is_public, expire_at, claim_infos, des, create_info, ip FROM core_file_qiniu_wait WHERE id = $1 AND file_claim_id < 1", id)
	if err != nil || data.ID < 1 {
		err = errors.New(fmt.Sprint("no data, ", err))
		return
	}
	return
}

// ArgsUpdateWaitID 更新文件ID参数
type ArgsUpdateWaitID struct {
	//wait ID
	ID int64 `db:"id" json:"id" check:"id"`
	//新的文件ID
	FileClaimID int64 `db:"file_claim_id" json:"fileClaimID" check:"id"`
}

// UpdateWaitID 更新文件ID
func UpdateWaitID(args *ArgsUpdateWaitID) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_file_qiniu_wait SET file_claim_id = :file_claim_id WHERE id = :id", args)
	return
}

// ArgsGetLastWaitByID 最终提取数据参数
type ArgsGetLastWaitByID struct {
	//wait ID
	ID int64 `db:"id" json:"id" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetLastWaitByID 最终提取数据
func GetLastWaitByID(args *ArgsGetLastWaitByID) (fileClaimID int64, err error) {
	//获取数据
	var data FieldsWait
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, file_claim_id FROM core_file_qiniu_wait WHERE id = $1 AND file_claim_id > 0 AND user_id = $2", args.ID, args.UserID)
	if err != nil {
		err = errors.New("get wait failed, " + err.Error() + fmt.Sprint(", wait id: ", args.ID, ", user id: ", args.UserID))
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_file_qiniu_wait", "id", map[string]interface{}{
		"id": args.ID,
	})
	if err != nil {
		err = errors.New("delete wait data, " + err.Error())
		return
	}
	fileClaimID = data.FileClaimID
	return
}
