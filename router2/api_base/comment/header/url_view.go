package Router2APIBaseCommentHeader

import (
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BaseFileSys2 "gitee.com/weeekj/weeekj_core/v5/base/filesys2"
	ClassComment "gitee.com/weeekj/weeekj_core/v5/class/comment"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	Router2Params "gitee.com/weeekj/weeekj_core/v5/router2/params"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
)

// URLView 通用查看评论数据包
func URLView(r *Router2Mid.RouterURLHeader, appendURL string, commentObj *ClassComment.Comment, permissions []string) {
	//获取评论列表
	r.POST(fmt.Sprint(appendURL, "/list"), func(c *Router2Mid.RouterURLHeaderC) {
		//日志
		c.LogAppend = "header url class comment get list"
		//权限检查
		if len(permissions) > 0 {
			if !Router2Mid.CheckPermission(c, permissions) {
				return
			}
		}
		//获取参数
		type paramsType struct {
			//分页
			Pages CoreSQLPages.ArgsDataList `json:"pages"`
			//上级ID
			ParentID int64 `json:"parentID" check:"id" empty:"true"`
			//绑定ID
			BindID int64 `json:"bindID" check:"id"`
			//评价类型
			// 0 好评 1 中立 2 差评
			LevelType int `db:"level_type" json:"levelType"`
			//分数范围
			LevelMin int `db:"level_min" json:"levelMin"`
			LevelMax int `db:"level_max" json:"levelMax"`
		}
		var params paramsType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取配置
		coreCommentShowListDesLen, _ := BaseConfig.GetDataInt("CoreCommentShowListDesLen")
		if coreCommentShowListDesLen < 1 {
			coreCommentShowListDesLen = 300
		}
		//获取数据
		dataList, dataCount, err := commentObj.GetList(&ClassComment.ArgsGetList{
			Pages:     params.Pages,
			CommentID: 0,
			ParentID:  params.ParentID,
			OrgID:     -1,
			UserID:    -1,
			BindID:    params.BindID,
			LevelType: params.LevelType,
			LevelMin:  params.LevelMin,
			LevelMax:  params.LevelMax,
			IsRemove:  false,
			Search:    "",
		})
		//重组数据
		type dataType struct {
			//基础
			ID int64 `db:"id" json:"id"`
			//创建时间
			CreateAt string `db:"create_at" json:"createAt"`
			//上级ID
			ParentID int64 `db:"parent_id" json:"parentID"`
			//用户信息
			UserName   string `json:"userName"`
			UserAvatar string `json:"userAvatar"`
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
				userData := UserCore.GetUserAndAvatar(v.UserID)
				newData = append(newData, dataType{
					ID:         v.ID,
					CreateAt:   CoreFilter.GetTimeToDefaultTime(v.CreateAt),
					ParentID:   v.ParentID,
					UserName:   userData.Name,
					UserAvatar: BaseFileSys2.GetPublicURLByClaimID(userData.Avatar),
					LevelType:  v.LevelType,
					Level:      v.Level,
					Title:      v.Title,
					Des:        CoreFilter.SubStrQuick(v.Des, coreCommentShowListDesLen),
					DesFiles:   BaseFileSys2.GetPublicURLsByClaimIDs(v.DesFiles),
					Params:     v.Params,
				})
			}
		}
		//反馈数据
		Router2Mid.ReportDataList(c, "get list", err, "", newData, dataCount)
	})
	//获取指定分类ID
	r.POST(fmt.Sprint(appendURL, "/id/:id"), func(c *Router2Mid.RouterURLHeaderC) {
		//日志
		c.LogAppend = "header url class comment get by id"
		//权限检查
		if len(permissions) > 0 {
			if !Router2Mid.CheckPermission(c, permissions) {
				return
			}
		}
		//获取参数
		id, b := Router2Params.GetID(c)
		if !b {
			return
		}
		//获取数据
		data := commentObj.GetByID(id)
		//重组数据
		type dataType struct {
			//基础
			ID int64 `db:"id" json:"id"`
			//创建时间
			CreateAt string `db:"create_at" json:"createAt"`
			//上级ID
			ParentID int64 `db:"parent_id" json:"parentID"`
			//用户信息
			UserName   string `json:"userName"`
			UserAvatar string `json:"userAvatar"`
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
		var newData dataType
		var err error
		if data.ID > 0 {
			userData := UserCore.GetUserAndAvatar(data.UserID)
			newData = dataType{
				ID:         data.ID,
				CreateAt:   CoreFilter.GetTimeToDefaultTime(data.CreateAt),
				ParentID:   data.ParentID,
				UserName:   userData.Name,
				UserAvatar: BaseFileSys2.GetPublicURLByClaimID(userData.Avatar),
				LevelType:  data.LevelType,
				Level:      data.Level,
				Title:      data.Title,
				Des:        data.Des,
				DesFiles:   BaseFileSys2.GetPublicURLsByClaimIDs(data.DesFiles),
				Params:     data.Params,
			}
		} else {
			err = errors.New("no data")
		}
		//反馈数据
		Router2Mid.ReportData(c, "get by id", err, "", newData)
	})
}
