package BaseFileSys

import (
	"time"
)

// Deprecated: 建议采用BaseFileSys2
//FieldsFileClaimVisit 认领文件访问数据参数
// 该表不是永久性数据存储，将用于数据分析调用和使用，一段时间后将自动删除相关数据，以确保数据快速写入
type FieldsFileClaimVisit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//认领文件ID
	ClaimID int64 `db:"claim_id" json:"claimID"`
	//实体文件ID
	FileID int64 `db:"file_id" json:"fileID"`
	//查看用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//创建IP
	CreateIP string `db:"create_ip" json:"createIP"`
}
