package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
)

func desEncryptString(src string, key []byte) (string, error) {
	if b, err := desEncrypt([]byte(src), key); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(b), nil
	}
}

func desEncrypt(src, key []byte) ([]byte, error) {
	if block, err := des.NewCipher(key); err != nil {
		return nil, err
	} else {
		src = pkcs5Padding(src, block.BlockSize())
		blockMode := cipher.NewCBCEncrypter(block, key)
		crypted := make([]byte, len(src))
		blockMode.CryptBlocks(crypted, src)
		return crypted, nil
	}
}

func desDecryptString(crypted string, key []byte) (string, error) {
	var b []byte
	var err error
	if b, err = hex.DecodeString(crypted); err != nil {
		return "", nil
	} else {
		if b, err = desDecrypt(b, key); err != nil {
			return "", err
		} else {
			return string(b), nil
		}
	}
}

func desDecrypt(crypted, key []byte) ([]byte, error) {
	if block, err := des.NewCipher(key); err != nil {
		return nil, err
	} else {
		blockMode := cipher.NewCBCDecrypter(block, key)
		src := make([]byte, len(crypted))
		blockMode.CryptBlocks(src, crypted)
		src = pkcs5UnPadding(src)
		return src, nil
	}
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
