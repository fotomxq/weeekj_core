package Router2Mid

import (
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

// RouterURLIOTC 普通级别头部
type RouterURLIOTC struct {
	//上下文
	Context *gin.Context
	//日志头部
	LogAppend string
	//上下文参数数据集合
	BodyByte []byte
}

func (t *RouterURLIOT) GET(urlPath string, handle func(*RouterURLIOTC)) {
	//映射URL
	t.BaseData.Routers.GET(urlPath, func(c *gin.Context) {
		getURLIOT(c, handle)
	})
}

func (t *RouterURLIOT) POST(urlPath string, handle func(*RouterURLIOTC)) {
	//映射URL
	t.BaseData.Routers.POST(urlPath, func(c *gin.Context) {
		getURLIOT(c, handle)
	})
}

func (t *RouterURLIOT) PUT(urlPath string, handle func(*RouterURLIOTC)) {
	//映射URL
	t.BaseData.Routers.PUT(urlPath, func(c *gin.Context) {
		getURLIOT(c, handle)
	})
}

func (t *RouterURLIOT) DELETE(urlPath string, handle func(*RouterURLIOTC)) {
	//映射URL
	t.BaseData.Routers.DELETE(urlPath, func(c *gin.Context) {
		getURLIOT(c, handle)
	})
}

// getURLIOT 方法集合处理封装
func getURLIOT(c *gin.Context, handle func(*RouterURLIOTC)) {
	//检查设备权限
	bodyByte, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, iotBanURL)
		c.Abort()
		return
	}
	args := ArgsIOTData{
		GroupMark: gjson.GetBytes(bodyByte, "keys.groupMark").String(),
		Code:      gjson.GetBytes(bodyByte, "keys.code").String(),
		NowTime:   gjson.GetBytes(bodyByte, "keys.nowTime").Int(),
		Rand:      gjson.GetBytes(bodyByte, "keys.rand").String(),
		Key:       gjson.GetBytes(bodyByte, "keys.key").String(),
		OrgID:     gjson.GetBytes(bodyByte, "keys.orgID").Int(),
	}
	if b := checkDeviceAndOrg(c, &args); !b {
		c.Redirect(http.StatusMovedPermanently, iotBanURL)
		c.Abort()
		return
	}
	//反馈
	handle(&RouterURLIOTC{
		Context:   c,
		LogAppend: "",
		//由于设备验证和获取参数冲突，gin的body不能连续获取两次，会触发EOF错误
		BodyByte: bodyByte,
	})
}
