package RouterAPIBaseSortOrg

import (
	BaseQiniu "gitee.com/weeekj/weeekj_core/v5/base/qiniu"
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	RouterMidOrg "gitee.com/weeekj/weeekj_core/v5/router/mid/org"
	RouterOrgCore "gitee.com/weeekj/weeekj_core/v5/router/org/core"
	RouterParams "gitee.com/weeekj/weeekj_core/v5/router/params"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// URLAll 组织层级的分类服务API
// 注意，必须登录和选择组织后才能使用本设计
// 除了指定routers外，还需要绑定分类对象和组织权限列
func URLAll(routers *gin.RouterGroup, sortObj *ClassSort.Sort, permissions []string) {
	//获取分类列表
	routers.POST("/list", func(c *gin.Context) {
		//权限检查
		if !RouterOrgCore.CheckPermissionByUser(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//分页
			Pages CoreSQLPages.ArgsDataList `json:"pages"`
			//标识码
			Mark string `json:"mark" check:"mark" empty:"true"`
			//上级ID
			ParentID int64 `json:"parentID" check:"id" empty:"true"`
			//搜索
			Search string `json:"search" check:"search" empty:"true"`
		}
		var params dataType
		if b := RouterParams.GetJSON(c, &params); !b {
			return
		}
		//获取组织
		orgData := RouterMidOrg.GetOrg(c)
		//获取数据
		dataList, dataCount, err := sortObj.GetList(&ClassSort.ArgsGetList{
			Pages:    params.Pages,
			BindID:   orgData.ID,
			Mark:     params.Mark,
			ParentID: params.ParentID,
			Search:   params.Search,
		})
		//扩展查询支持
		type newDataType struct {
			//主要数据集
			DataList []ClassSort.FieldsSort `json:"dataList"`
			//文件集
			// id => url
			FileList map[int64]string `json:"fileList"`
			//上级ID数据列
			ParentList []ClassSort.FieldsSort `json:"parentList"`
		}
		var newDataList newDataType
		if err == nil {
			newDataList.DataList = dataList
			var waitFiles, waitParents []int64
			for _, v := range dataList {
				if v.CoverFileID > 0 {
					waitFiles = append(waitFiles, v.CoverFileID)
				}
				for _, v := range v.DesFiles {
					waitFiles = append(waitFiles, v)
				}
				if v.ParentID > 0 {
					waitParents = append(waitParents, v.ParentID)
				}
			}
			if len(waitFiles) > 0 {
				newDataList.FileList, err = BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
					ClaimIDList: waitFiles,
					UserID:      0,
					OrgID:       0,
					IsPublic:    true,
				})
				if err != nil {
					//跳过错误
					err = nil
				}
			}
			if len(waitParents) > 0 {
				newDataList.ParentList, err = sortObj.GetByIDs(&ClassSort.ArgsGetIDs{
					IDs:    waitParents,
					BindID: 0,
					Limit:  999,
				})
				if err != nil {
					//跳过错误
					err = nil
				}
			}
		}
		//反馈数据
		RouterReport.DataList(c, "get sort list, ", "没有分类数据", err, newDataList, dataCount)
	})
	//获取一组分类
	routers.POST("/more", func(c *gin.Context) {
		//权限检查
		if !RouterOrgCore.CheckPermissionByUser(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//ID列
			IDs pq.Int64Array `json:"ids" check:"ids"`
		}
		var params dataType
		if b := RouterParams.GetJSON(c, &params); !b {
			return
		}
		//获取组织
		orgData := RouterMidOrg.GetOrg(c)
		//获取数据
		dataList, err := sortObj.GetByIDs(&ClassSort.ArgsGetIDs{
			IDs:    params.IDs,
			BindID: orgData.ID,
			Limit:  100,
		})
		//扩展查询支持
		type newDataType struct {
			//主要数据集
			DataList []ClassSort.FieldsSort `json:"dataList"`
			//文件集
			// id => url
			FileList map[int64]string `json:"fileList"`
			//上级ID数据列
			ParentList []ClassSort.FieldsSort `json:"parentList"`
		}
		var newDataList newDataType
		if err == nil {
			newDataList.DataList = dataList
			var waitFiles, waitParents []int64
			for _, v := range dataList {
				if v.CoverFileID > 0 {
					waitFiles = append(waitFiles, v.CoverFileID)
				}
				for _, v := range v.DesFiles {
					waitFiles = append(waitFiles, v)
				}
				if v.ParentID > 0 {
					waitParents = append(waitParents, v.ParentID)
				}
			}
			if len(waitFiles) > 0 {
				newDataList.FileList, err = BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
					ClaimIDList: waitFiles,
					UserID:      0,
					OrgID:       0,
					IsPublic:    true,
				})
				if err != nil {
					//跳过错误
					err = nil
				}
			}
			if len(waitParents) > 0 {
				newDataList.ParentList, err = sortObj.GetByIDs(&ClassSort.ArgsGetIDs{
					IDs:    waitParents,
					BindID: 0,
					Limit:  999,
				})
				if err != nil {
					//跳过错误
					err = nil
				}
			}
		}
		//反馈数据
		RouterReport.Data(c, "get sort more, ", "没有分类数据", err, newDataList)
	})
	//创建新的分类
	routers.PUT("", func(c *gin.Context) {
		//权限检查
		if !RouterOrgCore.CheckPermissionByUser(c, permissions) {
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
		if b := RouterParams.GetJSON(c, &params); !b {
			return
		}
		//获取组织
		orgData := RouterMidOrg.GetOrg(c)
		//获取数据
		data, err := sortObj.Create(&ClassSort.ArgsCreate{
			BindID:      orgData.ID,
			Mark:        params.Mark,
			ParentID:    params.ParentID,
			CoverFileID: params.CoverFileID,
			DesFiles:    params.DesFiles,
			Name:        params.Name,
			Des:         params.Des,
			Params:      params.Params,
		})
		if err == nil {
			RouterOrgCore.AddRecord(c, "创建新的分类(", params.Name, ")[", data.ID, "]")
		}
		//反馈数据
		RouterReport.ActionCreate(c, "create sort, ", "无法创建新的分类", err, data)
	})
	//修改分类
	routers.POST("/info", func(c *gin.Context) {
		//权限检查
		if !RouterOrgCore.CheckPermissionByUser(c, permissions) {
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
		if b := RouterParams.GetJSON(c, &params); !b {
			return
		}
		//获取组织
		orgData := RouterMidOrg.GetOrg(c)
		//获取数据
		err := sortObj.UpdateByID(&ClassSort.ArgsUpdateByID{
			ID:          params.ID,
			BindID:      orgData.ID,
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
			RouterOrgCore.AddRecord(c, "修改分类(", params.Name, ")[", params.ID, "]")
		}
		//反馈数据
		RouterReport.ActionUpdate(c, "update sort, ", "无法修改分类", err)
	})
	//删除分类
	routers.DELETE("/id", func(c *gin.Context) {
		//权限检查
		if !RouterOrgCore.CheckPermissionByUser(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//ID
			ID int64 `json:"id" check:"id"`
		}
		var params dataType
		if b := RouterParams.GetJSON(c, &params); !b {
			return
		}
		//获取组织
		orgData := RouterMidOrg.GetOrg(c)
		//获取数据
		err := sortObj.DeleteByID(&ClassSort.ArgsDeleteByID{
			ID:     params.ID,
			BindID: orgData.ID,
		})
		if err == nil {
			RouterOrgCore.AddRecord(c, "删除分类[", params.ID, "]")
		}
		//反馈数据
		RouterReport.ActionDelete(c, "delete sort, ", "无法删除分类", err)
	})
}
