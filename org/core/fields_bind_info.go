package OrgCoreCore

import (
	"github.com/lib/pq"
	"time"
)

// FieldsBindInfo 组织成员信息结构
type FieldsBindInfo struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//入职时间
	InAt time.Time `db:"in_at" json:"inAt"`
	//离职时间
	OutAt time.Time `db:"out_at" json:"outAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//性别
	// 0 男 1 女 2 未知
	Gender int `db:"gender" json:"gender"`
	//出生年月
	// date结构
	DateOfBirth time.Time `db:"date_of_birth" json:"dateOfBirth"`
	//联系电话
	Phone string `db:"phone" json:"phone"`
	//证件编号
	IDCard string `db:"id_card" json:"idCard"`
	//个人照片
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//其他照片信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//教育状态
	// 0 无教育; 1 小学; 2 初中; 3 高中; 4 专科技校; 5 本科; 6 研究生; 7 博士; 8 博士后
	EducationStatus int `db:"education_status" json:"educationStatus"`
	//婚姻状态
	MaritalStatus bool `db:"marital_status" json:"maritalStatus"`
	//家庭住址
	Address string `db:"address" json:"address"`
	//资质
	CertName string `db:"cert_name" json:"certName"`
	//保险
	InsuranceType string    `db:"insurance_type" json:"insuranceType"`
	InsuranceAt   time.Time `db:"insurance_at" json:"insuranceAt"`
	//银行信息
	BandName string `db:"band_name" json:"bindName"`
	BandSN   string `db:"band_sn" json:"bandSN"`
	//描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
}
