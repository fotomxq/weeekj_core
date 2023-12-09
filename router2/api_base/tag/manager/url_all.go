package Router2APIBaseTagManager

import (
	"fmt"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	Router2Params "gitee.com/weeekj/weeekj_core/v5/router2/params"
	Router2Record "gitee.com/weeekj/weeekj_core/v5/router2/record"
	"github.com/lib/pq"
)

// URLAll 管理层级操作接口汇总
func URLAll(r *Router2Mid.RouterURLUser, appendURL string, tagObj *ClassTag.Tag, permissions []string, haveView bool) {
	if haveView {
		//获取标签列表
		r.POST(fmt.Sprint(appendURL, "/list"), func(c *Router2Mid.RouterURLUserC) {
			//日志
			c.LogAppend = "manager url class tag get list"
			//权限检查
			if !Router2Mid.CheckPermission(c, permissions) {
				return
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
		r.POST(fmt.Sprint(appendURL, "/more"), func(c *Router2Mid.RouterURLUserC) {
			//日志
			c.LogAppend = "manager url class tag get more"
			//权限检查
			if !Router2Mid.CheckPermission(c, permissions) {
				return
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
	//创建新的标签
	r.PUT(fmt.Sprint(appendURL, ""), func(c *Router2Mid.RouterURLUserC) {
		//日志
		c.LogAppend = "manager url class tag create tag"
		//权限检查
		if !Router2Mid.CheckPermission(c, permissions) {
			return
		}
		//获取参数
		type dataType struct {
			//名称
			Name string `json:"name" check:"name"`
		}
		var params dataType
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		data, err := tagObj.Create(&ClassTag.ArgsCreate{
			BindID: 0,
			Name:   params.Name,
		})
		if err == nil {
			Router2Record.AddRecord(c, "class_tag", data.ID, "创建新的标签(", params.Name, ")[", data.ID, "]")
		}
		//反馈数据
		Router2Mid.ReportActionCreateNoData(c, "create failed", err, "")
	})
	//修改标签
	r.POST(fmt.Sprint(appendURL, "/info"), func(c *Router2Mid.RouterURLUserC) {
		//日志
		c.LogAppend = "manager url class tag update tag"
		//权限检查
		if !Router2Mid.CheckPermission(c, permissions) {
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
		if b := Router2Params.GetJSON(c, &params); !b {
			return
		}
		//获取数据
		err := tagObj.UpdateByID(&ClassTag.ArgsUpdateByID{
			ID:     params.ID,
			BindID: 0,
			Name:   params.Name,
		})
		if err == nil {
			Router2Record.AddRecord(c, "class_tag", params.ID, "修改标签(", params.Name, ")[", params.ID, "]")
		}
		//反馈数据
		Router2Mid.ReportActionUpdate(c, "update failed", err, "")
	})
	//删除标签
	r.DELETE(fmt.Sprint(appendURL, "/id"), func(c *Router2Mid.RouterURLUserC) {
		//日志
		c.LogAppend = "manager url class tag delete tag"
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
		err := tagObj.DeleteByID(&ClassTag.ArgsDeleteByID{
			ID:     params.ID,
			BindID: 0,
		})
		if err == nil {
			Router2Record.AddRecord(c, "class_tag", params.ID, "删除标签[", params.ID, "]")
		}
		//反馈数据
		Router2Mid.ReportActionDelete(c, "delete failed", err, "")
	})
}
