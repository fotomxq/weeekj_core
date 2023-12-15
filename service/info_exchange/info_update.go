package ServiceInfoExchange

import (
	"errors"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserSystemTip "github.com/fotomxq/weeekj_core/v5/user/system_tip"
	"github.com/lib/pq"
	"time"
)

// ArgsUpdateInfo 修改信息参数
type ArgsUpdateInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//过期时间
	ExpireAt string `db:"expire_at" json:"expireAt" check:"isoTime" empty:"true"`
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

// UpdateInfo 修改信息
func UpdateInfo(args *ArgsUpdateInfo) (err error) {
	var expireAt time.Time
	if args.ExpireAt == "" {
		expireAt = CoreFilter.GetNowTimeCarbon().AddDays(30).Time
	} else {
		expireAt, _ = CoreFilter.GetTimeByISO(args.ExpireAt)
	}
	if args.LimitCount < 1 {
		args.LimitCount = 0
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_info_exchange SET update_at = NOW(), publish_at = to_timestamp(0), audit_at = to_timestamp(0), expire_at = :expire_at, sort_id = :sort_id, tags = :tags, title = :title, title_des = :title_des, des = :des, cover_file_ids = :cover_file_ids, currency = :currency, price = :price, limit_count = :limit_count, address = :address, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
		"id":             args.ID,
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
	deleteInfoCache(args.ID)
	return
}

// ArgsPublishInfo 发布信息参数
type ArgsPublishInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// PublishInfo 发布信息
func PublishInfo(args *ArgsPublishInfo) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_info_exchange SET publish_at = NOW() WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", args)
	if err != nil {
		return
	}
	deleteInfoCache(args.ID)
	var serviceInfoExchangeAutoAudit bool
	serviceInfoExchangeAutoAudit, err = BaseConfig.GetDataBool("ServiceInfoExchangeAutoAudit")
	if err != nil {
		serviceInfoExchangeAutoAudit = false
		err = nil
	}
	if serviceInfoExchangeAutoAudit {
		_ = AuditInfo(&ArgsAuditInfo{
			ID:       args.ID,
			OrgID:    -1,
			IsAudit:  true,
			AuditDes: "",
		})
	}
	return
}

// ArgsAuditInfo 审核信息参数
type ArgsAuditInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否通过审核
	IsAudit bool `db:"is_audit" json:"isAudit" check:"bool"`
	//审核拒绝原因
	AuditDes string `db:"audit_des" json:"auditDes" check:"des" min:"1" max:"600" empty:"true"`
}

// AuditInfo 审核信息
func AuditInfo(args *ArgsAuditInfo) (err error) {
	//审核信息
	if args.IsAudit {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_info_exchange SET audit_at = NOW() WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_info_exchange SET audit_at = to_timestamp(0), audit_des = :audit_des WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	}
	if err != nil {
		return
	}
	//删除缓冲
	deleteInfoCache(args.ID)
	//发送用户消息
	infoData := getInfoByID(args.ID)
	if infoData.ID < 1 {
		err = errors.New("no data")
		return
	}
	UserSystemTip.SendSuccess(infoData.UserID, "帖子", infoData.ID, infoData.Title)
	//反馈
	return
}
