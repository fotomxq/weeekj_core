package IOTQuickRecord

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsRecord 临时记录表
type FieldsRecord struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//设备标识码
	DeviceCode string `db:"device_code" json:"deviceCode"`
	//设备ID
	// 匹配好的设备，保留直到设备领取数据
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
