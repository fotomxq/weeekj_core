package BlogExam

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetTopicGroupList 获取题库列表参数
type ArgsGetTopicGroupList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTopicGroupList 获取题库列表
func GetTopicGroupList(args *ArgsGetTopicGroupList) (dataList []FieldsTopicGroup, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "blog_exam_topic_group"
	var rawList []FieldsTopicGroup
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
		vData := GetTopicGroup(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// GetTopicGroupHaveAddVisit 获取题库详情并增加访问次数
func GetTopicGroupHaveAddVisit(id int64) (data FieldsTopicGroup) {
	data = GetTopicGroup(id)
	if CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		data.ID = 0
		return
	}
	_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE blog_exam_topic_group SET visit_count = visit_count + 1 WHERE id = :id", map[string]interface{}{
		"id": data.ID,
	})
	cacheMark := getTopicGroupCacheMark(id)
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// GetTopicGroup 获取题库详情
func GetTopicGroup(id int64) (data FieldsTopicGroup) {
	cacheMark := getTopicGroupCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, title, topic_ids, params, visit_count FROM blog_exam_topic_group WHERE id = $1", id)
	if err != nil || data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// ArgsCreateTopicGroup 创建题库参数
type ArgsCreateTopicGroup struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//题库名称
	Title string `db:"title" json:"title" check:"des" min:"1" max:"600"`
	//题目列
	TopicIDs pq.Int64Array `db:"topic_ids" json:"topicIDs" check:"ids"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateTopicGroup 创建题库
func CreateTopicGroup(args *ArgsCreateTopicGroup) (err error) {
	if args.TopicIDs == nil || len(args.TopicIDs) < 1 {
		args.TopicIDs = pq.Int64Array{}
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_exam_topic_group (org_id, title, topic_ids, params) VALUES (:org_id,:title,:topic_ids,:params)", args)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateTopicGroup 修改题库参数
type ArgsUpdateTopicGroup struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//题库名称
	Title string `db:"title" json:"title" check:"des" min:"1" max:"600"`
	//题目列
	TopicIDs pq.Int64Array `db:"topic_ids" json:"topicIDs" check:"ids"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateTopicGroup 修改题库
func UpdateTopicGroup(args *ArgsUpdateTopicGroup) (err error) {
	if args.TopicIDs == nil || len(args.TopicIDs) < 1 {
		args.TopicIDs = pq.Int64Array{}
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE blog_exam_topic_group SET update_at = NOW(), title = :title, topic_ids = :topic_ids, params = :params WHERE id = :id AND (org_id = :org_id OR :org_id < 1)", args)
	if err != nil {
		return
	}
	deleteTopicGroupCache(args.ID)
	return
}

// ArgsDeleteTopicGroup 删除题库参数
type ArgsDeleteTopicGroup struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteTopicGroup 删除题库
func DeleteTopicGroup(args *ArgsDeleteTopicGroup) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "blog_exam_topic_group", "id = :id AND (org_id = :org_id OR :org_id < 1)", args)
	if err != nil {
		return
	}
	deleteTopicGroupCache(args.ID)
	return
}

// 缓冲
func getTopicGroupCacheMark(id int64) string {
	return fmt.Sprint("blog:exam:topic:group:id:", id)
}

func deleteTopicGroupCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getTopicGroupCacheMark(id))
}

// 检查题目是否属于该组织
func checkTopicInOrg(topicIDs []int64, orgID int64) bool {
	if orgID < 0 {
		return true
	}
	for _, v := range topicIDs {
		vData := GetTopic(v)
		if vData.ID < 1 {
			return false
		}
		if CoreSQL.CheckTimeHaveData(vData.DeleteAt) {
			return false
		}
		if vData.OrgID != orgID {
			return false
		}
	}
	return true
}
