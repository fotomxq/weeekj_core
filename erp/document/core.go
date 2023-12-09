package ERPDocument

import (
	CoreSQL2 "gitee.com/weeekj/weeekj_core/v5/core/sql2"
	ERPCore "gitee.com/weeekj/weeekj_core/v5/erp/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

//ERP文档数据集
/**
1. 可定义不同的文档集类型
2. 授权后，成员可访问数据集合并添加文档数据
*/

var (
	//docComponentValObj 文档自定义组件
	docComponentValObj ERPCore.ComponentVal
	//配置sql
	configSQL CoreSQL2.Client
)

func Init() {
	//初始化节点内容对象
	docComponentValObj.TableName = "erp_document_doc_component"
	docComponentValObj.CacheName = "erp:document:doc:component:id:"
	//初始化sql
	configSQL.Init(&Router2SystemConfig.MainSQL, "erp_document_config")
}
