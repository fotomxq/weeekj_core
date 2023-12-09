package BaseSMS

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"testing"
	"time"
)

func TestInitSMS(t *testing.T) {
	TestInitConfig(t)
	TestCreateConfig(t)
}

func TestCreateSMS(t *testing.T) {
	var err error
	newUpdateHash := CoreFilter.GetRandNumber(10, 9999)
	_, err = CreateSMSCheck(&ArgsCreateSMSCheck{
		OrgID:      0,
		ConfigID:   newConfigData.ID,
		Token:      int64(newUpdateHash),
		NationCode: "86",
		Phone:      "17635705566",
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     1,
			Mark:   "user-mark",
			Name:   "user-name",
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestCheckSMS(t *testing.T) {
	//强制获取该验证码数据，并尝试验证
	var data []FieldsSMS
	if err := Router2SystemConfig.MainDB.Select(&data, "SELECT * FROM core_sms ORDER BY id DESC LIMIT 1 OFFSET 0"); err != nil {
		t.Error(err)
	} else {
		if len(data) < 1 {
			t.Error("len data less 1")
			return
		}
		val, b := data[0].Params.GetVal("val")
		if !b {
			val, b = data[0].Params.GetVal("code")
			if !b {
				t.Error("no have val or code")
				return
			}
		}
		if b := CheckSMS(&ArgsCheckSMS{
			ConfigID: newConfigData.ID,
			Token:    data[0].Token,
			Value:    val,
		}); !b {
			t.Error("check sms is false")
		}
	}
}

func TestRun(t *testing.T) {
	go Run()
	time.Sleep(time.Second * 5)
}

func TestClearSMS(t *testing.T) {
	TestDeleteConfig(t)
	TestClearConfig(t)
}
