package BlogExam

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetProductList 获取试卷列表参数
type ArgsGetProductList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetProductList 获取试卷列表
func GetProductList(args *ArgsGetProductList) (dataList []FieldsProduct, dataCount int64, err error) {
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
	tableName := "blog_exam_product"
	var rawList []FieldsProduct
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "start_at", "end_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := GetProduct(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// GetProduct 获取试卷详情
func GetProduct(id int64) (data FieldsProduct) {
	cacheMark := getProductCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, start_at, end_at, org_id, title, topic_ids, return_answer_now, params FROM blog_exam_product WHERE id = $1", id)
	if err != nil || data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// GetProductName 获取试卷名称
func GetProductName(id int64) string {
	data := GetProduct(id)
	return data.Title
}

// GetProductCountByOrgID 获取试卷数量
func GetProductCountByOrgID(orgID int64, afterAt time.Time) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM blog_exam_product WHERE ($1 < 0 OR org_id = $1) AND delete_at < to_timestamp(1000000) AND ($2 < to_timestamp(1000000) OR ($2 >= to_timestamp(1000000) AND create_at >= $2))", orgID, afterAt)
	return
}

// ArgsCreateProduct 创建试卷参数
type ArgsCreateProduct struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//开始时间
	StartAt string `db:"start_at" json:"startAt" check:"defaultTime"`
	//结束时间
	EndAt string `db:"end_at" json:"endAt" check:"defaultTime"`
	//标题
	Title string `db:"title" json:"title" check:"des" min:"1" max:"600"`
	//题目列
	TopicIDs pq.Int64Array `db:"topic_ids" json:"topicIDs" check:"ids"`
	//是否直接反馈正确答案
	ReturnAnswerNow bool `db:"return_answer_now" json:"returnAnswerNow" check:"bool"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateProduct 创建试卷
func CreateProduct(args *ArgsCreateProduct) (err error) {
	if !checkTopicInOrg(args.TopicIDs, args.OrgID) {
		err = errors.New("topic not in org")
		return
	}
	var startAt, endAt time.Time
	startAt, err = CoreFilter.GetTimeByDefault(args.StartAt)
	if err != nil {
		return
	}
	endAt, err = CoreFilter.GetTimeByDefault(args.EndAt)
	if err != nil {
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_exam_product (start_at, end_at, org_id, title, topic_ids, return_answer_now, params) VALUES (:start_at, :end_at, :org_id, :title, :topic_ids, :return_answer_now, :params)", map[string]interface{}{
		"start_at":          startAt,
		"end_at":            endAt,
		"org_id":            args.OrgID,
		"title":             args.Title,
		"topic_ids":         args.TopicIDs,
		"return_answer_now": args.ReturnAnswerNow,
		"params":            args.Params,
	})
	if err != nil {
		return
	}
	return
}

// ArgsDeleteProduct 删除试卷参数
type ArgsDeleteProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteProduct 删除试卷
func DeleteProduct(args *ArgsDeleteProduct) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "blog_exam_product", "id = :id AND (org_id = :org_id OR :org_id < 1)", args)
	if err != nil {
		return
	}
	deleteProductCache(args.ID)
	return
}

// 缓冲
func getProductCacheMark(id int64) string {
	return fmt.Sprint("blog:exam:product:id:", id)
}

func deleteProductCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getProductCacheMark(id))
}
