package BaseTempFile

import "time"

type FieldsFile struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//文件参数值
	// 可用于文件缓存，避免重复生成文件的负担
	FileParams string `db:"file_params" json:"fileParams"`
	//文件路径
	FileSrc string `db:"file_src" json:"fileSrc"`
	//文件SHA1
	FileSHA1 string `db:"file_sha1" json:"fileSHA1"`
	//文件名称
	Name string `db:"name" json:"name"`
	//文件类型
	FileType string `db:"file_type" json:"fileType"`
}
