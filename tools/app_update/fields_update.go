package ToolsAppUpdate

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsUpdate APP推送和发布表
type FieldsUpdate struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name" json:"name"`
	//升级内容
	Des string `db:"des" json:"des"`
	//描述附带文件
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//运行环境
	// android_phone / android_pad / ios_phone / ios_ipad / windows / osx / linux
	System string `db:"system" json:"system"`
	//环境的最低版本
	// 如果给与指定专供版本，则该设定无效
	// [7, 1, 4] => version 7.1.4
	// [0]则不限制
	SystemVerMin pq.Int64Array `db:"system_ver_min" json:"systemVerMin"`
	//环境的最高版本
	// 如果给与指定专供版本，则该设定无效
	// [7, 1, 4] => version 7.1.4
	// [0]则不限制
	SystemVerMax pq.Int64Array `db:"system_ver_max" json:"systemVerMax"`
	//专供版本
	// 该版本为专供特定环境的版本
	// {"7.1.4", "3.5.1"}
	SystemVer pq.StringArray `db:"system_ver" json:"systemVer"`
	//APP ID
	AppID int64 `db:"app_id" json:"appID"`
	//版本号
	// [7, 1, 4] => version 7.1.4
	Ver pq.Int64Array `db:"ver" json:"ver"`
	//app构建编号
	VerBuild string `db:"ver_build" json:"verBuild"`
	//下载文件ID或URL地址
	FileID      int64  `db:"file_id" json:"fileID"`
	DownloadURL string `db:"download_url" json:"downloadURL"`
	//文件大小
	AppSize int64 `db:"app_size" json:"appSize"`
	//文件MD5
	AppMD5 string `db:"app_md5" json:"appMD5"`
	//灰度发布
	// 如果>0，将根据上一个版本总数 / 灰度系数，随机抽中则推送，否则不推送
	GrayscaleRes bool `db:"grayscale_res" json:"grayscaleRes"`
	//是否跳过改版本？
	// 存在异常被标记后，禁止为用户推送该版本
	IsSkip bool `db:"is_skip" json:"isSkip"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
