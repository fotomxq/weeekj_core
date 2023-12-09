package OrgCoreCore

//获取一组组织的最高一级别的一组
// 从一组组织ID中，抽取出最高层级的部分，并反馈
type ArgsGetTopOrg struct {
	//一组组织ID
	OrgIDList []int64
}

func GetTopOrg(args *ArgsGetTopOrg) []int64 {
	//少于2个，直接反馈原始数据即可
	if len(args.OrgIDList) < 2 {
		return args.OrgIDList
	}
	//最高级别数据
	var topDataList []FieldsOrg
	//遍历数据
	for _, vID := range args.OrgIDList {
		if vID < 1 {
			continue
		}
		vOrgData, err := GetOrg(&ArgsGetOrg{
			ID: vID,
		})
		if err != nil {
			continue
		}
		//在topDataList里面查询所有parentID，如果存在符合的，则剔除
		//如果vOrgData.Parent不存在，则标记并加入
		isFind := false
		var newTopDataList []FieldsOrg
		for _, v2 := range topDataList {
			if v2.ParentID == vID {
				continue
			}
			if vOrgData.ParentID == v2.ID {
				isFind = true
			}
			newTopDataList = append(newTopDataList, v2)
		}
		if !isFind {
			newTopDataList = append(newTopDataList, vOrgData)
		}
		topDataList = newTopDataList
	}
	//构建反馈数据集合
	var res []int64
	for _, v := range topDataList {
		res = append(res, v.ID)
	}
	return res
}

//获取一组组织的最低级别
// 从一组组织ID中，抽取出最低级别的部分，并反馈
type ArgsGetLowOrg struct {
	//一组组织ID
	OrgIDList []int64
}

func GetLowOrg(args *ArgsGetLowOrg) []int64 {
	//少于2个，直接反馈原始数据即可
	if len(args.OrgIDList) < 2 {
		return args.OrgIDList
	}
	//最低级别数据集合
	var lowDataList []FieldsOrg
	//遍历数据
	for _, vID := range args.OrgIDList {
		if vID < 1 {
			continue
		}
		vOrgData, err := GetOrg(&ArgsGetOrg{
			ID: vID,
		})
		if err != nil {
			continue
		}
		//在lowDataList里面查询所有id，如果当前ID的上级ID存在，则剔除
		//如果vOrgData.Parent不存在，则标记并加入
		isFind := false
		var newLowDataList []FieldsOrg
		for _, v2 := range lowDataList {
			if v2.ID == vOrgData.ParentID {
				continue
			}
			if vOrgData.ID == v2.ParentID {
				isFind = true
			}
			newLowDataList = append(newLowDataList, v2)
		}
		if !isFind {
			newLowDataList = append(newLowDataList, vOrgData)
		}
		lowDataList = newLowDataList
	}
	//构建反馈数据集合
	var res []int64
	for _, v := range lowDataList {
		res = append(res, v.ID)
	}
	return res
}