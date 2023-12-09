package BlogExam

import (
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"strings"
	"time"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//参加考试
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.ProductID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "product_id = :product_id"
		maps["product_id"] = args.ProductID
	}
	if where == "" {
		where = "true"
	}
	tableName := "blog_exam_log"
	var rawList []FieldsLog
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "end_at", "score"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := GetLogByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// GetLogCountByOrgID 获取组织考试次数
func GetLogCountByOrgID(orgID int64) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM blog_exam_log WHERE ($1 < -1 OR org_id = $1)", orgID)
	if err != nil {
		return
	}
	return
}

// GetLogOrgCount 获取多少个组织参与的考试次数
func GetLogOrgCount() (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(org_id) FROM blog_exam_log GROUP BY org_id")
	if err != nil {
		return
	}
	return
}

// GetLogCountByUserID 获取用户考试次数
func GetLogCountByUserID(userID int64) (count int64) {
	cacheMark := getLogCountByUserIDCacheMark(userID)
	var err error
	count, err = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if err == nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM blog_exam_log WHERE user_id = $1", userID)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetInt64(cacheMark, count, 1800)
	return
}

// GetLogByID 获取指定日志
func GetLogByID(id int64) (data FieldsLog) {
	cacheMark := getLogCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, end_at, run_time, org_id, user_id, product_id, topic_ids, score, err_count, correct_count FROM blog_exam_log WHERE id = $1", id)
	if err != nil || data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// ArgsAppendLogAnswer 添加新的记录答题参数
type ArgsAppendLogAnswer struct {
	//题目ID
	TopicID int64 `json:"topicID" check:"id"`
	//答案
	Answer []string `json:"answer"`
}

// ArgsAppendLog 添加新的记录参数
type ArgsAppendLog struct {
	//创建时间
	CreateAt string `db:"create_at" json:"createAt" check:"defaultTime"`
	//结束时间
	EndAt string `db:"end_at" json:"endAt" check:"defaultTime"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//参加考试
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//参加的题目列
	TopicIDs pq.Int64Array `db:"topic_ids" json:"topicIDs" check:"ids" empty:"true"`
	//答案序列
	Answer []ArgsAppendLogAnswer `json:"answer"`
}

// AppendLog 添加新的记录
func AppendLog(args *ArgsAppendLog) (logID int64, allScore int, errCode string, err error) {
	//获取试卷
	if args.ProductID > 0 {
		productData := GetProduct(args.ProductID)
		if productData.ID < 1 || CoreSQL.CheckTimeHaveData(productData.DeleteAt) {
			errCode = "err_no_data"
			err = errors.New("product not exist")
			return
		}
		if len(args.TopicIDs) < 1 {
			for _, v := range productData.TopicIDs {
				args.TopicIDs = append(args.TopicIDs, v)
			}
		}
	}
	//检查答题数量
	if len(args.Answer) != len(args.TopicIDs) {
		errCode = "err_num"
		err = errors.New("answer not eq product topic len")
		return
	}
	if len(args.TopicIDs) < 1 {
		errCode = "err_num"
		err = errors.New("no topic ids")
		return
	}
	//获取参数
	var startAt, endAt time.Time
	startAt, err = CoreFilter.GetTimeByDefault(args.CreateAt)
	if err != nil {
		errCode = "err_time"
		return
	}
	endAt, err = CoreFilter.GetTimeByDefault(args.EndAt)
	if err != nil {
		errCode = "err_time"
		return
	}
	runTime := endAt.Unix() - startAt.Unix()
	//遍历数据，计算结果
	var score, errCount, correctCount int
	for _, v := range args.TopicIDs {
		vTopic := GetTopic(v)
		if vTopic.ID < 1 {
			continue
		}
		if CoreSQL.CheckTimeHaveData(vTopic.DeleteAt) {
			continue
		}
		allScore = vTopic.Score
		isOK := false
		for _, v2 := range args.Answer {
			if v2.TopicID == vTopic.ID {
				switch vTopic.TopicType {
				case 0:
					if len(v2.Answer) == len(vTopic.Answer) && len(v2.Answer) > 0 && v2.Answer[0] == vTopic.Answer[0] {
						isOK = true
					}
				case 1:
					isOK = true
					for _, v3 := range vTopic.Answer {
						isFind := false
						for _, v4 := range v2.Answer {
							if v3 == v4 {
								isFind = true
								break
							}
						}
						if !isFind {
							isOK = false
						}
					}
				case 2:
					if len(v2.Answer) == len(vTopic.Answer) && len(v2.Answer) > 0 && v2.Answer[0] == vTopic.Answer[0] {
						isOK = true
					}
				case 3:
					isOK = true
					for _, v3 := range v2.Answer {
						if !strings.Contains(vTopic.Des, fmt.Sprint("{", v3, "}")) {
							isOK = false
							break
						}
					}
				case 4:
					isOK = true
				default:
					break
				}
				break
			}
		}
		//如果正确，增加分数
		if isOK {
			correctCount += 1
			score = vTopic.Score
		} else {
			errCount += 1
		}
	}
	//重新计算分数
	var blogExamScore100 bool
	blogExamScore100, err = BaseConfig.GetDataBool("BlogExamScore100")
	if err != nil {
		blogExamScore100 = false
	}
	if blogExamScore100 {
		if allScore > 0 {
			score = int((float64(score) / float64(allScore)) * 10000)
		}
	}
	//记录数据
	logID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_exam_log (create_at, end_at, run_time, org_id, user_id, product_id, topic_ids, score, err_count, correct_count) VALUES (:create_at, :end_at, :run_time, :org_id, :user_id, :product_id, :topic_ids, :score, :err_count, :correct_count)", map[string]interface{}{
		"create_at":     startAt,
		"end_at":        endAt,
		"run_time":      runTime,
		"org_id":        args.OrgID,
		"user_id":       args.UserID,
		"product_id":    args.ProductID,
		"topic_ids":     args.TopicIDs,
		"score":         score,
		"err_count":     errCount,
		"correct_count": correctCount,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	logData := GetLogByID(logID)
	//删除缓冲
	deleteLogCache(logID)
	//推送nats
	CoreNats.PushDataNoErr("/blog/exam/log", "new", logID, "", logData)
	//反馈成功
	return
}

// 缓冲
func getLogCacheMark(id int64) string {
	return fmt.Sprint("blog:exam:log:id:", id)
}

func getLogCountByUserIDCacheMark(userID int64) string {
	return fmt.Sprint("blog:exam:log:user:", userID)
}

func deleteLogCache(id int64) {
	data := GetLogByID(id)
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(id))
	if data.UserID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getLogCountByUserIDCacheMark(data.UserID))
	}
}
