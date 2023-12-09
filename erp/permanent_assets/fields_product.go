package ERPPermanentAssets

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsProduct 固定资产
type FieldsProduct struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//使用年限
	UseExpireYear int64 `db:"use_expire_year" json:"useExpireYear"`
	//使用月份限
	UseExpireMonth int64 `db:"use_expire_month" json:"useExpireMonth"`
	//下一次盘点时间
	WaitCheckAt time.Time `db:"wait_check_at" json:"waitCheckAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//录入操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//计划盘点人
	CheckOrgBindID int64 `db:"check_org_bind_id" json:"checkOrgBindID"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//资产名称
	Name string `db:"name" json:"name"`
	//资产条码
	Code string `db:"code" json:"code"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述
	Des string `db:"des" json:"des"`
	//购买单价
	BuyPerPrice int64 `db:"buy_per_price" json:"buyPerPrice"`
	//购买总价
	BuyAllPrice int64 `db:"buy_all_price" json:"buyAllPrice"`
	//当前资产单价
	NowPerPrice int64 `db:"now_per_price" json:"nowPerPrice"`
	//当前总价值
	NowAllPrice int64 `db:"now_all_price" json:"nowAllPrice"`
	//当前数量
	Count int64 `db:"count" json:"count"`
	//正在使用数量
	UseCount int64 `db:"use_count" json:"useCount"`
	//存放地点
	SavePlace string `db:"save_place" json:"savePlace"`
	//约定使用部门名称
	// 可以使用log借出逻辑，或者在这里直接指定部门名称
	PlanUseOrgGroupName string `db:"plan_use_org_group_name" json:"planUseOrgGroupName"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
