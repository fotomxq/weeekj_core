package OrgSign

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

type ArgsGetSignAll struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgId" check:"id" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindId" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userId" check:"id" index:"true"`
}

// GetSignAll 获取所有签名
func GetSignAll(args *ArgsGetSignAll) (dataList []FieldsSign, err error) {
	_, err = signDB.GetList().GetAll(&BaseSQLTools.ArgsGetAll{
		ConditionFields: []BaseSQLTools.ArgsGetListSimpleConditionID{
			{
				Name: "org_id",
				Val:  args.OrgID,
			},
			{
				Name: "org_bind_id",
				Val:  args.OrgBindID,
			},
			{
				Name: "user_id",
				Val:  args.UserID,
			},
		},
		IsRemove: false,
	}, &dataList)
	if err != nil {
		return
	}
	return
}

// ArgsGetSignDefault 获取默认的签名参数
type ArgsGetSignDefault struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgId" check:"id" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindId" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userId" check:"id" index:"true"`
}

// GetSignDefault 获取默认的签名
func GetSignDefault(args *ArgsGetSignDefault) (data FieldsSign, err error) {
	err = signDB.GetInfo().GetInfoByFields(map[string]any{
		"org_id":      args.OrgID,
		"org_bind_id": args.OrgBindID,
		"user_id":     args.UserID,
		"is_default":  true,
	}, true, &data)
	if err != nil {
		return
	}
	return
}

// ArgsGetSignLastTemp 获取最后一次临时签字参数
type ArgsGetSignLastTemp struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgId" check:"id" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindId" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userId" check:"id" index:"true"`
}

// GetSignLastTemp 获取最后一次临时签字
func GetSignLastTemp(args *ArgsGetSignLastTemp) (data FieldsSign, err error) {
	//获取数据
	err = signDB.GetInfo().GetInfoByFields(map[string]any{
		"org_id":      args.OrgID,
		"org_bind_id": args.OrgBindID,
		"user_id":     args.UserID,
		"is_temp":     true,
	}, true, &data)
	if err != nil {
		return
	}
	//获取后自动删除
	err = DeleteSignByID(&ArgsDeleteSignByID{
		ID:        data.ID,
		OrgID:     data.OrgID,
		OrgBindID: data.OrgBindID,
		UserID:    data.UserID,
	})
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsGetSignLastTempAndDefault 获取最后一次临时签字并转化为默认参数
type ArgsGetSignLastTempAndDefault struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgId" check:"id" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindId" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userId" check:"id" index:"true"`
}

// GetSignLastTempAndDefault 获取最后一次临时签字并转化为默认
func GetSignLastTempAndDefault(args *ArgsGetSignLastTempAndDefault) (data FieldsSign, err error) {
	//获取数据
	err = signDB.GetInfo().GetInfoByFields(map[string]any{
		"org_id":      args.OrgID,
		"org_bind_id": args.OrgBindID,
		"user_id":     args.UserID,
		"is_temp":     true,
	}, true, &data)
	if err != nil {
		return
	}
	//转化数据
	type setData struct {
		//ID
		ID int64 `db:"id" json:"id" check:"id" index:"true"`
		//是否临时传递
		// 临时传递的签名将在使用后立即删除
		IsTemp bool `db:"is_temp" json:"isTemp" index:"true"`
	}
	_ = signDB.GetUpdate().UpdateByID(&setData{
		ID:     data.ID,
		IsTemp: false,
	})
	//设置默认
	updateDefaultSign(args.OrgID, args.OrgBindID, args.UserID)
	//反馈
	return
}

// GetSignByID 通过数据获取ID
func GetSignByID(id int64) (data FieldsSign, err error) {
	err = signDB.GetInfo().GetInfoByID(id, &data)
	if err != nil {
		return
	}
	return
}
