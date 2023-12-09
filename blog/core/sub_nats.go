package BlogCore

import (
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BlogUserReadMod "gitee.com/weeekj/weeekj_core/v5/blog/user_read/mod"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	OrgMapMod "gitee.com/weeekj/weeekj_core/v5/org/map/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//请求文章审核
	CoreNats.SubDataByteNoErr("/blog/core/audit", subNatsAudit)
	//请求阅读文章
	CoreNats.SubDataByteNoErr("/blog/core/read", subNatsRead)
	//商户地图创建后构建文章
	CoreNats.SubDataByteNoErr("/org/map/audit", subNatsOrgMapAudit)
}

// 请求审核文章
func subNatsAudit(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	appendLog := fmt.Sprint("blog core sub nats audit, content id: ", id)
	//检查是否开关
	blogAuditAuto, err := BaseConfig.GetDataBool("BlogAuditAuto")
	if err != nil {
		blogAuditAuto = false
	}
	if !blogAuditAuto {
		return
	}
	err = UpdateAudit(&ArgsUpdateAudit{
		ID:       id,
		OrgID:    -1,
		IsAudit:  true,
		AuditDes: "系统自动审核",
	})
	if err != nil {
		CoreLog.Warn(appendLog, "update audit failed, ", err)
	}
}

// 请求阅读文章
func subNatsRead(_ *nats.Msg, action string, _ int64, _ string, data []byte) {
	//读取参数
	userID := gjson.GetBytes(data, "userID").Int()
	contentID := gjson.GetBytes(data, "contentID").Int()
	//是否添加了数据
	isAdd := false
	//博客系统是否必须精确增加阅读数据？
	blogUserReadOnceAddCount, _ := BaseConfig.GetDataBool("BlogUserReadOnceAddCount")
	//识别action
	switch action {
	case "user":
		//获取用户只读一次的设计
		var blogUserReadOnce bool
		var err error
		blogUserReadOnce, err = BaseConfig.GetDataBool("BlogUserReadOnce")
		if err != nil {
			blogUserReadOnce = true
		}
		//检查用户访问日志数据是否存在
		if blogUserReadOnce {
			if BlogUserReadMod.CheckUserLogExist(userID, contentID) {
				break
			}
		}
		//记录日志
		contentData := getContentID(contentID)
		if contentData.ID < 1 {
			break
		}
		BlogUserReadMod.CreateLog(BlogUserReadMod.ArgsCreateLog{
			ChildOrgID: 0,
			UserID:     userID,
			FromMark:   "blog_content",
			FromName:   "",
			Name:       contentData.Title,
			IP:         "",
			ContentID:  contentData.ID,
			CreateAt:   "",
			LeaveAt:    "",
		})
		isAdd = true
	}
	//如果打开了精确增加阅读，则跳出后续
	if blogUserReadOnceAddCount && !isAdd {
		return
	}
	//更新文章次数
	_, _ = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE blog_core_content SET visit_count = visit_count + 1 WHERE id = :id", map[string]interface{}{
		"id": contentID,
	})
	//清理缓冲
	deleteContentCacheByID(contentID)
}

// 商户地图创建后构建帖子
func subNatsOrgMapAudit(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	//检查配置项
	open, _ := BaseConfig.GetDataBool("OrgMapAuditAutoCreateBlogCoreByUser")
	if !open {
		return
	}
	//获取地图信息
	mapData := OrgMapMod.GetMapByID(id)
	if mapData.ID < 1 {
		return
	}
	if mapData.UserID < 1 {
		return
	}
	//构建帖子信息
	var desFiles []int64
	if mapData.CoverFileID > 0 {
		desFiles = append(desFiles, mapData.CoverFileID)
	}
	blogData, err := CreateContent(&ArgsCreateContent{
		OrgID:       mapData.OrgID,
		UserID:      mapData.UserID,
		BindID:      0,
		ContentType: 4,
		Param1:      0,
		Param2:      0,
		Param3:      0,
		Key:         "",
		IsTop:       false,
		SortID:      0,
		Tags:        nil,
		Title:       mapData.Name,
		TitleDes:    CoreFilter.SubStrQuick(mapData.Des, 10),
		CoverFileID: mapData.CoverFileID,
		DesFiles:    desFiles,
		Des:         mapData.Des,
		Params: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "orgMapID",
				Val:  fmt.Sprint(mapData.ID),
			},
		},
	})
	if err != nil {
		CoreLog.Warn("blog core sub nats org map audit, create blog, ", err)
		return
	}
	if err = UpdatePublish(&ArgsUpdatePublish{
		ID:     blogData.ID,
		OrgID:  -1,
		UserID: -1,
	}); err != nil {
		CoreLog.Warn("blog core sub nats org map audit, publish blog, ", err)
		return
	}
	//if err = UpdateAudit(&ArgsUpdateAudit{
	//	ID:       blogData.ID,
	//	OrgID:    -1,
	//	IsAudit:  true,
	//	AuditDes: "系统自动通过审核",
	//}); err != nil {
	//	CoreLog.Warn("blog core sub nats org map audit, publish blog, ", err)
	//	return
	//}
}
