package CoreCurrency

import (
	"errors"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 货币模块
var (
	//人民币
	data = []DataCurrencyType{
		{
			ID:   86,
			Mark: "CNY",
			Name: "人民币",
		},
		{
			ID:   1,
			Mark: "USD",
			Name: "美元",
		},
		{
			ID:   61,
			Mark: "AUD",
			Name: "澳元",
		},
	}
)

type DataCurrencyType struct {
	//数字编码
	ID int
	//英文标识码，国际通用
	Mark string
	//名称
	Name string
}

// GetName 获取某个标识码的名称
func GetName(mark string) string {
	for _, v := range data {
		if v.Mark == mark {
			return v.Name
		}
	}
	return ""
}

// GetMarkByID 通过ID获取标识码
func GetMarkByID(id int) string {
	for _, v := range data {
		if v.ID == id {
			return v.Mark
		}
	}
	return ""
}

// GetID 获取ID
func GetID(mark string) int {
	for _, v := range data {
		if v.Mark == mark {
			return v.ID
		}
	}
	return 0
}

// CheckMark 检查标识码是否正确存在
func CheckMark(mark string) error {
	for _, v := range data {
		if v.Mark == mark {
			return nil
		}
	}
	return errors.New("currency mark is not exist")
}

// CheckID 检查是否支持
func CheckID(id int) error {
	if Router2SystemConfig.GlobConfig.Finance.NeedNoCheck {
		return nil
	}
	for _, v := range data {
		if v.ID == id {
			return nil
		}
	}
	return errors.New("currency id is not exist")
}
