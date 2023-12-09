package BaseQiniu

import (
	"errors"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

//ArgsGetCertificateToClient 获取一个上传凭证，给客户端使用参数
type ArgsGetCertificateToClient struct {
	//桶名称
	BucketName string
}

//GetCertificateToClient 获取一个上传凭证，给客户端使用
// 该方法为简化方案，可另外创建新的方法，
// 	注意，请在本模块内创建方法集合，不要在外面另外创建，以避免管理混乱
//param bucketName string 桶配置的标识码
//return string 临时凭证，有效期1小时
//return error 错误代码
func GetCertificateToClient(args *ArgsGetCertificateToClient) (string, error) {
	if !checkBucketName( args.BucketName) {
		return "", errors.New("bucket is not exist")
	}
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
		Scope: args.BucketName,
		//过期时间 1小时
		Expires: 3600,
		//回调反馈数据
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	//如果启动回调服务器
	certificateServerOn, err := getCallbackOn()
	if err != nil {
		return "", err
	}
	if certificateServerOn {
		putPolicy.CallbackURL = callBackURL
		putPolicy.CallbackBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`
		putPolicy.CallbackBodyType = "application/json"
	}
	mac := qbox.NewMac(ak, sk)
	upToken := putPolicy.UploadToken(mac)
	return upToken, nil
}
