package RouterSystem

import (
	"encoding/json"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreReg "github.com/fotomxq/weeekj_core/v5/core/reg"
	RouterGinSet "github.com/fotomxq/weeekj_core/v5/router/gin_set"
	RouterParams "github.com/fotomxq/weeekj_core/v5/router/params"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
)

type DataSystemRegType struct {
	Code string `json:"code"`
	Key  string `json:"key"`
}

type DataReportCodeType struct {
	Code string `json:"code"`
}

var (
	//Router 路由
	Router *gin.Engine
	//RegFileSrc 文件路径
	RegFileSrc = "reg.json"
)

// Reg reg处理
func Reg(appName string) bool {
	//初始化注册机
	CoreReg.Init(appName)
	//遍历目录进行验证
	if b, err := regVerify(); err != nil {
		//加载文件数据失败
		CoreLog.Error("reg system, ", err)
		return false
	} else {
		//验证失败，启动URL服务
		// 该部分将中断后续处理，函数将无法关闭
		if !b {
			return systemRegURL()
		}
		return true
	}
}

// 验证工作
func regVerify() (bool, error) {
	var err error
	//修正路径位置
	regDir := fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep, "conf", CoreFile.Sep, "reg")
	RegFileSrc = fmt.Sprint(regDir, CoreFile.Sep, "main.json")
	fileNames, err := CoreFile.GetFileList(regDir, []string{"json"}, true)
	if err != nil {
		return false, err
	}
	//获取本机序列
	localCode, err := CoreReg.GetCode()
	if err != nil {
		return false, err
	}
	for _, v := range fileNames {
		//获取reg.json文件
		dataByte, err := CoreFile.LoadFile(v)
		if err != nil {
			continue
		}
		var data DataSystemRegType
		if err = json.Unmarshal(dataByte, &data); err != nil {
			continue
		}
		//检查序列是否匹配
		if data.Code != localCode {
			continue
		}
		//进行验证
		if !CoreReg.Verify(data.Key) {
			continue
		}
		//全部完成
		return true, nil
	}
	//如果还在继续，说明文件不存在或序列号全部不匹配
	//检查文件是否存在
	if !CoreFile.IsFile(RegFileSrc) {
		//建立文件
		newData, err := json.Marshal(DataSystemRegType{})
		if err != nil {
			return false, err
		}
		if err = CoreFile.WriteFile(RegFileSrc, newData); err != nil {
			return false, err
		}
	}
	//获取reg.json文件
	dataByte, err := CoreFile.LoadFile(RegFileSrc)
	if err != nil {
		return false, err
	}
	//解析数据
	data := DataSystemRegType{}
	err = json.Unmarshal(dataByte, &data)
	if err != nil {
		return false, err
	}
	data.Code = localCode
	dataByte, err = json.Marshal(data)
	if err != nil {
		return false, err
	}
	if err := CoreFile.WriteFile(RegFileSrc, dataByte); err != nil {
		return false, err
	}
	//进行验证
	if !CoreReg.Verify(data.Key) {
		return false, nil
	}
	//反馈数据
	return true, nil
}

// URL服务
func systemRegURL() bool {
	//重定向Router
	Router = RouterGinSet.Router
	//构建全局指向URL
	// url: /
	Router.GET("/", func(c *gin.Context) {
		RouterReport.BaseError(c, "no_reg", "注册失效")
	})
	//提供本机序列号
	// action: /system/reg
	Router.GET("/reg/code", func(c *gin.Context) {
		localCode, err := CoreReg.GetCode()
		if err != nil {
			RouterReport.BaseError(c, "reg_error", "获取失败")
			return
		}
		data := DataReportCodeType{
			localCode,
		}
		RouterReport.BaseData(c, data)
	})
	//尝试修改注册序列号
	Router.POST("/reg", func(c *gin.Context) {
		type DataType struct {
			Key string `json:"key" filter:"Mark"`
		}
		params := DataType{}
		if b := RouterParams.GetJSON(c, &params); !b {
			return
		}
		localCode, err := CoreReg.GetCode()
		if err != nil {
			RouterReport.BaseError(c, "reg_error", "无法获取本机序列号")
			return
		}
		data := DataSystemRegType{
			localCode,
			params.Key,
		}
		dataByte, err := json.Marshal(data)
		if err != nil {
			RouterReport.BaseError(c, "save_error", "文件读取异常")
			return
		}
		//重新写入数据
		if err := CoreFile.WriteFile(RegFileSrc, dataByte); err != nil {
			RouterReport.BaseError(c, "save_error", "文件写入异常")
			return
		}
		//反馈成功
		RouterReport.BaseSuccess(c)
	})
	//启动gin
	return RouterGinSet.RunServer()
}
