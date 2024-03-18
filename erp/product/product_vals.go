package ERPProduct

// DataProductVal 标准化参数和反馈结构
type DataProductVal struct {
	//顺序序号
	OrderNum int64 `db:"order_num" json:"orderNum"`
	//插槽值
	SlotID int64 `db:"slot_id" json:"slotID" check:"id"`
	//值(字符串)
	DataValue string `db:"data_value" json:"dataValue"`
	//值(浮点数)
	DataValueNum float64 `db:"data_value_num" json:"dataValueNum"`
	//值(整数)
	DataValueInt int64 `db:"data_value_int" json:"dataValueInt"`
	//参数
	Params string `db:"params" json:"params"`
}

// ArgsGetProductVals 获取产品预设模板值参数
type ArgsGetProductVals struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
}

// GetProductVals 获取产品预设模板值
// 根据产品关联的分类、品牌，获取产品预设模板数据包
// 如已经存在值，则反馈具体值；否则反馈默认数据包
func GetProductVals(args *ArgsGetProductVals) (data []DataProductVal, err error) {
	return
}

// ArgsSetProductVals 设置产品数据参数
type ArgsSetProductVals struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//数据结构
	Vals []DataProductVal `db:"vals" json:"vals"`
}

// SetProductVals 设置产品数据
func SetProductVals(args *ArgsSetProductVals) (err error) {
	return
}

type ArgsClearProductVals struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// ClearProductVals 清空产品数据
func ClearProductVals(args *ArgsClearProductVals) (err error) {
	return
}
