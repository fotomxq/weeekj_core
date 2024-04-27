package Router2SystemInit

import (
	"errors"
	"fmt"
	AnalysisAny "github.com/fotomxq/weeekj_core/v5/analysis/any"
	AnalysisBindVisit "github.com/fotomxq/weeekj_core/v5/analysis/bind_visit"
	BaseApprover "github.com/fotomxq/weeekj_core/v5/base/approver"
	BaseAutoCode "github.com/fotomxq/weeekj_core/v5/base/auto_code"
	BaseBPM "github.com/fotomxq/weeekj_core/v5/base/bpm"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	BaseFileSys2 "github.com/fotomxq/weeekj_core/v5/base/filesys2"
	BaseIPAddr "github.com/fotomxq/weeekj_core/v5/base/ipaddr"
	BaseLookup "github.com/fotomxq/weeekj_core/v5/base/lookup"
	BaseMonitorGlob "github.com/fotomxq/weeekj_core/v5/base/monitor/glob"
	BaseMonitorPostgresql "github.com/fotomxq/weeekj_core/v5/base/monitor/postgresql"
	BasePython "github.com/fotomxq/weeekj_core/v5/base/python"
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	BaseTempFile "github.com/fotomxq/weeekj_core/v5/base/temp_file"
	BaseToken2 "github.com/fotomxq/weeekj_core/v5/base/token2"
	BaseUnit "github.com/fotomxq/weeekj_core/v5/base/unit"
	BaseVcode "github.com/fotomxq/weeekj_core/v5/base/vcode"
	BlogCore "github.com/fotomxq/weeekj_core/v5/blog/core"
	BlogExam "github.com/fotomxq/weeekj_core/v5/blog/exam"
	BlogStuRead "github.com/fotomxq/weeekj_core/v5/blog/stu_read"
	BlogUserRead "github.com/fotomxq/weeekj_core/v5/blog/user_read"
	EAMCore "github.com/fotomxq/weeekj_core/v5/eam/core"
	ERPAudit "github.com/fotomxq/weeekj_core/v5/erp/audit"
	ERPBudget "github.com/fotomxq/weeekj_core/v5/erp/budget"
	ERPDocument "github.com/fotomxq/weeekj_core/v5/erp/document"
	ERPPermanentAssets "github.com/fotomxq/weeekj_core/v5/erp/permanent_assets"
	ERPProduct "github.com/fotomxq/weeekj_core/v5/erp/product"
	ERPProductMall "github.com/fotomxq/weeekj_core/v5/erp/product_mall"
	ERPProject "github.com/fotomxq/weeekj_core/v5/erp/project"
	ERPPurchase "github.com/fotomxq/weeekj_core/v5/erp/purchase"
	ERPRequirement "github.com/fotomxq/weeekj_core/v5/erp/requirement"
	ERPSaleOut "github.com/fotomxq/weeekj_core/v5/erp/sale_out"
	ERPWarehouse "github.com/fotomxq/weeekj_core/v5/erp/warehouse"
	FinanceDeposit2 "github.com/fotomxq/weeekj_core/v5/finance/deposit2"
	FinanceLog "github.com/fotomxq/weeekj_core/v5/finance/log"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	FinanceReturnedMoney "github.com/fotomxq/weeekj_core/v5/finance/returned_money"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	IOTSensor "github.com/fotomxq/weeekj_core/v5/iot/sensor"
	MallCore "github.com/fotomxq/weeekj_core/v5/mall/core"
	MallRecommend "github.com/fotomxq/weeekj_core/v5/mall/recommend"
	MapCityData "github.com/fotomxq/weeekj_core/v5/map/city_data"
	MapRoom "github.com/fotomxq/weeekj_core/v5/map/room"
	MarketGivingNewUser "github.com/fotomxq/weeekj_core/v5/market/giving_new_user"
	Market2Log "github.com/fotomxq/weeekj_core/v5/market2/log"
	Market2ReferrerNewUser "github.com/fotomxq/weeekj_core/v5/market2/referrer_new_user"
	OrgCert2 "github.com/fotomxq/weeekj_core/v5/org/cert2"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	OrgMap "github.com/fotomxq/weeekj_core/v5/org/map"
	OrgShareSpace "github.com/fotomxq/weeekj_core/v5/org/share_space"
	OrgShareSpaceFileExcel "github.com/fotomxq/weeekj_core/v5/org/share_space_file_excel"
	OrgSubscription "github.com/fotomxq/weeekj_core/v5/org/subscription"
	OrgTime "github.com/fotomxq/weeekj_core/v5/org/time"
	OrgUser "github.com/fotomxq/weeekj_core/v5/org/user"
	OrgWorkTip "github.com/fotomxq/weeekj_core/v5/org/work_tip"
	RestaurantPurchase "github.com/fotomxq/weeekj_core/v5/restaurant/purchase_analysis"
	RestaurantRawMaterials "github.com/fotomxq/weeekj_core/v5/restaurant/raw_materials"
	RestaurantRecipe "github.com/fotomxq/weeekj_core/v5/restaurant/recipe"
	RestaurantWeeklyRecipeMarge "github.com/fotomxq/weeekj_core/v5/restaurant/weekly_recipe"
	RouterAPIRunBase "github.com/fotomxq/weeekj_core/v5/router/api/run_base"
	ServiceAD2 "github.com/fotomxq/weeekj_core/v5/service/ad2"
	ServiceCompany "github.com/fotomxq/weeekj_core/v5/service/company"
	ServiceHousekeeping "github.com/fotomxq/weeekj_core/v5/service/housekeeping"
	ServiceInfoExchange "github.com/fotomxq/weeekj_core/v5/service/info_exchange"
	ServiceOrder "github.com/fotomxq/weeekj_core/v5/service/order"
	ServiceOrderWait "github.com/fotomxq/weeekj_core/v5/service/order/wait"
	ServiceUserInfo "github.com/fotomxq/weeekj_core/v5/service/user_info"
	ServiceUserInfoCost "github.com/fotomxq/weeekj_core/v5/service/user_info_cost"
	TMSSelfOther "github.com/fotomxq/weeekj_core/v5/tms/self_other"
	TMSTransport "github.com/fotomxq/weeekj_core/v5/tms/transport"
	TMSUserRunning "github.com/fotomxq/weeekj_core/v5/tms/user_running"
	ToolsCommunication "github.com/fotomxq/weeekj_core/v5/tools/communication"
	UserChat "github.com/fotomxq/weeekj_core/v5/user/chat"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserLogin "github.com/fotomxq/weeekj_core/v5/user/login"
	UserLogin2 "github.com/fotomxq/weeekj_core/v5/user/login2"
	UserMessage "github.com/fotomxq/weeekj_core/v5/user/message"
	UserRecord2 "github.com/fotomxq/weeekj_core/v5/user/record2"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
)

