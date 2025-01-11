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

// GetSignByID 通过数据获取ID
func GetSignByID(id int64) (data FieldsSign, err error) {
	err = signDB.GetInfo().GetInfoByID(id, &data)
	if err != nil {
		return
	}
	return
}
