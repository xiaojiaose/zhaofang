package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

func DecryptWXData(sessionKey, encryptedData, iv string, v interface{}) error {
	// Base64解码
	key, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return err
	}

	// AES-128-CBC解密
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)

	// 去除填充
	cipherText, err = PKCS7UnPadding(cipherText)
	if err != nil {
		return err
	}

	// 解析JSON
	return json.Unmarshal(cipherText, v)
}

func PKCS7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid data length")
	}
	unPadding := int(data[length-1])
	if unPadding > length {
		return nil, errors.New("invalid padding size")
	}
	return data[:(length - unPadding)], nil
}
