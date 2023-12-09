package IOTDevice

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsAction 设备动作
type FieldsAction struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//动作对应任务的默认过期时间
	ExpireTime int64 `db:"expire_time" json:"expireTime"`
	//连接方式
	// mqtt_client 与设备直接连接，用于标准物联网设计
	// mqtt_group 设备分组与设备进行mqtt广播，可用于app通告方法等
	// none 交给业务模块进行处理，任务终端不做任何广播处理
	// 本系统默认支持的是mqtt，tcp建议采用微服务跨应用或组件方式构建，以避免系统级阻塞
	ConnectType string `db:"connect_type" json:"connectType"`
	//扩展参数
	Configs CoreSQLConfig.FieldsConfigsType `db:"configs" json:"configs"`
}
