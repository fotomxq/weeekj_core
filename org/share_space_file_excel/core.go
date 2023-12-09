package OrgShareSpaceFileExcel

//机构excel模块
/**
1. 采用系统专用结构体，约定子表和位置的值
2. 外部插入变更结构的样式
3. 支持导出excel文件，并同时支持在线查看和编辑操作
*/

var (
	//OpenSub 是否启动订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	if OpenSub {
		//消息列队
		subNats()
	}
}
