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

// URLView 组织层级的标签服务API
// 只读
// 注意，必须登录和选择组织后才能使用本设计
// 除了指定routers外，还需要绑定分类对象和组织权限列
func URLView(routers *gin.RouterGroup, tagObj *ClassTag.Tag, permissions []string) {
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
}
