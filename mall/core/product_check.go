package MallCore

import (
	"errors"
	"fmt"
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	UserSubscription "gitee.com/weeekj/weeekj_core/v5/user/subscription"
	UserTicket "gitee.com/weeekj/weeekj_core/v5/user/ticket"
	"github.com/lib/pq"
	"strings"
)

// ArgsCheckProductTicket 检查票据是否可用于商品参数
type ArgsCheckProductTicket struct {
	//商品ID
	// 不一定非要给该值
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//会员配置列
	// 会员配置分平台和商户，平台会员需参与活动才能使用，否则将禁止设置和后期使用
	UserSubs pq.Int64Array `db:"user_subs" json:"userSubs" check:"ids" empty:"true"`
	//票据ID列
	UserTicket pq.Int64Array `db:"user_ticket" json:"userTicket" check:"ids" empty:"true"`
	//购买数量
	// 如果为0，则只判断是否可用
	BuyCount int64 `db:"buy_count" json:"buyCount" check:"int64Than0" empty:"true"`
}

// CheckProductTicket 检查票据是否可用于商品
// 检查是否可用、授权给商户
func CheckProductTicket(args *ArgsCheckProductTicket) (errCode string, err error) {
	if len(args.UserTicket) < 1 {
		return
	}
	for _, v := range args.UserTicket {
		var configData UserTicket.FieldsConfig
		configData, err = UserTicket.GetConfigByID(&UserTicket.ArgsGetConfigByID{
			ID:    v,
			OrgID: args.OrgID,
		})
		if err != nil || configData.ID < 1 {
			errCode = "ticket_config_not_exist"
			err = errors.New(fmt.Sprint("user ticket config not exist, ", err))
			return
		}
	}
	if args.ID > 0 {
		//获取商品数据
		var productData FieldsCore
		productData, err = GetProduct(&ArgsGetProduct{
			ID:    args.ID,
			OrgID: -1,
		})
		if err != nil || productData.DeleteAt.Unix() > 1000000 {
			errCode = "ticket_product_not_exist"
			err = errors.New(fmt.Sprint("product not exist, ", err))
			return
		}
		//寻找会员配置，如果符合条件则反馈
		for _, v := range args.UserTicket {
			isFind := false
			for _, v2 := range productData.UserTicket {
				if v == v2 {
					isFind = true
					break
				}
			}
			if isFind {
				//符合条件，反馈
				err = nil
				return
			}
		}
		//获取分类数据
		var sortParam string
		sortParam, err = Sort.GetParam(&ClassSort.ArgsGetParam{
			ID:     args.SortID,
			BindID: args.OrgID,
			Mark:   "user_tickets",
		})
		if err == nil {
			//拆分数据结构体
			var sortParams []string
			sortParams = strings.Split(sortParam, ",")
			for _, v := range sortParams {
				var vInt64 int64
				vInt64, err = CoreFilter.GetInt64ByString(v)
				if err != nil {
					continue
				}
				for _, v2 := range args.UserSubs {
					if vInt64 == v2 {
						err = nil
						return
					}
				}
			}
		}
		//获取标签
		// 标案暂不支持该设计
		//反馈失败
		errCode = "ticket_not_find"
		err = errors.New("not find any config to use user tickets")
	}
	return
}

// ArgsCheckProductSub 检查会员是否可用于商品参数
type ArgsCheckProductSub struct {
	//商品ID
	// 不一定非要给该值
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//会员配置列
	// 会员配置分平台和商户，平台会员需参与活动才能使用，否则将禁止设置和后期使用
	UserSubs pq.Int64Array `db:"user_subs" json:"userSubs" check:"ids" empty:"true"`
	//票据ID列
	// 检查票据和会员的互斥关系
	UserTicket pq.Int64Array `db:"user_ticket" json:"userTicket" check:"ids" empty:"true"`
	//购买数量
	// 如果为0，则只判断是否可用
	BuyCount int64 `db:"buy_count" json:"buyCount" check:"int64Than0" empty:"true"`
}

// CheckProductSub 检查会员是否可用于商品
// 检查可用、授权给商户
// 分类和标签的相关数据，在分类和标签扩展参数user_subs中存放
func CheckProductSub(args *ArgsCheckProductSub) (errCode string, err error) {
	if len(args.UserSubs) < 1 {
		return
	}
	for _, v := range args.UserSubs {
		var configData UserSubscription.FieldsConfig
		configData, err = UserSubscription.GetConfigByID(&UserSubscription.ArgsGetConfigByID{
			ID:    v,
			OrgID: args.OrgID,
		})
		if err != nil || configData.ID < 1 {
			errCode = "sub_config_not_exist"
			err = errors.New(fmt.Sprint("user sub config not exist, ", err))
			return
		}
	}
	if args.ID > 0 {
		//获取商品数据
		var productData FieldsCore
		productData, err = GetProduct(&ArgsGetProduct{
			ID:    args.ID,
			OrgID: -1,
		})
		if err != nil || productData.DeleteAt.Unix() > 1000000 {
			errCode = "sub_product_not_exist"
			err = errors.New(fmt.Sprint("product not exist, ", err))
			return
		}
		//寻找会员配置，如果符合条件则反馈
		for _, v := range args.UserSubs {
			isFind := false
			for _, v2 := range productData.UserSubPrice {
				if v == v2.ID {
					isFind = true
					break
				}
			}
			if isFind {
				//符合条件，反馈
				err = nil
				return
			}
		}
		//获取分类数据
		var sortParam string
		sortParam, err = Sort.GetParam(&ClassSort.ArgsGetParam{
			ID:     args.SortID,
			BindID: args.OrgID,
			Mark:   "user_subs",
		})
		if err == nil {
			//拆分数据结构体
			var sortParams []string
			sortParams = strings.Split(sortParam, ",")
			for _, v := range sortParams {
				var vInt64 int64
				vInt64, err = CoreFilter.GetInt64ByString(v)
				if err != nil {
					continue
				}
				for _, v2 := range args.UserSubs {
					if vInt64 == v2 {
						err = nil
						return
					}
				}
			}
		}
		//获取标签
		// 标案暂不支持该设计
		//反馈失败
		errCode = "sub_not_find"
		err = errors.New("not find any config to use user sub")
	}
	return
}
