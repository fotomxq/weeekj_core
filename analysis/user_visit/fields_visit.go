package AnalysisUserVisit

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsVisit 访问记录表
type FieldsVisit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 访问的时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 如果存在数据，则表明该数据隶属于指定组织
	// 组织依可查看该数据
	OrgID int64 `db:"org_id" json:"orgID"`
	//数据来源
	// 来自哪个模块
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//关联的用户
	UserID int64 `db:"user_id" json:"userID"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//挖掘的电话号码
	Phone string `db:"phone" json:"phone"`
	//IP地址
	IP string `db:"ip" json:"ip"`
	//浏览器标识
	// 或设备标识
	Mark string `db:"mark" json:"mark"`
	//行为标记
	// insert 进入; out 离开; move 移动
	// buy_page 进入购物页面 ... 等，具体参考文档
	Action string `db:"action" json:"action"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
