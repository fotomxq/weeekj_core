package OrgSign

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

// ArgsCreateSign 创建签名参数
type ArgsCreateSign struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgId" check:"id" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindId" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userId" check:"id" index:"true"`
	//是否默认
	// 一个客体可拥有多个签名，但只能有一个默认签名
	IsDefault bool `db:"is_default" json:"isDefault" index:"true"`
	//签名类型
	// base64 Base64文件; file 文件系统ID
	SignType string `db:"sign_type" json:"signType" check:"des" min:"1" max:"50" index:"true"`
	//是否临时传递
	// 临时传递的签名将在使用后立即删除
	IsTemp bool `db:"is_temp" json:"isTemp" index:"true"`
	//签名数据
	SignData string `db:"sign_data" json:"signData" check:"des" min:"1" max:"-1" empty:"true"`
	//文件ID
	FileID int64 `db:"file_id" json:"fileID" check:"id" index:"true"`
}

// CreateSign 创建签名
func CreateSign(args *ArgsCreateSign) (errCode string, err error) {
	//检查数据是否正常
	if args.FileID < 1 && args.SignData == "" {
		errCode = "report_params_lost"
		err = errors.New("file id or sign data is required")
		return
	}
	//检查是否重复
	if args.FileID > 0 || args.SignData != "" {
		var b bool
		b, err = signDB.GetInfo().CheckInfoByFields(map[string]any{
			"file_id":   args.FileID,
			"sign_data": args.SignData,
		}, true)
		if err == nil && b {
			errCode = "err_have_replace"
			err = errors.New("sign data is exists")
			return
		}
		err = nil
	}
	//添加数据
	var newID int64
	newID, err = signDB.GetInsert().InsertRow(args)
	if err != nil {
		errCode = "err_insert"
		err = errors.New("create sign error: " + err.Error())
		return
	}
	//如果设置此内容为默认，则更新其他所有内容为非
	if args.IsDefault {
		ctx := signDB.GetUpdate().UpdatePrev()
		ctx = ctx.SetWhereStr("org_id = :org_id AND org_bind_id = :org_bind_id AND user_id = :user_id AND id != :id", map[string]any{
			"org_id":      args.OrgID,
			"org_bind_id": args.OrgBindID,
			"user_id":     args.UserID,
			"id":          newID,
		})
		ctx = ctx.SetFields([]string{"is_default"})
		_ = signDB.GetUpdate().UpdateDo(ctx, map[string]any{
			"is_default": false,
		})
	} else {
		//修正默认值数据
		updateDefaultSign(args.OrgID, args.OrgBindID, args.UserID)
	}
	//如果为临时传递
	if args.IsTemp {
		ctx := signDB.GetUpdate().UpdatePrev()
		ctx = ctx.SetWhereStr("org_id = :org_id AND org_bind_id = :org_bind_id AND user_id = :user_id AND id != :id", map[string]any{
			"org_id":      args.OrgID,
			"org_bind_id": args.OrgBindID,
			"user_id":     args.UserID,
			"id":          newID,
		})
		ctx = ctx.SetFields([]string{"delete_at"})
		_ = signDB.GetUpdate().UpdateDo(ctx, map[string]any{
			"delete_at": CoreFilter.GetNowTime(),
		})
	}
	//反馈
	return
}
