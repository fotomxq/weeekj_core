package BaseCache

type DataCache struct {
	//创建时间
	CreateTime int64
	//过期时间
	ExpireTime int64
	//数据标识码
	// 全局唯一，否则将覆盖该数据
	Mark string
	//存储值
	Value string
	//int类型将自动改写int64，转出时可强制缩短数据
	ValueInt64 int64
	//float64
	ValueFloat64 float64
	//bool
	ValueBool bool
	//byte
	ValueByte []byte
	//缓冲器数据集合
	ValueInterface interface{}
}
