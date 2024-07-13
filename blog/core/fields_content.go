package BlogCore

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsContent 文章核心模块
type FieldsContent struct {
	//ID
	ID int64 `db:"id" json:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0"`
	//文章类型
	// 0 普通文章; 1 挂靠视频; 3 第三方跳转; 4 组织地图
	// 挂靠或跳转内容，将在des中做特殊描述
	ContentType int `db:"content_type" json:"contentType"`
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt" default:"0"`
	//审核拒绝原因
	AuditDes string `db:"audit_des" json:"auditDes" check:"des" min:"1" max:"300" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" index:"true" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" index:"true" check:"id" empty:"true"`
	//扩展筛选项
	Param1 int64 `db:"param1" json:"param1"`
	Param2 int64 `db:"param2" json:"param2"`
	Param3 int64 `db:"param3" json:"param3"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserID int64 `db:"user_id" json:"userID" index:"true" check:"id" empty:"true"`
	//访问量
	VisitCount int64 `db:"visit_count" json:"visitCount"`
	//唯一标识码key
	// 作为id的补充，自动填写时，将自动生成随机字符串
	// 默认根据标题或标题拼音得出
	Key string `db:"key" json:"key" index:"true"`
	//归属关系
	// 删除后作为原始文档的子项目存在，key将自动失效
	ParentID int64 `db:"parent_id" json:"parentID" index:"true" check:"id" empty:"true"`
	//是否公开
	// 非公开数据将作为草稿或私有数据存在，只有管理员可以看到
	PublishAt time.Time `db:"publish_at" json:"publishAt" default:"0"`
	//是否置顶
	IsTop bool `db:"is_top" json:"isTop"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID" index:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//标题
	Title string `db:"title" json:"title" check:"des" min:"1" max:"300" empty:"true"`
	//小标题
	TitleDes string `db:"title_des" json:"titleDes" check:"des" min:"1" max:"300" empty:"true"`
	//封面文件
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" index:"true" check:"id" empty:"true"`
	//附加封面图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//内容
	Des string `db:"des" json:"des" check:"des" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
