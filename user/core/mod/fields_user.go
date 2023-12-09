package UserCoreMod

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsUserType 用户结构
type FieldsUserType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//状态
	// 0 -> ban后，用户可以正常登录，但一切内容无法使用
	// 1 -> audit后，用户登录后无法正常使用，但提示不太一样
	// 2 -> public 正常访问
	Status int `db:"status" json:"status"`
	//组织ID
	// 如果为空，则说明是平台的用户；否则为对应组织的用户
	// 所有获取的方法，都需要给与该ID参数，也可以留空，否则禁止获取
	OrgID int64 `db:"org_id" json:"orgID"`
	//姓名
	Name string `db:"name" json:"name"`
	//密码
	Password string `db:"password" json:"password"`
	//绑定手机号的国家代码
	NationCode string `db:"nation_code" json:"nationCode"`
	//手机号码，绑定后的手机
	Phone string `db:"phone" json:"phone"`
	//手机是否验证
	PhoneVerify time.Time `db:"phone_verify" json:"phoneVerify"`
	//邮箱，如果不存在手机则必须存在的
	Email string `db:"email" json:"email"`
	//邮箱是否验证
	EmailVerify time.Time `db:"email_verify" json:"emailVerify"`
	//用户名，如果不存在手机号和邮箱，则必须存在的
	Username string `db:"username" json:"username"`
	//用户头像
	Avatar int64 `db:"avatar" json:"avatar"`
	//上级关系
	Parents FieldsUserParents `db:"parents" json:"parents"`
	//用户组
	Groups FieldsUserGroupsType `db:"groups" json:"groups"`
	//信息结构体
	// 特殊情况下，如果不希望创建新的模块，可以使用该字段实现特定目标
	// 例如是否进入过某个页面，或是否进行过某些行为的标识码
	Infos CoreSQLConfig.FieldsConfigsType `db:"infos" json:"infos"`
	//登录结构体
	// 提供给第三方登录的数据接口
	Logins FieldsUserLoginsType `db:"logins" json:"logins"`
	//用户分类
	SortID int64 `db:"sort_id" json:"sortID"`
	//用户标签
	Tags pq.Int64Array `db:"tags" json:"tags"`
}

// FieldsUserParents 上下级关系处理
type FieldsUserParents []FieldsUserParent

func (t FieldsUserParents) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsUserParents) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsUserParent struct {
	//系统类型标识码
	// 用于指定不同类型的系统、模块的上下级关系
	// eg: transport / finance
	System string `db:"system" json:"system"`
	//上级ID
	ParentID int64 `db:"parentID" json:"parentID"`
	//权限标记
	Operate []string `db:"operate" json:"operate"`
}

func (t FieldsUserParent) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsUserParent) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsUserGroupsType 用户组信息结构
type FieldsUserGroupsType []FieldsUserGroupType

func (t FieldsUserGroupsType) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsUserGroupsType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsUserGroupType struct {
	//用户组ID
	GroupID int64 `db:"group_id" json:"groupID"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
}

// sql底层处理器
func (t FieldsUserGroupType) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsUserGroupType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// 登录结构体
type FieldsUserLoginsType []FieldsUserLoginType

// sql底层处理器
func (t FieldsUserLoginsType) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsUserLoginsType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsUserLoginType struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//值
	Val string `db:"val" json:"val"`
	//配置结构
	// 可任意建立，主要指一些特定的标记，如API反馈的原始数据等
	Config string `db:"config" json:"config"`
}

// sql底层处理器
func (t FieldsUserLoginType) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsUserLoginType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
