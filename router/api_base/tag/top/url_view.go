package RouterAPIBaseTagTop

import (
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	RouterMidAPI "gitee.com/weeekj/weeekj_core/v5/router/mid/api"
	RouterParams "gitee.com/weeekj/weeekj_core/v5/router/params"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// URLView 组织层级的标签服务API
// 只读
// 注意，必须登录后才能使用本设计
// 除了指定routers外，还需要绑定分类对象权限列
func URLView(routers *gin.RouterGroup, tagObj *ClassTag.Tag, permission string) {
	//获取标签列表
	routers.POST("/list", func(c *gin.Context) {
		if permission != "" {
			//权限检查
			if !RouterMidAPI.CheckUserPermission(c, permission) {
				return
			}
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
		//获取数据
		dataList, dataCount, err := tagObj.GetList(&ClassTag.ArgsGetList{
			Pages:  params.Pages,
			BindID: 0,
			Search: params.Search,
		})
		//反馈数据
		RouterReport.DataList(c, "get tag list, ", "没有标签数据", err, dataList, dataCount)
	})
	//获取一组标签
	routers.POST("/more", func(c *gin.Context) {
		if permission != "" {
			//权限检查
			if !RouterMidAPI.CheckUserPermission(c, permission) {
				return
			}
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
		//获取数据
		dataList, err := tagObj.GetByIDs(&ClassTag.ArgsGetIDs{
			IDs:    params.IDs,
			BindID: 0,
			Limit:  100,
		})
		//反馈数据
		RouterReport.Data(c, "get tag list, ", "没有标签数据", err, dataList)
	})
}
