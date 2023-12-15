package CoreRPCX

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
)

// 通用结构体声明
// 可用于server内部快速构建请求和接收结构
type ArgsEmpty struct {
}

type ArgsList struct {
	Pages  CoreSQLPages.ArgsDataList
	Search string
}

type ArgsString struct {
	String string
}

type ArgsMark struct {
	Mark string
}

type ArgsID struct {
	ID string
}

type ArgsIDHash struct {
	ID   string
	Hash string
}

type ArgsIDUser struct {
	ID     string
	UserID string
}

type ArgsIDFrom struct {
	ID   string
	From CoreSQLFrom.FieldsFrom
}

type ArgsOpen struct {
	Open bool
}

type ArgsRunTime struct {
	RunTime int64
}

type ArgsSetCache struct {
	Open  bool
	Limit int
}

type ArgsFrom struct {
	From CoreSQLFrom.FieldsFrom
}

type ArgsInt64 struct {
	Data int64
}

type ReplyEmpty struct {
}

type ReplyBool struct {
	Open bool
}

type ReplyString struct {
	Data string
}

type ReplyInt struct {
	Data int
}

type ReplyInt64 struct {
	Data int64
}

type ReplyFloat64 struct {
	Data float64
}

type ReplyByte struct {
	Data []byte
}
