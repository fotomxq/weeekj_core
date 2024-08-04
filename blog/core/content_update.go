package BlogCore

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsUpdateContent 修改词条参数
type ArgsUpdateContent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//扩展筛选项
	Param1 int64 `db:"param1" json:"param1"`
	Param2 int64 `db:"param2" json:"param2"`
	Param3 int64 `db:"param3" json:"param3"`
	//唯一标识码key
	// 作为id的补充，自动填写时，将自动生成随机字符串
	// 默认根据标题或标题拼音得出
	Key string `db:"key" json:"key" check:"mark" empty:"true"`
	//是否置顶
	IsTop bool `db:"is_top" json:"isTop" check:"bool"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//小标题
	TitleDes string `db:"title_des" json:"titleDes" check:"title" min:"1" max:"600" empty:"true"`
	//封面文件
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//附加封面图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"9000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateContent 修改词条
func UpdateContent(args *ArgsUpdateContent) (newData FieldsContent, err error) {
	//修正参数
	if len(args.DesFiles) < 1 {
		args.DesFiles = pq.Int64Array{}
	}
	//找到旧的数据
	var oldData FieldsContent
	err = Router2SystemConfig.MainDB.Get(&oldData, "SELECT id, create_at, update_at, delete_at, audit_at, audit_des, org_id, user_id, bind_id, param1, param2, param3, visit_count, key, parent_id, publish_at, is_top, sort_id, tags, title, title_des, cover_file_id, des_files, des, params FROM blog_core_content WHERE id = $1 AND ($2 < 0 OR org_id = $2) AND ($3 < 0 OR bind_id = $3) AND delete_at < to_timestamp(1000000) AND parent_id = 0", args.ID, args.OrgID, args.BindID)
	if err != nil || oldData.ID < 1 {
		err = errors.New(fmt.Sprint("not find old data, ", err))
		return
	}
	//生成key
	if args.Key == "" {
		args.Key = makeKey(args.Key, args.ID, args.Title)
	}
	if args.Key == "" {
		err = errors.New("key error")
		return
	}
	if len(oldData.Tags) < 1 {
		oldData.Tags = []int64{}
	}
	if len(oldData.DesFiles) < 1 {
		oldData.DesFiles = pq.Int64Array{}
	}
	//构建新的数据，历史存档
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_core_content (create_at, update_at, audit_at, audit_des, org_id, user_id, bind_id, param1, param2, param3, visit_count, key, parent_id, is_top, sort_id, tags, title, title_des, cover_file_id, des_files, des, params) VALUES (:create_at, :update_at, :audit_at, :audit_des, :org_id, :user_id, :bind_id, :param1, :param2, :param3, :visit_count, :key, :parent_id, :is_top, :sort_id, :tags, :title, :title_des, :cover_file_id, :des_files, :des, :params)", map[string]interface{}{
		"create_at":     oldData.CreateAt,
		"update_at":     oldData.UpdateAt,
		"audit_at":      oldData.AuditAt,
		"audit_des":     oldData.AuditDes,
		"org_id":        oldData.OrgID,
		"user_id":       oldData.UserID,
		"bind_id":       oldData.BindID,
		"param1":        oldData.Param1,
		"param2":        oldData.Param2,
		"param3":        oldData.Param3,
		"visit_count":   oldData.VisitCount,
		"key":           oldData.Key,
		"parent_id":     oldData.ID,
		"is_top":        oldData.IsTop,
		"sort_id":       oldData.SortID,
		"tags":          oldData.Tags,
		"title":         oldData.Title,
		"title_des":     oldData.TitleDes,
		"cover_file_id": oldData.CoverFileID,
		"des_files":     oldData.DesFiles,
		"des":           oldData.Des,
		"params":        oldData.Params,
	})
	if err != nil {
		return
	}
	//更新数据
	if len(args.Tags) < 1 {
		args.Tags = pq.Int64Array{}
	}
	if len(args.DesFiles) < 1 {
		args.DesFiles = pq.Int64Array{}
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE blog_core_content SET update_at = NOW(), publish_at = to_timestamp(0), audit_at = to_timestamp(0), bind_id = :bind_id, param1 = :param1, param2 = :param2, param3 = :param3, key = :key, is_top = :is_top, sort_id = :sort_id, tags = :tags, title = :title, title_des = :title_des, cover_file_id = :cover_file_id, des = :des, des_files = :des_files, params = :params WHERE id = :id", map[string]interface{}{
		"id":            args.ID,
		"key":           args.Key,
		"bind_id":       oldData.BindID,
		"param1":        args.Param1,
		"param2":        args.Param2,
		"param3":        args.Param3,
		"is_top":        args.IsTop,
		"sort_id":       args.SortID,
		"tags":          args.Tags,
		"title":         args.Title,
		"title_des":     args.TitleDes,
		"cover_file_id": args.CoverFileID,
		"des_files":     args.DesFiles,
		"des":           args.Des,
		"params":        args.Params,
	})
	if err != nil {
		return
	}
	deleteContentCacheByID(args.ID)
	newData = getContentID(args.ID)
	if newData.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsMoveSort 修改指定分类下文章迁移到另外一个分类参数
type ArgsMoveSort struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id"`
	//目标分类ID
	DestSortID int64 `db:"dest_sort_id" json:"destSortID" check:"id"`
}

