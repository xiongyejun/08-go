package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"hash"
)

type EncryptedType int

const (
	RC4 EncryptedType = iota
	AGILE
)

type IEncryptedType interface {
	initData() (err error) // 从encryptionInfo读取需要的数据
	CheckPassword(psw string) (err error)
}

// 2.3.3
type EncryptionVerifier struct {
	SaltSize              uint32 // It MUST be 0x00000010
	Salt                  []byte
	EncryptedVerifier     []byte
	VerifierHashSize      uint32
	EncryptedVerifierHash []byte //  (variable) RC4-20个字节。AES-32个字节
}

// 2.3.3
func (me *EncryptionVerifier) getEncryptionVerifier(src []byte, index int, providerType uint32) (endIndex int, err error) {
	if me.SaltSize, err = byteToUint32(src[index : index+4]); err != nil {
		return
	}
	index += 4

	me.Salt = src[index : index+int(me.SaltSize)]
	index += int(me.SaltSize)

	me.EncryptedVerifier = src[index : index+int(me.SaltSize)]
	index += int(me.SaltSize)

	if me.VerifierHashSize, err = byteToUint32(src[index : index+4]); err != nil {
		return
	}
	index += 4

	//		enum ProviderType
	//		{
	//			Any = 0x00000000,
	//			RC4 = 0x00000001,
	//			AES = 0x00000018
	//		}
	var encryptedVerifierHashLen int
	if providerType == 0x00000018 {
		encryptedVerifierHashLen = 32
	} else if providerType == 0x00000001 {
		encryptedVerifierHashLen = 20
	}

	me.EncryptedVerifierHash = src[index : index+encryptedVerifierHashLen]
	index += encryptedVerifierHashLen

	return index, nil
}

// 密码验证
func (me *EncryptionVerifier) evPasswordVerifier(encryptionKey []byte, sha hash.Hash) (err error) {
	fmt.Println("len=", len(encryptionKey))
	// Decrypt the encrypted verifier 解密加密验证器
	var decryptedVerifier []byte
	if decryptedVerifier, err = aesDecrypt(me.EncryptedVerifier, encryptionKey); err != nil {
		return errors.New("decryptedVerifier:" + err.Error())
	}
	decryptedVerifier = decryptedVerifier[:16]

	var decryptedVerifierHash []byte
	if decryptedVerifierHash, err = aesDecrypt(me.EncryptedVerifierHash, encryptionKey); err != nil {
		return errors.New("decryptedVerifierHash:" + err.Error())
	}
	// Hash the decrypted verifier (2.3.4.9)
	if _, err = sha.Write(decryptedVerifier); err != nil {
		return errors.New("sha.Write:" + err.Error())
	}
	checkHash := sha.Sum(nil)

	if bytes.Compare(checkHash, decryptedVerifierHash) != 0 {
		return errors.New("密码不正确。")
	}

	fmt.Println("test")
	return nil
}

func aesDecrypt(crypted, key []byte) (b []byte, err error) {
	if block, err := aes.NewCipher(key); err != nil {
		return nil, err
	} else {
		blockMode := cipher.NewCBCDecrypter(block, key)
		src := make([]byte, len(crypted))
		blockMode.CryptBlocks(src, crypted)
		src = pkcs5UnPadding(src)
		return src, nil
	}
}
func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func H(sha hash.Hash, b1, b2 []byte) (b []byte, err error) {
	sha.Reset()
	b = append(b1, b2...)
	if _, err := sha.Write(b); err != nil {
		return nil, err
	}
	return sha.Sum(nil), nil
}
