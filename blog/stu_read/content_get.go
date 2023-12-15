package BlogStuRead

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetContentList 获取文章列表参数
type ArgsGetContentList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//文章类型
	// 0 普通文章; 1 挂靠视频; 3 第三方跳转
	// 挂靠或跳转内容，将在des中做特殊描述
	ContentType int `db:"content_type" json:"contentType"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetContentList 获取文章列表
func GetContentList(args *ArgsGetContentList) (dataList []FieldsContent, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ContentType > -1 {
		where = where + " AND content_type = :content_type"
		maps["content_type"] = args.ContentType
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR title_des ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "blog_stu_read_content"
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
		[]string{"id", "create_at", "update_at", "delete_at", "visit_count"},
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

// GetContentByIDHaveAddVisit 获取文章并增加阅读次数
func GetContentByIDHaveAddVisit(id int64) (data FieldsContent) {
	//获取文章
	data, _ = GetContentByID(&ArgsGetContentByID{
		ID:    id,
		OrgID: -1,
	})
	if data.ID < 1 {
		return
	}
	//更新访问次数
	_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE blog_stu_read_content SET visit_count = visit_count + 1 WHERE id = :id", map[string]interface{}{
		"id": data.ID,
	})
	//删除缓冲
	deleteContentCacheByID(id)
	//保存缓冲
	data.VisitCount += 1
	//反馈
	return
}

// ArgsGetContentByID 获取ID参数
type ArgsGetContentByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetContentByID 获取ID
func GetContentByID(args *ArgsGetContentByID) (data FieldsContent, err error) {
	data = getContentID(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data.ID = 0
		err = errors.New("no data")
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
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, content_type, org_id, title, title_des, cover_file_id, des_files, des, params, visit_count FROM blog_stu_read_content WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 604800)
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
