package FinanceMargePay

//ArgsCreatePay 创建聚合支付请求参数
type ArgsCreatePay struct {
	//支付ID
	// 二维码将识别对应的支付ID，并根据端的区别，自动修正支付请求的端数据包，方便调用支付请求处理
	PayID int64 `db:"pay_id" json:"payID" check:"id"`
}

//CreatePay 创建聚合支付请求
// 1. 创建正常支付请求，不要做客户端确认
// 2. 通过支付请求，创建本聚合请求，出具二维码让客户完成后续支付即可
func CreatePay(args *ArgsCreatePay) (data string, err error) {
	return
}

//ArgsPayClient 提供端请求，修正支付请求并反馈数据包参数
type ArgsPayClient struct {
}

//PayClient 提供端请求，修正支付请求并反馈数据包
/**
1. 提供端数据包
2. 根据数据包，修正支付请求
3. 反馈支付请求的客户端所需资源
*/
func PayClient(args *ArgsPayClient) (err error) {
	return
}

//ArgsPayment 完成二维码的支付操作
type ArgsPayment struct {
}

//Payment 完成二维码的支付
func Payment(args *ArgsPayment) (err error) {
	return
}
