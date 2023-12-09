package UserCore

//FieldsPermissionType 权限表结构
type FieldsPermissionType struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//是否允许组织授权
	AllowOrg bool `db:"allow_org" json:"allowOrg"`
}