package CoreFilter

//GetMaxRand 从一定额度下随机抽取
// 可用于红包的随机计算等
// maxCount int64 允许给予的总额度
// haveCount int64 已经发放的额度总额
// sendCount int64 已经发放的数量
// sendLimit int64 发出的数量限制
func GetMaxRand(maxCount int64, haveCount int64, sendCount int64, sendLimit int64) (result int64) {
	//如果额度不足
	if sendCount >= sendLimit {
		return
	}
	if haveCount >= maxCount {
		return
	}
	//如果仅剩余1个名额，则直接发放
	if sendCount+1 == sendLimit {
		result = maxCount - haveCount
		return
	}
	//如果多余1个名额，则计算随机数
	result = int64(GetRandNumber(1, int(maxCount-haveCount-(sendLimit-sendCount))))
	return
}
