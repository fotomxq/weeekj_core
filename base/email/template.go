package BaseEmail

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"strings"
	"time"
)

// ArgsGetTemplateList 获取模板列表参数
type ArgsGetTemplateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTemplateList 获取模板列表
func GetTemplateList(args *ArgsGetTemplateList) (dataList []FieldsTemplate, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR content ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "core_email_template"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, server_ids, title, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

type ArgsGetTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func GetTemplate(args *ArgsGetTemplate) (data FieldsTemplate, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, server_ids, title, content, params FROM core_email_template WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	if err == nil && data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsSendTemplate 采用模板发送邮件参数
type ArgsSendTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//替换的数据集合
	ReplaceData CoreSQLConfig.FieldsConfigsType `json:"replaceData"`
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom `json:"createInfo"`
	//计划送达时间
	SendAt time.Time `json:"sendAt"`
	//发送邮箱地址列
	ToEmailList []string `json:"toEmailList"`
}

// SendTemplate 采用模板发送邮件
func SendTemplate(args *ArgsSendTemplate) (err error) {
	//获取模版
	var templateData FieldsTemplate
	err = Router2SystemConfig.MainDB.Get(&templateData, "SELECT id, server_ids, title, content, params FROM core_email_template WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	if err != nil || templateData.ID < 1 {
		err = errors.New(fmt.Sprint("no data, ", err))
		return
	}
	//修正内容
	var content = templateData.Content + ""
	for _, v := range args.ReplaceData {
		content = strings.ReplaceAll(content, v.Mark, v.Val)
	}
	//遍历邮箱服务
	var findServerID int64 = 0
	if len(templateData.ServerIDs) == 1 {
		findServerID = templateData.ServerIDs[0]
	} else {
		key := CoreFilter.GetRandNumber(0, len(templateData.ServerIDs)-1)
		for vKey, vServerID := range templateData.ServerIDs {
			if vKey == key {
				findServerID = vServerID
				break
			}
		}
	}
	if findServerID < 1 {
		if len(templateData.ServerIDs) > 0 {
			findServerID = templateData.ServerIDs[0]
		} else {
			err = errors.New("no find server id")
			return
		}
	}
	for _, vToEmail := range args.ToEmailList {
		//根据邮箱配置，推送数据
		_, err = Send(&ArgsSend{
			OrgID:      templateData.OrgID,
			ServerID:   findServerID,
			CreateInfo: args.CreateInfo,
			SendAt:     args.SendAt,
			ToEmail:    vToEmail,
			Title:      templateData.Title,
			Content:    content,
			IsHtml:     true,
		})
		if err != nil {
			CoreLog.Warn("core email send by template failed, ", err)
			err = nil
			continue
		}
	}
	return
}

// ArgsCreateTemplate 创建模板参数
type ArgsCreateTemplate struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//计划使用的邮箱配置列
	// 多个配置将随机抽取一个发送
	ServerIDs pq.Int64Array `db:"server_ids" json:"serverIDs" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//内容
	// 将强制邮件采用HTML模式发送，此处存放HTML内容
	// 相关变量根据模块的约定执行
	Content string `db:"content" json:"content" check:"des" min:"1" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateTemplate 创建模板
func CreateTemplate(args *ArgsCreateTemplate) (err error) {
	//检查配置列
	if err = checkServerOrg(args.ServerIDs, args.OrgID); err != nil {
		return
	}
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_email_template (org_id, server_ids, title, content, params) VALUES (:org_id,:server_ids,:title,:content,:params)", args)
	return
}

// ArgsUpdateTemplate 修改模板参数
type ArgsUpdateTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//计划使用的邮箱配置列
	// 多个配置将随机抽取一个发送
	ServerIDs pq.Int64Array `db:"server_ids" json:"serverIDs" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//内容
	// 将强制邮件采用HTML模式发送，此处存放HTML内容
	// 相关变量根据模块的约定执行
	Content string `db:"content" json:"content" check:"des" min:"1" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateTemplate 修改模板
func UpdateTemplate(args *ArgsUpdateTemplate) (err error) {
	//检查配置列
	if err = checkServerOrg(args.ServerIDs, args.OrgID); err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_email_template SET update_at = NOW(), server_ids = :server_ids, title = :title, content = :content, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteTemplate 删除模板参数
type ArgsDeleteTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteTemplate 删除模板
func DeleteTemplate(args *ArgsDeleteTemplate) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "core_email_template", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// 检查配置组织关系
func checkServerOrg(serverIDs pq.Int64Array, orgID int64) (err error) {
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
	}
	if len(serverIDs) > 0 {
		var dataList []dataType
		err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM core_email_server WHERE id = ANY($1) AND ($2 < 1 OR org_id = $2)", serverIDs, orgID)
		if err != nil || len(dataList) < 1 {
			err = errors.New(fmt.Sprint("no data, ids: ", serverIDs, ", err: ", err))
			return
		}
		if len(dataList) != len(serverIDs) {
			err = errors.New("have no org data")
			return
		}
		return
	} else {
		return
	}
}
