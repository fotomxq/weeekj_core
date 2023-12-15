package BaseEmail

import (
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

func runSendEmail() {
	var err error
	//检查阻断器
	if !runBlocker.CheckPass() {
		return
	}
	//加载未发送数据
	limit := 1000
	step := 0
	for {
		var dataList []FieldsEmailType
		err = Router2SystemConfig.MainDB.Select(
			&dataList,
			"SELECT id, create_at, update_at, server_id, create_info, send_at, is_success, is_failed, fail_message, title, content, content_type, to_email, delete_at FROM core_email WHERE delete_at < to_timestamp(1000000) AND is_success=false AND is_failed=false AND send_at<=NOW() LIMIT $1 OFFSET $2",
			limit, step,
		)
		if err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		//遍历数据依次发送并标记结果
		for _, v := range dataList {
			//获取server数据
			serverData, err := GetServerByID(&ArgsGetServerByID{
				ID:    v.ServerID,
				OrgID: -1,
			})
			if err != nil {
				if err2 := updateFailed(v.ID, fmt.Sprint("email run error, cannot get server by id, ", err.Error(), ", id: ", v.ServerID)); err2 != nil {
					CoreLog.Error("email run error, get server data, ", err2.Error())
				}
				continue
			}
			//根据server发送数据
			isSuccess := false
			isFailed := false
			failedMessage := ""
			if serverData.IsSSL {
				if err := sendSSLMail(serverData, v); err != nil {
					isFailed = true
					failedMessage = "send ssl mail is failed, err: " + err.Error()
				} else {
					isSuccess = true
				}
			} else {
				if err := sendMail(serverData, v); err != nil {
					isFailed = true
					failedMessage = "send mail is failed, err: " + err.Error()
				} else {
					isSuccess = true
				}
			}
			//如果成功则标记成功
			if isSuccess {
				if _, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_email SET is_success=true WHERE id=:id", map[string]interface{}{
					"id": v.ID,
				}); err != nil {
					CoreLog.Error("email run error, cannot update mail success, id: ", v.ID, ", err: ", err.Error())
				}
			}
			if isFailed {
				if err := updateFailed(v.ID, failedMessage); err != nil {
					CoreLog.Error("email run error, update failed, " + err.Error())
				}
			}
			//强制延迟10毫秒，消峰处理
			time.Sleep(time.Millisecond * 10)
		}
		//下一页
		step += limit
	}
}
