package BaseQiniu

import (
	"encoding/json"
	BaseFileSys "gitee.com/weeekj/weeekj_core/v5/base/filesys"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

//七牛云反馈处理

// 根据数据集合，创建数据
// param data []byte 解析数据集合
// param ip string IP地址
// param from string 来源
// param fromID string 来源ID
// param userID string 用户ID
// param expireTime int64 过期时间
// return FieldsFileType
// return error
type ArgsCreateByReport struct {
	//数据包
	Data []byte
	//IP
	IP string
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//过期时间
	ExpireAt time.Time
	//文件类型
	FileType string
	//扩展信息
	ClaimInfos []CoreSQLConfig.FieldsConfigType
	//备注
	Des string
}

func CreateByReport(args *ArgsCreateByReport) (fileData BaseFileSys.FieldsFileClaimType, errCode string, err error) {
	//记录结果集
	var info reportCallBackDataType
	info, err = getCallBackObject(string(args.Data))
	if err != nil {
		return
	}
	//解析数据并创建文件块，等待认领
	infos := []CoreSQLConfig.FieldsConfigType{
		{
			Mark: "bucket",
			Val:  info.Bucket,
		},
	}
	fileData, _, errCode, err = BaseFileSys.Create(&BaseFileSys.ArgsCreate{
		CreateIP:   args.IP,
		CreateInfo: args.CreateInfo,
		FileSize:   info.FileSize,
		FileType:   args.FileType,
		FileHash:   info.Hash,
		FileSrc:    "",
		ExpireAt:   args.ExpireAt,
		FromInfo:   CoreSQLFrom.FieldsFrom{System: "qiniu", Mark: info.Key},
		Infos:      infos,
		ClaimInfos: args.ClaimInfos,
		Des:        args.Des,
	})
	return
}

// 反馈信息组结构
type reportCallBackDataType struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	FileSize int64  `json:"fsize"`
	Bucket   string `json:"bucket"`
	Name     string `json:"name"`
}

// 尝试解析数据
// {"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}
func getCallBackObject(data string) (reportCallBackDataType, error) {
	res := reportCallBackDataType{}
	err := json.Unmarshal([]byte(data), &res)
	return res, err
}
