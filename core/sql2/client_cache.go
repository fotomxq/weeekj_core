package CoreSQL2

import (
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

func (t *ClientCtx) getCacheMark() string {
	return "sql:core:" + CoreFilter.GetSha1Str(t.client.TableName) + "." + CoreFilter.GetSha1Str(fmt.Sprint(t.query, t.appendArgs))
}

func (t *ClientCtx) getCacheData(data interface{}) error {
	err := t.client.cacheObj.GetStruct(t.getCacheMark(), data)
	if err != nil {
		return err
	}
	return nil
}

func (t *ClientCtx) setCacheData(data interface{}) {
	t.client.cacheObj.SetStruct(t.getCacheMark(), data, CoreCache.CacheTime1Hour)
}
