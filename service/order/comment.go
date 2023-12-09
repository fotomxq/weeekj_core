package ServiceOrder

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateCommentBuyer 买家评价参数
type ArgsUpdateCommentBuyer struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
	//货物来源
	GoodFrom CoreSQLFrom.FieldsFrom `db:"good_from" json:"goodFrom"`
	//买家评论ID
	CommentBuyerID int64 `db:"comment_buyer_id" json:"commentBuyerID" check:"id"`
}

// UpdateCommentBuyer 买家评价
func UpdateCommentBuyer(args *ArgsUpdateCommentBuyer) (err error) {
	//获取数据
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  0,
		UserID: 0,
	})
	if err != nil {
		return
	}
	//找到货物，和设定评论
	for k, v := range orderData.Goods {
		if !v.From.CheckEg(args.GoodFrom) {
			continue
		}
		if orderData.Goods[k].CommentBuyer {
			err = errors.New("have comment")
			return
		}
		orderData.Goods[k].CommentBuyer = true
		orderData.Goods[k].CommentBuyerID = args.CommentBuyerID
	}
	//修改数据
	var newLog string
	newLog, err = getLogData(args.UserID, 0, "cancel", args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), goods = :goods, logs = logs || :log WHERE id = :id AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
		"id":      args.ID,
		"user_id": args.UserID,
		"log":     newLog,
		"goods":   orderData.Goods,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}

// ArgsUpdateCommentSeller 卖家评价参数
type ArgsUpdateCommentSeller struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
	//货物来源
	GoodFrom CoreSQLFrom.FieldsFrom `db:"good_from" json:"goodFrom"`
	//是否卖家ID
	CommentSellerID int64 `db:"comment_seller_id" json:"commentSellerID" check:"id"`
}

// UpdateCommentSeller 卖家评价
func UpdateCommentSeller(args *ArgsUpdateCommentSeller) (err error) {
	//获取数据
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  0,
		UserID: 0,
	})
	if err != nil {
		return
	}
	//找到货物，和设定评论
	for k, v := range orderData.Goods {
		if !v.From.CheckEg(args.GoodFrom) {
			continue
		}
		if orderData.Goods[k].CommentSeller {
			err = errors.New("have comment")
			return
		}
		orderData.Goods[k].CommentSeller = true
		orderData.Goods[k].CommentSellerID = args.CommentSellerID
	}
	//修改数据
	var newLog string
	newLog, err = getLogData(0, args.OrgBindID, "comment_seller", args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), goods = :goods, logs = logs || :log WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
		"goods":  orderData.Goods,
		"log":    newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}
