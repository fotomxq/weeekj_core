package CoreHttp2

import (
	"encoding/json"
	"errors"
	"io"
)

// ResultBody 获取原生数据
func (t *Core) ResultBody() io.ReadCloser {
	//如果获取失败，反馈错误
	if t.Err != nil {
		return nil
	}
	//退出后关闭连接
	defer func() {
		_ = t.Response.Body.Close()
	}()
	//反馈数据
	return t.Response.Body
}

// ResultBytes 获取byte数据
func (t *Core) ResultBytes() (result []byte, err error) {
	if t.Err != nil {
		return []byte{}, t.Err
	}
	if t.Response.Body == nil {
		return []byte{}, errors.New("no data")
	}
	result, t.Err = io.ReadAll(t.Response.Body)
	if t.Err != nil {
		return []byte{}, t.Err
	}
	return
}

// ResultString 将结果转为字符串
func (t *Core) ResultString() (result string, err error) {
	var d []byte
	d, err = t.ResultBytes()
	if err != nil {
		return
	}
	result = string(d)
	return
}

// ResultJSON 将结果转为json
func (t *Core) ResultJSON(result any) (err error) {
	var d []byte
	d, t.Err = t.ResultBytes()
	if t.Err != nil {
		return t.Err
	}
	t.Err = json.Unmarshal(d, result)
	if t.Err != nil {
		return t.Err
	}
	return t.Err
}
