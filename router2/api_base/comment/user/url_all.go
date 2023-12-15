package Router2APIBaseCommentUser

import (
	"fmt"
	BaseFileSys2 "github.com/fotomxq/weeekj_core/v5/base/filesys2"
	ClassComment "github.com/fotomxq/weeekj_core/v5/class/comment"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	Router2Params "github.com/fotomxq/weeekj_core/v5/router2/params"
	Router2Record "github.com/fotomxq/weeekj_core/v5/router2/record"
	"github.com/lib/pq"
)

// URLAll 用户层级操作接口汇总
func URLAll(r *Router2Mid.RouterURLUser, appendURL string, commentObj *ClassComment.Comment, permissions []string) {
	//获取评论列表
	r.POST(fmt.Sprint(appendURL, "/list"), func(c *Router2Mid.RouterURLUserC) {
		//日志
		c.LogAppend = "user url class comment get list"
		//权限检查
		if !Router2Mid.CheckPermission(c, permissions) {
			return
		}
		//获取参数
		type paramsType struct {
			//分页
			Pages CoreSQLPages.ArgsDataList `json:"pages"`
			//评论ID
			// 评论被删除后出现，指向新的评论
			// 下级评论的上级江全部改为该新的ID
			CommentID int64 `db:"comment_id" json:"commentID" check:"id" empty:"true"`
			//上级ID
			ParentID int64 `json:"parentID" check:"id" empty:"true"`
			//绑定组织
			// 该组织根据资源来源设定
			// 如果是平台资源，则为0
			OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
			//绑定ID
			BindID int64 `json:"bindID" check:"id" empty:"true"`
			//评价类型
			// 0 好评 1 中立 2 差评
			LevelType int `db:"level_type" json:"levelType"`
			//分数范围
			LevelMin int `db:"level_min" json:"levelMin"`
			LevelMax int `db:"level_max" json:"levelMax"`
			//是否删除
			IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
			//搜索
			Search string `json:"search" check:"search" empty:"true"`
		}
		var params paramsType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		dataList, dataCount, err := commentObj.GetList(&ClassComment.ArgsGetList{
			Pages:     params.Pages,
			CommentID: params.CommentID,
			ParentID:  params.ParentID,
			OrgID:     params.OrgID,
			UserID:    c.UserID,
			BindID:    params.BindID,
			LevelType: params.LevelType,
			LevelMin:  params.LevelMin,
			LevelMax:  params.LevelMax,
			IsRemove:  params.IsRemove,
			Search:    params.Search,
		})
		//重组数据
		type dataType struct {
			//基础
			ID int64 `db:"id" json:"id"`
			//创建时间
			CreateAt string `db:"create_at" json:"createAt"`
			//删除时间
			DeleteAt string `db:"delete_at" json:"deleteAt"`
			//评论ID
			// 评论被删除后出现，指向新的评论
			// 下级评论的上级江全部改为该新的ID
			CommentID int64 `db:"comment_id" json:"commentID"`
			//上级ID
			ParentID int64 `db:"parent_id" json:"parentID"`
			//绑定组织
			// 该组织根据资源来源设定
			// 如果是平台资源，则为0
			OrgID int64 `db:"org_id" json:"orgID"`
			//绑定内容
			BindID int64 `db:"bind_id" json:"bindID"`
			//评价类型
			// 0 好评 1 中立 2 差评
			LevelType int `db:"level_type" json:"levelType"`
			//分数
			Level int `db:"level" json:"level"`
			//标题
			Title string `db:"title" json:"title"`
			//内容
			Des string `db:"des" json:"des"`
			//介绍图文
			DesFiles []string `db:"des_files" json:"desFiles"`
			//扩展参数
			Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
		}
		var newData []dataType
		if err == nil {
			for _, v := range dataList {
				newData = append(newData, dataType{
					ID:        v.ID,
					CreateAt:  CoreFilter.GetTimeToDefaultTime(v.CreateAt),
					DeleteAt:  CoreFilter.GetTimeToDefaultTime(v.DeleteAt),
					CommentID: v.CommentID,
					ParentID:  v.ParentID,
					OrgID:     v.OrgID,
					BindID:    v.BindID,
					LevelType: v.LevelType,
					Level:     v.Level,
					Title:     v.Title,
					Des:       v.Des,
					DesFiles:  BaseFileSys2.GetPublicURLsByClaimIDs(v.DesFiles),
					Params:    v.Params,
				})
			}
		}
		//反馈数据
		Router2Mid.ReportDataList(c, "get list", err, "", newData, dataCount)
	})
	//创建新的评论
	r.PUT(fmt.Sprint(appendURL, ""), func(c *Router2Mid.RouterURLUserC) {
		//日志
		c.LogAppend = "user url class comment create sort"
		//权限检查
		if !Router2Mid.CheckPermission(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//上级ID
			// 评论的上下级关系，一旦建立无法修改
			ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
			//绑定组织
			// 该组织根据资源来源设定
			// 如果是平台资源，则为0
			OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
			//绑定内容
			BindID int64 `db:"bind_id" json:"bindID" check:"id"`
			//评价类型
			// 0 好评 1 中立 2 差评
			LevelType int `db:"level_type" json:"levelType"`
			//分数
			Level int `db:"level" json:"level"`
			//标题
			Title string `db:"title" json:"title" check:"title" min:"1" max:"100"`
			//内容
			Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
			//介绍图文
			DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
		}
		var params dataType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		data, err := commentObj.Create(&ClassComment.ArgsCreate{
			ParentID:  params.ParentID,
			OrgID:     params.OrgID,
			UserID:    c.UserID,
			BindID:    params.BindID,
			LevelType: params.Level,
			Level:     params.Level,
			Title:     params.Title,
			Des:       params.Des,
			DesFiles:  params.DesFiles,
			Params:    []CoreSQLConfig.FieldsConfigType{},
		})
		if err == nil {
			Router2Record.AddRecord(c, "class_comment", data.ID, "创建新的评论(", params.Title, ")[", data.ID, "]")
		}
		//反馈数据
		Router2Mid.ReportActionCreateNoData(c, "create failed", err, "")
	})
	//删除评论
	r.DELETE(fmt.Sprint(appendURL, "/id"), func(c *Router2Mid.RouterURLUserC) {
		//日志
		c.LogAppend = "user url class comment delete sort"
		//权限检查
		if !Router2Mid.CheckPermission(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//ID
			ID int64 `json:"id" check:"id"`
		}
		var params dataType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		err := commentObj.DeleteByID(&ClassComment.ArgsDeleteByID{
			ID:     params.ID,
			OrgID:  -1,
			UserID: c.UserID,
		})
		if err == nil {
			Router2Record.AddRecord(c, "class_comment", params.ID, "删除评论[", params.ID, "]")
		}
		//反馈数据
		Router2Mid.ReportActionDelete(c, "delete failed", err, "")
	})
}
