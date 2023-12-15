package ERPDocument

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ERPCore "github.com/fotomxq/weeekj_core/v5/erp/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsCreateDoc struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
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

func CreateDoc(args *ArgsCreateDoc) (errCode string, err error) {
	//声明隐藏字段并且赋值
	var searchDes string
	for _, v := range args.DataList {
		searchDes += v.Val
	}
	//获取配置
	configData := GetConfig(args.ConfigID, args.OrgID)
	if configData.ID < 1 {
		errCode = "err_erp_document_config_not_exist"
		err = errors.New("no config data")
		return
	}
	//检查是否发布
	if !checkConfigPublish(args.ConfigID, args.OrgID) {
		errCode = "err_erp_document_config_not_exist"
		err = errors.New("config no publish")
		return
	}
	//写入数据
	var newID int64
	newID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_document_doc (config_id, org_id, name, des, cover_file_id, search_des, params) VALUES (:config_id, :org_id, :name, :des, :cover_file_id, :search_des, :params)", map[string]interface{}{
		"config_id":     configData.ID,
		"org_id":        args.OrgID,
		"name":          args.Name,
		"des":           args.Des,
		"cover_file_id": args.CoverFileID,
		"search_des":    searchDes,
		"params":        args.Params,
	})
	if err != nil {
		errCode = "err_insert"
		err = errors.New(fmt.Sprint("create doc failed, config id: ", configData.ID, ", name: ", args.Name, ", err: ", err))
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
			BindID:   newID,
			DataList: configData.ComponentList,
		})
		if err != nil {
			errCode = "err_update"
			err = errors.New(fmt.Sprint("set component val, doc id: ", newID, ", config id: ", configData.ID, ", err: ", err))
			return
		}
	case "doc":
	case "excel":
		//获取excel配置
		err = setExcelByConfigID(configData.ID, newID)
		if err != nil {
			errCode = "err_update"
			err = errors.New(fmt.Sprint("set excel val, doc id: ", newID, ", config id: ", configData.ID, ", err: ", err))
			return
		}
	}
	//反馈
	return
}
