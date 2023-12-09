package FinancePhysicalPay

import "sync"

/**
以实物代替财务完成对标的物的支付：
- 商户可设置标的物抵扣额度，以及限制
- 记录每次抵扣的信息
 */

var(
	//置换物锁定机制，避免超限置换
	logLock sync.Mutex
)
