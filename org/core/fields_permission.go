package OrgCoreCore

//FieldsPermission 组织权限
type FieldsPermission struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//分组标识码
	FuncMark string `db:"func_mark" json:"funcMark"`
	//名称
	Name string `db:"name" json:"name"`
}
