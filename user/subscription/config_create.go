package UserSubscription

import (
	"errors"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsCreateConfig 创建新的订阅配置参数
type ArgsCreateConfig struct {
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
	// 1-100百分比
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

// CreateConfig 创建新的订阅配置
func CreateConfig(args *ArgsCreateConfig) (data FieldsConfig, err error) {
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
		var maxConfig int64
		maxConfig, err = OrgCore.Config.GetConfigValInt64(&ClassConfig.ArgsGetConfig{
			BindID:    args.OrgID,
			Mark:      "UserSubscriptionOrgMax",
			VisitType: "admin",
		})
		if err != nil {
			maxConfig = 0
		}
		if maxConfig > 0 {
			var count int64
			count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "user_sub_config", "id", "org_id = :org_id AND delete_at < to_timestamp(1000000)", map[string]interface{}{
				"org_id": args.OrgID,
			})
			if err == nil && count > 0 {
				if count >= maxConfig {
					err = errors.New("org have too many config")
					return
				}
			}
		}
	}
	//检查mark
	if args.Mark != "" {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_sub_config WHERE org_id = $1 AND mark = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.Mark)
		if err == nil && data.ID > 0 {
			err = errors.New("config mark is exist")
			return
		}
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_sub_config", "INSERT INTO user_sub_config (org_id, mark, time_type, time_n, currency, price, price_old, title, des, cover_file_id, des_files, user_groups, exemption_discount, exemption_price, exemption_min_price, limits, exemption_time, style_id, params) VALUES (:org_id, :mark, :time_type, :time_n, :currency, :price, :price_old, :title, :des, :cover_file_id, :des_files, :user_groups, :exemption_discount, :exemption_price, :exemption_min_price, :limits, :exemption_time, :style_id, :params)", args, &data)
	return
}
