package CoreLog

func runMake() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			Error("log run make, ", r)
		}
	}()
	if allowGin {
		ginLog.MakeNextHandle()
	}
	if allowDefault {
		globLog.MakeNextHandle()
		errLog.MakeNextHandle()
		warnLog.MakeNextHandle()
		appLog.MakeNextHandle()
		mqttLog.MakeNextHandle()
	}
}