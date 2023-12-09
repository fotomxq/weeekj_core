package OrgShareSpace

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type FieldsFile struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//协同人列
	ShareOrgBindIDs FieldsFileShareOrgBindList `db:"share_org_bind_ids" json:"shareOrgBindIDs"`
	//目录ID
	DirID int64 `db:"dir_id" json:"dirID"`
	//名称
	Name string `db:"name" json:"name"`
	//文件系统
	System string `db:"system" json:"system"`
	//文件ID
	FileID int64 `db:"file_id" json:"fileID"`
	//文件尺寸
	FileSize int64 `db:"file_size" json:"fileSize"`
}

type FieldsFileShareOrgBindList []FieldsFileShareOrgBind

// Value sql底层处理器
func (t FieldsFileShareOrgBindList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsFileShareOrgBindList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsFileShareOrgBind struct {
	//所属人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//模式
	// 0 仅查看; 1 查看和编辑
	Mode int `db:"mode" json:"mode"`
}

// Value sql底层处理器
func (t FieldsFileShareOrgBind) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsFileShareOrgBind) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
