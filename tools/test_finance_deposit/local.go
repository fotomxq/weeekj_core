package TestFinanceDeposit

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	"testing"
)

// 本地组件
var (
	LocalConfig      FinanceDeposit.FieldsConfigType
	LocalDepositData FinanceDeposit.FieldsDepositType
)

//用于资金池的处理套件

// 创建新的资金池
func LocalCreateConfig(t *testing.T, mark string, takeOut bool, name, des string) {
	//创建资金池配置
	data, err := FinanceDeposit.SetConfig(&FinanceDeposit.ArgsSetConfig{
		Mark:             mark,
		Name:             name,
		Des:              des,
		Currency:         86,
		TakeOut:          takeOut,
		TakeLimit:        0,
		OnceSaveMinLimit: 0,
		OnceSaveMaxLimit: 0,
		OnceTakeMinLimit: 0,
		OnceTakeMaxLimit: 0,
		Configs:          nil,
	})
	if err != nil {
		t.Error(err)
	} else {
		//t.Log(data)
		LocalConfig = data
	}
}

// 检查目标资金量
func LocalGetPriceByFrom(createInfo CoreSQLFrom.FieldsFrom, fromInfo CoreSQLFrom.FieldsFrom, configMark string) {
	data, err := FinanceDeposit.GetByFrom(&FinanceDeposit.ArgsGetByFrom{
		CreateInfo: createInfo,
		FromInfo:   fromInfo,
		ConfigMark: configMark,
	})
	if err != nil {
		//如果不存在，则构建0的单位结构
		//不会发出错误信息，需要自行通过data数据结构发觉
	} else {
		//t.Log(data)
	}
	LocalDepositData = data
}

// LocalSetPrice 设置目标资金量
func LocalSetPrice(t *testing.T, createInfo CoreSQLFrom.FieldsFrom, fromInfo CoreSQLFrom.FieldsFrom, configMark string, price int64) {
	//尝试获取，如果存在则修正updateHash
	dataDeposit, err := FinanceDeposit.GetByFrom(&FinanceDeposit.ArgsGetByFrom{
		CreateInfo: createInfo,
		FromInfo:   fromInfo,
		ConfigMark: configMark,
	})
	updateHash := ""
	if err == nil {
		updateHash = dataDeposit.UpdateHash
	}
	//执行操作
	data, errCode, err := FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
		UpdateHash:      updateHash,
		CreateInfo:      createInfo,
		FromInfo:        fromInfo,
		ConfigMark:      configMark,
		AppendSavePrice: price,
	})
	if err != nil {
		t.Error("set finance deposit, code: ", errCode, ", err: ", err)
	} else {
		//t.Log(data)
		LocalDepositData = data
	}
}

// 删除资金池配置
func LocalDeleteConfig(t *testing.T, mark string) {
	if err := FinanceDeposit.DeleteConfigByMark(&FinanceDeposit.ArgsDeleteConfigByMark{
		Mark: mark,
	}); err != nil {
		t.Error(err)
	}
}
