package TMSSelfOtherKuai100

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 读取公司名录数据
type dataTMSCompanyChild struct {
	//公司名称
	Name string `json:"name"`
	//编码
	Code string `json:"code"`
	//业务类型
	TMSType string `json:"tmsType"`
}

type dataTMSCompany struct {
	DataList []dataTMSCompanyChild `json:"dataList"`
}

// 读取公司名录数据
func loadTMSCompanyJSON() (err error) {
	var dataByte []byte
	dataByte, err = CoreFile.LoadFile(fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep, "data", CoreFile.Sep, "tms_self_other_kuai100.json"))
	if err != nil {
		err = errors.New("load city_data json file, " + err.Error())
		return
	}
	var jsonData dataTMSCompany
	err = json.Unmarshal(dataByte, &jsonData)
	if err != nil {
		return
	}
	globTMSCompanyData = jsonData.DataList
	return
}
