package CorePostgres

import (
	"errors"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
)

// Install 自动安装配置文件
func (t *Client) Install() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	if t.InstallDir == "" {
		return
	}
	var files []string
	files, err = CoreFile.GetFileList(t.InstallDir, []string{"sql"}, true)
	if err != nil {
		err = nil
		return
	}
	for _, v := range files {
		var sqlByte []byte
		sqlByte, err = CoreFile.LoadFile(v)
		if err != nil {
			err = errors.New("install sql load install file, " + err.Error())
			return
		}
		_, err = t.DB.Exec(string(sqlByte))
		if err != nil {
			err = errors.New(fmt.Sprint("install sql exec sql, sql file: ", v, ", sql data: ", string(sqlByte), ", err: ", err))
			return
		}
	}
	return
}
