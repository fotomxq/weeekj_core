package Router2Params

type ParamsListType struct {
	Page   int64  `json:"page" check:"page"`
	Max    int64  `json:"max" check:"max"`
	Sort   string `json:"sort" check:"sort"`
	Desc   bool   `json:"desc" check:"desc"`
	Search string `json:"search" check:"search" empty:"true"`
}

func GetDataList(context any) (ParamsListType, bool) {
	c := getContext(context)
	//获取参数
	var params ParamsListType
	if b := GetJSON(c, &params); !b {
		return ParamsListType{}, false
	}
	//反馈
	return params, true
}