func Init() (err error) {
	//postgresql数据库监控服务
	BaseMonitorPostgresql.OpenSub = OpenSub
	BaseMonitorPostgresql.Init()
	//配置模块
	if err = BaseConfig.Init(); err != nil {
		fmt.Println("base config, " + err.Error())
		//err = errors.New("base config, " + err.Error())
		//return
		err = nil
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
	//管理单元
	BaseUnit.Init()
	//数据字典
	BaseLookup.Init()
	//临时文件
	BaseTempFile.OpenSub = OpenSub
	BaseTempFile.Init()
	//过期通知
	BaseExpireTip.OpenSub = OpenSub
	BaseExpireTip.Init()
	//会话
	BaseToken2.OpenSub = OpenSub
	BaseToken2.Init()
	//python模块
	BasePython.OpenSub = OpenSub
	BasePython.Init()
	//系统检测
	// 注意所有进程会强制订阅nats处理机制
	BaseMonitorGlob.Init()
	//文件系统
	BaseFileSys2.OpenSub = OpenSub
	BaseFileSys2.Init()
	//BPM
	BaseBPM.OpenSub = OpenSub
	BaseBPM.Init()
	//服务
	BaseService.OpenSub = OpenSub
	BaseService.Init()
	//审批
	BaseApprover.Init()
	//编码生成器
	BaseAutoCode.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//统计
	///////////////////////////////////////////////////////////////////////////////////
	//访问统计
	AnalysisBindVisit.OpenSub = OpenSub
	AnalysisBindVisit.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//工具
	///////////////////////////////////////////////////////////////////////////////////
	//通讯
	ToolsCommunication.OpenSub = OpenSub
	ToolsCommunication.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//地图模块
	///////////////////////////////////////////////////////////////////////////////////
	//城市模块
	if err = MapCityData.Init(); err != nil {
		err = errors.New("map city data, " + err.Error())
		return
	}
	//房间服务
	MapRoom.OpenSub = OpenSub
	MapRoom.OpenAnalysis = OpenAnalysis
	MapRoom.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//用户服务
	///////////////////////////////////////////////////////////////////////////////////
	//登陆服务
	UserLogin.Init(fmt.Sprint(AppName, AppVersion))

	///////////////////////////////////////////////////////////////////////////////////
	//用户
	///////////////////////////////////////////////////////////////////////////////////
	//用户基础
	UserCore.OpenSub = OpenSub
	UserCore.OpenAnalysis = OpenAnalysis
	UserCore.Init()
	//用户订阅
	UserSubscription.OpenSub = OpenSub
	UserSubscription.Init()
	//用户消息
	UserMessage.OpenSub = OpenSub
	UserMessage.OpenAnalysis = OpenAnalysis
	UserMessage.Init()
	//用户聊天室
	UserChat.Init()
	//用户记录
	UserRecord2.OpenSub = OpenSub
	UserRecord2.Init()
	//用户登录
	UserLogin2.OpenSub = OpenSub
	UserLogin2.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//行政
	///////////////////////////////////////////////////////////////////////////////////
	//组织核心
	OrgCoreCore.OpenSub = OpenSub
	OrgCoreCore.OpenAnalysis = OpenAnalysis
	OrgCoreCore.Init()
	//组织订阅
	OrgSubscription.OpenSub = OpenSub
	OrgSubscription.Init()
	//组织用户聚合数据
	OrgUser.OpenSub = OpenSub
	OrgUser.Init()
	//组织工作提醒
	OrgWorkTip.OpenSub = OpenSub
	OrgWorkTip.Init()
	//共享空间
	OrgShareSpace.OpenSub = OpenSub
	OrgShareSpace.Init()
	//空间excel文件
	OrgShareSpaceFileExcel.OpenSub = OpenSub
	OrgShareSpaceFileExcel.Init()
	//组织证件
	OrgCert2.OpenSub = OpenSub
	OrgCert2.Init()
	//组织地图
	OrgMap.OpenSub = OpenSub
	OrgMap.OpenAnalysis = OpenSub
	OrgMap.Init()
	//考勤打卡
	OrgTime.OpenSub = OpenSub
	OrgTime.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//设备
	///////////////////////////////////////////////////////////////////////////////////
	//设备
	IOTDevice.OpenSub = OpenSub
	IOTDevice.Init()
	//传感器
	IOTSensor.OpenSub = OpenSub
	IOTSensor.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//财务
	///////////////////////////////////////////////////////////////////////////////////
	//财务支付
	FinancePay.OpenSub = OpenSub
	FinancePay.Init()
	//财务日志
	FinanceLog.OpenSub = OpenSub
	FinanceLog.Init()
	//回款
	FinanceReturnedMoney.OpenSub = OpenSub
	FinanceReturnedMoney.Init()
	//财务储蓄
	FinanceDeposit2.OpenAnalysis = OpenSub
	FinanceDeposit2.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//服务
	///////////////////////////////////////////////////////////////////////////////////
	//订单
	ServiceOrderWait.OpenSub = OpenSub
	ServiceOrderWait.Init()
	ServiceOrder.OpenAnalysis = OpenAnalysis
	ServiceOrder.OpenSub = OpenSub
	ServiceOrder.Init()
	//档案
	ServiceUserInfo.OpenSub = OpenSub
	ServiceUserInfo.OpenAnalysis = OpenAnalysis
	ServiceUserInfo.Init()
	//信息档案消耗模块
	ServiceUserInfoCost.OpenSub = OpenSub
	ServiceUserInfoCost.Init()
	//服务配单
	ServiceHousekeeping.OpenSub = OpenSub
	ServiceHousekeeping.Init()
	//信息交互模块
	ServiceInfoExchange.OpenAnalysis = OpenAnalysis
	ServiceInfoExchange.OpenSub = OpenSub
	ServiceInfoExchange.Init()
	//广告
	ServiceAD2.OpenAnalysis = OpenAnalysis
	ServiceAD2.Init()
	//公司
	ServiceCompany.OpenSub = OpenSub
	ServiceCompany.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//配送
	///////////////////////////////////////////////////////////////////////////////////
	//配送
	TMSTransport.OpenSub = OpenSub
	TMSTransport.Init()
	//跑腿
	TMSUserRunning.OpenSub = OpenSub
	TMSUserRunning.Init()
	//第三方平台
	err = TMSSelfOther.Init()
	if err != nil {
		err = errors.New(fmt.Sprint("init tms self other, ", err))
		return
	}

	///////////////////////////////////////////////////////////////////////////////////
	//ERP
	///////////////////////////////////////////////////////////////////////////////////
	//ERP产品
	ERPProduct.OpenSub = OpenSub
	ERPProduct.Init()
	//ERP产品商品
	ERPProductMall.OpenSub = OpenSub
	ERPProductMall.Init()
	//出货单
	ERPSaleOut.OpenSub = OpenSub
	ERPSaleOut.Init()
	//审批
	ERPAudit.OpenSub = OpenSub
	ERPAudit.Init()
	//文档
	ERPDocument.Init()
	//仓储
	ERPWarehouse.OpenAnalysis = OpenAnalysis
	ERPWarehouse.OpenSub = OpenSub
	ERPWarehouse.Init()
	//固定资产
	ERPPermanentAssets.OpenSub = OpenSub
	ERPPermanentAssets.Init()
	//项目管理
	ERPProject.Init()
	//采购申请
	ERPRequirement.Init()
	//预算管理
	ERPBudget.Init()
	//ERP采购订单
	ERPPurchase.Init()
	//EAM
	EAMCore.OpenSub = OpenSub
	EAMCore.Init()
	///////////////////////////////////////////////////////////////////////////////////
	//商城
	///////////////////////////////////////////////////////////////////////////////////
	//商城
	MallCore.OpenSub = OpenSub
	MallCore.Init()
	//商品推荐
	MallRecommend.OpenSub = OpenSub
	MallRecommend.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//营销系统
	///////////////////////////////////////////////////////////////////////////////////
	//营销核心
	MarketGivingNewUser.OpenSub = OpenSub
	MarketGivingNewUser.Init()
	//营销核心2
	Market2Log.OpenSub = OpenSub
	Market2Log.Init()
	Market2ReferrerNewUser.OpenSub = OpenSub
	Market2ReferrerNewUser.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//博客
	///////////////////////////////////////////////////////////////////////////////////
	//博客文章
	BlogCore.OpenSub = OpenSub
	BlogCore.Init()
	//博客考试
	BlogExam.OpenSub = OpenSub
	BlogExam.Init()
	//博客学习
	BlogStuRead.OpenSub = OpenSub
	BlogStuRead.Init()
	//博客文章阅读记录
	BlogUserRead.OpenSub = OpenSub
	BlogUserRead.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//统计
	///////////////////////////////////////////////////////////////////////////////////
	AnalysisAny.Init()

	///////////////////////////////////////////////////////////////////////////////////
	//外部
	///////////////////////////////////////////////////////////////////////////////////
	//系统配置维护
	RouterAPIRunBase.Init()
	///////////////////////////////////////////////////////////////////////////////////
	//电子餐厅
	///////////////////////////////////////////////////////////////////////////////////
	//电子菜品
	RestaurantRecipe.Init()
	//周食谱
	RestaurantWeeklyRecipeMarge.Init()
	//原材料采购台账
	RestaurantPurchase.Init()
	//原材料库
	RestaurantRawMaterials.Init()

	//启动完成提示
	fmt.Println("main router init success.")
	//反馈
	return
}
