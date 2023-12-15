package TMSSelfOtherKuai100

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreHttp2 "github.com/fotomxq/weeekj_core/v5/core/http2"
)

// PollQueryData 实时快递单号查询数据
type PollQueryData struct {
	//消息内容
	Message string `json:"message"`
	//当前状态
	State string `json:"state"`
	//通讯状态
	Status string `json:"status"`
	//快递明细状态
	Condition string `json:"condition"`
	//是否已经签收
	IsCheck string `json:"ischeck"`
	//快递公司编号
	Com string `json:"com"`
	//快递单号
	Nu string `json:"nu"`
	//日志
	Data []PollQueryDataLog `json:"data"`
}

type PollQueryDataLog struct {
	//内容
	Context string `json:"context"`
	//时间
	Time string `json:"time"`
	//格式化后的时间
	FTime string `json:"ftime"`
	//状态
	Status string `json:"status"`
	//状态代码
	StatusCode string `json:"status_code"`
	//分区代码
	AreaCode string `json:"area_code"`
	//分区名称
	AreaName string `json:"area_name"`
	//分区经纬度
	AreaCenter string `json:"area_center"`
	//当前快递点
	Location string `json:"location"`
	//行政区拼音
	AreaPinYin string `json:"area_pin_yin"`
}

// PollQuery 实时快递单号查询
// 参考文档: https://api.kuaidi100.com/document/5f0ffb5ebc8da837cbd8aefc
func PollQuery(orgID int64, companySN, tmsSN, phone, addrFrom, addrTo string) (result PollQueryData, errCode string, err error) {
	//请求数据集合
	customer := getCustomer(orgID)
	if customer == "" {
		errCode = "err_config"
		err = errors.New("customer is empty")
		return
	}
	type paramsType struct {
		Com   string `json:"com"`
		Num   string `json:"num"`
		Phone string `json:"phone"`
		From  string `json:"from"`
		To    string `json:"to"`
	}
	params := paramsType{
		Com:   companySN,
		Num:   tmsSN,
		Phone: phone,
		From:  addrFrom,
		To:    addrTo,
	}
	var sign string
	sign, errCode, err = getSign(orgID, fmt.Sprint(params.Com, params.Num, params.Phone, params.From, params.To))
	if err != nil {
		return
	}
	type bodyType struct {
		Customer string     `json:"customer"`
		Sign     string     `json:"sign"`
		Params   paramsType `json:"params"`
		Order    string     `json:"order"`
	}
	body := bodyType{
		Customer: customer,
		Sign:     sign,
		Params:   params,
		Order:    "desc",
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		errCode = "err_json"
		return
	}
	client := CoreHttp2.Core{
		Url:          "https://poll.kuaidi100.com/poll/query.do",
		Method:       CoreHttp2.MethodPost,
		Body:         bodyJSON,
		Params:       nil,
		IsRandHeader: false,
		IsProxy:      false,
		ProxyIP:      "",
		Client:       nil,
		Request:      nil,
		Response:     nil,
		Err:          nil,
		RequestErr:   nil,
		ResponseErr:  nil,
	}
	client.Resp()
	if client.Err != nil {
		err = client.Err
		return
	}
	client.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client.DoResp()
	if client.Err != nil {
		err = client.Err
		return
	}
	//将数据解析为json结构
	err = client.ResultJSON(&result)
	if err != nil {
		return
	}
	//反馈
	return
}
