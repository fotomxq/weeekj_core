package OrgCoreCore

import (
	"errors"
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetBindAttrAll 获取指定成员的信息参数
type ArgsGetBindAttrAll struct {
	//所属企业ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true" index:"true"`
}

// GetBindAttrAll 获取指定成员的信息
func GetBindAttrAll(args *ArgsGetBindAttrAll) (dataList []FieldsBindAttr, err error) {
	_, err = orgBindAttrDB.Select().SetDeleteQuery("delete_at", false).SetIDQuery("org_id", args.OrgID).SetIDQuery("user_id", args.UserID).SetIDQuery("org_bind_id", args.OrgBindID).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  9999,
		Sort: "id",
		Desc: false,
	}).ResultAndCount(&dataList)
	if err != nil {
		return
	}
	for k, v := range dataList {
		dataList[k] = getBindAttr(v.ID)
	}
	return
}

// ArgsGetBindAttr 获取指定的信息参数
type ArgsGetBindAttr struct {
	//所属企业ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true" index:"true"`
	//标识码
	AttrCode string `db:"attr_code" json:"attrCode" index:"true"`
}

// GetBindAttr 获取指定的信息
func GetBindAttr(args *ArgsGetBindAttr) (data FieldsBindAttr, err error) {
	if args.UserID < 1 && args.OrgBindID < 1 {
		err = errors.New("args error")
		return
	}
	if args.UserID > 0 {
		err = orgBindAttrDB.Get().SetIDQuery("org_id", args.OrgID).SetIDQuery("user_id", args.UserID).SetStringQuery("attr_code", args.AttrCode).NeedLimit().Result(&data)
	} else {
		err = orgBindAttrDB.Get().SetIDQuery("org_id", args.OrgID).SetIDQuery("org_bind_id", args.OrgBindID).SetStringQuery("attr_code", args.AttrCode).NeedLimit().Result(&data)
	}
	if err != nil {
		return
	}
	data = getBindAttr(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsSetBindAttrs 批量修改成员信息参数
type ArgsSetBindAttrs struct {
	//所属企业ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true" index:"true"`
	//信息列表
	Attrs []ArgsSetBindAttrsChild `json:"attrs"`
}

type ArgsSetBindAttrsChild struct {
	//标识码
	AttrCode string `db:"attr_code" json:"attrCode" index:"true"`
	//值
	AttrValue string `db:"attr_value" json:"attrValue"`
	//整数
	AttrInt int64 `db:"attr_int" json:"attrInt"`
	//浮点数
	AttrFloat float64 `db:"attr_float" json:"attrFloat"`
}

// SetBindAttrs 批量修改成员信息
func SetBindAttrs(args *ArgsSetBindAttrs) (err error) {
	for _, v := range args.Attrs {
		err = SetBindAttr(&ArgsSetBindAttr{
			OrgID:     args.OrgID,
			UserID:    args.UserID,
			OrgBindID: args.OrgBindID,
			AttrCode:  v.AttrCode,
			AttrValue: v.AttrValue,
			AttrInt:   v.AttrInt,
			AttrFloat: v.AttrFloat,
		})
		if err != nil {
			return
		}
	}
	return
}

// ArgsSetBindAttr 设置指定成员信息参数
type ArgsSetBindAttr struct {
	//所属企业ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true" index:"true"`
	//标识码
	AttrCode string `db:"attr_code" json:"attrCode" index:"true"`
	//值
	AttrValue string `db:"attr_value" json:"attrValue"`
	//整数
	AttrInt int64 `db:"attr_int" json:"attrInt"`
	//浮点数
	AttrFloat float64 `db:"attr_float" json:"attrFloat"`
}

// SetBindAttr 设置指定成员信息
func SetBindAttr(args *ArgsSetBindAttr) (err error) {
	var data FieldsBindAttr
	data, err = GetBindAttr(&ArgsGetBindAttr{
		OrgID:     args.OrgID,
		UserID:    args.UserID,
		OrgBindID: args.OrgBindID,
		AttrCode:  args.AttrCode,
	})
	if err != nil {
		err = orgBindAttrDB.Insert().SetDefaultInsertFields().Add(map[string]any{
			"org_id":      args.OrgID,
			"user_id":     args.UserID,
			"org_bind_id": args.OrgBindID,
			"attr_code":   args.AttrCode,
			"attr_value":  args.AttrValue,
			"attr_int":    args.AttrInt,
			"attr_float":  args.AttrFloat,
		}).ExecAndCheckID()
	} else {
		err = orgBindAttrDB.Update().SetFields([]string{"attr_value", "attr_int", "attr_float"}).AddWhereID(data.ID).NamedExec(map[string]any{
			"attr_value": args.AttrValue,
			"attr_int":   args.AttrInt,
			"attr_float": args.AttrFloat,
		})
	}
	if err != nil {
		return
	}
	deleteBindAttrCache(data.ID)
	return
}

// ArgsClearBindAttr 清空成员信息参数
type ArgsClearBindAttr struct {
	//所属企业ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true" index:"true"`
}

// ClearBindAttr 清空成员信息
func ClearBindAttr(args *ArgsClearBindAttr) (err error) {
	var dataList []FieldsBindAttr
	dataList, err = GetBindAttrAll(&ArgsGetBindAttrAll{})
	if err != nil {
		err = nil
	}
	if args.UserID > 0 {
		err = orgBindAttrDB.Delete().NeedSoft(true).AddWhereOrgID(args.OrgID).AddWhereUserID(args.UserID).ExecNamed(nil)
	} else {
		err = orgBindAttrDB.Delete().NeedSoft(true).AddWhereOrgID(args.OrgID).SetWhereAnd("org_bind_id", args.OrgBindID).ExecNamed(nil)
	}
	if err != nil {
		return
	}
	for _, v := range dataList {
		deleteBindAttrCache(v.ID)
	}
	return
}

func getBindAttr(id int64) (data FieldsBindAttr) {
	cacheMark := getBindAttrCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := orgBindAttrDB.Get().SetDefaultFields().GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	return
}

// 缓冲
func getBindAttrCacheMark(id int64) string {
	return fmt.Sprint("org:bind:attr:id.", id)
}

func deleteBindAttrCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBindAttrCacheMark(id))
}
