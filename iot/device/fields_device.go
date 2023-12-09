package IOTDevice

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsDevice 设备
type FieldsDevice struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//状态
	// 0 public 公共可用 / 1 private 私有 / 2 ban 停用
	Status int `db:"status" json:"status"`
	//在线状态
	IsOnline bool `db:"is_online" json:"isOnline"`
	//最后一次通讯时间
	LastAt time.Time `db:"last_at" json:"lastAt"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面列
	// 第一张作为封面
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles"`
	//描述信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID"`
	//设备编号
	// 同一个分组下，必须唯一
	Code string `db:"code" json:"code"`
	//连接密钥
	// 设备连接使用的唯一密钥
	// 设备需使用该key+code+时间戳+随机码混合计算，作为握手的识别码
	Key string `db:"key" json:"key"`
	//注册地
	// 如果设置将优先使用设备注册地，而不是管辖注册地
	Address string `db:"address" json:"address"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
