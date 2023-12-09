package BaseToken

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsTokenType 会话结构体
type FieldsTokenType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//key
	// 钥匙，用于配对
	Key string `db:"key" json:"key"`
	//配对密钥
	// key + password和其他信息混合后，验证是否正确，该内容只有在刚开始下放给用户
	// 如果不存在则跳过
	Password string `db:"password" json:"password"`
	//来源结构体
	// system为来源系统
	// id为来源ID
	// name为名称，用于快速识别来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//登陆渠道信息
	// 渠道必须是唯一的，如果重叠，则自动剔除旧的
	// system为来源系统，例如来源于小程序，weixin-xiaochengxu; 或来源于手机客户端, phone-a...
	LoginInfo CoreSQLFrom.FieldsFrom `db:"login_info" json:"loginInfo"`
	//IP地址
	IP string `db:"ip" json:"ip"`
	//是否记住我
	IsRemember bool `db:"is_remember" json:"isRemember"`
}
