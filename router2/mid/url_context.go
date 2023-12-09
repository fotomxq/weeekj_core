package Router2Mid

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// GetContext 识别和获取头部上下文
func GetContext(context any) *gin.Context {
	publicC, ok := context.(*RouterURLPublicC)
	if ok {
		return publicC.Context
	} else {
		headerC, ok := context.(*RouterURLHeaderC)
		if ok {
			return headerC.Context
		} else {
			userC, ok := context.(*RouterURLUserC)
			if ok {
				return userC.Context
			} else {
				orgC, ok := context.(*RouterURLOrgC)
				if ok {
					return orgC.Context
				} else {
					roleC, ok := context.(*RouterURLRoleC)
					if ok {
						return roleC.Context
					} else {
						iotC, ok := context.(*RouterURLIOTC)
						if ok {
							return iotC.Context
						} else {
							//给与的上下文结构无法识别
							panic(fmt.Sprint("get context failed, ", context))
							return &gin.Context{}
						}
					}
				}
			}
		}
	}
}

// GetContextBodyByte 尝试获取上下文的body byte
// 由于设备验证和获取参数冲突，gin的body不能连续获取两次，会触发EOF错误
func GetContextBodyByte(context any) (dataByte []byte, b bool) {
	iotC, ok := context.(*RouterURLIOTC)
	if ok {
		dataByte = iotC.BodyByte
		b = true
		return
	} else {
		return
	}
}

// DataGetContextData 识别和获取头部上下文并带数据结果
type DataGetContextData struct {
	//日志头部
	LogAppend string
	//用户ID
	UserID int64
	//组织ID
	OrgID int64
	//组织成员ID
	OrgBindID int64
}

// GetContextData 识别和获取头部上下文并带数据
func GetContextData(context any) (c *gin.Context, result DataGetContextData) {
	publicC, ok := context.(*RouterURLPublicC)
	if ok {
		c = publicC.Context
		result.LogAppend = publicC.LogAppend
	} else {
		headerC, ok := context.(*RouterURLHeaderC)
		if ok {
			c = headerC.Context
			result.LogAppend = headerC.LogAppend
		} else {
			userC, ok := context.(*RouterURLUserC)
			if ok {
				c = userC.Context
				result.LogAppend = userC.LogAppend
				result.UserID = userC.UserID
			} else {
				orgC, ok := context.(*RouterURLOrgC)
				if ok {
					c = orgC.Context
					result.LogAppend = orgC.LogAppend
					result.UserID = orgC.UserID
					result.OrgID = orgC.OrgID
					result.OrgBindID = orgC.OrgBindID
				} else {
					roleC, ok := context.(*RouterURLRoleC)
					if ok {
						c = roleC.Context
						result.LogAppend = roleC.LogAppend
						result.UserID = roleC.UserID
					} else {
						iotC, ok := context.(*RouterURLIOTC)
						if ok {
							c = iotC.Context
							result.LogAppend = iotC.LogAppend
						} else {
							c = &gin.Context{}
							result.LogAppend = ""
							//给与的上下文结构无法识别
							panic(fmt.Sprint("get context failed, ", context))
						}
					}
				}
			}
		}
	}
	return
}
