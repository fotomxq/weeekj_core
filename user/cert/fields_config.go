package UserCert

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/lib/pq"
	"time"
)

//FieldsConfig 证件配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//默认过期时间
	// 到达时间后，将自动标记请求为过期并删除
	// 单位：秒
	DefaultExpireTime int64 `db:"default_expire_time" json:"defaultExpireTime"`
	//绑定组织
	// 该组织根据资源来源设定
	// 如果是平台资源，则为0
	OrgID int64 `db:"org_id" json:"orgID"`
	//证件名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//图文
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//开关
	// 关闭后将禁止新增申请
	AllowOpen bool `db:"allow_open" json:"allowOpen"`
	//提交间隔
	// 提交之间的时间间隔，也可以理解为审核失败后再次提交的时间间隔
	// 单位：秒
	PostInterval int64 `db:"post_interval" json:"postInterval"`
	//证件申请步骤
	Steps FieldsConfigSteps `db:"steps" json:"steps"`
}

//证件步骤
// 在配置内部嵌入，每个步骤都可设置并列入否
// 写入顺序将影响呈现给用户的步骤
type FieldsConfigSteps []FieldsConfigStep

type FieldsConfigStep struct {
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//图文
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//并列标识码
	// 相同的标识码会被并列处理
	Mark string `db:"mark" json:"mark"`
	//步骤所需材料
	// 写入顺序影响上下文顺序
	Contents FieldsContents `db:"contents" json:"contents"`
}

//扩展结构
type FieldsContents []FieldsContent

//sql底层处理器
func (t FieldsContents) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsContents) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsContent struct {
	//标识码
	// 系统不会做唯一判断，该标识码用于前端判定处理，如特殊的页面需做定制化处理的
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//是否必须填写
	Must bool `db:"must" json:"must"`
	//值类型
	// 0 string / 1 bool / 2 int / 3 int64 / 4 float64
	// 5 time 时间 / 6 daytime 带有日期的时间 / 7 unix 时间戳
	// 8 fileID 文件ID / 9 fileIDList 文件ID列
	// 10 userID 用户ID / 11 userIDList 用户ID列
	ValType int `db:"val_type" json:"valType"`
	//正则表达式范围
	ValCheck string `db:"val_check" json:"valCheck"`
	//值/默认值
	Val string `db:"val" json:"val"`
}

//sql底层处理器
func (t FieldsContent) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsContent) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
