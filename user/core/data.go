package UserCore

// Deprecated: 准备废弃
// DataUserDataType 解决一体化工具，用于反馈用户的整合信息结构
type DataUserDataType struct {
	//用户基本信息结构
	Info FieldsUserType `json:"info"`
	//用户组信息
	Groups []FieldsGroupType `json:"groups"`
	//权限组
	Permissions []string `json:"permissions"`
	//文件列
	FileList map[int64]string `json:"fileList"`
}
