package OrgSign

import "errors"

// ArgsSelectSignDefault 更换默认选择的签名参数
type ArgsSelectSignDefault struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" index:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgId" check:"id" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindId" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userId" check:"id" index:"true"`
}

// SelectSignDefault 更换默认选择的签名
func SelectSignDefault(args *ArgsSelectSignDefault) (err error) {
	var defaultData FieldsSign
	defaultData, err = GetSignDefault(&ArgsGetSignDefault{
		OrgID:     args.OrgID,
		OrgBindID: args.OrgBindID,
		UserID:    args.UserID,
	})
	if err == nil && defaultData.ID > 0 {
		if defaultData.ID != args.ID {
			type setNoDefaultType struct {
				//ID
				ID int64 `db:"id" json:"id" check:"id" index:"true"`
				//是否默认
				// 一个客体可拥有多个签名，但只能有一个默认签名
				IsDefault bool `db:"is_default" json:"isDefault" index:"true"`
			}
			err = signDB.GetUpdate().UpdateByID(&setNoDefaultType{
				ID:        defaultData.ID,
				IsDefault: false,
			})
			if err != nil {
				return
			}
		}
	} else {
		var data FieldsSign
		data, err = GetSignByID(args.ID)
		if err != nil {
			err = errors.New("no data found")
			return
		}
		args.OrgID = data.OrgID
		args.OrgBindID = data.OrgBindID
		args.UserID = data.UserID
	}
	type setType struct {
		//ID
		ID int64 `db:"id" json:"id" check:"id" index:"true"`
		//组织ID
		OrgID int64 `db:"org_id" json:"orgId" check:"id" index:"true"`
		//组织成员ID
		OrgBindID int64 `db:"org_bind_id" json:"orgBindId" check:"id" index:"true"`
		//用户ID
		UserID int64 `db:"user_id" json:"userId" check:"id" index:"true"`
		//是否默认
		// 一个客体可拥有多个签名，但只能有一个默认签名
		IsDefault bool `db:"is_default" json:"isDefault" index:"true"`
	}
	err = signDB.GetUpdate().UpdateByIDAndCheckOrgBindUser(&setType{
		ID:        args.ID,
		OrgID:     args.OrgID,
		OrgBindID: args.OrgBindID,
		UserID:    args.UserID,
		IsDefault: true,
	})
	if err != nil {
		return
	}
	return
}

// 自动修正默认签名
// 如果没有默认时才会生效
func updateDefaultSign(orgID int64, orgBindID int64, userID int64) {
	//获取所有数据
	dataList, err := GetSignAll(&ArgsGetSignAll{
		OrgID:     orgID,
		OrgBindID: orgBindID,
		UserID:    userID,
	})
	if err == nil && len(dataList) > 0 {
		//如果存在数据，且没有其他默认时，则更换一个作为默认签名
		haveDefault := false
		for _, vData := range dataList {
			if vData.IsDefault {
				haveDefault = true
				break
			}
		}
		if !haveDefault {
			_ = SelectSignDefault(&ArgsSelectSignDefault{
				ID:        dataList[0].ID,
				OrgID:     orgID,
				OrgBindID: orgBindID,
				UserID:    userID,
			})
		}
	}
}
