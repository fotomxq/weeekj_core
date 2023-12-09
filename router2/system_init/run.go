package Router2SystemInit

import (
	"fmt"
	BaseLog "gitee.com/weeekj/weeekj_core/v5/base/log"
	BaseOtherCheck "gitee.com/weeekj/weeekj_core/v5/base/other_check"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	RouterGinSet "gitee.com/weeekj/weeekj_core/v5/router/gin_set"
	RouterSystem "gitee.com/weeekj/weeekj_core/v5/router/system"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Run 核心run
// 请阻塞本函数，不要使用go运行，否则应用会跳出
func Run(urlHandle func()) {
	///////////////////////////////////////////////////////////////////////////////////
	//系统底层
	// 必须引用的关键组件
	///////////////////////////////////////////////////////////////////////////////////
	//日志
	go CoreLog.Run()
	go BaseLog.Run()

	//启动debug
	if Router2SystemConfig.Debug {
		go runDebug()
	}

	///////////////////////////////////////////////////////////////////////////////////
	//路由
	///////////////////////////////////////////////////////////////////////////////////
	//构建gin通用基础
	RouterGinSet.PageDefaultRoot()
	//启动静态路由默认设置
	if OpenAPIDefaultSet {
		//映射路由
		apiRouter := RouterGinSet.Router
		//找不到页面操作设计
		apiRouter.NoRoute(func(c *gin.Context) {
			//动态验证模块
			data, err := BaseOtherCheck.GetURL(&BaseOtherCheck.ArgsGetURL{
				URL: c.Request.RequestURI,
			})
			if err == nil {
				c.String(http.StatusOK, data)
				c.Abort()
				return
			}
			//其他页面自动反馈
			c.String(http.StatusNotFound, "")
			c.Abort()
			return
		})
	}
	//启动URL模块
	urlHandle()
	//启动服务
	go runServer()
	//系统阻塞和关闭服务
	RouterSystem.Close()
}

// runServer 启动服务
func runServer() {
	closeSysSrc := fmt.Sprint(CoreFile.BaseSrc, CoreFile.Sep, "close_sys")
	if b := RouterGinSet.RunServer(); !b {
		//无法启动服务
		//创建close_sys文件，主程序发现后将在1秒内关闭整个服务
		if err := CoreFile.WriteFile(closeSysSrc, []byte{}); err != nil {
			CoreLog.Error("create close sys pls file failed, ", err)
			return
		}
		//等待几秒后再次尝试
		time.Sleep(time.Second * 4)
		if b := RouterGinSet.RunServer(); !b {
			//还是失败了
		}
	}
}
