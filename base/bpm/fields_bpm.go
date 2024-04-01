package BaseBPM

import "time"

// FieldsBPM BPM核心载体
type FieldsBPM struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//所属主题
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id"`
	//流程节点数量
	NodeCount int `db:"node_count" json:"nodeCount" check:"int64Than0"`
	//流程json文件内容
	JSONNode string `db:"json_node" json:"jsonNode"`
}
