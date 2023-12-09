package Router2Core

var (
	AppName    = "weeekj"
	AppVersion = "2.0.0"
)

//Init 初始化处理模块
// 必须在主进程直接启动，否则请建立子进程，并在尾部做拦截程序处理，避免程序跳出
func Init(appName, appVer string, urlFunc func(), initFunc func()) (err error) {
	return
}
