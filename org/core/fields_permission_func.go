package OrgCoreCore

import "github.com/lib/pq"

//FieldsPermissionFunc 权限业务分组
type FieldsPermissionFunc struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//所需业务
	ParentMarks pq.StringArray `db:"parent_marks" json:"parentMarks"`
}
