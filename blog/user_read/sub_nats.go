package BlogUserRead

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//记录阅读
	CoreNats.SubDataByteNoErr("/blog/user_read/new", func(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
		logAppend := "blog user read sub nats new, "
		//获取参数
		var args ArgsCreateLog
		if err := CoreNats.ReflectDataByte(data, &args); err != nil {
			CoreLog.Error(logAppend, "get data, ", err)
			return
		}
		//添加数据
		if err := CreateLog(&args); err != nil {
			CoreLog.Error(logAppend, "insert log, ", err)
			return
		}
	})
	//删除文章统计
	CoreNats.SubDataByteNoErr("/blog/core/delete", func(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
		deleteLogByContentID(id)
	})
	//删除用户
	CoreNats.SubDataByteNoErr("/user/core/delete", func(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "blog_user_read_log", "user_id = :user_id", map[string]interface{}{
			"user_id": id,
		})
	})
}
