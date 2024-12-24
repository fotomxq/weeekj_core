package BaseLookup

import "errors"

// ArgsUpdateDomain 更新主题域参数
type ArgsUpdateDomain struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// UpdateDomain 更新主题域
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

// SetDomainOrGet 设置新的领域
// 如果领域名称重复，则直接反馈信息
func SetDomainOrGet(name string) (id int64, err error) {
	//检查是否存在
	data := GetDomainByName(name)
	if data.ID > 0 {
		id = data.ID
		return
	}
	//添加数据
	data, err = CreateDomain(&ArgsCreateDomain{
		Name: name,
	})
	if err != nil {
		return
	}
	id = data.ID
	if id < 1 {
		err = errors.New("id error")
		return
	}
	//反馈
	return
}
