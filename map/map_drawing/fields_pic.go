package MapMapDrawing

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsPic 地图的图片
type FieldsPic struct {
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
	//主图ID
	ParentID int64 `db:"parent_id" json:"parentID"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//核心地图
	FileID int64 `db:"file_id" json:"fileID"`
	//修正图片高度和宽度
	FixHeight int `db:"fix_height" json:"fixHeight"`
	FixWidth  int `db:"fix_width" json:"fixWidth"`
	//按钮文字
	ButtonName string `db:"button_name" json:"buttonName"`
	//绑定电子围栏
	BindAreaID int64 `db:"bind_area_id" json:"bindAreaID"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
