package BasePython

import (
	"time"
)

// fieldsWait 等待处理的模块函数库
// 数据库会保存处理情况，方便一些特殊的处理方法回调，同时避免反复调用处理的问题
type fieldsWait struct {
	//ID
	// 将采用系统来源/时间月+时间日/ID，作为文件存储路径
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 满足一定条件后将用于自动删除数据的处理机制
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//超时时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//是否处理完成
	IsFinish bool `db:"is_finish" json:"isFinish"`
	//系统来源
	System string `db:"system" json:"system"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//唯一标识码
	// 在系统来源内唯一，如果存在BindID则优先需bindID也唯一
	// 可用于参数混合，避免重复处理数据
	Mark string `db:"mark" json:"mark"`
	//推送数据流
	Param []byte `db:"param" json:"param"`
}
