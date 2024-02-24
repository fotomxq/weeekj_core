package BaseLookup

type ArgsUpdateDomain struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

func UpdateDomain(args *ArgsUpdateDomain) (err error) {
	//更新数据
	err = domainDB.Update().SetFields([]string{"name"}).NeedUpdateTime().AddWhereID(args.ID).NeedSoft(true).NamedExec(map[string]any{
		"name": args.Name,
	})
	if err != nil {
		return
	}
	deleteDomainCache(args.ID)
	return
}
