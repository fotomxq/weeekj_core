package ClassComment

import (
	"database/sql"
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "gitee.com/weeekj/weeekj_core/v5/core/sql/analysis"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//Comment 通用评论模块
/**
可以模块引用到
*/
type Comment struct {
	//主表
	TableName string
	//用户是否可以为某个绑定关系创建多个评论
	UserMoreComment bool
	//用户是否可以编辑评论
	UserEditComment bool
	//用户是否可以删除评论
	UserDeleteComment bool
	//组织是否可以删除评论
	OrgDeleteComment bool
	//系统来源
	// 用于推送nats等操作
	System string
}

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//评论ID
	// 评论被删除后出现，指向新的评论
	// 下级评论的上级江全部改为该新的ID
	CommentID int64 `db:"comment_id" json:"commentID" check:"id" empty:"true"`
	//上级ID
	ParentID int64 `json:"parentID" check:"id" empty:"true"`
	//绑定组织
	// 该组织根据资源来源设定
	// 如果是平台资源，则为0
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//评价类型
	// 0 好评 1 中立 2 差评
	LevelType int `db:"level_type" json:"levelType"`
	//分数范围
	LevelMin int `db:"level_min" json:"levelMin"`
	LevelMax int `db:"level_max" json:"levelMax"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取列表
func (t *Comment) GetList(args *ArgsGetList) (dataList []FieldsComment, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.CommentID > 0 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.CommentID
	}
	if args.ParentID > 0 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > 0 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.BindID > 0 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.LevelType > 0 {
		where = where + " AND level_type = :level_type"
		maps["level_type"] = args.LevelType
	}
	if args.LevelMin > 0 {
		where = where + " AND level >= :level_min"
		maps["level_min"] = args.LevelMin
	}
	if args.LevelMax > 0 {
		where = where + " AND level <= :level_max"
		maps["level_max"] = args.LevelMax
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		t.TableName,
		"id",
		"SELECT id, create_at, delete_at, comment_id, org_id, user_id, bind_id, parent_id, level_type, level, title, des, des_files, params "+"FROM "+t.TableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "delete_at", "level"},
	)
	return
}

// GetByID 获取评论ID
func (t *Comment) GetByID(id int64) (data FieldsComment) {
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, delete_at, comment_id, org_id, user_id, bind_id, parent_id, level_type, level, title, des, des_files, params "+"FROM "+t.TableName+" WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	return
}

// GetCountByFrom 获取评论人数
func (t *Comment) GetCountByFrom(bindID int64) (count int64) {
	where := "bind_id = :bind_id"
	maps := map[string]interface{}{
		"bind_id": bindID,
	}
	where = CoreSQL.GetDeleteSQL(false, where)
	var err error
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, t.TableName, "id", where, maps)
	if err != nil {
		count = 0
		return
	}
	return
}

// ArgsGetAnalysisAvg 统计指定范围的评价平均值参数
type ArgsGetAnalysisAvg struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//购买人
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
}

// DataGetAnalysisAvg 统计指定范围的评价平均值数据
type DataGetAnalysisAvg struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//数据
	Count int `db:"count" json:"count"`
}

// GetAnalysisAvg 统计指定范围的评价平均值
func (t *Comment) GetAnalysisAvg(args *ArgsGetAnalysisAvg) (dataList []DataGetAnalysisAvg, err error) {
	where := "((org_id = :org_id OR :org_id < 1) OR (bind_id = :bind_id OR :bind_id < 1) OR (user_id = :user_id OR :user_id < 1)) AND parent_id = 0"
	maps := map[string]interface{}{
		"org_id":  args.OrgID,
		"bind_id": args.BindID,
		"user_id": args.UserID,
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", AVG(level) as count FROM "+t.TableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

type ArgsGetAnalysisAvgOne struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
}

type DataGetAnalysisAvgOne struct {
	//数据
	Count int `db:"count" json:"count"`
}

// GetAnalysisAvgOne 统计指定绑定的评价平均值
func (t *Comment) GetAnalysisAvgOne(args *ArgsGetAnalysisAvgOne) (resultData int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("recover failed,", r))
		}
	}()
	createAt := CoreFilter.GetNowTimeCarbon()
	switch args.TimeType {
	case "year":
		createAt = createAt.SubYear()
	case "month":
		createAt = createAt.SubMonth()
	case "day":
		createAt = createAt.SubDay()
	case "hour":
		createAt = createAt.SubHour()
	}
	where := "(org_id = :org_id OR :org_id < 1) AND parent_id = 0 AND (:bind_id < 1 OR bind_id = :bind_id) AND create_at > :create_at"
	maps := map[string]interface{}{
		"org_id":    args.OrgID,
		"bind_id":   args.BindID,
		"create_at": createAt.Time,
	}
	var data DataGetAnalysisAvgOne
	err = CoreSQL.GetOne(
		Router2SystemConfig.MainDB.DB,
		&data,
		"SELECT AVG(level) as count "+"FROM "+t.TableName+" WHERE "+where,
		maps,
	)
	resultData = data.Count
	return
}

// ArgsCreate 创建新的评论参数
type ArgsCreate struct {
	//上级ID
	// 评论的上下级关系，一旦建立无法修改
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//绑定组织
	// 该组织根据资源来源设定
	// 如果是平台资源，则为0
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属用户
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//绑定内容
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//评价类型
	// 0 好评 1 中立 2 差评
	LevelType int `db:"level_type" json:"levelType"`
	//分数
	Level int `db:"level" json:"level"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"100"`
	//内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//介绍图文
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Create 创建新的评论
func (t *Comment) Create(args *ArgsCreate) (data FieldsComment, err error) {
	//检查参数
	if err = t.checkLevelType(args.LevelType); err != nil {
		return
	}
	//检查上级和当前的绑定来源是否一致
	if args.ParentID > 0 {
		var parentData FieldsComment
		err = Router2SystemConfig.MainDB.Get(&parentData, "SELECT id "+"FROM "+t.TableName+" WHERE id = $1 AND bind_id = $2", args.ParentID, args.BindID)
		if err != nil {
			err = errors.New("parent not exist, " + err.Error())
			return
		} else {
			if parentData.ID < 1 {
				err = errors.New("parent not exist")
				return
			}
		}
		//检查通过，继续执行后续
	}
	//检查评论数量
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id "+"FROM "+t.TableName+" WHERE org_id = $1 AND bind_id = $2 AND user_id = $3", args.OrgID, args.BindID, args.UserID)
	if err == nil && data.ID > 0 {
		if t.UserMoreComment {
			err = errors.New("user have comment")
			return
		}
	}
	//写入数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, t.TableName, "INSERT "+"INTO "+t.TableName+" (comment_id, org_id, user_id, bind_id, parent_id, level_type, level, title, des, des_files, params) VALUES (0, :org_id, :user_id, :bind_id, :parent_id, :level_type, :level, :title, :des, :des_files, :params)", args, &data)
	if err != nil {
		return
	}
	//通知新的评论
	CoreNats.PushDataNoErr("/class/comment", "new", data.ID, t.System, map[string]interface{}{
		"parentID":  args.ParentID,
		"orgID":     args.OrgID,
		"userID":    args.UserID,
		"bindID":    args.BindID,
		"levelType": args.LevelType,
		"level":     args.Level,
	})
	//反馈
	return
}

// ArgsUpdate 修改评论参数
type ArgsUpdate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//绑定组织
	// 用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属用户
	// 用于验证
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//评价类型
	// 0 好评 1 中立 2 差评
	LevelType int `db:"level_type" json:"levelType"`
	//分数
	Level int `db:"level" json:"level"`
	//标题
	Title string `db:"title" json:"title" check:"name" empty:"true"`
	//内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//介绍图文
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Update 修改评论
func (t *Comment) Update(args *ArgsUpdate) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	if args.UserID > 0 && !t.UserEditComment {
		err = errors.New("user cannot edit comment")
		return
	}
	//检查参数
	if err = t.checkLevelType(args.LevelType); err != nil {
		return
	}
	//获取原始数据
	var data FieldsComment
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id, user_id, bind_id, parent_id "+"FROM "+t.TableName+" WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR user_id = $3)", args.ID, args.OrgID, args.UserID)
	if err != nil {
		return
	}
	//开始事务
	tx := Router2SystemConfig.MainDB.MustBegin()
	//重构新的数据
	var stmt *sqlx.NamedStmt
	stmt, err = tx.PrepareNamed("INSERT " + "INTO " + t.TableName + " (comment_id, org_id, user_id, bind_id, parent_id, level_type, level, title, des, des_files, params) VALUES (0, :org_id, :user_id, :bind_id, :parent_id, :level_type, :level, :title, :des, :des_files, :params)  RETURNING id;")
	var lastID int64
	lastID, err = CoreSQL.LastRowsAffectedCreate(tx, stmt, map[string]interface{}{
		"org_id":     data.OrgID,
		"user_id":    data.UserID,
		"bind_id":    data.BindID,
		"parent_id":  data.ParentID,
		"level_type": args.LevelType,
		"level":      args.Level,
		"title":      args.Title,
		"des":        args.Des,
		"des_files":  args.DesFiles,
		"params":     args.Params,
	}, err)
	if err != nil {
		return
	}
	//修改本数据
	var result sql.Result
	result, err = tx.NamedExec("UPDATE "+t.TableName+" SET delete_at = NOW(), comment_id = :comment_id WHERE id = :id", map[string]interface{}{
		"id":         data.ID,
		"comment_id": lastID,
	})
	err = CoreSQL.LastRowsAffected(tx, result, err)
	if err != nil {
		return
	}
	//修改子数据
	_, err = tx.NamedExec("UPDATE "+t.TableName+" SET parent_id = :new_parent_id WHERE parent_id = :parent_id", map[string]interface{}{
		"parent_id":     data.ID,
		"new_parent_id": lastID,
	})
	err = CoreSQL.LastRows(tx, err)
	if err != nil {
		return
	}
	//运行sql
	err = tx.Commit()
	return
}

// ArgsDeleteByID 删除评论ID参数
type ArgsDeleteByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//绑定组织
	// 用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属用户
	// 用于验证
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// DeleteByID 删除评论ID
func (t *Comment) DeleteByID(args *ArgsDeleteByID) (err error) {
	if args.OrgID > 0 && !t.OrgDeleteComment {
		err = errors.New("org cannot delete comment")
		return
	}
	if args.UserID > 0 && !t.UserDeleteComment {
		err = errors.New("user cannot delete comment")
		return
	}
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, t.TableName, "id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", args)
	return
}

// ArgsDeleteByUser 删除用户所有评论参数
type ArgsDeleteByUser struct {
	//所属用户
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//绑定组织
	// 用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteByUser 删除用户所有评论
func (t *Comment) DeleteByUser(args *ArgsDeleteByUser) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, t.TableName, "user_id = :user_id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteByBind 删除来源所有评论参数
type ArgsDeleteByBind struct {
	//绑定内容
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//绑定组织
	// 用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteByBind 删除来源所有评论
func (t *Comment) DeleteByBind(args *ArgsDeleteByBind) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, t.TableName, "bind_id = :bind_id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteByOrg 删除组织所有数据参数
type ArgsDeleteByOrg struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteByOrg 删除组织所有数据
func (t *Comment) DeleteByOrg(args *ArgsDeleteByOrg) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, t.TableName, "org_id = :org_id", args)
	return
}

// checkLevelType 检查评级
func (t *Comment) checkLevelType(levelType int) (err error) {
	switch levelType {
	case 0:
	case 1:
	case 2:
	default:
		err = errors.New("level type error")
		return
	}
	return
}
