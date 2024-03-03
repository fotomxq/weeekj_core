package BaseBPM

import "time"

// FieldsEvent 节点事件注册
type FieldsEvent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	//删除时间
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//所属主题分类
	ThemeCategoryID int64 `db:"theme_category_id" json:"themeCategoryId" check:"id"`
	//所属主题
	// 插槽可用于的主题域
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id"`
	//事件编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300"`
	//事件类型
	// nats - NATS事件
	EventType string `db:"event_type" json:"eventType" check:"intThan0"`
	//事件地址
	// nats - 触发的地址
	EventURL string `db:"event_url" json:"eventURL" check:"des" min:"1" max:"600"`
	//事件固定参数
	// nats - 事件附带的固定参数，如果为空则根据流程阶段事件触发填入
	EventParams string `db:"event_params" json:"eventParams" check:"des" min:"1" max:"1000" empty:"true"`
}
