package UserCore

import (
	"bytes"
	"fmt"
	BaseQiniu "gitee.com/weeekj/weeekj_core/v5/base/qiniu"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/issue9/identicon"
	"github.com/nats-io/nats.go"
	"image/color"
	"image/jpeg"
	"time"
)

// 为用户创建自动化头像
func subNatsCreateAvatar(_ *nats.Msg, _ string, userID int64, _ string, _ []byte) {
	appendLog := "sub nats create user avatar, "
	//获取用户数据
	userData := getUserByID(userID)
	if userData.ID < 1 {
		return
	}
	//初始化图形化
	back := color.RGBA{R: 255, G: 250, B: 250, A: 100}
	fore := color.RGBA{R: 135, G: 206, B: 250, A: 100}
	fores := []color.Color{fore, fore, fore}
	avatarObj, err := identicon.New(128, back, fores...)
	if err != nil {
		CoreLog.Error(appendLog, "create new avatar by identicon, ", err)
		return
	}
	//构建图片
	img := avatarObj.Make([]byte(fmt.Sprint(userID, userData.Name, userData.CreateAt)))
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		CoreLog.Error(appendLog, "get png encode by buf, ", err)
		return
	}
	fileData, errCode, err := BaseQiniu.Upload(buf.Bytes(), "", "jpg", "0.0.0.0", CoreSQLFrom.FieldsFrom{
		System: "user",
		ID:     userData.ID,
		Mark:   "",
		Name:   userData.Name,
	}, userData.ID, userData.OrgID, true, time.Time{}, []CoreSQLConfig.FieldsConfigType{}, "")
	if err != nil {
		if _, err := CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), avatar = :avatar WHERE id = :id", map[string]interface{}{
			"id":     userData.ID,
			"avatar": -1,
		}); err != nil {
			CoreLog.Error(appendLog, "update user avatar, user id: ", userData.ID, ", err: ", err)
			return
		}
		CoreLog.Error(appendLog, "upload file to qiniu, user id: ", userData.ID, ", code: ", errCode, ", err: ", err)
	} else {
		if _, err := CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), avatar = :avatar WHERE id = :id", map[string]interface{}{
			"id":     userData.ID,
			"avatar": fileData.ID,
		}); err != nil {
			CoreLog.Error(appendLog, "update user avatar, user id: ", userData.ID, ", err: ", err)
			return
		}
	}
}
