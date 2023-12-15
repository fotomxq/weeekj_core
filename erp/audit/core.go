package ERPAudit

import (
	ERPCore "github.com/fotomxq/weeekj_core/v5/erp/core"
	"sync"
)

//ERP核心模块
/**
流程化管理主要面向企业ERP系统，提供一套标准化的流程管理组件
1. 根据订单或自建ERP核心流程
2. 设计和约定商品，作为自动创建流程
3. 设计流程
4. 每个流程和节点有标记，可外挂模块实现追踪、统计需求
*/

var (
	//componentValObj 节点内容对象
	componentValObj ERPCore.ComponentVal
	//审批创建锁定
	auditCreateLock sync.Mutex
	//审批编辑锁定
	auditUpdateLock sync.Mutex
	//配置节点限制
	limitConfigStepCount = 10
	//OpenSub 订阅
	OpenSub = false
)

func Init() {
	//初始化节点内容对象
	componentValObj.TableName = "erp_audit_step_child_component"
	componentValObj.CacheName = "erp:audit:step:child:component:id:"
	//nats
	if OpenSub {
		subNats()
	}
}
