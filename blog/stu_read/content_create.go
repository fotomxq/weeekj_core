package BlogStuRead

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsCreateContent 创建新的词条参数
type ArgsCreateContent struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//文章类型
	// 0 普通文章; 1 挂靠视频; 3 第三方跳转
	// 挂靠或跳转内容，将在des中做特殊描述
	ContentType int `db:"content_type" json:"contentType" check:"intThan0" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//小标题
	TitleDes string `db:"title_des" json:"titleDes" check:"title" min:"1" max:"600" empty:"true"`
	//封面文件
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//附加封面图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateContent 创建新的词条
func CreateContent(args *ArgsCreateContent) (data FieldsContent, err error) {
	//修正参数
	if len(args.DesFiles) < 1 {
		args.DesFiles = pq.Int64Array{}
	}
	//检查文章类型
	if err = checkContentType(args.ContentType); err != nil {
		return
	}
	//构建新的数据
	var contentID int64
	contentID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_stu_read_content (content_type, org_id, title, title_des, cover_file_id, des_files, des, params, visit_count) VALUES (:content_type, :org_id, :title, :title_des, :cover_file_id, :des_files, :des, :params, 0)", args)
	if err != nil {
		return
	}
	//获取数据
	data = getContentID(contentID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//反馈
	return
}
