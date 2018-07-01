package main

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"hash"
)

// 密钥数据
type KeyData struct {
	XMLName         xml.Name `xml:"keyData"`
	SaltSize        uint     `xml:"saltSize,attr"`
	BlockSize       uint     `xml:"blockSize,attr"`
	KeyBits         uint     `xml:"keyBits,attr"`
	HashSize        uint     `xml:"hashSize,attr"`
	CipherAlgorithm string   `xml:"cipherAlgorithm,attr"`
	CipherChaining  string   `xml:"cipherChaining,attr"`
	HashAlgorithm   string   `xml:"hashAlgorithm,attr"`
	SaltValue       string   `xml:"saltValue,attr"`
}

// 数据完整性
type DataIntegrity struct {
	XMLName            xml.Name `xml:"dataIntegrity"`
	EncryptedHmacKey   string   `xml:"encryptedHmacKey,attr"`
	EncryptedHmacValue string   `xml:"encryptedHmacValue,attr"`
}

// 加密密钥
type EncryptedKey struct {
	XMLName                    xml.Name `xml:"encryptedKey"`
	SaltSize                   uint     `xml:"saltSize,attr"`
	BlockSize                  uint     `xml:"blockSize,attr"`
	KeyBits                    uint     `xml:"keyBits,attr"`
	HashSize                   uint     `xml:"hashSize,attr"`
	CipherAlgorithm            string   `xml:"cipherAlgorithm,attr"`
	CipherChaining             string   `xml:"cipherChaining,attr"`
	HashAlgorithm              string   `xml:"hashAlgorithm,attr"`
	SaltValue                  string   `xml:"saltValue,attr"`
	SpinCount                  uint     `xml:"spinCount,attr"`
	EncryptedVerifierHashInput string   `xml:"encryptedVerifierHashInput,attr"`
	EncryptedVerifierHashValue string   `xml:"encryptedVerifierHashValue,attr"`
	EncryptedKeyValue          string   `xml:"encryptedKeyValue,attr"`
}

// 密钥加密器
type KeyEncryptor struct {
	XMLName xml.Name      `xml:"keyEncryptor"`
	URI     string        `xml:"uri,attr"`
	EK      *EncryptedKey `xml:"encryptedKey"`
}
type KeyEncryptors struct {
	KES []*KeyEncryptor `xml:"keyEncryptor"`
}

// 2.3.4.10
type Encryption struct {
	XMLName xml.Name       `xml:"encryption"`
	KD      *KeyData       `xml:"keyData"`
	DI      *DataIntegrity `xml:"dataIntegrity"`
	KE      *KeyEncryptors `xml:"keyEncryptors"`
}

// 2.3.4.11	Encryption Key Generation (Agile Encryption)
type agile struct {
	b []byte
	E *Encryption `xml:"encryption,attr"`

	encryptionKey []byte

	EncryptionVerifier
	encryptedKeyValue []byte

	sha hash.Hash
}

func (me *agile) CheckPassword(password string) (err error) {
	// 生成加密密钥
	if err = me.getEncryptionKey(string2Unicode(password)); err != nil {
		return
	}

	// 验证密码
	return me.passwordVerifier()
}

// 从encryptionInfo读取需要的数据
func (me *agile) initData() (err error) {
	// Reserved (4 bytes): A value that MUST be 0x00000040.
	var tmp uint32
	if tmp, err = byteToUint32(me.b[4:8]); tmp != 0x00000040 || err != nil {
		return errors.New("Agile Encryption Reserved (4:8 bytes): A value that MUST be 0x00000040.")
	}
	// 从xml中读取所需要的数据
	if err = me.parseXml(); err != nil {
		return
	} else {
		if me.Salt, err = base64.StdEncoding.DecodeString(me.E.KE.KES[0].EK.SaltValue); err != nil {
			return
		}
		if me.EncryptedVerifier, err = base64.StdEncoding.DecodeString(me.E.KE.KES[0].EK.EncryptedVerifierHashInput); err != nil {
			return
		}
		if me.EncryptedVerifierHash, err = base64.StdEncoding.DecodeString(me.E.KE.KES[0].EK.EncryptedVerifierHashValue); err != nil {
			return
		}
		if me.encryptedKeyValue, err = base64.StdEncoding.DecodeString(me.E.KE.KES[0].EK.EncryptedKeyValue); err != nil {
			return
		}

		if me.E.KE.KES[0].EK.HashAlgorithm == "SHA512" {
			me.sha = sha512.New()
		} else if me.E.KE.KES[0].EK.HashAlgorithm == "SHA1" {
			me.sha = sha1.New()
		} else {
			return errors.New("未设置的HashAlgorithm：" + me.E.KE.KES[0].EK.HashAlgorithm)
		}

	}

	return nil
}

