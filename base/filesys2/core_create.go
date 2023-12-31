package BaseFileSys2

import CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"

type argsCreateCore struct {
	//创建IP
	CreateIP string `db:"create_ip" json:"createIP"`
	//文件原始创建人
	//创建组织
	// 可选，指定后该文件归属于组织，用户ID将只是指引，没有操作权限
	OrgID int64 `db:"org_id" json:"orgID"`
	//创建用户
	// 必须指定创建的用户，如果组织失效，则文件将自动归属于用户
	UserID int64 `db:"user_id" json:"userID"`
	//文件尺寸
	FileSize int64 `db:"file_size" json:"fileSize"`
	//文件类型
	FileType string `db:"file_type" json:"fileType"`
	//文件hash
	// 默认采用sha256作为标准
	FileHash string `db:"file_hash" json:"fileHash"`
	//文件路径
	FileSrc string `db:"file_src" json:"fileSrc"`
	//存储方式
	// local 本地化单一服务器存储; qiniu 七牛云存储
	SaveSystem string `db:"save_system" json:"saveSystem"`
	//存储块
	SaveMark string `db:"save_mark" json:"saveMark"`
	//第三方服务是否确认
	SaveSuccess bool `db:"save_success" json:"saveSuccess"`
	//其他扩展信息
	Infos CoreSQLConfig.FieldsConfigsType `db:"infos" json:"infos"`
}

func createCore(args *argsCreateCore) (err error) {
	err = coreDB.Insert().SetFields([]string{"create_ip", "org_id", "user_id", "file_size", "file_type", "file_hash", "file_src", "save_system", "save_mark", "save_success", "infos"}).ExecAndCheckID()
	if err != nil {
		return
	}
	return
}