// MoveSort 修改指定分类下文章迁移到另外一个分类
func MoveSort(args *ArgsMoveSort) (err error) {
	//检查分类是否属于该商户？
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM blog_core_sort WHERE id = $1 AND bind_id = $2", args.SortID, args.OrgID)
	if err != nil || id < 1 {
		err = errors.New("org not own sort ")
		return
	}
	var destID int64
	err = Router2SystemConfig.MainDB.Get(&destID, "SELECT id FROM blog_core_sort WHERE id = $1 AND bind_id = $2", args.DestSortID, args.OrgID)
	if err != nil || destID < 1 {
		err = errors.New("org not own dest sort ")
		return
	}
	//删除缓冲
	var ids []FieldsContent
	err = Router2SystemConfig.MainDB.Select(&ids, "SELECT id FROM blog_core_content WHERE sort_id = $1", args.SortID)
	if err == nil {
		for _, v := range ids {
			deleteContentCacheByID(v.ID)
		}
	}
	//修改操作
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE blog_core_content SET sort_id = :dest_sort_id WHERE sort_id = :sort_id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	return
}

// ArgsUpdatePublish 发布文章参数
type ArgsUpdatePublish struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// UpdatePublish 发布文章
func UpdatePublish(args *ArgsUpdatePublish) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE blog_core_content SET publish_at = NOW(), audit_at = to_timestamp(0) WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND (:user_id < 1 OR user_id = :user_id)", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteContentCacheByID(args.ID)
	//请求审核
	pushAudit(args.ID)
	//反馈
	return
}

// ArgsUpdateAudit 审核文章参数
type ArgsUpdateAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否通过审核
	IsAudit bool `db:"is_audit" json:"isAudit" check:"bool"`
	//审核内容
	AuditDes string `db:"audit_des" json:"auditDes" check:"des" min:"1" max:"600" empty:"true"`
}

// UpdateAudit 审核文章
func UpdateAudit(args *ArgsUpdateAudit) (err error) {
	if args.IsAudit {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE blog_core_content SET audit_at = NOW() WHERE id = :id AND (org_id = :org_id OR :org_id < 1)", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
		})
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE blog_core_content SET audit_at = to_timestamp(0), audit_des = :audit_des WHERE id = :id AND (org_id = :org_id OR :org_id < 1)", map[string]interface{}{
			"id":        args.ID,
			"org_id":    args.OrgID,
			"audit_des": args.AuditDes,
		})
	}
	if err != nil {
		return
	}
	//如果审核通过，推送审核结束
	if args.IsAudit {
		pushAuditDone(args.ID)
	}
	//删除缓冲
	deleteContentCacheByID(args.ID)
	//反馈
	return
}

// UpdateContentTop 设置文章置顶
func UpdateContentTop(id int64, isTop bool) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE blog_core_content SET is_top = :is_top WHERE id = :id", map[string]interface{}{
		"is_top": isTop,
		"id":     id,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteContentCacheByID(id)
	return
}
