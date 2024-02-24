package BaseUnit

import "errors"

type ArgsCreateUnit struct {
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"100"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

func CreateUnit(args *ArgsCreateUnit) (data FieldsUnit, err error) {
	//确保唯一性
	data = GetUnitByCode(args.Code)
	if data.ID > 0 {
		data = FieldsUnit{}
		err = errors.New("code replace")
		return
	}
	//创建数据
	err = unitDB.Insert().SetFields([]string{"code", "name"}).Add(map[string]interface{}{
		"code": args.Code,
		"name": args.Name,
	}).ExecAndResultData(&data)
	if err != nil {
		return
	}
	return
}
