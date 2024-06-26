package ServiceCompany

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//用户绑定了手机号
	CoreNats.SubDataByteNoErr("user_core_new_phone", "/user/core/new_phone", subNatsUserNewPhone)
	//注册服务
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "服务公司变动通知",
		Description:  "",
		EventSubType: "all",
		Code:         "service_company",
		EventType:    "nats",
		EventURL:     "/service/company",
		//TODO:待补充
		EventParams: "",
	})
}

// 用户绑定了手机号
func subNatsUserNewPhone(_ *nats.Msg, _ string, userID int64, _ string, data []byte) {
	//获取参数
	nationCode := gjson.GetBytes(data, "nationCode").String()
	phone := gjson.GetBytes(data, "phone").String()
	if nationCode == "" && phone == "" {
		return
	}
	//检查和获取绑定关系
	var count int64
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_company_bind WHERE nation_code = $1 AND phone = $2 AND user_id < 1", nationCode, phone)
	if err != nil || count < 1 {
		return
	}
	//修改所有绑定关系
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_company_bind SET user_id = :user_id WHERE nation_code = :nation_code AND phone = :phone", map[string]interface{}{
		"user_id":     userID,
		"nation_code": nationCode,
		"phone":       phone,
	})
	if err != nil {
		CoreLog.Error("service company sub nats user new phone, update all user, ", err)
		return
	}
	//删除缓冲
	var dataList []FieldsBind
	err = Router2SystemConfig.MainDB.Get(&dataList, "SELECT id FROM service_company_bind WHERE nation_code = $1 AND phone = $2", nationCode, phone)
	if err == nil && len(dataList) > 0 {
		for _, v := range dataList {
			deleteBindCache(v.ID)
		}
	}
}
