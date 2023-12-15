package BaseQiniu

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseFileSys "github.com/fotomxq/weeekj_core/v5/base/filesys"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreHttp "github.com/fotomxq/weeekj_core/v5/core/http"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"strings"
	"time"
)

// ArgsUploadBySrc 将本地文件上传处理，之后自动删除临时数据参数
type ArgsUploadBySrc struct {
	//文件本地路径
	Src string
	//存储块名称
	BucketName string
	//文件类型
	FileType string
	//IP地址
	IP string
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//创建用户
	UserID int64
	//创建组织
	OrgID int64
	//是否为公开的文件
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//过期时间
	ExpireAt time.Time
	//扩展信息
	ClaimInfos []CoreSQLConfig.FieldsConfigType
	//描述
	Des string
}

// UploadBySrc 将本地文件上传处理，之后自动删除临时数据
func UploadBySrc(args *ArgsUploadBySrc) (fileData BaseFileSys.FieldsFileClaimType, errCode string, err error) {
	var localFileData []byte
	localFileData, err = CoreFile.LoadFile(args.Src)
	if err != nil {
		err = errors.New("load upload file data, " + err.Error())
		errCode = "load_file"
		return
	}
	fileData, errCode, err = Upload(localFileData, args.BucketName, args.FileType, args.IP, args.CreateInfo, args.UserID, args.OrgID, args.IsPublic, args.ExpireAt, args.ClaimInfos, args.Des)
	if err2 := CoreFile.DeleteF(args.Src); err2 != nil {
		//无法删除该文件，不记录
		// 该操作可交给文件处理模块或其他维护模块，对临时文件进行定期清理处理
		// 临时文件空间应该为共享的空间，而不是完全封闭独立的空间
	}
	return
}

// ArgsUploadByURL 通过网络文件推送数据参数
type ArgsUploadByURL struct {
	//远程文件URL
	URL string
	//存储块名称
	BucketName string
	//文件类型
	FileType string
	//IP地址
	IP string
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//创建用户
	UserID int64
	//创建组织
	OrgID int64
	//是否为公开的文件
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//过期时间
	ExpireAt time.Time
	//扩展信息
	ClaimInfos []CoreSQLConfig.FieldsConfigType
	//描述
	Des string
}

// UploadByURL 通过网络文件推送数据
func UploadByURL(args *ArgsUploadByURL) (fileData BaseFileSys.FieldsFileClaimType, errCode string, err error) {
	//加载远程文件
	var fileByte []byte
	fileByte, err = CoreHttp.GetData(args.URL, nil, "", false)
	if err != nil {
		errCode = "http_failed"
		return
	}
	//如果文件格式为空，则自动解析
	if args.FileType == "" {
		urlNames := CoreFilter.GetURLNameType(args.URL)
		args.FileType = urlNames["type"]
	}
	//推送上传
	fileData, errCode, err = Upload(fileByte, args.BucketName, args.FileType, args.IP, args.CreateInfo, args.UserID, args.OrgID, args.IsPublic, args.ExpireAt, args.ClaimInfos, args.Des)
	return
}

