package MallLog

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	MallCore "gitee.com/weeekj/weeekj_core/v5/mall/core"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//添加一条日志
	CoreNats.SubDataByteNoErr("/mall/log/new", func(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
		//解析参数
		userID := gjson.GetBytes(data, "userID").Int()
		ip := gjson.GetBytes(data, "ip").String()
		orgID := gjson.GetBytes(data, "orgID").Int()
		productID := gjson.GetBytes(data, "productID").Int()
		action := int(gjson.GetBytes(data, "action").Int())
		//添加日志
		AppendLog(userID, ip, orgID, productID, action)
	})
	//商品评论订阅
	CoreNats.SubDataByteNoErr("/class/comment", subNatsNewComment)
}

// 商品评论订阅
func subNatsNewComment(_ *nats.Msg, action string, id int64, mark string, _ []byte) {
	//必须是商品的评论，否则退出
	if mark != "mall_core_product" {
		return
	}
	//必须是新增评论
	if action != "new" {
		return
	}
	//获取评论
	commentData := MallCore.Comment.GetByID(id)
	if commentData.ID < 1 {
		return
	}
	//获取商品信息
	productData, _ := MallCore.GetProduct(&MallCore.ArgsGetProduct{
		ID:    commentData.BindID,
		OrgID: -1,
	})
	if productData.ID < 1 {
		return
	}
	//记录数据
	AppendLog(commentData.UserID, "", commentData.OrgID, productData.ID, 1)
}
