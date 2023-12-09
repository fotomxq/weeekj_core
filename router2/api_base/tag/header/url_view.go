package Router2APIBaseTagHeader

import (
	"fmt"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	Router2Params "gitee.com/weeekj/weeekj_core/v5/router2/params"
	"github.com/lib/pq"
)

// URLView 通用查看标签列表
func URLView(r *Router2Mid.RouterURLHeader, appendURL string, tagObj *ClassTag.Tag, permissions []string) {
	//获取标签列表
	r.POST(fmt.Sprint(appendURL, "/list"), func(c *Router2Mid.RouterURLHeaderC) {
		//日志
		c.LogAppend = "header url class tag get list"
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
			//搜索
			Search string `json:"search" check:"search" empty:"true"`
		}
		var params paramsType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		dataList, dataCount, err := tagObj.GetList(&ClassTag.ArgsGetList{
			Pages:  params.Pages,
			BindID: params.BindID,
			Search: params.Search,
		})
		//重组数据
		type dataType struct {
			//基础
			ID int64 `db:"id" json:"id"`
			//名称
			Name string `db:"name" json:"name"`
		}
		var newData []dataType
		if err == nil {
			for _, v := range dataList {
				newData = append(newData, dataType{
					ID:   v.ID,
					Name: v.Name,
				})
			}
		}
		//反馈数据
		Router2Mid.ReportDataList(c, "get list", err, "", newData, dataCount)
	})
	//获取一组标签
	r.POST(fmt.Sprint(appendURL, "/more"), func(c *Router2Mid.RouterURLHeaderC) {
		//日志
		c.LogAppend = "header url class tag get more"
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
		dataList, err := tagObj.GetByIDs(&ClassTag.ArgsGetIDs{
			IDs:    params.IDs,
			BindID: params.BindID,
			Limit:  100,
		})
		//重组数据
		type dataType struct {
			//基础
			ID int64 `db:"id" json:"id"`
			//名称
			Name string `db:"name" json:"name"`
		}
		var newData []dataType
		if err == nil {
			for _, v := range dataList {
				newData = append(newData, dataType{
					ID:   v.ID,
					Name: v.Name,
				})
			}
		}
		//反馈数据
		Router2Mid.ReportData(c, "get more", err, "", newData)
	})
}
