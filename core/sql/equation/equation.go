package CoreSQLEquation

//方程式设计
// 该设计将写入sql的同时，可自行对内容进行检测
// 提供模版后，可写入数据并进行验证，反馈结果
// 该方程式设计主要面向传统对等、范围内容的判断处理

type FieldsEquations []FieldsEquation

type FieldsEquation struct {
	//等式ID
	// 该ID由前端拟定，确保唯一即可
	ID int `db:"id" json:"id"`
	//上级ID
	// 嵌套关系处理
	ParentID int `db:"parent_id" json:"parentID"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//等式类型
	// 0 > 大于; 1 >= 大于等于; 2 = 等于; 3 < 小于; 4 <= 小于等于
	// 5 ? 查询; 6 ?? 忽略大小写查询
	Eq int `db:"eq" json:"eq"`
	//条件处理方式
	// 0 + 相加 / 1 - 相减 / 2 x 相乘 / 3 "/" 相除
	Conditions int `db:"conditions" json:"conditions"`
	//值
	Val string `db:"val" json:"val"`
}
