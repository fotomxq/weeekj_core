package BaseSMS

import (
	"errors"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

//给腾讯发送短信
func createSMSToTencent(configData *FieldsConfig, smsData *FieldsSMS) error {
	//初始化短信前置
	credential := common.NewCredential(
		configData.AppID,
		configData.AppKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	cpf.SignMethod = "HmacSHA1"
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
	request := sms.NewSendSmsRequest()
	appID, b := configData.Params.GetVal("appID")
	if !b {
		return errors.New("tencent sms not have app id")
	}
	request.SmsSdkAppId = common.StringPtr(appID)
	request.SignName = common.StringPtr(configData.TemplateSign)
	request.SenderId = common.StringPtr("")
	request.SessionContext = common.StringPtr(fmt.Sprint(smsData.ID))
	request.ExtendCode = common.StringPtr("")
	var params []string
	if len(configData.TemplateParams) > 0 {
		for i := 0; i < len(configData.TemplateParams); i++ {
			for _, v2 := range smsData.Params {
				if configData.TemplateParams[i].Mark == v2.Mark {
					params = append(params, v2.Val)
					break
				}
			}
		}
	}
	if len(params) > 0 {
		request.TemplateParamSet = common.StringPtrs(params)
	}
	request.TemplateId = common.StringPtr(configData.TemplateID)
	if smsData.NationCode == "" {
		smsData.NationCode = "86"
	}
	request.PhoneNumberSet = common.StringPtrs([]string{"+" + smsData.NationCode + smsData.Phone})
	response, err := client.SendSms(request)
	if err != nil {
		return errors.New(fmt.Sprintf("An API error has returned: %s", err))
	}
	for _, v := range response.Response.SendStatusSet {
		if *v.Code != "" {
			//CoreLog.Info("sms failed: ", *v.Code, ", msg: ", *v.Message)
		}
	}
	return nil
}
