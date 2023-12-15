package BaseToken

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

// 生成一个不重复的key
func getKeyNoReplace(limitKeyLen int) (string, error) {
	var key string
	var err error
	tryStep := 1
	maxTry := 10
	for {
		key, err = CoreFilter.GetRandStr3(limitKeyLen)
		if err != nil {
			return "", errors.New("cannot get new rand data, " + err.Error())
		}
		//检查是否重复
		_, err = GetByKey(&ArgsGetByKey{
			Key: key,
		})
		if err == nil {
			tryStep += 1
			if tryStep > maxTry {
				return "", errors.New("try is too many, pls delete some key or set limit more")
			}
			continue
		}
		return key, nil
	}
}
