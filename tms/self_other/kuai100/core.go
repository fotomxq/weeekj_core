package TMSSelfOtherKuai100

//快递100查询方案

var (
	//公司名录数据
	globTMSCompanyData []dataTMSCompanyChild
)

// Init 初始化
func Init() (err error) {
	//初始化公司名录数据
	err = loadTMSCompanyJSON()
	if err != nil {
		return
	}
	return
}
