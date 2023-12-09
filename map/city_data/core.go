package MapCityData

//城市数据包

var (
	//全局城市数据结构体
	globCityAreaData []DataProvinceData
)

// Init 初始化
func Init() (err error) {
	//加载城市数据
	err = loadCityData()
	if err != nil {
		return
	}
	//反馈
	return
}
