package ServiceInfoExchange

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateInfo 创建信息参数
type ArgsCreateInfo struct {
	//信息类型
	// none 普通类型; recruitment 招聘信息; rent 租房信息; map 商户地图信息
	InfoType string `db:"info_type" json:"infoType" check:"mark"`
	//过期时间
	ExpireAt string `db:"expire_at" json:"expireAt" check:"isoTime" empty:"true"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//副标题
	TitleDes string `db:"title_des" json:"titleDes" check:"title" min:"1" max:"600" empty:"true"`
	//商品描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面ID
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs" check:"ids" empty:"true"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency" empty:"true"`
	//费用
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//报名人数限制
	// <1 不限制
	LimitCount int64 `db:"limit_count" json:"limitCount" check:"int64Than0" empty:"true"`
	//唯一送货地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address" check:"address_data" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

// CreateInfo 创建信息
func CreateInfo(args *ArgsCreateInfo) (data FieldsInfo, err error) {
	var infoID int64
	infoID, err = createInfo(args)
	if err != nil {
		return
	}
	data = getInfoByID(infoID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

func createInfo(args *ArgsCreateInfo) (infoID int64, err error) {
	if args.Tags == nil {
		args.Tags = []int64{}
	}
	if args.CoverFileIDs == nil {
		args.CoverFileIDs = []int64{}
	}
	var expireAt time.Time
	if args.ExpireAt == "" {
		expireAt = CoreFilter.GetNowTimeCarbon().AddDays(30).Time
	} else {
		expireAt, _ = CoreFilter.GetTimeByISO(args.ExpireAt)
	}
	if args.LimitCount < 1 {
		args.LimitCount = 0
	}
	infoID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO service_info_exchange (audit_des, info_type, expire_at, org_id, user_id, sort_id, tags, title, title_des, des, cover_file_ids, currency, price, limit_count, order_id, address, params) VALUES ('',:info_type,:expire_at,:org_id,:user_id,:sort_id,:tags,:title,:title_des,:des,:cover_file_ids,:currency,:price,:limit_count,0,:address,:params)", map[string]interface{}{
		"info_type":      args.InfoType,
		"expire_at":      expireAt,
		"org_id":         args.OrgID,
		"user_id":        args.UserID,
		"sort_id":        args.SortID,
		"tags":           args.Tags,
		"title":          args.Title,
		"title_des":      args.TitleDes,
		"des":            args.Des,
		"cover_file_ids": args.CoverFileIDs,
		"currency":       args.Currency,
		"price":          args.Price,
		"limit_count":    args.LimitCount,
		"address":        args.Address,
		"params":         args.Params,
	})
	if err != nil {
		return
	}
	return
}