// 2.3.4.11	Encryption Key Generation加密密钥生成
func (me *agile) getEncryptionKey(pwd []byte) (err error) {
	if me.encryptionKey, err = H(me.sha, me.Salt, pwd); err != nil {
		return
	}
	var i uint = 0
	for ; i < me.E.KE.KES[0].EK.SpinCount; i++ {
		// Hn = H(iterator + Hn-1)
		// iterator is an unsigned 32-bit value that is initially set to 0x00000000 and then incremented
		if me.encryptionKey, err = H(me.sha, uintToByte(i), me.encryptionKey); err != nil {
			return
		}
	}

	return nil
}

// 2.3.4.13
func (me *agile) createVerifier() (encryptedValue []byte, err error) {
	hashInputBlockKey := []byte{0xfe, 0xa7, 0xd2, 0x76, 0x3b, 0x4b, 0x9e, 0x79}
	// 解密EncryptedVerifier
	var plaintextEncryptedVerifier []byte
	if plaintextEncryptedVerifier, err = me.cryptor(hashInputBlockKey, me.EncryptedVerifier, false); err != nil {
		return
	}

	// 获取加密的encryptedHashValue
	var tmp int = 0
	tmp1 := int(me.E.KE.KES[0].EK.HashSize)
	tmp2 := int(me.E.KE.KES[0].EK.BlockSize)
	tmp = ((tmp1 + tmp2 - 1) / tmp2) * tmp2

	var encryptedHashValue []byte
	if encryptedHashValue, err = H(me.sha, plaintextEncryptedVerifier, nil); err != nil {
		return
	}
	encryptedHashValue = appendByte(encryptedHashValue, tmp, 0)

	// 加密encryptedHashValue
	hashValueBlockKey := []byte{0xd7, 0xaa, 0x0f, 0x6d, 0x30, 0x61, 0x34, 0x4e}
	if encryptedValue, err = me.cryptor(hashValueBlockKey, encryptedHashValue, true); err != nil {
		return
	}

	return
}

// bEncryption	true	加密
// bEncryption	false	解密
func (me *agile) cryptor(blockKey, valueInput []byte, bEncryption bool) (bResult []byte, err error) {
	// 先用密码最后的hash值encryptionKey，与blockKey进行hash
	if bResult, err = H(me.sha, me.encryptionKey, blockKey); err != nil {
		return
	}
	// 得到密key
	secretKey := appendByte(bResult, int(me.E.KE.KES[0].EK.KeyBits/8), 0x36)
	// 得到密iv
	iv := appendByte(me.Salt, int(me.E.KE.KES[0].EK.BlockSize), 0x36)

	if bEncryption {
		// 加密
		if bResult, err = aesEncrypt(valueInput, secretKey, iv); err != nil {
			return
		}
	} else {
		// 解密
		if bResult, err = aesDecrypt(valueInput, secretKey, iv); err != nil {
			return
		}
	}

	return
}

// 解析xml
func (me *agile) parseXml() (err error) {
	me.E = new(Encryption)
	me.E.KD = new(KeyData)
	me.E.DI = new(DataIntegrity)
	me.E.KE = new(KeyEncryptors)

	if err = xml.Unmarshal(me.b[8:], me.E); err != nil {
		return
	}
	return nil
}

// 密码验证
func (me *agile) passwordVerifier() (err error) {
	var encryptedHashValueNew []byte
	if encryptedHashValueNew, err = me.createVerifier(); err != nil {
		return
	}

	if bytes.Compare(encryptedHashValueNew, me.EncryptedVerifierHash) != 0 {
		return errors.New("密码错误。")
	}

	//	secretKeyBlockKey := []byte{0x14, 0x6e, 0x0b, 0xe7, 0xab, 0xac, 0xd0, 0xd6}
	//	// 解密EncryptedVerifier
	//	var decryptedSecretKey []byte
	//	if decryptedSecretKey, err = me.cryptor(secretKeyBlockKey, me.encryptedKeyValue, false); err != nil {
	//		return
	//	}

	return nil
}
