package BlogCore

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"strings"
	"time"
)

// ArgsGetContentList 获取列表参数
type ArgsGetContentList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//文章类型
	// 0 普通文章; 1 挂靠视频; 3 第三方跳转
	// 挂靠或跳转内容，将在des中做特殊描述
	ContentType int `db:"content_type" json:"contentType" check:"intThan0" empty:"true"`
	//分类ID
	// > -1 为包含；否则不包含。0为没有设定
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//标签或关系
	TagsOr bool `json:"tagsOr" check:"bool"`
	//是否已经发布
	NeedIsPublish bool `db:"need_is_publish" json:"needIsPublish" check:"bool"`
	IsPublish     bool `db:"is_publish" json:"isPublish" check:"bool"`
	//归属关系
	// 删除后作为原始文档的子项目存在，key将自动失效
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//是否需要审核
	NeedAudit bool `json:"needAudit" check:"bool"`
	//是否已经审核
	IsAudit bool `json:"isAudit" check:"bool"`
	//是否置顶
	NeedIsTop bool `json:"needIsTop" check:"bool" empty:"true"`
	IsTop     bool `db:"is_top" json:"isTop" check:"bool" empty:"true"`
	//扩展选项范围查询
	Param1Min int64 `json:"param1Min"`
	Param1Max int64 `json:"param1Max"`
	Param2Min int64 `json:"param2Min"`
	Param2Max int64 `json:"param2Max"`
	Param3Min int64 `json:"param3Min"`
	Param3Max int64 `json:"param3Max"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetContentList 获取列表
func GetContentList(args *ArgsGetContentList) (dataList []FieldsContent, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.BindID > 0 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.ContentType > -1 {
		where = where + " AND content_type = :content_type"
		maps["content_type"] = args.ContentType
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.TagsOr {
		if len(args.Tags) > 0 {
			var tagStrs []string
			for _, v := range args.Tags {
				tagStrs = append(tagStrs, fmt.Sprint("tags && :tags_", v))
				maps[fmt.Sprint("tags_", v)] = pq.Int64Array{v}
			}
			whereOr := strings.Join(tagStrs, " OR ")
			where = where + " AND (" + whereOr + ")"
		}
	} else {
		if len(args.Tags) > 0 {
			where = where + " AND tags @> :tags"
			maps["tags"] = args.Tags
		}
	}
	if args.NeedAudit {
		if args.IsAudit {
			where = where + " AND audit_at > to_timestamp(1000000)"
		} else {
			where = where + " AND audit_at <= to_timestamp(1000000)"
		}
	}
	if args.NeedIsPublish {
		if args.IsPublish {
			where = where + " AND publish_at > to_timestamp(1000000)"
		} else {
			where = where + " AND publish_at <= to_timestamp(1000000)"
		}
	}
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.NeedIsTop {
		if args.IsTop {
			where = where + " AND is_top = true"
		} else {
			where = where + " AND is_top = false"
		}
	}
	if args.Param1Min > 0 && args.Param1Max > 0 {
		where = where + " AND param1 >= :param1_min AND param1 <= :param1_max"
		maps["param1_min"] = args.Param1Min
		maps["param1_max"] = args.Param1Max
	}
	if args.Param2Min > 0 && args.Param2Max > 0 {
		where = where + " AND param2 >= :param2_min AND param2 <= :param2_max"
		maps["param2_min"] = args.Param2Min
		maps["param2_max"] = args.Param2Max
	}
	if args.Param3Min > 0 && args.Param3Max > 0 {
		where = where + " AND param3 >= :param3_min AND param3 <= :param3_max"
		maps["param3_min"] = args.Param3Min
		maps["param3_max"] = args.Param3Max
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR title_des ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "blog_core_content"
	var rawList []FieldsContent
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "audit_at", "publish_at", "visit_count"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getContentID(v.ID)
		if vData.ID < 1 {
			continue
		}
		vData.Des = ""
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

//V3修改获取列表，支持一组用户发布的文章检索

// ArgsGetContentListV3 获取列表参数
type ArgsGetContentListV3 struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserIDs pq.Int64Array `db:"user_ids" json:"userIDs" check:"ids" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//分类ID
	// > -1 为包含；否则不包含。0为没有设定
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//标签或关系
	TagsOr bool `json:"tagsOr" check:"bool"`
	//是否已经发布
	NeedIsPublish bool `db:"need_is_publish" json:"needIsPublish" check:"bool"`
	IsPublish     bool `db:"is_publish" json:"IsPublish" check:"bool"`
	//归属关系
	// 删除后作为原始文档的子项目存在，key将自动失效
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//是否需要审核
	NeedAudit bool `json:"needAudit" check:"bool"`
	//是否已经审核
	IsAudit bool `json:"isAudit" check:"bool"`
	//是否置顶
	NeedIsTop bool `json:"needIsTop" check:"bool" empty:"true"`
	IsTop     bool `db:"is_top" json:"isTop" check:"bool" empty:"true"`
	//扩展选项范围查询
	Param1Min int64 `json:"param1Min"`
	Param1Max int64 `json:"param1Max"`
	Param2Min int64 `json:"param2Min"`
	Param2Max int64 `json:"param2Max"`
	Param3Min int64 `json:"param3Min"`
	Param3Max int64 `json:"param3Max"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetContentListV3 获取列表
func GetContentListV3(args *ArgsGetContentListV3) (dataList []FieldsContent, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if len(args.UserIDs) > -1 {
		where = where + " AND user_id = ANY(:user_id)"
		maps["user_id"] = args.UserIDs
	}
	if args.BindID > 0 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.TagsOr {
		if len(args.Tags) > 0 {
			var tagStrs []string
			for _, v := range args.Tags {
				tagStrs = append(tagStrs, fmt.Sprint("tags && :tags_", v))
				maps[fmt.Sprint("tags_", v)] = pq.Int64Array{v}
			}
			whereOr := strings.Join(tagStrs, " OR ")
			where = where + " AND (" + whereOr + ")"
		}
	} else {
		if len(args.Tags) > 0 {
			where = where + " AND tags @> :tags"
			maps["tags"] = args.Tags
		}
	}
	if args.NeedAudit {
		if args.IsAudit {
			where = where + " AND audit_at > to_timestamp(1000000)"
		} else {
			where = where + " AND audit_at <= to_timestamp(1000000)"
		}
	}
	if args.NeedIsPublish {
		if args.IsPublish {
			where = where + " AND publish_at > to_timestamp(1000000)"
		} else {
			where = where + " AND publish_at <= to_timestamp(1000000)"
		}
	}
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.NeedIsTop {
		if args.IsTop {
			where = where + " AND is_top = true"
		} else {
			where = where + " AND is_top = false"
		}
	}
	if args.Param1Min > 0 && args.Param1Max > 0 {
		where = where + " AND param1 >= :param1_min AND param1 <= :param1_max"
		maps["param1_min"] = args.Param1Min
		maps["param1_max"] = args.Param1Max
	}
	if args.Param2Min > 0 && args.Param2Max > 0 {
		where = where + " AND param2 >= :param2_min AND param2 <= :param2_max"
		maps["param2_min"] = args.Param2Min
		maps["param2_max"] = args.Param2Max
	}
	if args.Param3Min > 0 && args.Param3Max > 0 {
		where = where + " AND param3 >= :param3_min AND param3 <= :param3_max"
		maps["param3_min"] = args.Param3Min
		maps["param3_max"] = args.Param3Max
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR title_des ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "blog_core_content"
	var rawList []FieldsContent
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "audit_at", "publish_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getContentID(v.ID)
		vData.Des = ""
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// ArgsGetContentListV4 获取列表参数
type ArgsGetContentListV4 struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserIDs pq.Int64Array `db:"user_ids" json:"userIDs" check:"ids" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//分类ID
	// > -1 为包含；否则不包含。0为没有设定
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//标签或关系
	TagsOr bool `json:"tagsOr" check:"bool"`
	//是否已经发布
	NeedIsPublish bool `db:"need_is_publish" json:"needIsPublish" check:"bool"`
	IsPublish     bool `db:"is_publish" json:"IsPublish" check:"bool"`
	//归属关系
	// 删除后作为原始文档的子项目存在，key将自动失效
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//是否需要审核
	NeedAudit bool `json:"needAudit" check:"bool"`
	//是否已经审核
	IsAudit bool `json:"isAudit" check:"bool"`
	//是否置顶
	NeedIsTop bool `json:"needIsTop" check:"bool" empty:"true"`
	IsTop     bool `db:"is_top" json:"isTop" check:"bool" empty:"true"`
	//扩展选项范围查询
	Param1Min int64 `json:"param1Min"`
	Param1Max int64 `json:"param1Max"`
	Param2Min int64 `json:"param2Min"`
	Param2Max int64 `json:"param2Max"`
	Param3Min int64 `json:"param3Min"`
	Param3Max int64 `json:"param3Max"`
	//发布时间范围
	PublishBetween CoreSQL2.ArgsTimeBetween `json:"publishBetween"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetContentListV4 获取列表
func GetContentListV4(args *ArgsGetContentListV4) (dataList []FieldsContent, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if len(args.UserIDs) > 0 {
		where = where + " AND user_id = ANY(:user_id)"
		maps["user_id"] = args.UserIDs
	}
	if args.BindID > 0 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.TagsOr {
		if len(args.Tags) > 0 {
			var tagStrs []string
			for _, v := range args.Tags {
				tagStrs = append(tagStrs, fmt.Sprint("tags && :tags_", v))
				maps[fmt.Sprint("tags_", v)] = pq.Int64Array{v}
			}
			whereOr := strings.Join(tagStrs, " OR ")
			where = where + " AND (" + whereOr + ")"
		}
	} else {
		if len(args.Tags) > 0 {
			where = where + " AND tags @> :tags"
			maps["tags"] = args.Tags
		}
	}
	if args.NeedAudit {
		if args.IsAudit {
			where = where + " AND audit_at > to_timestamp(1000000)"
		} else {
			where = where + " AND audit_at <= to_timestamp(1000000)"
		}
	}
	if args.NeedIsPublish {
		if args.IsPublish {
			where = where + " AND publish_at > to_timestamp(1000000)"
		} else {
			where = where + " AND publish_at <= to_timestamp(1000000)"
		}
	}
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.NeedIsTop {
		if args.IsTop {
			where = where + " AND is_top = true"
		} else {
			where = where + " AND is_top = false"
		}
	}
	if args.Param1Min > 0 && args.Param1Max > 0 {
		where = where + " AND param1 >= :param1_min AND param1 <= :param1_max"
		maps["param1_min"] = args.Param1Min
		maps["param1_max"] = args.Param1Max
	}
	if args.Param2Min > 0 && args.Param2Max > 0 {
		where = where + " AND param2 >= :param2_min AND param2 <= :param2_max"
		maps["param2_min"] = args.Param2Min
		maps["param2_max"] = args.Param2Max
	}
	if args.Param3Min > 0 && args.Param3Max > 0 {
		where = where + " AND param3 >= :param3_min AND param3 <= :param3_max"
		maps["param3_min"] = args.Param3Min
		maps["param3_max"] = args.Param3Max
	}
	if args.PublishBetween.MinTime != "" {
		where = where + " AND publish_at >= :publish_at_min"
		maps["publish_at_min"] = CoreFilter.GetTimeByDefaultNoErr(args.PublishBetween.MinTime)
	}
	if args.PublishBetween.MaxTime != "" {
		where = where + " AND publish_at <= :publish_at_max"
		maps["publish_at_max"] = CoreFilter.GetTimeByDefaultNoErr(args.PublishBetween.MaxTime)
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR title_des ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "blog_core_content"
	var rawList []FieldsContent
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&CoreSQLPages.ArgsDataList{
			Page: args.Pages.Page,
			Max:  args.Pages.Max,
			Sort: args.Pages.Sort,
			Desc: args.Pages.Desc,
		},
		[]string{"id", "create_at", "update_at", "delete_at", "audit_at", "publish_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getContentID(v.ID)
		vData.Des = ""
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// ArgsGetContentSortCount 统计某个分类多少文章参数
type ArgsGetContentSortCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//分类ID
	// > -1 为包含；否则不包含。0为没有设定
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//指定更新日期之后的数据
	// 如果给空，则忽略
	AfterAt string `json:"afterAt" check:"isoTime" empty:"true"`
}

// GetContentSortCount 统计某个分类多少文章
func GetContentSortCount(args *ArgsGetContentSortCount) (count int64) {
	//获取数据
	var err error
	if args.AfterAt == "" {
		count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "blog_core_content", "id", "org_id = :org_id AND (sort_id = :sort_id OR :sort_id < 1) AND delete_at < to_timestamp(1000000) AND publish_at > to_timestamp(1000000) AND audit_at > to_timestamp(1000000) AND parent_id = 0", args)
	} else {
		var afterAt time.Time
		afterAt, err = CoreFilter.GetTimeByISO(args.AfterAt)
		if err != nil {
			return
		}
		count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "blog_core_content", "id", "org_id = :org_id AND (sort_id = :sort_id OR :sort_id < 1) AND update_at >= :after_at AND delete_at < to_timestamp(1000000) AND publish_at > to_timestamp(1000000) AND audit_at > to_timestamp(1000000) AND parent_id = 0", map[string]interface{}{
			"org_id":   args.OrgID,
			"sort_id":  args.SortID,
			"after_at": afterAt,
		})
	}
	if err != nil {
		count = 0
	}
	return
}

// ArgsGetContentTagsCount 获取标签有多少篇文章参数
type ArgsGetContentTagsCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//标签或关系
	TagsOr bool `json:"tagsOr"`
	//指定更新日期之后的数据
	// 如果给空，则忽略
	AfterAt string `json:"afterAt" check:"isoTime" empty:"true"`
}

// GetContentTagsCount 获取标签有多少篇文章
// 多个标签为或的关系
func GetContentTagsCount(args *ArgsGetContentTagsCount) (count int64) {
	var err error
	where := "org_id = :org_id AND delete_at < to_timestamp(1000000) AND publish_at > to_timestamp(1000000) AND audit_at > to_timestamp(1000000) AND parent_id = 0"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.TagsOr {
		if len(args.Tags) > 0 {
			var tagStrs []string
			for _, v := range args.Tags {
				tagStrs = append(tagStrs, fmt.Sprint("tags @> :tags_", v))
				maps[fmt.Sprint("tags_", v)] = pq.Int64Array{v}
			}
			whereOr := strings.Join(tagStrs, " OR ")
			where = where + " AND (" + whereOr + ")"
		}
	} else {
		if len(args.Tags) > 0 {
			where = where + " AND tags @> :tags"
			maps["tags"] = args.Tags
		}
	}
	if args.AfterAt == "" {
		count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "blog_core_content", "id", where, maps)
	} else {
		var afterAt time.Time
		afterAt, err = CoreFilter.GetTimeByISO(args.AfterAt)
		if err != nil {
			return
		}
		where = where + " AND update_at >= :after_at"
		maps["after_at"] = afterAt
		count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "blog_core_content", "id", where, maps)
	}
	if err != nil {
		count = 0
	}
	return
}

// ArgsGetContentByID 获取ID参数
type ArgsGetContentByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否需要审核
	IsAudit bool `db:"is_audit" json:"isAudit" check:"bool"`
	//是否需要检查发布状态
	IsPublish bool `db:"is_publish" json:"IsPublish" check:"bool"`
	//阅读用户ID
	// 可以留空
	ReadUserID int64 `json:"readUserID" check:"id" empty:"true"`
}

// GetContentByID 获取ID
func GetContentByID(args *ArgsGetContentByID) (data FieldsContent, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM blog_core_content WHERE (org_id = $1 OR $1 < 1) AND id = $2 AND delete_at < to_timestamp(1000000) AND ($3 = FALSE OR ($3 = TRUE AND publish_at > to_timestamp(1000000))) AND ($4 < 1 OR user_id = $4) AND ($5 = FALSE OR ($5 = TRUE AND audit_at > to_timestamp(1000000)))", args.OrgID, args.ID, args.IsPublish, args.UserID, args.IsAudit)
	if err != nil {
		return
	}
	data = getContentID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	PushRead(data.OrgID, args.ReadUserID, data.ID)
	return
}

// GetContentByIDNoErr 无错误获取文章信息
func GetContentByIDNoErr(id int64, orgID int64) (data FieldsContent) {
	data = getContentID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsContent{}
		return
	}
	return
}

// 获取文章
func getContentID(id int64) (data FieldsContent) {
	cacheMark := getContentCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, content_type, audit_at, audit_des, org_id, user_id, bind_id, param1, param2, param3, visit_count, key, parent_id, publish_at, is_top, sort_id, tags, title, title_des, cover_file_id, des_files, des, params FROM blog_core_content WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 604800)
	return
}

// ArgsGetContentCountByUser 获取用户发布文章数量参数
type ArgsGetContentCountByUser struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetContentCountByUser 获取用户发布文章数量
func GetContentCountByUser(args *ArgsGetContentCountByUser) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM blog_core_content WHERE user_id = $1 AND delete_at < to_timestamp(1000000) AND audit_at > to_timestamp(1000000) AND publish_at > to_timestamp(1000000)", args.UserID)
	return
}

// GetContentCountByOrgID 获取组织发布文章数量
func GetContentCountByOrgID(orgID int64, contentType int, afterAt time.Time) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM blog_core_content WHERE ($1 < 0 OR org_id = $1) AND delete_at < to_timestamp(1000000) AND audit_at > to_timestamp(1000000) AND publish_at > to_timestamp(1000000) AND ($2 < 0 OR content_type = $2) AND ($3 < to_timestamp(1000000) OR ($3 >= to_timestamp(1000000) AND create_at >= $3))", orgID, contentType, afterAt)
	return
}

// ArgsGetContentByKey 获取指定mey参数
type ArgsGetContentByKey struct {
	//唯一标识码key
	// 作为id的补充，自动填写时，将自动生成随机字符串
	// 默认根据标题或标题拼音得出
	Key string `db:"key" json:"key" check:"mark"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否需要审核
	IsAudit bool `db:"is_audit" json:"isAudit" check:"bool"`
	//是否需要检查发布状态
	IsPublish bool `db:"is_publish" json:"IsPublish" check:"bool"`
	//阅读用户ID
	// 可以留空
	ReadUserID int64 `json:"readUserID" check:"id" empty:"true"`
}

