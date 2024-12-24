package BaseLookup

import "errors"

// ArgsUpdateLookup 更新编码参数
type ArgsUpdateLookup struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//领域ID
	DomainID int64 `db:"domain_id" json:"domainID" check:"id" empty:"true"`
	//管理单元
	UnitID int64 `db:"unit_id" json:"unitID" check:"id" empty:"true"`
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"100"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// UpdateLookup 更新编码
// 注意，主要创建过的编码，无论是否删除都会被占用无法使用
func UpdateLookup(args *ArgsUpdateLookup) (err error) {
	//同一个领域下，编码不可重复
	if args.DomainID > 0 {
		var data FieldsLookup
		_ = lookupDB.DB.GetPostgresql().Get(&data, "SELECT id FROM base_lookup_child WHERE code = $1", args.Code)
		if data.ID > 0 && data.ID != args.ID {
			err = errors.New("code repeat")
			return
		}
	}
	//更新数据
	err = lookupDB.Update().SetFields([]string{"domain_id", "unit_id", "code", "name"}).NeedUpdateTime().AddWhereID(args.ID).NeedSoft(true).NamedExec(map[string]any{
		"domain_id": args.DomainID,
		"unit_id":   args.UnitID,
		"code":      args.Code,
		"name":      args.Name,
	})
	if err != nil {
		return
	}
	deleteLookupCache(args.ID)
	return
}

// ArgsSetLookupList 批量设置一组编码参数
type ArgsSetLookupList struct {
	//领域ID
	DomainID int64 `db:"domain_id" json:"domainID" check:"id" empty:"true"`
	//管理单元
	UnitID int64 `db:"unit_id" json:"unitID" check:"id" empty:"true"`
	//数据列
	DataList []ArgsSetLookupListChild `json:"dataList"`
}

type ArgsSetLookupListChild struct {
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"100"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// SetLookupList 批量设置一组编码
func SetLookupList(args *ArgsSetLookupList) (err error) {
	allCode, _ := GetLookupAll(args.DomainID, args.UnitID)
	var readyCodes []string
	for _, v := range allCode {
		for _, v2 := range args.DataList {
			if v.Code == v2.Code {
				err = UpdateLookup(&ArgsUpdateLookup{
					ID:       v.ID,
					DomainID: args.DomainID,
					UnitID:   args.UnitID,
					Code:     v2.Code,
					Name:     v2.Name,
				})
				if err != nil {
					err = errors.New("update error, code: " + v.Code)
				}
			}
		}
		vData := GetLookupCode(v.Code)
		if vData.ID > 0 {
		}
		readyCodes = append(readyCodes, v.Code)
	}
	return
}
