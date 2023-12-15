package ServiceCompany

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsCompanyAudit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//hash值
	// 唯一的数据，可用于查询对应组织，取代ID直接查询；或用于第三方系统同步数据处理用
	Hash string `db:"hash" json:"hash"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用途
	// client 客户; supplier 供应商; partners 合作商; service 服务商
	UseType string `db:"use_type" json:"useType"`
	//绑定组织ID
	BindOrgID int64 `db:"bind_org_id" json:"bindOrgID"`
	//绑定用户ID
	// 主绑定关系具备所有能力，类似组织的拥有人
	BindUserID int64 `db:"bind_user_id" json:"bindUserID"`
	//名称
	Name string `db:"name" json:"name"`
	//公司营业执照编号
	SN string `db:"sn" json:"sn"`
	//描述
	Des string `db:"des" json:"des"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//所属城市
	City int `db:"city" json:"city"`
	//街道详细信息
	Address string `db:"address" json:"address"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
