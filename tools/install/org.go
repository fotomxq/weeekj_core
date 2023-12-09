package ToolsInstall

import (
	"errors"
	"fmt"
	ClassConfig "gitee.com/weeekj/weeekj_core/v5/class/config"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
)

// InstallOrg 安装组织数据包
func InstallOrg() (err error) {
	//配置文件名称
	configConfigDefaultFileName := fmt.Sprint("org", CoreFile.Sep, "config_default.json")
	//检查配置文件是否存在
	if checkConfigFile(configConfigDefaultFileName) {
		//声明结构
		type dataInstallValueType struct {
			//标识码
			Mark string `db:"mark" json:"mark"`
			//名称
			Name string `db:"name" json:"name"`
			//是否可以公开
			AllowPublic bool `db:"allow_public" json:"allowPublic"`
			//是否允许组织查看该配置
			AllowSelfView bool `db:"allow_self_view" json:"allowSelfView"`
			//是否允许组织自己修改
			AllowSelfSet bool `db:"allow_self_set" json:"allowSelfSet"`
			//结构
			// 0 string / 1 bool / 2 int / 3 int64 / 4 float64
			// 结构也可用于前端判定某个特殊的样式，如时间样式、过期时间样式等，程序内不做任何限定，只是标记
			ValueType string `db:"value_type" json:"valueType"`
			//正则表达式
			ValueCheck string `db:"value_check" json:"valueCheck"`
			//默认值
			ValueDefault string `db:"value_default" json:"valueDefault"`
		}
		type dataInstallConfig struct {
			//增量配置表
			DataList []dataInstallValueType `json:"dataList"`
		}
		//获取文件数据
		var dataConfig dataInstallConfig
		if err := loadConfigFile(configConfigDefaultFileName, &dataConfig); err != nil {
			return nil
		}
		//处理数据
		for _, v := range dataConfig.DataList {
			var valueType int
			switch v.ValueType {
			case "string":
				valueType = 0
			case "bool":
				valueType = 1
			case "int":
				valueType = 2
			case "int64":
				valueType = 3
			case "float64":
				valueType = 4
			}
			if err = OrgCore.Config.Default.SetConfigDefault(&ClassConfig.ArgsSetConfigDefault{
				Mark:          v.Mark,
				Name:          v.Name,
				AllowPublic:   v.AllowPublic,
				AllowSelfSet:  v.AllowSelfSet,
				AllowSelfView: v.AllowSelfView,
				ValueType:     valueType,
				ValueCheck:    v.ValueCheck,
				ValueDefault:  v.ValueDefault,
			}); err != nil {
				err = errors.New(fmt.Sprint("安装组织配置失败, ", err, ", args: ", v))
				return
			}
		}
	}
	//第二代商户配置处理机制
	// 读取org_config文件下所有文件，写入处理
	newConfigDirSrc := fmt.Sprint(configDir, "org_config")
	configList, _ := CoreFile.GetFileList(newConfigDirSrc, []string{"json"}, false)
	for _, vConfigFileName := range configList {
		//声明结构
		type dataInstallValueType struct {
			//标识码
			Mark string `db:"mark" json:"mark"`
			//名称
			Name string `db:"name" json:"name"`
			//是否可以公开
			AllowPublic bool `db:"allow_public" json:"allowPublic"`
			//是否允许组织查看该配置
			AllowSelfView bool `db:"allow_self_view" json:"allowSelfView"`
			//是否允许组织自己修改
			AllowSelfSet bool `db:"allow_self_set" json:"allowSelfSet"`
			//结构
			// 0 string / 1 bool / 2 int / 3 int64 / 4 float64
			// 结构也可用于前端判定某个特殊的样式，如时间样式、过期时间样式等，程序内不做任何限定，只是标记
			ValueType string `db:"value_type" json:"valueType"`
			//正则表达式
			ValueCheck string `db:"value_check" json:"valueCheck"`
			//默认值
			ValueDefault string `db:"value_default" json:"valueDefault"`
		}
		type dataInstallConfig struct {
			//增量配置表
			DataList []dataInstallValueType `json:"dataList"`
		}
		//获取文件数据
		var dataConfig dataInstallConfig
		if err := loadConfigFile(fmt.Sprint("org_config", CoreFile.Sep, vConfigFileName), &dataConfig); err != nil {
			continue
		}
		//处理数据
		for _, v := range dataConfig.DataList {
			var valueType int
			switch v.ValueType {
			case "string":
				valueType = 0
			case "bool":
				valueType = 1
			case "int":
				valueType = 2
			case "int64":
				valueType = 3
			case "float64":
				valueType = 4
			}
			if err = OrgCore.Config.Default.SetConfigDefault(&ClassConfig.ArgsSetConfigDefault{
				Mark:          v.Mark,
				Name:          v.Name,
				AllowPublic:   v.AllowPublic,
				AllowSelfSet:  v.AllowSelfSet,
				AllowSelfView: v.AllowSelfView,
				ValueType:     valueType,
				ValueCheck:    v.ValueCheck,
				ValueDefault:  v.ValueDefault,
			}); err != nil {
				err = errors.New(fmt.Sprint("安装组织配置失败, ", err, ", args: ", v))
				return
			}
		}
	}
	//配置文件名称
	configPermissionFuncFileName := fmt.Sprint("org", CoreFile.Sep, "permission_func_mark.json")
	//检查配置文件是否存在
	if checkConfigFile(configPermissionFuncFileName) {
		//安装业务声明
		//声明结构
		type dataInstallPermissionFunc struct {
			//标识码
			Mark string `db:"mark" json:"mark"`
			//名称
			Name string `db:"name" json:"name"`
			//描述
			Des string `db:"des" json:"des"`
			//所需业务
			ParentMarks []string `db:"parent_marks" json:"parentMarks"`
		}
		type dataInstallPermissionFuncList struct {
			//增量配置表
			DataList []dataInstallPermissionFunc `json:"dataList"`
		}
		//获取文件数据
		var permissionFuncConfig dataInstallPermissionFuncList
		if err := loadConfigFile(configPermissionFuncFileName, &permissionFuncConfig); err != nil {
			return err
		}
		for _, v := range permissionFuncConfig.DataList {
			if err = OrgCore.SetPermissionFunc(&OrgCore.ArgsSetPermissionFunc{
				Mark:        v.Mark,
				Name:        v.Name,
				Des:         v.Des,
				ParentMarks: v.ParentMarks,
			}); err != nil {
				err = errors.New("安装组织业务失败, " + err.Error())
				return
			}
		}
	}
	//配置文件名称
	configPermissionFileName := fmt.Sprint("org", CoreFile.Sep, "permission.json")
	//检查配置文件是否存在
	if checkConfigFile(configPermissionFileName) {
		//安装权限
		type DataPermission struct {
			//标识码
			Mark string `json:"mark"`
			//名称
			Name string `json:"name"`
			//模块标识码
			FuncMark string `json:"funcMark"`
		}
		type dataInstallPermission struct {
			//增量配置表
			DataList []DataPermission `json:"dataList"`
		}
		//获取文件数据
		var dataPermission dataInstallPermission
		if err := loadConfigFile(configPermissionFileName, &dataPermission); err != nil {
			return err
		}
		for _, v := range dataPermission.DataList {
			if err = OrgCore.SetPermission(&OrgCore.ArgsSetPermission{
				Mark:     v.Mark,
				FuncMark: v.FuncMark,
				Name:     v.Name,
			}); err != nil {
				err = errors.New("安装组织权限失败, " + err.Error())
				return
			}
		}
	}
	return
}
