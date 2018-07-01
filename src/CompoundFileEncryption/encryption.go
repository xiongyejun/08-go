package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"hash"
)

// 1.3.3加密方式
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
	if me.SaltSize != 0x10 {
		return 0, errors.New("me.SaltSize必须等于0x10.")
	}

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

func aesEncrypt(src, key, IV []byte) ([]byte, error) {
	if block, err := aes.NewCipher(key); err != nil {
		return nil, err
	} else {
		//		src = pkcs5Padding(src, block.BlockSize())
		blockMode := cipher.NewCBCEncrypter(block, IV)
		crypted := make([]byte, len(src))
		blockMode.CryptBlocks(crypted, src)
		return crypted, nil
	}
}
func aesDecrypt(crypted, key, IV []byte) (b []byte, err error) {
	if block, err := aes.NewCipher(key); err != nil {
		return nil, err
	} else {
		blockMode := cipher.NewCBCDecrypter(block, IV)

		src := make([]byte, len(crypted))
		blockMode.CryptBlocks(src, crypted)
		//		src = pkcs5UnPadding(src)
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
	//	不用 b = append(b1,b2)，防止b1地址的内容被修改
	b = make([]byte, len(b1)+len(b2))
	copy(b, b1)
	copy(b[len(b1):], b2)
	if _, err := sha.Write(b); err != nil {
		return nil, err
	}
	return sha.Sum(nil), nil
}
