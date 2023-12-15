package RouterAPIBaseTagOrg

import (
	ClassTag "github.com/fotomxq/weeekj_core/v5/class/tag"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	RouterMidOrg "github.com/fotomxq/weeekj_core/v5/router/mid/org"
	RouterOrgCore "github.com/fotomxq/weeekj_core/v5/router/org/core"
	RouterParams "github.com/fotomxq/weeekj_core/v5/router/params"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// URLAll 组织层级的标签服务API
// 注意，必须登录和选择组织后才能使用本设计
// 除了指定routers外，还需要绑定分类对象和组织权限列
func URLAll(routers *gin.RouterGroup, tagObj *ClassTag.Tag, permissions []string) {
	//获取标签列表
	routers.POST("/list", func(c *gin.Context) {
		//权限检查
		if !RouterOrgCore.CheckPermissionByUser(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//分页
			Pages CoreSQLPages.ArgsDataList `json:"pages"`
			//搜索标签
			Search string `json:"search" check:"search" empty:"true"`
		}
		var params dataType
		if b := RouterParams.GetJSON(c, &params); !b {
			return
		}
		//获取组织
		orgData := RouterMidOrg.GetOrg(c)
		//获取数据
		dataList, dataCount, err := tagObj.GetList(&ClassTag.ArgsGetList{
			Pages:  params.Pages,
			BindID: orgData.ID,
			Search: params.Search,
		})
		//反馈数据
		RouterReport.DataList(c, "get tag list, ", "没有标签数据", err, dataList, dataCount)
	})
	//获取一组标签
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
		dataList, err := tagObj.GetByIDs(&ClassTag.ArgsGetIDs{
			IDs:    params.IDs,
			BindID: orgData.ID,
			Limit:  100,
		})
		//反馈数据
		RouterReport.Data(c, "get tag more, ", "没有标签数据", err, dataList)
	})
	//创建新的标签
	routers.PUT("", func(c *gin.Context) {
		//权限检查
		if !RouterOrgCore.CheckPermissionByUser(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//名称
			Name string `json:"name" check:"name"`
		}
		var params dataType
		if b := RouterParams.GetJSON(c, &params); !b {
			return
		}
		//获取组织
		orgData := RouterMidOrg.GetOrg(c)
		//获取数据
		data, err := tagObj.Create(&ClassTag.ArgsCreate{
			BindID: orgData.ID,
			Name:   params.Name,
		})
		if err == nil {
			RouterOrgCore.AddRecord(c, "创建新的标签(", params.Name, ")[", data.ID, "]")
		}
		//反馈数据
		RouterReport.ActionCreate(c, "create tag, ", "无法创建新的标签", err, data)
	})
	//修改标签
	routers.POST("/info", func(c *gin.Context) {
		//权限检查
		if !RouterOrgCore.CheckPermissionByUser(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//ID
			ID int64 `json:"id" check:"id"`
			//名称
			Name string `json:"name" check:"name"`
		}
		var params dataType
		if b := RouterParams.GetJSON(c, &params); !b {
			return
		}
		//获取组织
		orgData := RouterMidOrg.GetOrg(c)
		//获取数据
		err := tagObj.UpdateByID(&ClassTag.ArgsUpdateByID{
			ID:     params.ID,
			BindID: orgData.ID,
			Name:   params.Name,
		})
		if err == nil {
			RouterOrgCore.AddRecord(c, "修改标签(", params.Name, ")[", params.ID, "]")
		}
		//反馈数据
		RouterReport.ActionUpdate(c, "update tag, ", "无法修改标签", err)
	})
	//删除标签
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
		err := tagObj.DeleteByID(&ClassTag.ArgsDeleteByID{
			ID:     params.ID,
			BindID: orgData.ID,
		})
		if err == nil {
			RouterOrgCore.AddRecord(c, "删除标签[", params.ID, "]")
		}
		//反馈数据
		RouterReport.ActionDelete(c, "delete tag, ", "无法删除标签", err)
	})
}
