package UserCore

import "github.com/lib/pq"

//FieldsGroupType 用户组结构表
type FieldsGroupType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 如果为空，则说明是平台的用户；否则为对应组织的用户
	// 所有获取的方法，都需要给与该ID参数，也可以留空，否则禁止获取
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//关联的权限组
	Permissions pq.StringArray `db:"permissions" json:"permissions"`
}
