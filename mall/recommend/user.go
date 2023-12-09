package MallRecommend

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	MallCore "gitee.com/weeekj/weeekj_core/v5/mall/core"
)

// GetRecommendDefaultList 获取某个用户在某个商户下的推荐商品列
// 本方法会从列队中提取数据，且排除要排除掉的商品ID
// 满足条件后反馈数据，如果数据不足的，将用销量倒叙投放数据
func GetRecommendDefaultList(orgID int64, userID int64, noHaveProductID int64, limit int) (dataList []MallCore.FieldsCore, err error) {
	//如果需要补全数据，则推送请求计算推荐数据包
	CoreNats.PushDataNoErr("/mall/recommend/user", "", userID, "", nil)
	//获取销量倒叙排名商品，补全数据
	appendByCountList, _, _ := MallCore.GetProductList(&MallCore.ArgsGetProductList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  int64(limit - len(dataList)),
			Sort: "buy_count",
			Desc: true,
		},
		OrgID:              orgID,
		ProductType:        -1,
		NeedIsVirtual:      false,
		IsVirtual:          false,
		SortID:             -1,
		Tags:               []int64{},
		NeedIsPublish:      true,
		IsPublish:          true,
		PriceMin:           -1,
		PriceMax:           -1,
		NeedHaveVIP:        false,
		HaveVIP:            false,
		Tickets:            []int64{},
		NeedHaveCount:      false,
		HaveCount:          false,
		NeedHaveIntegral:   false,
		HaveIntegral:       false,
		ParentID:           0,
		TransportID:        -1,
		WarehouseProductID: -1,
		IsRemove:           false,
		Search:             "",
	})
	if len(appendByCountList) > 0 {
		for _, v := range appendByCountList {
			dataList = append(dataList, v)
		}
	}
	//反馈
	return
}

// 为用户构建数据
func putUserData(userID int64) {

}
