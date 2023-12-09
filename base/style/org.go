package BaseStyle

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// 获取组织的样式库
type ArgsGetOrgList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//样式ID
	StyleID int64 `db:"style_id" json:"styleID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetOrgList(args *ArgsGetOrgList) (dataList []FieldsOrg, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.StyleID > 0 {
		where = where + " AND style_id = :style_id"
		maps["style_id"] = args.StyleID
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_style_org",
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, style_id, components, title, cover_file_id, des_files FROM core_style_org WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// 获取ID
type ArgsGetOrgByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func GetOrgByID(args *ArgsGetOrgByID) (data FieldsOrg, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, style_id, components, title, des, cover_file_id, des_files, params FROM core_style_org WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	return
}

// ArgsGetOrgByStyleMark 获取指定页面Mark的样式数据参数
type ArgsGetOrgByStyleMark struct {
	//样式mark
	Mark string `db:"mark" json:"mark" check:"mark"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

type DataGetOrgByStyleMarkComponents struct {
	//组件ID
	ComponentID int64 `json:"componentID"`
	//关联标识码
	// 必填
	// 页面内独特的代码片段，声明后可直接定义该组件的默认参数形式
	Mark string `db:"mark" json:"mark"`
	//组件名称
	Name string `db:"name" json:"name"`
	//组件介绍
	Des string `db:"des" json:"des"`
	//组件封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//组件描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

type DataGetOrgByStyleMark struct {
	//样式基础
	Name string `json:"name"`
	//关联标识码
	// 用于识别代码片段
	Mark string `db:"mark" json:"mark"`
	//样式使用渠道
	// app APP；wxx 小程序等，可以任意定义，模块内不做限制
	SystemMark string `db:"system_mark" json:"systemMark"`
	//默认标题
	// 标题是展示给用户的，样式库名称和该标题不是一个
	Title string `db:"title" json:"title"`
	//默认描述
	Des string `db:"des" json:"des"`
	//默认封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//默认描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//组件列
	Components []DataGetOrgByStyleMarkComponents `json:"components"`
}

// GetOrgByStyleMark 获取指定页面Mark的样式数据
// 如果组织没有自定义，则按照默认构建数据并反馈
func GetOrgByStyleMark(args *ArgsGetOrgByStyleMark) (data DataGetOrgByStyleMark, err error) {
	var styleData FieldsStyle
	styleData, err = GetStyleByMark(&ArgsGetStyleByMark{
		Mark: args.Mark,
	})
	if err != nil {
		return
	}
	var orgData FieldsOrg
	err = Router2SystemConfig.MainDB.Get(&orgData, "SELECT id, org_id, style_id, components, title, des, cover_file_id, des_files, params FROM core_style_org WHERE style_id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", styleData.ID, args.OrgID)
	if err == nil {
		data = DataGetOrgByStyleMark{
			Name:        styleData.Name,
			Mark:        styleData.Mark,
			SystemMark:  styleData.SystemMark,
			Title:       styleData.Title,
			Des:         styleData.Des,
			CoverFileID: styleData.CoverFileID,
			DesFiles:    styleData.DesFiles,
			Params:      styleData.Params,
			Components:  []DataGetOrgByStyleMarkComponents{},
		}
		if len(orgData.Components) > 0 {
			var ids []int64
			for key := 0; key < len(orgData.Components); key++ {
				ids = append(ids, orgData.Components[key].ComponentID)
			}
			var components []FieldsComponent
			components, err = GetComponentMore(&ArgsGetComponentMore{
				IDs:        ids,
				HaveRemove: false,
			})
			if err != nil {
				return
			}
			for key := 0; key < len(orgData.Components); key++ {
				for key2 := 0; key2 < len(components); key2++ {
					if orgData.Components[key].ComponentID == components[key2].ID {
						data.Components = append(data.Components, DataGetOrgByStyleMarkComponents{
							ComponentID: components[key2].ID,
							Mark:        components[key2].Mark,
							Name:        components[key2].Name,
							Des:         components[key2].Des,
							CoverFileID: components[key2].CoverFileID,
							DesFiles:    components[key2].DesFiles,
							Params:      orgData.Components[key].Params,
						})
						break
					}
				}
			}
		}
		return
	} else {
		var components []FieldsComponent
		components, err = GetComponentMore(&ArgsGetComponentMore{
			IDs:        styleData.Components,
			HaveRemove: false,
		})
		if err != nil {
			return
		}
		var orgComponents []DataGetOrgByStyleMarkComponents
		for key := 0; key < len(styleData.Components); key++ {
			for key2 := 0; key2 < len(components); key2++ {
				if styleData.Components[key] == components[key2].ID {
					orgComponents = append(orgComponents, DataGetOrgByStyleMarkComponents{
						ComponentID: components[key2].ID,
						Mark:        components[key2].Mark,
						Name:        components[key2].Name,
						Des:         components[key2].Des,
						CoverFileID: components[key2].CoverFileID,
						DesFiles:    components[key2].DesFiles,
						Params:      components[key2].Params,
					})
					break
				}
			}
		}
		data = DataGetOrgByStyleMark{
			Name:        styleData.Name,
			Mark:        styleData.Mark,
			SystemMark:  styleData.SystemMark,
			Title:       styleData.Title,
			Des:         styleData.Des,
			CoverFileID: styleData.CoverFileID,
			DesFiles:    styleData.DesFiles,
			Params:      styleData.Params,
			Components:  orgComponents,
		}
	}
	return
}

type ArgsGetOrgByStyleID struct {
	//样式表ID
	StyleID int64 `db:"style_id" json:"styleID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

func GetOrgByStyleID(args *ArgsGetOrgByStyleID) (data DataGetOrgByStyleMark, err error) {
	var styleData FieldsStyle
	styleData, err = GetStyleByID(&ArgsGetStyleByID{
		ID: args.StyleID,
	})
	if err != nil {
		return
	}
	var orgData FieldsOrg
	err = Router2SystemConfig.MainDB.Get(&orgData, "SELECT id, create_at, update_at, delete_at, org_id, style_id, components, title, des, cover_file_id, des_files, params FROM core_style_org WHERE style_id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", styleData.ID, args.OrgID)
	if err == nil {
		data = DataGetOrgByStyleMark{
			Name:        styleData.Name,
			Mark:        styleData.Mark,
			SystemMark:  styleData.SystemMark,
			Title:       styleData.Title,
			Des:         styleData.Des,
			CoverFileID: styleData.CoverFileID,
			DesFiles:    styleData.DesFiles,
			Params:      styleData.Params,
			Components:  []DataGetOrgByStyleMarkComponents{},
		}
		if len(orgData.Components) > 0 {
			var ids []int64
			for key := 0; key < len(orgData.Components); key++ {
				ids = append(ids, orgData.Components[key].ComponentID)
			}
			var components []FieldsComponent
			components, err = GetComponentMore(&ArgsGetComponentMore{
				IDs:        ids,
				HaveRemove: false,
			})
			if err != nil {
				return
			}
			for key := 0; key < len(orgData.Components); key++ {
				for key2 := 0; key2 < len(components); key2++ {
					if orgData.Components[key].ComponentID == components[key2].ID {
						data.Components = append(data.Components, DataGetOrgByStyleMarkComponents{
							ComponentID: components[key2].ID,
							Mark:        components[key2].Mark,
							Name:        components[key2].Name,
							Des:         components[key2].Des,
							CoverFileID: components[key2].CoverFileID,
							DesFiles:    components[key2].DesFiles,
							Params:      orgData.Components[key].Params,
						})
						break
					}
				}
			}
		}
		return
	} else {
		var components []FieldsComponent
		components, err = GetComponentMore(&ArgsGetComponentMore{
			IDs:        styleData.Components,
			HaveRemove: false,
		})
		if err != nil {
			return
		}
		var orgComponents []DataGetOrgByStyleMarkComponents
		for key := 0; key < len(styleData.Components); key++ {
			for key2 := 0; key2 < len(components); key2++ {
				if styleData.Components[key] == components[key2].ID {
					orgComponents = append(orgComponents, DataGetOrgByStyleMarkComponents{
						ComponentID: components[key2].ID,
						Mark:        components[key2].Mark,
						Name:        components[key2].Name,
						Des:         components[key2].Des,
						CoverFileID: components[key2].CoverFileID,
						DesFiles:    components[key2].DesFiles,
						Params:      components[key2].Params,
					})
					break
				}
			}
		}
		data = DataGetOrgByStyleMark{
			Name:        styleData.Name,
			Mark:        styleData.Mark,
			SystemMark:  styleData.SystemMark,
			Title:       styleData.Title,
			Des:         styleData.Des,
			CoverFileID: styleData.CoverFileID,
			DesFiles:    styleData.DesFiles,
			Params:      styleData.Params,
			Components:  orgComponents,
		}
	}
	return
}

// ArgsSetOrgStyle 创建新的样式关联参数
type ArgsSetOrgStyle struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//样式表ID
	StyleID int64 `db:"style_id" json:"styleID" check:"id"`
	//组件列
	Components FieldsOrgComponents `db:"components" json:"components"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetOrgStyle 创建新的样式关联
func SetOrgStyle(args *ArgsSetOrgStyle) (data FieldsOrg, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, style_id, components, title, des, cover_file_id, des_files, params FROM core_style_org WHERE delete_at < to_timestamp(1000000) AND style_id = $1 AND org_id = $2", args.StyleID, args.OrgID)
	if err != nil {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_style_org", "INSERT INTO core_style_org (org_id, style_id, components, title, des, cover_file_id, des_files, params) VALUES (:org_id, :style_id, :components, :title, :des, :cover_file_id, :des_files, :params)", args, &data)
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_style_org SET update_at = NOW(), components = :components, title = :title, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, params = :params WHERE id = :id AND org_id = :org_id", map[string]interface{}{
			"id":            data.ID,
			"org_id":        args.OrgID,
			"components":    args.Components,
			"title":         args.Title,
			"des":           args.Des,
			"cover_file_id": args.CoverFileID,
			"des_files":     args.DesFiles,
			"params":        args.Params,
		})
	}
	return
}

// ArgsDeleteOrgStyle 删除样式关联参数
type ArgsDeleteOrgStyle struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteOrgStyle 删除样式关联
func DeleteOrgStyle(args *ArgsDeleteOrgStyle) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "core_style_org", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
