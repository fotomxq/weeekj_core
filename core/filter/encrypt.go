package CoreFilter

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

//加密摘要模块

// GetSha1Str 获取字符串SHA1摘要
// 该模块返回string类型
// param str string 要加密的字符串
// return string SHA1值，加密失败则返回空字符串
func GetSha1Str(str string) string {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(str))
	if err != nil {
		return ""
	}
	sha := hasher.Sum(nil)
	return hex.EncodeToString(sha)
}

// GetSha1Str2 获取字符串SHA1摘要
// 该模块返回string类型
// param str string 要加密的字符串
// return string SHA1值，加密失败则返回空字符串
func GetSha1Str2(str string) (string, error) {
	res, err := GetSha1([]byte(str))
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// GetSha1ByString 获取字符串的SHA1值
// param content string 要计算的字符串
// return string 计算出的SHA1值
// return error
func GetSha1ByString(content string) (string, error) {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(content))
	if err != nil {
		return "", err
	}
	sha := hasher.Sum(nil)
	return hex.EncodeToString(sha), nil
}

// GetSha1 获取字符串SHA1摘要
// 该模块返回[]byte类型
// param str []byte 要加密的字符串
// return []byte SHA1值，加密失败则返回空字符串
func GetSha1(str []byte) ([]byte, error) {
	if len(str) < 1 {
		return nil, errors.New("Encrypt get sha1 , str is empty.")
	}
	hasher := sha1.New()
	_, err := hasher.Write(str)
	if err != nil {
		return nil, err
	}
	sha := hasher.Sum(nil)
	dest := make([]byte, hex.EncodedLen(len(sha)))
	if sha == nil {
		return nil, errors.New("Encrypt get sha1 dest is nil.")
	}
	_ = hex.Encode(dest, sha)
	return dest, nil
}

// GetSha256Str 获取sha256字符串
// param str string 要计算的字符串
// return string 字符串
// return error 错误信息
func GetSha256Str(str string) (string, error) {
	res, err := GetSha256([]byte(str))
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// GetSha256 计算sha256
// param str []byte 要加密的字符串
// return []byte SHA256值，加密失败则返回空字符串
func GetSha256(str []byte) ([]byte, error) {
	if len(str) < 1 {
		return nil, errors.New("Encrypt get sha256 , str is empty.")
	}
	hasher := sha256.New()
	_, err := hasher.Write(str)
	if err != nil {
		return nil, err
	}
	sha := hasher.Sum(nil)
	dest := make([]byte, hex.EncodedLen(len(sha)))
	if sha == nil {
		return nil, errors.New("Encrypt get sha256 dest is nil.")
	}
	_ = hex.Encode(dest, sha)
	return dest, nil
}

// GetMd5StrByStr 获取MD5字符串
func GetMd5StrByStr(str string) string {
	//计算MD5
	md5Byte := md5.Sum([]byte(str))
	//转为字符串
	return hex.EncodeToString(md5Byte[:])
}
