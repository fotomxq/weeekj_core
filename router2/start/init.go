package Router2Core

import (
	"errors"
	"fmt"
	AnalysisAny "github.com/fotomxq/weeekj_core/v5/analysis/any"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseIPAddr "github.com/fotomxq/weeekj_core/v5/base/ipaddr"
	BaseTempFile "github.com/fotomxq/weeekj_core/v5/base/temp_file"
	BaseVcode "github.com/fotomxq/weeekj_core/v5/base/vcode"
	BlogCore "github.com/fotomxq/weeekj_core/v5/blog/core"
	BlogExam "github.com/fotomxq/weeekj_core/v5/blog/exam"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	IOTSensor "github.com/fotomxq/weeekj_core/v5/iot/sensor"
	MallCore "github.com/fotomxq/weeekj_core/v5/mall/core"
	MapCityData "github.com/fotomxq/weeekj_core/v5/map/city_data"
	MapRoom "github.com/fotomxq/weeekj_core/v5/map/room"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	OrgSubscription "github.com/fotomxq/weeekj_core/v5/org/subscription"
	RouterAPIRunBase "github.com/fotomxq/weeekj_core/v5/router/api/run_base"
	ServiceHousekeeping "github.com/fotomxq/weeekj_core/v5/service/housekeeping"
	ServiceInfoExchange "github.com/fotomxq/weeekj_core/v5/service/info_exchange"
	ServiceOrder "github.com/fotomxq/weeekj_core/v5/service/order"
	ServiceUserInfo "github.com/fotomxq/weeekj_core/v5/service/user_info"
	ServiceUserInfoCost "github.com/fotomxq/weeekj_core/v5/service/user_info_cost"
	TMSTransport "github.com/fotomxq/weeekj_core/v5/tms/transport"
	ToolsCommunication "github.com/fotomxq/weeekj_core/v5/tools/communication"
	UserChat "github.com/fotomxq/weeekj_core/v5/user/chat"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserLogin "github.com/fotomxq/weeekj_core/v5/user/login"
	UserMessage "github.com/fotomxq/weeekj_core/v5/user/message"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
)

func moduleInit() (err error) {
	//配置模块
	if err = BaseConfig.Init(); err != nil {
		err = errors.New("base config, " + err.Error())
		err = nil
		//return
	}
	if err = BaseIPAddr.Init(); err != nil {
		err = errors.New("base ipaddr, " + err.Error())
		return
	}
	//图形验证码
	if err = BaseVcode.Init(); err != nil {
		err = errors.New("base vcode, " + err.Error())
		return
	}
	//临时文件
	BaseTempFile.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//工具
	///////////////////////////////////////////////////////////////////////////////////
	//通讯
	ToolsCommunication.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//地图模块
	///////////////////////////////////////////////////////////////////////////////////
	//城市模块
	if err = MapCityData.Init(); err != nil {
		err = errors.New("map city data, " + err.Error())
		return
	}

	///////////////////////////////////////////////////////////////////////////////////
	//用户服务
	///////////////////////////////////////////////////////////////////////////////////
	//登陆服务
	UserLogin.Init(fmt.Sprint(AppName, AppVersion))

	///////////////////////////////////////////////////////////////////////////////////
	//用户
	///////////////////////////////////////////////////////////////////////////////////
	//用户基础
	UserCore.Init()
	//用户订阅
	UserSubscription.Init()
	//用户消息
	UserMessage.Init()
	//用户聊天室
	UserChat.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//行政
	///////////////////////////////////////////////////////////////////////////////////
	//组织核心
	OrgCoreCore.Init()
	//组织订阅
	OrgSubscription.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//设备
	///////////////////////////////////////////////////////////////////////////////////
	//设备
	IOTDevice.Init()
	//传感器
	IOTSensor.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//地图
	///////////////////////////////////////////////////////////////////////////////////
	//房间服务
	MapRoom.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//服务
	///////////////////////////////////////////////////////////////////////////////////
	//订单
	ServiceOrder.Init()
	//档案
	ServiceUserInfo.Init()
	//信息档案消耗模块
	ServiceUserInfoCost.Init()
	//服务配单
	ServiceHousekeeping.Init()
	//信息交互模块
	ServiceInfoExchange.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//配送
	///////////////////////////////////////////////////////////////////////////////////
	//配送
	TMSTransport.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//商城
	///////////////////////////////////////////////////////////////////////////////////
	//商城
	MallCore.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//博客
	///////////////////////////////////////////////////////////////////////////////////
	BlogCore.Init()
	BlogExam.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//统计
	///////////////////////////////////////////////////////////////////////////////////
	AnalysisAny.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//外部
	///////////////////////////////////////////////////////////////////////////////////
	//系统配置维护
	RouterAPIRunBase.Init()
	//反馈
	return
}
