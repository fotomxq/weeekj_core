package BaseFileSys

//检查一组认领文件是否属于此人
type ArgsCheckFileListAndUserID struct {
	//要验证的引用文件ID列
	ClaimList []int64
	//用户ID
	UserID int64
}

func CheckFileListAndUserID(args *ArgsCheckFileListAndUserID) bool {
	for _, v := range args.ClaimList {
		if b := CheckFileIDAndUserID(&ArgsCheckFileIDAndUserID{
			ClaimID: v,
			UserID:  args.UserID,
		}); !b {
			return false
		}
	}
	return true
}

//检查认领ID是否属于此人
type ArgsCheckFileIDAndUserID struct {
	//引用文件ID
	ClaimID int64
	//用户ID
	UserID int64
}

func CheckFileIDAndUserID(args *ArgsCheckFileIDAndUserID) bool {
	return CheckFileIDAndCreateInfo(&ArgsCheckFileIDAndCreateInfo{
		ClaimID: args.ClaimID,
		UserID:  args.UserID,
		OrgID:   0,
	})
}

//检查认领ID列是否属于此组织
type ArgsCheckFileIDsAndOrgID struct {
	//文件引用序列
	ClaimIDs []int64 `json:"claimIDs"`
	//组织ID
	OrgID int64 `json:"orgID"`
}

func CheckFileIDsAndOrgID(args *ArgsCheckFileIDsAndOrgID) (b bool) {
	for _, v := range args.ClaimIDs {
		_, err := GetFileClaimByID(&ArgsGetFileClaimByID{
			ClaimID: v,
			UserID:  0,
			OrgID:   args.OrgID,
		})
		if err != nil {
			return
		}
	}
	b = true
	return
}

//检查认领ID是否属于此来源
type ArgsCheckFileIDsAndCreateInfo struct {
	//引用文件ID
	ClaimIDs []int64
	//用户ID
	// 可选，用于检测
	UserID int64 `json:"userID"`
	//组织ID
	// 可选，用于检测
	OrgID int64 `json:"orgID"`
}

func CheckFileIDsAndCreateInfo(args *ArgsCheckFileIDsAndCreateInfo) (b bool) {
	for _, v := range args.ClaimIDs {
		b = CheckFileIDAndCreateInfo(&ArgsCheckFileIDAndCreateInfo{
			ClaimID: v,
			UserID:  args.UserID,
			OrgID:   args.OrgID,
		})
		if !b {
			return
		}
	}
	b = true
	return
}

//检查认领ID是否属于来源
type ArgsCheckFileIDAndCreateInfo struct {
	//引用文件ID
	ClaimID int64
	//用户ID
	// 可选，用于检测
	UserID int64 `json:"userID"`
	//组织ID
	// 可选，用于检测
	OrgID int64 `json:"orgID"`
}

func CheckFileIDAndCreateInfo(args *ArgsCheckFileIDAndCreateInfo) (b bool) {
	data, err := GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: args.ClaimID,
		UserID:  args.UserID,
		OrgID:   args.OrgID,
	})
	if err != nil {
		return false
	}
	if data.ID < 1 {
		return false
	}
	return true
}
