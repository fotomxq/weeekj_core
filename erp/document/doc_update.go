package ERPDocument

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ERPCore "github.com/fotomxq/weeekj_core/v5/erp/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateDoc 更新配置信息参数
type ArgsUpdateDoc struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述
	// 根据文档格式决定，默认采用html富文本形式记录数据
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//内容
	DataList []ERPCore.ArgsComponentValSetOnlyUpdate `json:"dataList"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateDoc 更新配置信息
func UpdateDoc(args *ArgsUpdateDoc) (errCode string, err error) {
	//声明隐藏字段并且赋值
	var searchDes string
	for _, v := range args.DataList {
		searchDes += v.Val
	}
	//获取数据
	data := getDocByID(args.ID)
	if data.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	configData := getConfigByID(data.ConfigID)
	if configData.ID < 1 {
		errCode = "err_erp_document_config_not_exist"
		err = errors.New("no config data")
		return
	}
	//检查是否发布
	if !checkConfigPublish(data.ConfigID, args.OrgID) {
		err = errors.New("config no publish")
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_document_doc SET update_at = NOW(), name = :name, des = :des, cover_file_id = :cover_file_id, search_des = :search_des, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":            args.ID,
		"org_id":        args.OrgID,
		"name":          args.Name,
		"des":           args.Des,
		"cover_file_id": args.CoverFileID,
		"search_des":    searchDes,
		"params":        args.Params,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	//设置附加数据
	switch configData.DocType {
	case "custom":
		//修改组件列
		for k, v := range configData.ComponentList {
			for _, v2 := range args.DataList {
				if v.Key == v2.Key {
					configData.ComponentList[k].Val = v2.Val
					configData.ComponentList[k].Params = v2.Params
					break
				}
			}
		}
		err = docComponentValObj.SetMore(&ERPCore.ArgsSetMore{
			BindID:   args.ID,
			DataList: configData.ComponentList,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
	case "doc":
	case "excel":
		//获取excel配置
		// 通过其他接口存储
	}
	//删除缓冲
	deleteDocCache(args.ID)
	//反馈
	return
}
