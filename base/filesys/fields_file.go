package BaseFileSys

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

// Deprecated: 建议采用BaseFileSys2
// FieldsFileType 文件主要结构体
type FieldsFileType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//验证Hash
	UpdateHash string `db:"update_hash" json:"updateHash"`
	//创建IP
	CreateIP string `db:"create_ip" json:"createIP"`
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//文件尺寸
	FileSize int64 `db:"file_size" json:"fileSize"`
	//文件类型
	FileType string `db:"file_type" json:"fileType"`
	//文件hash
	// 默认采用sha256作为标准
	FileHash string `db:"file_hash" json:"fileHash"`
	//文件路径
	FileSrc string `db:"file_src" json:"fileSrc"`
	//文件来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//其他扩展信息
	Infos CoreSQLConfig.FieldsConfigsType `db:"infos" json:"infos"`
}
