package BaseLookup

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

func UpdateLookup(args *ArgsUpdateLookup) (err error) {
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
