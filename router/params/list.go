package RouterParams

import "github.com/gin-gonic/gin"

type ParamsListType struct {
	Page   int64  `json:"page" check:"page"`
	Max    int64  `json:"max" check:"max"`
	Sort   string `json:"sort" check:"sort"`
	Desc   bool   `json:"desc" check:"desc"`
	Search string `json:"search" check:"search" empty:"true"`
}

func GetDataList(c *gin.Context) (ParamsListType, bool) {
	//获取参数
	var params ParamsListType
	if b := GetJSON(c, &params); !b {
		return ParamsListType{}, false
	}
	//反馈
	return params, true
}