// Upload 上传新的文件
func Upload(data []byte, bucketName string, fileType string, ip string, createInfo CoreSQLFrom.FieldsFrom, userID, orgID int64, isPublic bool, expireAt time.Time, claimInfos []CoreSQLConfig.FieldsConfigType, des string) (fileData BaseFileSys.FieldsFileClaimType, errCode string, err error) {
	//获取上传凭证
	var upToken string
	upToken, bucketName, errCode, err = GetCertSimple(bucketName, 0)
	if err != nil {
		return
	}
	//上传文件数据
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuabei
	// 找到配置
	configBuckets, err := BaseConfig.GetDataString("FileQiNiuBucketList")
	if err != nil {
		errCode = "qiniu_bucket_config"
		err = errors.New("get config, " + err.Error())
		return
	}
	configBucketsArr := strings.Split(configBuckets, "|")
	var configZones string
	configZones, err = BaseConfig.GetDataString("FileQiNiuZones")
	if err == nil {
		configZonesArr := strings.Split(configZones, "|")
		isFind := false
		for k, v := range configBucketsArr {
			if v != bucketName {
				continue
			}
			for k2, v2 := range configZonesArr {
				if k == k2 {
					isFind = true
					switch v2 {
					case "ZoneHuabei":
						cfg.Zone = &storage.ZoneHuabei
					case "ZoneHuanan":
						cfg.Zone = &storage.ZoneHuanan
					case "ZoneHuadong":
						cfg.Zone = &storage.ZoneHuadong
					case "ZoneBeimei":
						cfg.Zone = &storage.ZoneBeimei
					case "ZoneXinjiapo":
						cfg.Zone = &storage.ZoneXinjiapo
					case "Zone_as0":
						cfg.Zone = &storage.Zone_as0
					case "Zone_na0":
						cfg.Zone = &storage.Zone_na0
					case "Zone_z0":
						cfg.Zone = &storage.Zone_z0
					case "Zone_z1":
						cfg.Zone = &storage.Zone_z1
					case "Zone_z2":
						cfg.Zone = &storage.Zone_z2
					}
					break
				}
			}
			if isFind {
				break
			}
		}
	}
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	//开始上传文件
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			//"x:name": "github logo",
		},
	}
	dataLen := int64(len(data))
	var newFileName string
	newFileName, err = getFileName(data, fileType)
	if err != nil {
		errCode = "file_info"
		err = errors.New("get file info of names, " + err.Error())
		return
	}
	err = formUploader.Put(context.Background(), &ret, upToken, newFileName, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		errCode = "qiniu_put"
		err = errors.New("upload qiniu file, formUploader put file is error, bucket: " + bucketName + ", " + err.Error())
		return
	}
	if ret.Hash == "" {
		ret.Hash, err = CoreFilter.GetSha256Str(string(data))
		if err != nil {
			errCode = "hash_256"
			err = errors.New("get file hash 256, " + err.Error())
			return
		}
	}
	if ret.Hash == "" || ret.Key == "" {
		err = errors.New("upload qiniu file, ret hash or key is empty")
		errCode = "qiniu_info"
		return
	}
	//上传成功后，获取反馈数据
	infos := []CoreSQLConfig.FieldsConfigType{
		{
			Mark: "bucket",
			Val:  bucketName,
		},
	}
	fileData, _, errCode, err = BaseFileSys.Create(&BaseFileSys.ArgsCreate{
		CreateIP:   ip,
		CreateInfo: createInfo,
		UserID:     userID,
		OrgID:      orgID,
		IsPublic:   isPublic,
		FileSize:   dataLen,
		FileType:   fileType,
		FileHash:   ret.Hash,
		FileSrc:    "",
		ExpireAt:   expireAt,
		FromInfo:   CoreSQLFrom.FieldsFrom{System: "qiniu", Mark: ret.Key},
		Infos:      infos,
		ClaimInfos: claimInfos,
		Des:        des,
	})
	return
}

// GetCertSimple 获取上传凭证
func GetCertSimple(bucketName string, waitID int64) (upToken string, newBucketName, errCode string, err error) {
	if bucketName == "" {
		bucketName, err = BaseConfig.GetDataString("FileQiniuDefaultBucketName")
		if err != nil {
			errCode = "config_bucket_default"
			return
		}
	}
	//获取对应空间的上传凭证
	upToken, err = getCertificateToServer(bucketName, waitID)
	if err != nil {
		errCode = "qiniu_token"
		err = errors.New("upload qiniu file, " + err.Error())
		return
	}
	newBucketName = bucketName
	return
}

// getCertificateToServer 获取上传凭证
// 服务端专用
func getCertificateToServer(bucketName string, waitID int64) (string, error) {
	ak, sk, err := getKey()
	if err != nil {
		return "", err
	}
	callBackURL, err := getCallBack()
	if err != nil {
		return "", err
	}
	putPolicy := storage.PutPolicy{
		//设备所属的文件领域
		Scope: bucketName,
		//过期时间 1小时
		Expires: 3600,
		//回调反馈数据
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)","mimeType":"$(mimeType)","ext":"${ext}","tokenid":` + fmt.Sprint(waitID) + `}`,
	}
	//如果启动回调服务器
	certificateServerOn, err := getCallbackOn()
	if err != nil {
		return "", err
	}
	if certificateServerOn && waitID > 0 {
		putPolicy.CallbackURL = callBackURL
		putPolicy.CallbackBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)","mimeType":"$(mimeType)","ext":"$(ext)","tokenid":` + fmt.Sprint(waitID) + `}`
		putPolicy.CallbackBodyType = "application/json"
		putPolicy.SaveKey = fmt.Sprint(CoreFilter.GetNowTime().Format("2006-01-02_15-04-05"), "_", "$(etag)$(ext)")
	}
	mac := qbox.NewMac(ak, sk)
	upToken := putPolicy.UploadToken(mac)
	return upToken, nil
}
