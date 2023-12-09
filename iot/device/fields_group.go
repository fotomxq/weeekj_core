package IOTDevice

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsGroup 设备分组
type FieldsGroup struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//分区标识码
	// 全局必须唯一
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面列
	// 第一张作为封面
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles"`
	//支持动作ID组
	Action pq.Int64Array `db:"action" json:"action"`
	//心跳超时时间
	// 超出时间没有通讯则判定掉线
	// 单位: 秒
	ExpireTime int64 `db:"expire_time" json:"expireTime"`
	//设备的预计使用场景
	// 0 public 公共设备 / 1 private 私有设备
	// 如果>1 则为自定义设置，具体由设备驱动识别处理
	UseType int `db:"use_type" json:"useType"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
