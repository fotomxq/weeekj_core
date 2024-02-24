package BaseLookup

type ArgsCreateDomain struct {
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

func CreateDomain(args *ArgsCreateDomain) (data FieldsDomain, err error) {
	//创建数据
	err = domainDB.Insert().SetFields([]string{"name"}).Add(map[string]interface{}{
		"name": args.Name,
	}).ExecAndResultData(&data)
	if err != nil {
		return
	}
	return
}
