package BaseLookup

import "errors"

// ArgsCreateLookup 创建编码参数
type ArgsCreateLookup struct {
	//是否为系统预设
	IsSys bool `db:"is_sys" json:"isSys" check:"bool"`
	//领域ID
	DomainID int64 `db:"domain_id" json:"domainID" check:"id" empty:"true"`
	//管理单元
	UnitID int64 `db:"unit_id" json:"unitID" check:"id" empty:"true"`
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"100"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// CreateLookup 创建编码
// 注意，主要创建过的编码，无论是否删除都会被占用无法使用
func CreateLookup(args *ArgsCreateLookup) (data FieldsLookup, err error) {
	//同一个领域下，编码不可重复
	if args.DomainID > 0 {
		_ = lookupDB.DB.GetPostgresql().Get(&data, "SELECT id FROM base_lookup_child WHERE code = $1", args.Code)
		if data.ID > 0 {
			err = errors.New("code repeat")
			return
		}
	}
	//创建数据
	err = lookupDB.Insert().SetFields([]string{"is_sys", "domain_id", "unit_id", "code", "name"}).Add(map[string]interface{}{
		"is_sys":    args.IsSys,
		"domain_id": args.DomainID,
		"unit_id":   args.UnitID,
		"code":      args.Code,
		"name":      args.Name,
	}).ExecAndResultData(&data)
	if err != nil {
		return
	}
	return
}
