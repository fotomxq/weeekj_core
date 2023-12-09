package ToolsAppUpdate

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"strings"
)

//APP升级服务模块

// ArgsGetUpdateLastVer 获取指定渠道的最新版本号参数
type ArgsGetUpdateLastVer struct {
	//运行环境
	// android_phone / android_pad / ios_phone / ios_ipad / windows / osx / linux
	// 或者特定品牌的定制
	System string `db:"system" json:"system"`
	//APP标识码
	AppMark string `db:"app_mark" json:"appMark"`
}

// GetUpdateLastVer 获取指定渠道的最新版本号
func GetUpdateLastVer(args *ArgsGetUpdateLastVer) (varData string, err error) {
	var data FieldsUpdate
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT tools_app_update.id as id, tools_app_update.ver as ver FROM tools_app_update, tools_app_update_app WHERE tools_app_update.system = $1 AND tools_app_update_app.app_mark = $2 AND tools_app_update_app.id = tools_app_update.app_id ORDER BY id DESC LIMIT 1", args.System, args.AppMark)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no update")
		return
	}
	for _, v := range data.Ver {
		if varData == "" {
			varData = fmt.Sprint(v)
		} else {
			varData = fmt.Sprint(varData, ".", v)
		}
	}
	return
}

// ArgsCheckUpdate 检查是否需要升级参数
type ArgsCheckUpdate struct {
	//运行环境
	// android_phone / android_pad / ios_phone / ios_ipad / windows / osx / linux
	// 或者特定品牌的定制
	System string `db:"system" json:"system"`
	//环境的最低版本
	// 如果给与指定专供版本，则该设定无效
	// [7, 1, 4] => version 7.1.4
	SystemVersion string `db:"system_version" json:"systemVersion"`
	//APP标识码
	AppMark string `db:"app_mark" json:"appMark"`
	//版本号
	// [7, 1, 4] => version 7.1.4
	Version string `db:"version" json:"version"`
}

// CheckUpdate 检查是否需要升级
func CheckUpdate(args *ArgsCheckUpdate) (data FieldsUpdate, needUpdate bool) {
	//systemVerStr := strings.Split(args.SystemVersion, ".")
	//var systemVer []int64
	//for k := 0; k < len(systemVerStr); k++ {
	//	v, err := CoreFilter.GetInt64ByString(systemVerStr[k])
	//	if err != nil {
	//		continue
	//	}
	//	systemVer = append(systemVer, v)
	//}
	//verStr := strings.Split(args.Version, ".")
	//var ver []int64
	//for k := 0; k < len(verStr); k++ {
	//	v, err := CoreFilter.GetInt64ByString(verStr[k])
	//	if err != nil {
	//		continue
	//	}
	//	ver = append(ver, v)
	//}
	//找到符合条件的app应用
	var appData FieldsApp
	if err := Router2SystemConfig.MainDB.Get(&appData, "SELECT id FROM tools_app_update_app WHERE app_mark = $1", args.AppMark); err != nil {
		return
	}
	//找到符合条件的更新
	var page int64 = 1
	for {
		if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, system, system_ver_min, system_ver_max, system_ver, file_id, download_url, app_size, app_md5, name, des, des_files, ver, ver_build, grayscale_res, params FROM tools_app_update WHERE system = $1 AND app_id = $2 AND is_skip = false ORDER BY id DESC LIMIT 1 OFFSET $3", args.System, appData.ID, page-1); err != nil || data.ID < 1 {
			break
		}
		//检查该数据的版本是否存在特殊指定？
		//if len(data.SystemVer) > 0 {
		//for _, v := range data.SystemVer {
		//	var vSInt64 []int64
		//	//拆解v，改为整数数据后进行对比
		//	vS := strings.Split(v, ".")
		//	for _, v2 := range vS {
		//		v2Int64, err := CoreFilter.GetInt64ByString(v2)
		//		if err != nil {
		//			v2Int64 = 0
		//		}
		//		vSInt64 = append(vSInt64, v2Int64)
		//	}
		//	if b := CoreFilter.CheckArrayEq(systemVer, vSInt64); b {
		//		needUpdate = true
		//		_ = appendCount(&argsAppendCount{
		//			OrgID:    data.OrgID,
		//			AppID:    data.AppID,
		//			UpdateID: data.ID,
		//		})
		//		return
		//	}
		//}
		//如果没有相等，则禁止更新
		//return
		//}
		//if len(data.SystemVerMin) > 0 {
		//	//如果左侧较小，则说明低于最低需求，需跳过该版本继续查询
		//	if b := CoreFilter.CheckArray(data.SystemVerMin, systemVer); b {
		//		page += 1
		//		continue
		//	}
		//}
		//if len(data.SystemVerMax) > 0 {
		//	//如果超出版本需求，则需继续查询
		//	if b := CoreFilter.CheckArray(data.SystemVerMax, systemVer); !b {
		//		page += 1
		//		continue
		//	}
		//}
		//检查该版本是否满足升级条件
		//if b := CoreFilter.CheckArrayEq(data.Ver, ver); b {
		//	//禁止升级，跳出
		//	return
		//}
		var vVerInt64 []string
		for k := 0; k < len(data.Ver); k++ {
			vVerInt64 = append(vVerInt64, CoreFilter.GetStringByInt64(data.Ver[k]))
		}
		vVer := strings.Join(vVerInt64, ".")
		if args.Version == vVer {
			//禁止升级，跳出
			return
		}
		break
	}
	//全部满足条件，则直接反馈数据
	needUpdate = true
	_ = appendCount(&argsAppendCount{
		OrgID:    data.OrgID,
		AppID:    data.AppID,
		UpdateID: data.ID,
	})
	return
}

