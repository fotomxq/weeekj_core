package UserSubscription

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsUpdateConfig 修改订阅配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//时间类型
	// 0 小时 1 天 2 周 3 月 4 年
	TimeType int `db:"time_type" json:"timeType" check:"intThan0" empty:"true"`
	//时间长度
	TimeN int `db:"time_n" json:"timeN" check:"intThan0" empty:"true"`
	//开通价格
	Currency int   `db:"currency" json:"currency" check:"intThan0" empty:"true"`
	Price    int64 `db:"price" json:"price" check:"price"`
	//折扣前费用，用于展示
	PriceOld int64 `db:"price_old" json:"priceOld" check:"price"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//关联的用户组
	// 只有为平台配置时，该数据才可修改病会生效
	UserGroups pq.Int64Array `db:"user_groups" json:"userGroups" check:"ids" empty:"true"`
	//默认减免的费用比例、费用金额
	ExemptionPrice int64 `db:"exemption_price" json:"exemptionPrice" check:"price" empty:"true"`
	// 1-100% 百分比
	ExemptionDiscount int64 `db:"exemption_discount" json:"exemptionDiscount"`
	//费用低于多少时，将失效
	// 依赖于订单的总金额判断
	ExemptionMinPrice int64 `db:"exemption_min_price" json:"exemptionMinPrice" check:"price" empty:"true"`
	//限制设计
	// 允许设置多个条件，如1天限制一次的同时、30天能使用10次
	Limits FieldsLimits `db:"limits" json:"limits"`
	//周期价格
	ExemptionTime FieldsExemptionTimes `db:"exemption_time" json:"exemptionTime"`
	//样式ID
	// 关联到样式库后，本记录的图片和文本将交给样式库布局实现
	StyleID int64 `db:"style_id" json:"styleID" check:"id" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConfig 修改订阅配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	//组织的订阅不能授权用户分组
	if args.OrgID > 0 {
		//重新处理用户组
		if len(args.UserGroups) > 0 {
			args.UserGroups = []int64{}
		}
		//检查用户组归属权
		if err = checkUserGroup(args.OrgID, args.UserGroups); err != nil {
			return
		}
		//重新获取用户组数据包
		if len(args.UserGroups) > 0 {
			var data FieldsConfig
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, user_groups FROM user_sub_config WHERE org_id = $1 AND id != $2 AND mark = $3 AND delete_at < to_timestamp(1000000)", args.OrgID, args.ID, args.Mark)
			if err != nil {
				return
			}
			args.UserGroups = data.UserGroups
		}
	}
	//更新的mark如果和之前不一致
	var data FieldsConfig
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_sub_config WHERE org_id = $1 AND id != $2 AND mark = $3 AND delete_at < to_timestamp(1000000)", args.OrgID, args.ID, args.Mark)
	if err == nil && data.ID > 0 {
		err = errors.New("mark is exist")
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_sub_config SET update_at = NOW(), mark = :mark, time_type = :time_type, time_n = :time_n, currency = :currency, price = :price, price_old = :price_old, title = :title, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, user_groups = :user_groups, exemption_discount = :exemption_discount, exemption_price = :exemption_price, exemption_min_price = :exemption_min_price, limits = :limits, exemption_time = :exemption_time, style_id = :style_id, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
