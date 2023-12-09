package BaseEarlyWarning

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

//送达人关系
type FieldsToType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//绑定用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//姓名
	Name string `db:"name" json:"name"`
	//备注和描述
	Des string `db:"des" json:"des"`
	//联系电话国家区号
	PhoneNationCode string `db:"phone_nation_code" json:"phoneNationCode"`
	//联系电话
	Phone string `db:"phone" json:"phone"`
	//联系邮箱
	Email string `db:"email" json:"email"`
}

//模版
type FieldsTemplateType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//模版标识码
	// 用于程序内的识别
	// 如果客户修改会造成识别异常或丢弃告警模块
	Mark string `db:"mark" json:"mark"`
	//名称
	// 后台区分使用
	Name string `db:"name" json:"name"`
	//默认超出时间
	// 模版超出时间和送达人超出时间，会自动取最短值计算
	DefaultExpireTime string `db:"default_expire_time" json:"defaultExpireTime"`
	//标题
	// email专用，短信无效
	// 同样可带有[1]结构变量
	Title string `db:"title" json:"title"`
	//内容
	// 例如: 你好！尊敬的[1]用户，欢迎在[2]时间访问[3]系统。
	Content string `db:"content" json:"content"`
	//短消息对应的模版ID
	// 例如，腾讯云需审核过的模版ID，才能使用短信服务
	TemplateID string `db:"template_id" json:"templateID"`
	//内容和可绑定的关系组
	// 对应Content中的变量和值
	// 例如：[]string{"[1]","[2]",...}
	BindData FieldsTemplateBindData `db:"bind_data" json:"bindData"`
}

type FieldsTemplateBindData []string

func (t FieldsTemplateBindData) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsTemplateBindData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

//模版和送达人关系结构体
type FieldsBindType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//优先级
	// 从小到大排序，先给最前面的发送，超时后自动到下一级；
	// 同一级别将视为并列发送处理
	Level int `db:"level" json:"level"`
	//进入下一级条件
	// or 只要有一个同意则完成，否则到下一级
	// and 必须所有人同意才能完成，否则到下一级
	// none 本集自动终结，不再给下一个级别通知
	// 注意，同一个级别的任意修改操作，将造成同级别所有关系条件调整，以确保一致性
	LevelMode string `db:"level_mode" json:"levelMode"`
	//通知下一级的等待时间
	NextWaitTime string `db:"next_wait_time" json:"nextWaitTime"`
	//关系人
	ToID int64 `db:"to_id" json:"toID"`
	//模版ID
	TemplateID int64 `db:"template_id" json:"templateID"`
	//送出方式
	NeedPhone bool `db:"need_phone" json:"needPhone"`
	NeedSMS   bool `db:"need_sms" json:"needSMS"`
	NeedEmail bool `db:"need_email" json:"needEmail"`
	NeedAPP   bool `db:"need_app" json:"needAPP"`
}

//等待送达信息
type FieldsWaitType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//绑定关系
	BindID int64 `db:"bind_id" json:"bindID"`
	//关联级别
	Level int `db:"level" json:"level"`
	//送达人
	ToID int64 `db:"to_id" json:"toID"`
	//已经全部送出
	IsSend bool `db:"is_send" json:"isSend"`
	//是否已读
	IsRead bool `db:"is_read" json:"isRead"`
	//已经超时处理
	// 给下一级发送完成消息后，自动标记该数据
	ExpireFinish bool `db:"expire_finish" json:"expireFinish"`
	//超时时间
	// 如果未读，超出后自动根据送达下一级发送消息
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//采用模版
	TemplateID int64 `db:"template_id" json:"templateID"`
	//消息数据
	// 来自模版的数据结构
	Content string `db:"content" json:"content"`
	//消息变量
	BindData FieldsWaitBindData `db:"bind_data" json:"bindData"`
	//送出方式
	// 电话联络
	NeedPhone bool `db:"need_phone" json:"needPhone"`
	// 短信推送
	NeedSMS bool `db:"need_sms" json:"needSMS"`
	// 邮件推送
	NeedEmail bool `db:"need_email" json:"needEmail"`
	// APP通知
	NeedAPP bool `db:"need_app" json:"needAPP"`
}

type FieldsWaitBindData map[string]string

func (t FieldsWaitBindData) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsWaitBindData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}