// ArgsGetUpdateList 获取更新列表参数
type ArgsGetUpdateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//运行环境
	// android_phone / android_pad / ios_phone / ios_ipad / windows / osx / linux
	System string `db:"system" json:"system" check:"mark" empty:"true"`
	//APP ID
	AppID int64 `db:"app_id" json:"appID" check:"id" empty:"true"`
	//搜索
	Search string `db:"search" json:"search" check:"search" empty:"true"`
}

// GetUpdateList 获取更新列表
func GetUpdateList(args *ArgsGetUpdateList) (dataList []FieldsUpdate, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.System != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "system = :system"
		maps["system"] = args.System
	}
	if args.AppID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "app_id = :app_id"
		maps["app_id"] = args.AppID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "tools_app_update"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, name, des, des_files, system, system_ver_min, system_ver_max, system_ver, app_id, ver, ver_build, file_id, download_url, app_size, app_md5, grayscale_res, is_skip, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetUpdateID 获取升级版本参数
type ArgsGetUpdateID struct {
	//ID
	ID int64 `db:"id" json:"id"`
}

// GetUpdateID 获取升级版本
func GetUpdateID(args *ArgsGetUpdateID) (data FieldsUpdate, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, name, des, des_files, system, system_ver_min, system_ver_max, system_ver, app_id, ver, ver_build, file_id, download_url, app_size, app_md5, grayscale_res, is_skip, params FROM tools_app_update WHERE id = $1", args.ID)
	return
}

