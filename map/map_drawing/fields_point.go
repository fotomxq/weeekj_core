package MapMapDrawing

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	"time"
)

// FieldsPoint 点定位
type FieldsPoint struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//主图
	PicID int64 `db:"pic_id" json:"picID"`
	//图坐标
	PicPoint CoreSQLGPS.FieldsPoint `db:"pic_point" json:"picPoint"`
	//GPS坐标
	GPSPoint CoreSQLGPS.FieldsPoint `db:"gps_point" json:"gpsPoint"`
	//误差半径
	Radius float64 `db:"radius" json:"radius"`
	//显示图标
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//显示位置点系统图标
	CoverIcon string `db:"cover_icon" json:"coverIcon"`
	//图标颜色
	CoverRGB string `db:"cover_rgb" json:"coverRGB"`
	//绑定Mark
	// 绑定的系统，如pic 主图系统; room 房间; device 设备
	BindMark string `db:"bind_mark" json:"bindMark"`
	//绑定ID
	// 该点绑定的房间，可作为联动处理
	BindID int64 `db:"bind_id" json:"bindID"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
