package OrgUserMod

import (
	"fmt"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
)

// PushUpdateUserData 请求更新用户的聚合数据
func PushUpdateUserData(orgID int64, userID int64) {
	CoreNats.PushDataNoErr("org_user_post_update", "/org/user/post_update", "", userID, fmt.Sprint(orgID), nil)
}