// GetContentByKey 获取指定key
func GetContentByKey(args *ArgsGetContentByKey) (data FieldsContent, err error) {
	//查询数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM blog_core_content WHERE org_id = $1 AND key = $2 AND delete_at < to_timestamp(1000000) AND ($3 = FALSE OR ($3 = TRUE AND publish_at > to_timestamp(1000000))) AND ($4 < 1 OR user_id = $4) AND ($5 = FALSE OR ($5 = TRUE AND audit_at > to_timestamp(1000000)))", args.OrgID, args.Key, args.IsPublish, args.UserID, args.IsAudit)
	if err != nil {
		return
	}
	data = getContentID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	PushRead(data.OrgID, args.ReadUserID, data.ID)
	return
}

type ArgsCheckContentByKey struct {
	//唯一标识码key
	// 作为id的补充，自动填写时，将自动生成随机字符串
	// 默认根据标题或标题拼音得出
	Key string `db:"key" json:"key" check:"mark"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

func CheckContentByKey(args *ArgsCheckContentByKey) (data FieldsContent, err error) {
	//查询数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM blog_core_content WHERE org_id = $1 AND key = $2", args.OrgID, args.Key)
	return
}

// ArgsGetContentMore 获取一组IDs参数
type ArgsGetContentMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetContentMore 获取一组IDs
func GetContentMore(args *ArgsGetContentMore) (dataList []FieldsContent, err error) {
	var rawList []FieldsContent
	err = CoreSQLIDs.GetIDsOrgAndDelete(&rawList, "blog_core_content", "id", args.IDs, args.OrgID, args.HaveRemove)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getContentID(v.ID)
		vData.Des = ""
		dataList = append(dataList, vData)
	}
	return
}

func GetContentMoreMap(args *ArgsGetContentMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsOrgTitleAndDelete("blog_core_content", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// GetContentTitle 获取文章名称
func GetContentTitle(id int64) string {
	data := getContentID(id)
	if data.ID < 1 {
		return ""
	}
	return data.Title
}
