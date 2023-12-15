package ServiceInfoExchange

import (
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	OrgMapMod "github.com/fotomxq/weeekj_core/v5/org/map/mod"
	ServiceOrderMod "github.com/fotomxq/weeekj_core/v5/service/order/mod"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//商户地图创建后构建帖子
	CoreNats.SubDataByteNoErr("/org/map/audit", subNatsOrgMapAudit)
	//通知等待订单创建完成
	CoreNats.SubDataByteNoErr("/service/order/create_wait_finish", subNatsOrderWaitFinish)
	//通知订单审核并完成支付
	CoreNats.SubDataByteNoErr("/service/order/next", subNatsOrderNext)
	//发生评论行为
	CoreNats.SubDataByteNoErr("/class/comment", subNatsComment)
}

// 商户地图创建后构建帖子
func subNatsOrgMapAudit(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	//检查配置项
	open, _ := BaseConfig.GetDataBool("OrgMapAuditAutoCreateServiceInfoExchangeByUser")
	if !open {
		return
	}
	//获取地图信息
	mapData := OrgMapMod.GetMapByID(id)
	if mapData.ID < 1 {
		return
	}
	if mapData.UserID < 1 {
		return
	}
	//构建帖子信息
	infoID, err := createInfo(&ArgsCreateInfo{
		InfoType:     "org_map",
		ExpireAt:     "",
		OrgID:        0,
		UserID:       mapData.UserID,
		SortID:       0,
		Tags:         nil,
		Title:        mapData.Name,
		TitleDes:     CoreFilter.SubStrQuick(mapData.Des, 10),
		Des:          mapData.Des,
		CoverFileIDs: nil,
		Currency:     0,
		Price:        0,
		LimitCount:   0,
		Address:      CoreSQLAddress.FieldsAddress{},
		Params: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "orgMapID",
				Val:  fmt.Sprint(mapData.ID),
			},
		},
	})
	if err != nil {
		CoreLog.Warn("service info exchange sub nats org map audit, create info, ", err)
		return
	}
	if err = PublishInfo(&ArgsPublishInfo{
		ID:     infoID,
		OrgID:  -1,
		UserID: -1,
	}); err != nil {
		CoreLog.Warn("service info exchange sub nats org map audit, publish info, ", err)
		return
	}
}

// 通知等待订单创建完成
func subNatsOrderWaitFinish(_ *nats.Msg, _ string, id int64, _ string, data []byte) {
	//获取信息结构
	infoData := getInfoByWaitOrderID(id)
	if infoData.ID < 1 {
		return
	}
	//解析结构体
	orderID := gjson.GetBytes(data, "orderID").Int()
	if orderID < 1 {
		return
	}
	//修改订单ID
	if err := updateInfoOrderID(infoData.ID, orderID); err != nil {
		CoreLog.Warn("service info exchange sub nats order wait finish, info id: ", infoData.ID, ", update info order id: ", orderID, ", err: ", err)
		return
	}
}

// 通知订单支付完成
func subNatsOrderNext(_ *nats.Msg, _ string, orderID int64, _ string, _ []byte) {
	//获取信息
	infoData := getInfoByOrderID(orderID)
	if infoData.ID < 1 {
		return
	}
	//标记订单完成
	if err := updateInfoOrderFinish(infoData.ID, orderID); err != nil {
		CoreLog.Warn("service info exchange sub nats order wait finish, info id: ", infoData.ID, ", update info order finish by id: ", orderID, ", err: ", err)
		return
	}
	//反馈订单完成
	ServiceOrderMod.UpdateFinish(orderID, "信息交互订单完成")
}

// 发生评论
func subNatsComment(_ *nats.Msg, _ string, _ int64, systemMark string, data []byte) {
}
