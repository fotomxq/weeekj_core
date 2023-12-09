package BasePython

import (
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
)

//python处理对接模块
/**
1. 本模块用于和python进行无缝对接
2. 方法分为同步和异步，同步会等待处理结果并反馈；异步会在处理完成后调用回调函数
3. 必须配合专用的py模块实现对接
*/
var (
	//OpenSub 是否启动订阅
	OpenSub = false
	//处理文件路径
	dirSrc = ""
)

func Init() {
	dirSrc = fmt.Sprint(CoreFile.BaseDir(), CoreFile.Sep, "temp", CoreFile.Sep, "python", CoreFile.Sep)
	if OpenSub {
		//消息列队
		subNats()
	}
}
