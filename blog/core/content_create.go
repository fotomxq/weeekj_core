package BlogCore

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"github.com/mozillazg/go-pinyin"
	"strings"
)

// ArgsCreateContent 创建新的词条参数
type ArgsCreateContent struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//文章类型
	// 0 普通文章; 1 挂靠视频; 3 第三方跳转; 4 组织地图
	// 挂靠或跳转内容，将在des中做特殊描述
	ContentType int `db:"content_type" json:"contentType" check:"intThan0" empty:"true"`
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

// CreateContent 创建新的词条
func CreateContent(args *ArgsCreateContent) (data FieldsContent, err error) {
	//检查文章类型
	if err = checkContentType(args.ContentType); err != nil {
		return
	}
	//生成key
	if args.Key == "" {
		args.Key = makeKey(args.Key, 0, args.Title)
	}
	if args.Key == "" {
		err = errors.New("key error")
		return
	}
	//修正参数
	if len(args.Tags) < 1 {
		args.Tags = pq.Int64Array{}
	}
	if len(args.DesFiles) < 1 {
		args.DesFiles = pq.Int64Array{}
	}
	//构建新的数据
	var contentID int64
	contentID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_core_content (content_type, audit_at, audit_des, org_id, user_id, bind_id, param1, param2, param3, key, parent_id, is_top, sort_id, tags, title, title_des, cover_file_id, des_files, des, params) VALUES (:content_type, to_timestamp(0),'',:org_id, :user_id, :bind_id, :param1, :param2, :param3, :key, 0, :is_top, :sort_id, :tags, :title, :title_des, :cover_file_id, :des_files, :des, :params)", args)
	if err != nil {
		return
	}
	//获取数据
	data = getContentID(contentID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//新的文章
	pushCreate(contentID)
	//请求审核
	pushAudit(contentID)
	//反馈
	return
}

// 生成文章key
func makeKey(key string, oldID int64, title string) string {
	//锁定机制，避免重叠数据
	makeKeyLock.Lock()
	defer makeKeyLock.Unlock()
	//初始化id
	var id int64
	//检查key
	if key != "" {
		err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM blog_core_content WHERE key = $1", key)
		if err != nil || id < 1 {
			return key
		}
		if oldID > 0 && id == oldID {
			return key
		}
	}
	//修正key
	key = strings.Replace(key, " ", "_", -1)
	//生成拼音
	newKeys := pinyin.LazyPinyin(title, pinyin.NewArgs())
	var newKey string
	for k := 0; k < len(newKeys); k++ {
		newKey = newKey + newKeys[k]
	}
	//检查key是否存在
	if key != "" {
		err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM blog_core_content WHERE key = $1", newKey)
		if err != nil || id < 1 {
			return newKey
		}
	}
	step := 1
	for {
		step += 1
		if step > 100 {
			break
		}
		var vID int64
		newKey2 := newKey + fmt.Sprint(step)
		if step > 5 {
			newKey2 = fmt.Sprint(newKey, "_", CoreFilter.GetRandStr4(10))
			err := Router2SystemConfig.MainDB.Get(&vID, "SELECT id FROM blog_core_content WHERE key = $1", newKey2)
			if err == nil && vID > 0 {
				continue
			}
		} else {
			err := Router2SystemConfig.MainDB.Get(&vID, "SELECT id FROM blog_core_content WHERE key = $1", newKey2)
			if err == nil && vID > 0 {
				continue
			}
		}
		return newKey2
	}
	//反馈失败
	return ""
}