// ArgsCreateUpdate 创建新的版本参数
type ArgsCreateUpdate struct {
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//升级内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//描述附带文件
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//运行环境
	// android_phone / android_pad / ios_phone / ios_ipad / windows / osx / linux
	System string `db:"system" json:"system" check:"mark"`
	//环境的最低版本
	// 如果给与指定专供版本，则该设定无效
	// [7, 1, 4] => version 7.1.4
	SystemVerMin pq.Int64Array `db:"system_ver_min" json:"systemVerMin"`
	//环境的最高版本
	// 如果给与指定专供版本，则该设定无效
	// [7, 1, 4] => version 7.1.4
	SystemVerMax pq.Int64Array `db:"system_ver_max" json:"systemVerMax"`
	//专供版本
	// 该版本为专供特定环境的版本
	// {"7.1.4", "3.5.1"}
	SystemVer pq.StringArray `db:"system_ver" json:"systemVer"`
	//APP ID
	AppID int64 `db:"app_id" json:"appID" check:"id"`
	//版本号
	// [7, 1, 4] => version 7.1.4
	Ver pq.Int64Array `db:"ver" json:"ver"`
	//app构建编号
	VerBuild string `db:"ver_build" json:"verBuild"`
	//下载文件ID或URL地址
	FileID      int64  `db:"file_id" json:"fileID" check:"id" empty:"true"`
	DownloadURL string `db:"download_url" json:"downloadURL"`
	//文件大小
	AppSize int64 `db:"app_size" json:"appSize"`
	//文件MD5
	AppMD5 string `db:"app_md5" json:"appMD5"`
	//灰度发布
	// 如果>0，将根据上一个版本总数 / 灰度系数，随机抽中则推送，否则不推送
	GrayscaleRes bool `db:"grayscale_res" json:"grayscaleRes" check:"bool" empty:"true"`
	//是否跳过改版本？
	// 存在异常被标记后，禁止为用户推送该版本
	IsSkip bool `db:"is_skip" json:"isSkip" check:"bool" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateUpdate 创建新的版本
func CreateUpdate(args *ArgsCreateUpdate) (data FieldsUpdate, err error) {
	var appData FieldsApp
	err = Router2SystemConfig.MainDB.Get(&appData, "SELECT id FROM tools_app_update_app WHERE id = $1", args.AppID)
	if err != nil || appData.ID < 1 {
		err = errors.New(fmt.Sprint("app not exist, ", err))
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tools_app_update", "INSERT INTO tools_app_update (org_id, name, des, des_files, system, system_ver_min, system_ver_max, system_ver, app_id, ver, ver_build, file_id, download_url, app_size, app_md5, grayscale_res, is_skip, params) VALUES (:org_id,:name,:des,:des_files,:system,:system_ver_min,:system_ver_max,:system_ver,:app_id,:ver,:ver_build,:file_id,:download_url,:app_size,:app_md5,:grayscale_res,:is_skip,:params)", args, &data)
	return
}

// ArgsUpdateUpdate 修改版本信息参数
type ArgsUpdateUpdate struct {
	//APP ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//升级内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//描述附带文件
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//运行环境
	// android_phone / android_pad / ios_phone / ios_ipad / windows / osx / linux
	System string `db:"system" json:"system" check:"mark"`
	//环境的最低版本
	// 如果给与指定专供版本，则该设定无效
	// [7, 1, 4] => version 7.1.4
	SystemVerMin pq.Int64Array `db:"system_ver_min" json:"systemVerMin"`
	//环境的最高版本
	// 如果给与指定专供版本，则该设定无效
	// [7, 1, 4] => version 7.1.4
	SystemVerMax pq.Int64Array `db:"system_ver_max" json:"systemVerMax"`
	//专供版本
	// 该版本为专供特定环境的版本
	// {"7.1.4", "3.5.1"}
	SystemVer pq.StringArray `db:"system_ver" json:"systemVer"`
	//版本号
	// [7, 1, 4] => version 7.1.4
	Ver pq.Int64Array `db:"ver" json:"ver"`
	//app构建编号
	VerBuild string `db:"ver_build" json:"verBuild"`
	//下载文件ID或URL地址
	FileID      int64  `db:"file_id" json:"fileID" check:"id" empty:"true"`
	DownloadURL string `db:"download_url" json:"downloadURL"`
	//文件大小
	AppSize int64 `db:"app_size" json:"appSize"`
	//文件MD5
	AppMD5 string `db:"app_md5" json:"appMD5"`
	//灰度发布
	// 如果>0，将根据上一个版本总数 / 灰度系数，随机抽中则推送，否则不推送
	GrayscaleRes bool `db:"grayscale_res" json:"grayscaleRes" check:"bool" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateUpdate 修改版本信息
func UpdateUpdate(args *ArgsUpdateUpdate) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tools_app_update SET org_id = :org_id, name = :name, des = :des, des_files = :des_files, system = :system, system_ver_min = :system_ver_min, system_ver_max = :system_ver_max, system_ver = :system_ver, ver = :ver, ver_build = :ver_build, file_id = :file_id, download_url = :download_url, app_size = :app_size, app_md5 = :app_md5, grayscale_res = :grayscale_res, params = :params WHERE id = :id AND org_id = :org_id", args)
	return
}

// ArgsDeleteUpdate 删除版本参数
type ArgsDeleteUpdate struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteUpdate 删除版本
func DeleteUpdate(args *ArgsDeleteUpdate) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tools_app_update SET is_skip = true WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
