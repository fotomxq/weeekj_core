package IOTMQTT

import (
	"encoding/json"
	"errors"
	"fmt"
	BaseQiniu "github.com/fotomxq/weeekj_core/v5/base/qiniu"
	ToolsAppUpdate "github.com/fotomxq/weeekj_core/v5/tools/app_update"
	"strings"
)

// PushAppUpdate 推送设备更新数据
func PushAppUpdate(groupMark, deviceCode string, appMark string, data ToolsAppUpdate.FieldsUpdate) (err error) {
	//重组数据
	type dataType struct {
		//APP下载地址
		DownloadURL string `json:"downloadURL"`
		//运行环境
		// android_phone / android_pad / ios_phone / ios_ipad / windows / osx / linux
		System string `db:"system" json:"system"`
		//环境的最低版本
		// 如果给与指定专供版本，则该设定无效
		// [7, 1, 4] => version 7.1.4
		// [0]则不限制
		SystemVerMin string `json:"systemVerMin"`
		//环境的最高版本
		// 如果给与指定专供版本，则该设定无效
		// [7, 1, 4] => version 7.1.4
		// [0]则不限制
		SystemVerMax string `json:"systemVerMax"`
		//版本号
		Ver string `json:"ver"`
		//app构建编号
		VerBuild string `json:"verBuild"`
		//APP大小
		// 字节
		AppSize int64 `json:"appSize"`
		//文件MD5摘要值
		MD5 string `json:"md5"`
		//标题
		Name string `json:"name"`
		//介绍文字
		Des string `json:"des"`
		//介绍文件列
		DesFiles []string `json:"desFiles"`
	}
	var newData dataType
	newData.DownloadURL = data.DownloadURL
	newData.System = data.System
	var systemVerMin []string
	for i := 0; i < len(data.SystemVerMax); i++ {
		systemVerMin = append(systemVerMin, fmt.Sprint(data.Ver[i]))
	}
	newData.SystemVerMin = strings.Join(systemVerMin, ".")
	var systemVerMax []string
	for i := 0; i < len(data.SystemVerMax); i++ {
		systemVerMax = append(systemVerMax, fmt.Sprint(data.Ver[i]))
	}
	newData.SystemVerMax = strings.Join(systemVerMax, ".")
	var vers []string
	for i := 0; i < len(data.Ver); i++ {
		vers = append(vers, fmt.Sprint(data.Ver[i]))
	}
	newData.Ver = strings.Join(vers, ".")
	newData.VerBuild = data.VerBuild
	newData.AppSize = data.AppSize
	newData.MD5 = data.AppMD5
	newData.Name = data.Name
	newData.Des = data.Des
	//如果下载文件不存在，则通过下载文件获取数据
	var fileIDs []int64
	if data.DownloadURL == "" && data.FileID > 0 {
		if data.FileID > 0 {
			fileIDs = append(fileIDs, data.FileID)
		}
	}
	if len(data.DesFiles) > 0 {
		for _, v := range data.DesFiles {
			fileIDs = append(fileIDs, v)
		}
	}
	fileURLs, _ := BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
		ClaimIDList: fileIDs,
		UserID:      0,
		OrgID:       0,
		IsPublic:    true,
	})
	if data.DownloadURL == "" && data.FileID > 0 {
		for k, v := range fileURLs {
			if k == data.FileID {
				newData.DownloadURL = v
				break
			}
		}
	}
	for i := 0; i < len(data.DesFiles); i++ {
		for k, v := range fileURLs {
			if k == data.DesFiles[i] {
				newData.DesFiles = append(newData.DesFiles, v)
			}
		}
	}
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(newData)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	if groupMark == "" || deviceCode == "" {
		topic := fmt.Sprint("app/update/system/", data.System, "/", appMark)
		err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	} else {
		topic := fmt.Sprint("app/update/group/", groupMark, "/code/", deviceCode)
		err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	}
	return
}
