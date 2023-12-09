package BaseQiniu

import (
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

//获取管理对象
func getManager() (*storage.BucketManager, error) {
	ak, sk, err := getKey()
	if err != nil {
		return bucketManager, err
	}
	mac := qbox.NewMac(ak, sk)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	return bucketManager, nil
}
