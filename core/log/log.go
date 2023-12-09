package CoreLog

//Error 发送日志组合
// 系统错误
func Error(args ...interface{}) {
	errLog.LogHandle.Error(args...)
}

func Warn(args ...interface{}) {
	warnLog.LogHandle.Warn(args...)
}

func Mqtt(args ...interface{}) {
	mqttLog.LogHandle.Info(args...)
}

func MqttWarn(args ...interface{}) {
	mqttLog.LogHandle.Warn(args...)
}

func MqttError(args ...interface{}) {
	mqttLog.LogHandle.Error(args...)
}

func Info(args ...interface{}) {
	globLog.LogHandle.Info(args...)
}

func Debug(args ...interface{}) {
	globLog.LogHandle.Debug(args...)
}

func App(args ...interface{}) {
	appLog.LogHandle.Info(args...)
}
