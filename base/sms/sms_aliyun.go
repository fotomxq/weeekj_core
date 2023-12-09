package BaseSMS

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

//阿里云发送短信封装
func createSMSToAliyun(configData *FieldsConfig, smsData *FieldsSMS) error {
	//获取阿里云配置
	aliyunRegionID, b := configData.Params.GetVal("aliyunRegionID")
	if !b {
		return errors.New("aliyun not have aliyunRegionID config")
	}
	//构建签名
	client, err := dysmsapi.NewClientWithAccessKey(aliyunRegionID, configData.AppID, configData.AppKey)
	if err != nil {
		return errors.New("client is error, " + err.Error())
	}
	//构建请求头
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = smsData.Phone
	request.SignName = configData.TemplateSign
	request.TemplateCode = configData.TemplateID
	if len(smsData.Params) > 0 {
		request.TemplateParam = "{"
		for _, v := range smsData.Params {
			if request.TemplateParam == "{" {
				request.TemplateParam = fmt.Sprint(request.TemplateParam, "\"", v.Mark, "\":\"", v.Val, "\"")
			} else {
				request.TemplateParam = fmt.Sprint(request.TemplateParam, ",\"", v.Mark, "\":\"", v.Val, "\"")
			}
		}
		request.TemplateParam = request.TemplateParam + "}"
	}
	//发起请求
	response, err := client.SendSms(request)
	if err != nil {
		return errors.New("send sms failed, " + err.Error())
	}
	//反馈数据分析
	if response.Code != "OK" {
		return errors.New("response is not ok, code: " + response.Code + ", BizId: " + response.BizId + ", RequestId: " + response.RequestId)
	}
	//反馈成功
	return nil
}
