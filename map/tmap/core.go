package MapTMap

import (
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	"io/ioutil"
	"net/http"
)

//天地图相关服务设计
// 注意天地图坐标系为：CGCS2000，该坐标和WGS84一致。前端需要手动进行转化操作。具体转化方法可使用chatGPT生成。
// 模块内提供了方法：WGS84ToGcj02可以将坐标转为国内常用的坐标系：GCJ02

// getKey 获取key
func getKey() (key string) {
	return BaseConfig.GetDataStringNoErr("MapTMapWEBKey")
}

// postURL 底层服务相应处理模块
// 将自动追加密钥数据参数
func postURL(url string) (bodyData []byte, err error) {
	//追加key
	url = fmt.Sprint(url, "&tk=", getKey())
	//替换为您的实际API URL以及请求参数
	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("User-Agent", "TianDiTu-Go-Client")
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	bodyData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//反馈
	return
}
