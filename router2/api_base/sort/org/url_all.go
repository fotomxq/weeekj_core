package Router2APIBaseSortOrg

import (
	"fmt"
	BaseFileSys2 "gitee.com/weeekj/weeekj_core/v5/base/filesys2"
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	Router2Params "gitee.com/weeekj/weeekj_core/v5/router2/params"
	Router2Record "gitee.com/weeekj/weeekj_core/v5/router2/record"
	"github.com/lib/pq"
)

// URLAll 组织层级的分类服务API
// 注意，必须登录和选择组织后才能使用本设计
// 除了指定routers外，还需要绑定分类对象和组织权限列
func URLAll(r *Router2Mid.RouterURLOrg, appendURL string, sortObj *ClassSort.Sort, permissions []string, haveView bool) {
	if haveView {
		//获取分类列表
		r.POST(fmt.Sprint(appendURL, "/list"), func(c *Router2Mid.RouterURLOrgC) {
			//日志
			c.LogAppend = "om url class sort get list"
			//权限检查
			if !Router2Mid.CheckPermission(c, permissions) {
				return
			}
			//获取参数
			type paramsType struct {
				//分页
				Pages CoreSQLPages.ArgsDataList `json:"pages"`
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
				BindID:   c.OrgID,
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
		r.POST(fmt.Sprint(appendURL, "/more"), func(c *Router2Mid.RouterURLOrgC) {
			//日志
			c.LogAppend = "om url class sort get more"
			//权限检查
			if !Router2Mid.CheckPermission(c, permissions) {
				return
			}
			//获取参数
			type paramsType struct {
				//ID列
				IDs pq.Int64Array `json:"ids" check:"ids"`
			}
			var params paramsType
			if b := Router2Params.GetJSON(c, &params); !b {
				return
			}
			//获取数据
			dataList, err := sortObj.GetByIDs(&ClassSort.ArgsGetIDs{
				IDs:    params.IDs,
				BindID: c.OrgID,
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
		r.POST(fmt.Sprint(appendURL, "/id"), func(c *Router2Mid.RouterURLOrgC) {
			//日志
			c.LogAppend = "om url class sort get by id"
			//权限检查
			if !Router2Mid.CheckPermission(c, permissions) {
				return
			}
			//获取参数
			type paramsType struct {
				//ID
				ID int64 `db:"id" json:"id" check:"id"`
			}
			var params paramsType
			if b := Router2Params.GetJSON(c, &params); !b {
				return
			}
			//获取数据
			data, err := sortObj.GetByID(&ClassSort.ArgsGetByID{
				ID:     params.ID,
				BindID: c.OrgID,
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
	//创建新的分类
	r.PUT(fmt.Sprint(appendURL, ""), func(c *Router2Mid.RouterURLOrgC) {
		//日志
		c.LogAppend = "om url class sort create sort"
		//权限检查
		if !Router2Mid.CheckPermission(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//分组标识码
			Mark string `db:"mark" json:"mark"  check:"mark" empty:"true"`
			//上级ID
			ParentID int64 `json:"parentID" check:"id" empty:"true"`
			//封面图
			CoverFileID int64 `json:"coverFileID" check:"id" empty:"true"`
			//介绍图文
			DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
			//名称
			Name string `json:"name" check:"name"`
			//描述
			Des string `json:"des" check:"des" min:"1" max:"3000" empty:"true"`
			//扩展参数
			Params CoreSQLConfig.FieldsConfigsType `json:"params"`
		}
		var params dataType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		data, err := sortObj.Create(&ClassSort.ArgsCreate{
			BindID:      c.OrgID,
			Mark:        params.Mark,
			ParentID:    params.ParentID,
			CoverFileID: params.CoverFileID,
			DesFiles:    params.DesFiles,
			Name:        params.Name,
			Des:         params.Des,
			Params:      params.Params,
		})
		if err == nil {
			Router2Record.AddRecord(c, "class_sort", data.ID, "创建新的分类(", params.Name, ")[", data.ID, "]")
		}
		//反馈数据
		Router2Mid.ReportActionCreateNoData(c, "create failed", err, "")
	})
	//修改分类
	r.POST(fmt.Sprint(appendURL, "/info"), func(c *Router2Mid.RouterURLOrgC) {
		//日志
		c.LogAppend = "om url class sort update sort"
		//权限检查
		if !Router2Mid.CheckPermission(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//ID
			ID int64 `json:"id" check:"id"`
			//分组标识码
			Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
			//上级ID
			ParentID int64 `json:"parentID" check:"id" empty:"true"`
			//排序
			Sort int `json:"sort"`
			//封面图
			CoverFileID int64 `json:"coverFileID" check:"id" empty:"true"`
			//介绍图文
			DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
			//名称
			Name string `json:"name" check:"name"`
			//描述
			Des string `json:"des" check:"des" min:"1" max:"3000" empty:"true"`
			//扩展参数
			Params CoreSQLConfig.FieldsConfigsType `json:"params"`
		}
		var params dataType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		err := sortObj.UpdateByID(&ClassSort.ArgsUpdateByID{
			ID:          params.ID,
			BindID:      c.OrgID,
			Mark:        params.Mark,
			ParentID:    params.ParentID,
			Sort:        params.Sort,
			CoverFileID: params.CoverFileID,
			DesFiles:    params.DesFiles,
			Name:        params.Name,
			Des:         params.Des,
			Params:      params.Params,
		})
		if err == nil {
			Router2Record.AddRecord(c, "class_sort", params.ID, "修改分类(", params.Name, ")[", params.ID, "]")
		}
		//反馈数据
		Router2Mid.ReportActionUpdate(c, "update failed", err, "")
	})
	//删除分类
	r.DELETE(fmt.Sprint(appendURL, "/id"), func(c *Router2Mid.RouterURLOrgC) {
		//日志
		c.LogAppend = "om url class sort delete sort"
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
		err := sortObj.DeleteByID(&ClassSort.ArgsDeleteByID{
			ID:     params.ID,
			BindID: c.OrgID,
		})
		if err == nil {
			Router2Record.AddRecord(c, "class_sort", params.ID, "删除分类[", params.ID, "]")
		}
		//反馈数据
		Router2Mid.ReportActionDelete(c, "delete failed", err, "")
	})
}
