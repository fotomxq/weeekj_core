package Router2APIBaseSortHeader

import (
	"fmt"
	BaseFileSys2 "github.com/fotomxq/weeekj_core/v5/base/filesys2"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	Router2Params "github.com/fotomxq/weeekj_core/v5/router2/params"
	"github.com/lib/pq"
)

// URLView 通用查看分类数据包
func URLView(r *Router2Mid.RouterURLHeader, appendURL string, sortObj *ClassSort.Sort, permissions []string) {
	//获取分类列表
	r.POST(fmt.Sprint(appendURL, "/list"), func(c *Router2Mid.RouterURLHeaderC) {
		//日志
		c.LogAppend = "header url class sort get list"
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
			//绑定ID
			BindID int64 `json:"bindID" check:"id" empty:"true"`
			//标识码
			Mark string `json:"mark" check:"mark" empty:"true"`
			//上级ID
			ParentID int64 `json:"parentID" check:"id" empty:"true"`
			//搜索
			Search string `json:"search" check:"search" empty:"true"`
		}
		var params paramsType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		dataList, dataCount, err := sortObj.GetList(&ClassSort.ArgsGetList{
			Pages:    params.Pages,
			BindID:   params.BindID,
			Mark:     params.Mark,
			ParentID: params.ParentID,
			Search:   params.Search,
		})
		//重组数据
		type dataType struct {
			//基础
			ID int64 `db:"id" json:"id"`
			//分组标识码
			// 用于一些特殊的显示场景做区分，可以重复
			Mark string `db:"mark" json:"mark"`
			//上级ID
			ParentID int64 `db:"parent_id" json:"parentID"`
			//排序
			Sort int `db:"sort" json:"sort"`
			//封面图
			CoverFileURL string `json:"coverFileURL"`
			//名称
			Name string `db:"name" json:"name"`
		}
		var newData []dataType
		if err == nil {
			for _, v := range dataList {
				newData = append(newData, dataType{
					ID:           v.ID,
					Mark:         v.Mark,
					ParentID:     v.ParentID,
					Sort:         v.Sort,
					CoverFileURL: BaseFileSys2.GetPublicURLByClaimID(v.CoverFileID),
					Name:         v.Name,
				})
			}
		}
		//反馈数据
		Router2Mid.ReportDataList(c, "get list", err, "", newData, dataCount)
	})
	//获取一组分类
	r.POST(fmt.Sprint(appendURL, "/more"), func(c *Router2Mid.RouterURLHeaderC) {
		//日志
		c.LogAppend = "header url class sort get more"
		//权限检查
		if len(permissions) > 0 {
			if !Router2Mid.CheckPermission(c, permissions) {
				return
			}
		}
		//获取参数
		type paramsType struct {
			//ID列
			IDs pq.Int64Array `json:"ids" check:"ids"`
			//绑定ID
			BindID int64 `json:"bindID" check:"id" empty:"true"`
		}
		var params paramsType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		dataList, err := sortObj.GetByIDs(&ClassSort.ArgsGetIDs{
			IDs:    params.IDs,
			BindID: params.BindID,
			Limit:  100,
		})
		//重组数据
		type dataType struct {
			//基础
			ID int64 `db:"id" json:"id"`
			//分组标识码
			// 用于一些特殊的显示场景做区分，可以重复
			Mark string `db:"mark" json:"mark"`
			//上级ID
			ParentID int64 `db:"parent_id" json:"parentID"`
			//排序
			Sort int `db:"sort" json:"sort"`
			//封面图
			CoverFileURL string `json:"coverFileURL"`
			//名称
			Name string `db:"name" json:"name"`
			//介绍图文
			DesFiles []string `db:"des_files" json:"desFiles"`
			//描述
			Des string `db:"des" json:"des"`
			//扩展参数
			Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
		}
		var newData []dataType
		if err == nil {
			for _, v := range dataList {
				newData = append(newData, dataType{
					ID:           v.ID,
					Mark:         v.Mark,
					ParentID:     v.ParentID,
					Sort:         v.Sort,
					CoverFileURL: BaseFileSys2.GetPublicURLByClaimID(v.CoverFileID),
					Name:         v.Name,
					DesFiles:     BaseFileSys2.GetPublicURLsByClaimIDs(v.DesFiles),
					Des:          v.Des,
					Params:       v.Params,
				})
			}
		}
		//反馈数据
		Router2Mid.ReportData(c, "get more", err, "", newData)
	})
	//获取指定分类ID
	r.POST(fmt.Sprint(appendURL, "/id"), func(c *Router2Mid.RouterURLHeaderC) {
		//日志
		c.LogAppend = "header url class sort get by id"
		//权限检查
		if len(permissions) > 0 {
			if !Router2Mid.CheckPermission(c, permissions) {
				return
			}
		}
		//获取参数
		type paramsType struct {
			//ID
			ID int64 `db:"id" json:"id" check:"id"`
			//绑定ID
			BindID int64 `json:"bindID" check:"id" empty:"true"`
		}
		var params paramsType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		data, err := sortObj.GetByID(&ClassSort.ArgsGetByID{
			ID:     params.ID,
			BindID: params.BindID,
		})
		//重组数据
		type dataType struct {
			//基础
			ID int64 `db:"id" json:"id"`
			//创建时间
			CreateAt string `db:"create_at" json:"createAt"`
			//更新时间
			UpdateAt string `db:"update_at" json:"updateAt"`
			//来源ID
			// 可以是某个组织，或特定的其他系统ID
			BindID int64 `db:"bind_id" json:"bindID"`
			//分组标识码
			// 用于一些特殊的显示场景做区分，可以重复
			Mark string `db:"mark" json:"mark"`
			//上级ID
			ParentID int64 `db:"parent_id" json:"parentID"`
			//排序
			Sort int `db:"sort" json:"sort"`
			//封面图
			CoverFileID  int64  `db:"cover_file_id" json:"coverFileID"`
			CoverFileURL string `json:"coverFileURL"`
			//介绍图文
			DesFileIDs []int64          `json:"desFileIDs"`
			DesFiles   map[int64]string `db:"des_files" json:"desFiles"`
			//名称
			Name string `db:"name" json:"name"`
			//描述
			Des string `db:"des" json:"des"`
			//扩展参数
			Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
		}
		var newData dataType
		if err == nil {
			newData = dataType{
				ID:           data.ID,
				CreateAt:     CoreFilter.GetTimeToDefaultTime(data.CreateAt),
				UpdateAt:     CoreFilter.GetTimeToDefaultTime(data.UpdateAt),
				BindID:       data.BindID,
				Mark:         data.Mark,
				ParentID:     data.ParentID,
				Sort:         data.Sort,
				CoverFileID:  data.CoverFileID,
				CoverFileURL: BaseFileSys2.GetPublicURLByClaimID(data.CoverFileID),
				DesFileIDs:   data.DesFiles,
				DesFiles:     BaseFileSys2.GetPublicURLMapsByClaimIDsTo(data.DesFiles),
				Name:         data.Name,
				Des:          data.Des,
				Params:       data.Params,
			}
		}
		//反馈数据
		Router2Mid.ReportData(c, "get by id", err, "", newData)
	})
}
