package Router2SystemConfig

import (
	"errors"
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
)

func LoadPostgres() (err error) {
	// 连接数据库
	timeZone := Cfg.Section("db").Key("time_zone").Value()
	if timeZone == "" {
		timeZone = "PRC"
	}
	MainDB.MaxConnect, _ = Cfg.Section("db").Key("postgresql_max_connect").Int()
	if MainDB.MaxConnect < 1 {
		MainDB.MaxConnect = 90
	}
	MainDB.ConnectExpireSec, _ = Cfg.Section("db").Key("postgresql_expire_sec").Int()
	if MainDB.ConnectExpireSec < 1 {
		MainDB.ConnectExpireSec = 30
	}
	err = MainDB.Init(PostgresURL, fmt.Sprint(RootDir, CoreFile.Sep, "install"), timeZone)
	if err != nil {
		err = errors.New(fmt.Sprint("无法连接postgres数据库, " + err.Error()))
		return
	}
	return
}
