package ServiceInfoExchangeMod

import (
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsInfo 信息核心
type FieldsInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime" empty:"true"`
	//发布时间
	PublishAt time.Time `db:"publish_at" json:"publishAt"`
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//审核拒绝原因
	AuditDes string `db:"audit_des" json:"auditDes" check:"des" min:"1" max:"600" empty:"true"`
	//信息类型
	// none 普通类型; recruitment 招聘信息; rent 租房信息; thing 物品交易
	InfoType string `db:"info_type" json:"infoType" check:"mark"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//副标题
	TitleDes string `db:"title_des" json:"titleDes" check:"title" min:"1" max:"600" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面ID
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs" check:"ids" empty:"true"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency" empty:"true"`
	//费用
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//报名人数限制
	// <1 不限制
	LimitCount int64 `db:"limit_count" json:"limitCount" check:"int64Than0"`
	//关联的订单
	OrderID     int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	WaitOrderID int64 `db:"wait_order_id" json:"waitOrderID" check:"id" empty:"true"`
	//订单是否完成
	OrderFinish bool `db:"order_finish" json:"orderFinish"`
	//唯一送货地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
