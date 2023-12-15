package BaseWeixinWXXClientCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	BaseWeixinWXXClientECB "github.com/fotomxq/weeekj_core/v5/base/weixin/wxx/client/ecb"
	"io"
	"sort"
	"strings"
)

// SignByMD5 多参数通过MD5签名
func SignByMD5(data map[string]string, key string) (string, error) {

	var query []string
	for k, v := range data {
		query = append(query, k+"="+v)
	}

	sort.Strings(query)
	query = append(query, "key="+key)
	str := strings.Join(query, "&")

	str, err := MD5(str)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(str), nil
}

// MD5 加密
func MD5(str string) (string, error) {
	hs := md5.New()
	if _, err := hs.Write([]byte(str)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hs.Sum(nil)), nil
}

// PKCS5UnPadding 反补
// Golang AES没有64位的块, 如果采用PKCS5, 那么实质上就是采用PKCS7
func PKCS5UnPadding(plaintext []byte) ([]byte, error) {
	ln := len(plaintext)

	// 去掉最后一个字节 unPadding 次
	unPadding := int(plaintext[ln-1])

	if unPadding > ln {
		return []byte{}, errors.New("数据不正确")
	}

	return plaintext[:(ln - unPadding)], nil
}

// PKCS5Padding 补位
// Golang AES没有64位的块, 如果采用PKCS5, 那么实质上就是采用PKCS7
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// Validate 对数据包进行签名校验，确保数据的完整性。
//
// @rawData 不包括敏感信息的原始数据字符串，用于计算签名。
// @signature 使用 sha1( rawData + sessionkey ) 得到字符串，用于校验用户信息
// @ssk 微信 session_key
func Validate(rawData, ssk, signature string) bool {
	r := sha1.Sum([]byte(rawData + ssk))

	return signature == hex.EncodeToString(r[:])
}

// CBCDecrypt CBC解密数据
//
// @ssk 通过 Login 向微信服务端请求得到的 session_key
// @data 小程序通过 api 得到的加密数据(encryptedData)
// @iv 小程序通过 api 得到的初始向量(iv)
func CBCDecrypt(ssk, data, iv string) (bts []byte, err error) {
	var key []byte
	key, err = base64.StdEncoding.DecodeString(ssk)
	if err != nil {
		err = errors.New(fmt.Sprint("get ssk, ", err))
		return
	}
	var cipherText []byte
	cipherText, err = base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = errors.New(fmt.Sprint("get data, ", err))
		return
	}
	var rawIV []byte
	rawIV, err = base64.StdEncoding.DecodeString(iv)
	if err != nil {
		err = errors.New(fmt.Sprint("get iv, ", err))
		return
	}
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		err = errors.New(fmt.Sprint("get block, ", err))
		return
	}
	size := aes.BlockSize
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipherText.
	if len(cipherText) < size {
		err = errors.New("cipher too short")
		return
	}
	// CBC mode always works in whole blocks.
	if len(cipherText)%size != 0 {
		err = errors.New("cipher is not a multiple of the block size")
		return
	}
	mode := cipher.NewCBCDecrypter(block, rawIV[:size])
	plaintext := make([]byte, len(cipherText))
	mode.CryptBlocks(plaintext, cipherText)
	return PKCS5UnPadding(plaintext)
}

// CBCEncrypt CBC加密数据
func CBCEncrypt(key, data string) (ciphertext []byte, err error) {
	dk, err := hex.DecodeString(key)
	if err != nil {
		return
	}

	plaintext := []byte(data)

	if len(plaintext)%aes.BlockSize != 0 {
		err = errors.New("plaintext is not a multiple of the block size")
		return
	}

	block, err := aes.NewCipher(dk)
	if err != nil {
		return
	}

	ciphertext = make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return PKCS5Padding(ciphertext, block.BlockSize()), nil
}

// AesECBDecrypt CBC解密数据
//
// @ciphertext 加密数据
// @key 商户支付密钥
func AesECBDecrypt(ciphertext, key []byte) (plaintext []byte, err error) {

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("cipher too short")
	}
	// ECB mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("cipher is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	BaseWeixinWXXClientECB.NewDecrypter(block).CryptBlocks(ciphertext, ciphertext)

	return PKCS5UnPadding(ciphertext)
}
