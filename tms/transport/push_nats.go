package TMSTransport

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// 请求统计配送员信息
func pushNatsAnalysisBind(orgBindID int64) {
	CoreNats.PushDataNoErr("/tms/transport/analysis_bind", "", orgBindID, "", nil)
}
