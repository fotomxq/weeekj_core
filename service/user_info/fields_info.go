package ServiceUserInfo

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsInfo 信息主表
type FieldsInfo struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//死亡标记
	DieAt time.Time `db:"die_at" json:"dieAt"`
	//出院标记
	OutAt time.Time `db:"out_at" json:"outAt"`
	//组织ID
	// 允许平台方的0数据，该数据可能来源于其他领域
	OrgID int64 `db:"org_id" json:"orgID"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
	UserID int64 `db:"user_id" json:"userID"`
	//从属关系
	BindID int64 `db:"bind_id" json:"bindID"`
	//从属关系的类型
	// 0 子女; 1 亲属非子女; 2 好友; 3 其他
	BindType int `db:"bind_type" json:"bindType"`
	//姓名
	Name string `db:"name" json:"name"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//性别
	// 0 男 1 女 2 未知
	Gender int `db:"gender" json:"gender"`
	//证件编号
	IDCard string `db:"id_card" json:"idCard"`
	//身份证正反面照片
	IDCardFrontFileID int64 `db:"id_card_front_file_id" json:"idCardFrontFileID"`
	IDCardBackFileID  int64 `db:"id_card_back_file_id" json:"idCardBackFileID"`
	//联系电话
	Phone string `db:"phone" json:"phone"`
	//个人照片
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//其他照片信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//家庭住址
	Address string `db:"address" json:"address"`
	//出生年月
	// date结构
	DateOfBirth time.Time `db:"date_of_birth" json:"dateOfBirth"`
	//婚姻状态
	MaritalStatus bool `db:"marital_status" json:"maritalStatus"`
	//教育状态
	// 0 无教育; 1 小学; 2 初中; 3 高中; 4 专科技校; 5 本科; 6 研究生; 7 博士; 8 博士后
	EducationStatus int `db:"education_status" json:"educationStatus"`
	//职业
	Profession string `db:"profession" json:"profession"`
	//级别
	// 看护级别; 其他级别
	Level int `db:"level" json:"level"`
	//主要紧急联系人
	EmergencyContact      string `db:"emergency_contact" json:"emergencyContact"`
	EmergencyContactPhone string `db:"emergency_contact_phone" json:"emergencyContactPhone"`
	//关联的文件
	DocID int64 `db:"doc_id" json:"docID"`
	//描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//负责人1
	Director1 int64 `db:"director_1" json:"director1" check:"id" empty:"true"`
	//负责人2
	Director2 int64 `db:"director_2" json:"director2" check:"id" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
