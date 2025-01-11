package OrgSign

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

// ArgsDeleteSignByID 删除签名参数
type ArgsDeleteSignByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" index:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgId" check:"id" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindId" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userId" check:"id" index:"true"`
}

// DeleteSignByID 删除签名
func DeleteSignByID(args *ArgsDeleteSignByID) (err error) {
	//获取数据
	var data FieldsSign
	data, err = GetSignByID(args.ID)
	if err != nil {
		return
	}
	//检查权限
	if !CoreFilter.EqID2(args.OrgID, args.OrgID) || !CoreFilter.EqID2(args.OrgBindID, data.OrgBindID) || !CoreFilter.EqID2(args.UserID, data.UserID) {
		err = errors.New("org_id not match")
		return
	}
	//删除数据
	err = signDB.GetDelete().DeleteByID(args.ID)
	if err != nil {
		return
	}
	//修正默认值数据
	updateDefaultSign(data.OrgID, data.OrgBindID, data.UserID)
	//反馈
	return
}
