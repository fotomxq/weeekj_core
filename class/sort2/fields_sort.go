package ClassSort2

import (
	"github.com/lib/pq"
	"time"
)

// FieldsSort 分组分类
type FieldsSort struct {
	//基础
	ID int64 `db:"id" json:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0"`
	//来源ID
	// 可以是某个组织，或特定的其他系统ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true" index:"true"`
	//分组标识码
	// 用于一些特殊的显示场景做区分，可以重复
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true" index:"true"`
	//上级ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true" index:"true"`
	//排序
	Sort int `db:"sort" json:"sort" index:"true"`
	//封面图
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//介绍图文
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
}
