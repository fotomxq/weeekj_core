package BlogExam

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	BaseFileUpload "github.com/fotomxq/weeekj_core/v5/base/fileupload"
	CoreExcel "github.com/fotomxq/weeekj_core/v5/core/excel"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ToolsLoadExcel "github.com/fotomxq/weeekj_core/v5/tools/load_excel"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"strings"
	"time"
)

// ArgsGetTopicList 获取题目列表参数
type ArgsGetTopicList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type DataGetTopicList struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//描述
	DesCut string `json:"desCut"`
	//题目类型
	// 0 单选； 1 多选； 2 判断； 3 填空题； 4 问答题
	TopicType int `db:"topic_type" json:"topicType"`
	//正确得分
	Score int `db:"score" json:"score"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// GetTopicList 获取题目列表
func GetTopicList(args *ArgsGetTopicList) (dataList []DataGetTopicList, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%' OR answer_analysis ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "blog_exam_topic"
	var rawList []FieldsTopic
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := GetTopic(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, DataGetTopicList{
			ID:        vData.ID,
			CreateAt:  vData.CreateAt,
			UpdateAt:  vData.UpdateAt,
			DeleteAt:  vData.DeleteAt,
			OrgID:     vData.OrgID,
			DesCut:    CoreFilter.SubStrQuick(vData.Des, 30),
			TopicType: vData.TopicType,
			Score:     vData.Score,
			Params:    vData.Params,
		})
	}
	//反馈
	return
}

// GetTopic 获取题目信息
func GetTopic(id int64) (data FieldsTopic) {
	cacheMark := getTopicCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, des, topic_type, score, options, answer, answer_analysis, params FROM blog_exam_topic WHERE id = $1", id)
	if err != nil || data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// GetRandTopic 随机抽取指定数量的题目
func GetRandTopic(orgID int64, limit int) (dataList []FieldsTopic) {
	if limit < 1 {
		limit = 1
	}
	if limit > 1000 {
		limit = 1000
	}
	var rawList []FieldsTopic
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM blog_exam_topic WHERE org_id = $1 AND delete_at < to_timestamp(1000000) ORDER BY random() LIMIT $2", orgID, limit)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := GetTopic(v.ID)
		if vData.ID < 1 || CoreSQL.CheckTimeHaveData(vData.DeleteAt) {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetTopicByIDs 获取一组题目
func GetTopicByIDs(ids []int64, haveRemove bool) (dataList []FieldsTopic) {
	if len(ids) < 1 {
		return
	}
	if len(ids) > 1000 {
		return
	}
	for _, v := range ids {
		isFind := false
		for _, v2 := range dataList {
			if v2.ID == v {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		vData := GetTopic(v)
		if vData.ID < 1 {
			continue
		}
		if !haveRemove && CoreSQL.CheckTimeHaveData(vData.DeleteAt) {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsCreateTopic 创建新的题目参数
type ArgsCreateTopic struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000"`
	//题目类型
	// 0 单选； 1 多选； 2 判断； 3 填空题； 4 问答题
	TopicType int `db:"topic_type" json:"topicType"`
	//正确得分
	Score int `db:"score" json:"score" check:"intThan0"`
	//选项
	// 单选、多选、判断
	Options FieldsTopicOptions `db:"options" json:"options"`
	//正确选项
	// 可能是多个，用于支持单选、多选、判断、填空
	// 注意填空题此处如果设置，则需des配合填写{marK}的字符，方便直到具体是哪个mark
	Answer pq.StringArray `db:"answer" json:"answer" check:"marks" empty:"true"`
	//解析
	AnswerAnalysis string `db:"answer_analysis" json:"answerAnalysis" check:"des" min:"1" max:"10000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateTopic 创建新的题目
func CreateTopic(args *ArgsCreateTopic) (data FieldsTopic, err error) {
	if !checkTopic(args.TopicType, args.Options, args.Answer, args.Des) {
		err = errors.New("check topic")
		return
	}
	var newID int64
	newID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_exam_topic (org_id, des, topic_type, score, options, answer, answer_analysis, params) VALUES (:org_id,:des,:topic_type,:score,:options,:answer,:answer_analysis,:params)", args)
	if err != nil {
		return
	}
	data = GetTopic(newID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateTopicMore 批量创建考题参数
type ArgsCreateTopicMore struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//题目列表
	DataList []ArgsCreateTopicMoreChild `json:"dataList"`
}

type ArgsCreateTopicMoreChild struct {
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000"`
	//题目类型
	// 0 单选； 1 多选； 2 判断； 3 填空题； 4 问答题
	TopicType int `db:"topic_type" json:"topicType"`
	//正确得分
	Score int `db:"score" json:"score" check:"intThan0"`
	//选项
	// 单选、多选、判断
	Options FieldsTopicOptions `db:"options" json:"options"`
	//正确选项
	// 可能是多个，用于支持单选、多选、判断、填空
	// 注意填空题此处如果设置，则需des配合填写{marK}的字符，方便直到具体是哪个mark
	Answer pq.StringArray `db:"answer" json:"answer" check:"marks" empty:"true"`
	//解析
	AnswerAnalysis string `db:"answer_analysis" json:"answerAnalysis" check:"des" min:"1" max:"10000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateTopicMore 批量创建考题
func CreateTopicMore(args *ArgsCreateTopicMore) (dataList []FieldsTopic, errCode string, err error) {
	for _, v := range args.DataList {
		var vData FieldsTopic
		vData, err = CreateTopic(&ArgsCreateTopic{
			OrgID:          args.OrgID,
			Des:            v.Des,
			TopicType:      v.TopicType,
			Score:          v.Score,
			Options:        v.Options,
			Answer:         v.Answer,
			AnswerAnalysis: v.AnswerAnalysis,
			Params:         v.Params,
		})
		if err != nil {
			return
		}
		dataList = append(dataList, vData)
	}
	return
}

// UploadExcelAndGetTopicArgs 上传和反馈excel数据结构，用于批量添加考题
func UploadExcelAndGetTopicArgs(c *gin.Context, args *BaseFileUpload.ArgsUploadToTemp) (dataList []ArgsCreateTopicMoreChild, errCode string, err error) {
	//获取excel文件
	var excelData *excelize.File
	var waitDeleteFile string
	excelData, waitDeleteFile, errCode, err = ToolsLoadExcel.UploadFileAndGetExcelData(c, args)
	if err != nil {
		return
	}
	//遍历数据并构建参数
	sheetMaps := excelData.GetSheetMap()
	var sheetName string
	for _, v := range sheetMaps {
		sheetName = v
		break
	}
	excelVals := CoreExcel.GetSheetRows(excelData, sheetName)
	if len(excelVals) < 1 {
		errCode = "err_excel"
		err = errors.New("excel vals is empty")
		return
	}
	step := 0
	for {
		step += 1
		if step >= len(excelVals) {
			break
		}
		rows := excelVals[step]
		if len(rows) < 12 {
			break
		}
		if rows[1] == "" {
			break
		}
		topicType := 0
		options := FieldsTopicOptions{}
		var answer []string
		switch rows[1] {
		case "单选":
			topicType = 0
			answer = []string{rows[9]}
			options = uploadExcelAndGetTopicArgsGetRows(rows)
		case "多选":
			topicType = 1
			answer = strings.Split(rows[9], ",")
			if len(answer) < 2 {
				answer = strings.Split(rows[9], "，")
				if len(answer) < 2 {
					answer = strings.Split(rows[9], "|")
				}
			}
			options = uploadExcelAndGetTopicArgsGetRows(rows)
		case "判断":
			topicType = 2
			answer = []string{rows[9]}
			if rows[3] != "" {
				options = append(options, FieldsTopicOption{
					Mark: "正确",
					Des:  rows[3],
				})
			}
			if rows[4] != "" {
				options = append(options, FieldsTopicOption{
					Mark: "错误",
					Des:  rows[4],
				})
			}
		case "填空":
			topicType = 3
			answer = []string{rows[9]}
		case "问答":
			topicType = 4
			answer = []string{rows[9]}
		default:
			continue
		}
		score, _ := CoreFilter.GetIntByString(rows[10])
		vData := ArgsCreateTopicMoreChild{
			Des:            rows[1],
			TopicType:      topicType,
			Score:          score,
			Options:        options,
			Answer:         answer,
			AnswerAnalysis: rows[11],
			Params:         CoreSQLConfig.FieldsConfigsType{},
		}
		dataList = append(dataList, vData)
	}
	//删除临时文件
	err = CoreFile.DeleteF(waitDeleteFile)
	if err != nil {
		CoreLog.Error("blog exam topic upload excel and get topic args, delete temp file: ", waitDeleteFile, ", err: ", err)
		err = nil
	}
	//反馈
	return
}

func uploadExcelAndGetTopicArgsGetRows(rows []string) (options []FieldsTopicOption) {
	if rows[3] != "" {
		options = append(options, FieldsTopicOption{
			Mark: "A",
			Des:  rows[3],
		})
	}
	if rows[4] != "" {
		options = append(options, FieldsTopicOption{
			Mark: "B",
			Des:  rows[4],
		})
	}
	if rows[5] != "" {
		options = append(options, FieldsTopicOption{
			Mark: "C",
			Des:  rows[5],
		})
	}
	if rows[6] != "" {
		options = append(options, FieldsTopicOption{
			Mark: "D",
			Des:  rows[6],
		})
	}
	if rows[7] != "" {
		options = append(options, FieldsTopicOption{
			Mark: "E",
			Des:  rows[7],
		})
	}
	if rows[8] != "" {
		options = append(options, FieldsTopicOption{
			Mark: "F",
			Des:  rows[8],
		})
	}
	return
}

// ArgsUpdateTopic 修改题目参数
type ArgsUpdateTopic struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000"`
	//正确得分
	Score int `db:"score" json:"score" check:"intThan0"`
	//选项
	// 单选、多选、判断
	Options FieldsTopicOptions `db:"options" json:"options"`
	//正确选项
	// 可能是多个，用于支持单选、多选、判断、填空
	// 注意填空题此处如果设置，则需des配合填写{marK}的字符，方便直到具体是哪个mark
	Answer pq.StringArray `db:"answer" json:"answer" check:"marks" empty:"true"`
	//解析
	AnswerAnalysis string `db:"answer_analysis" json:"answerAnalysis" check:"des" min:"1" max:"10000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateTopic 修改题目
func UpdateTopic(args *ArgsUpdateTopic) (err error) {
	data := GetTopic(args.ID)
	if !checkTopic(data.TopicType, args.Options, args.Answer, args.Des) {
		err = errors.New("check topic")
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE blog_exam_topic SET update_at = NOW(), des = :des, score = :score, options = :options, answer = :answer, answer_analysis = :answer_analysis, params = :params WHERE id = :id AND (org_id = :org_id OR :org_id < 1)", args)
	if err != nil {
		return
	}
	deleteTopicCache(args.ID)
	return
}

// ArgsDeleteTopic 删除题目参数
type ArgsDeleteTopic struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteTopic 删除题目
func DeleteTopic(args *ArgsDeleteTopic) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "blog_exam_topic", "id = :id AND (org_id = :org_id OR :org_id < 1)", args)
	if err != nil {
		return
	}
	deleteTopicCache(args.ID)
	return
}

// 缓冲
func getTopicCacheMark(id int64) string {
	return fmt.Sprint("blog:exam:topic:id:", id)
}

func deleteTopicCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getTopicCacheMark(id))
}

// checkTopic 检查题目设置合理性
// topicType: 0 单选； 1 多选； 2 判断； 3 填空题； 4 问答题
func checkTopic(topicType int, options FieldsTopicOptions, answer []string, des string) bool {
	switch topicType {
	case 0:
		if len(answer) > 1 {
			return false
		}
		for _, v := range answer {
			isFind := false
			for _, v2 := range options {
				if v == v2.Mark {
					isFind = true
				}
			}
			if !isFind {
				return false
			}
		}
	case 1:
		for _, v := range answer {
			isFind := false
			for _, v2 := range options {
				if v == v2.Mark {
					isFind = true
				}
			}
			if !isFind {
				return false
			}
		}
	case 2:
		if len(answer) > 2 {
			return false
		}
	case 3:
		if len(options) > 0 {
			return false
		}
		for _, v := range answer {
			if !strings.Contains(des, fmt.Sprint("{", v, "}")) {
				return false
			}
		}
	case 4:
	default:
		if len(options) > 0 {
			return false
		}
		if len(answer) > 0 {
			return false
		}
		//不支持的类型
		return false
	}
	//反馈成功
	return true
}
