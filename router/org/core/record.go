package RouterOrgCore

import (
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	OrgRecord "gitee.com/weeekj/weeekj_core/v5/org/record"
	RouterMidOrg "gitee.com/weeekj/weeekj_core/v5/router/mid/org"
	RouterUserRecord "gitee.com/weeekj/weeekj_core/v5/router/user/record"
	"github.com/gin-gonic/gin"
	"strings"
)

// AddRecord 行政日志行为统一处理
// 注意必须具备上下文关系，否则请提前自行创建
// createInfo 创建日志模块来源，Mark: 系统内部的行为识别编码，例如任务终端中的修改目标状态，也可以用于其他的识别编码
func AddRecord(c *gin.Context, content ...interface{}) {
	//获取绑定数据
	bindData := RouterMidOrg.GetOrgBindData(c)
	//生成mark
	mark := strings.Replace(c.Request.URL.String(), "/", "_", -1)
	mark = strings.Replace(mark, ":", "", -1)
	//重组消息
	var contentMsg string
	for _, v := range content {
		contentMsg = contentMsg + fmt.Sprint(v)
	}
	//创建行政日志
	if err := OrgRecord.Create(&OrgRecord.ArgsCreate{
		OrgID:       bindData.OrgID,
		BindID:      bindData.ID,
		ContentMark: mark,
		Content:     contentMsg,
	}); err != nil {
		CoreLog.Error("create work record, ", err)
	}
	//创建用户日志
	RouterUserRecord.CreateByC(c, content)
}
