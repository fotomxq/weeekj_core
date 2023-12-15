package OrgSubscription

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsConfig 可选的功能
type FieldsConfig struct {
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
	//对应组织功能
	FuncList pq.StringArray `db:"func_list" json:"funcList"`
	//时间类型
	// 0 小时 1 天 2 周 3 月 4 年
	TimeType int `db:"time_type" json:"timeType"`
	//时间长度
	TimeN int `db:"time_n" json:"timeN"`
	//开通价格
	Currency int   `db:"currency" json:"currency"`
	Price    int64 `db:"price" json:"price"`
	//折扣前费用，用于展示
	PriceOld int64 `db:"price_old" json:"priceOld"`
	//标题
	Title string `db:"title" json:"title"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//样式ID
	// 关联到样式库后，本记录的图片和文本将交给样式库布局实现
	StyleID int64 `db:"style_id" json:"styleID"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
