package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
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

	sha      hash.Hash
	blockKey []byte
	keySalt  []byte
	iv       []byte // Initialization vectors
}

func (me *agile) CheckPassword(password string) (err error) {
	b := []byte(password)
	var bb []byte = make([]byte, len(b)*2)
	for i := range b {
		bb[i*2] = b[i]
		bb[i*2+1] = 0
	}
	// 生成加密密钥
	if err = me.getEncryptionKey(bb); err != nil {
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
		//		fmt.Printf("%#v\r\n\r\n", me.E)
		//		fmt.Printf("%#v\r\n\r\n", me.E.KD)
		//		fmt.Printf("%#v\r\n\r\n", me.E.DI)
		//		fmt.Printf("%#v\r\n\r\n", me.E.KE.KES[0])
		//		fmt.Printf("%#v\r\n\r\n", me.E.KE.KES[0].EK)

		if me.Salt, err = base64.StdEncoding.DecodeString(me.E.KE.KES[0].EK.SaltValue); err != nil {
			return
		}

		fmt.Println(me.E.KE.KES[0].EK.EncryptedVerifierHashInput)
		if me.EncryptedVerifier, err = base64.StdEncoding.DecodeString(me.E.KE.KES[0].EK.EncryptedVerifierHashInput); err != nil {
			return
		}
		if me.EncryptedVerifierHash, err = base64.StdEncoding.DecodeString(me.E.KE.KES[0].EK.EncryptedVerifierHashValue); err != nil {
			return
		}

		if me.Salt, err = base64.StdEncoding.DecodeString(me.E.KE.KES[0].EK.SaltValue); err != nil {
			return
		}

		if me.E.KE.KES[0].EK.HashAlgorithm == "SHA512" {
			me.sha = sha512.New()
		}

	}

	return nil
}

// 2.3.4.11	Encryption Key Generation加密密钥生成
func (me *agile) getEncryptionKey(pwd []byte) (err error) {
	fmt.Printf("SaltValue=% x\r\n", me.Salt)
	if me.encryptionKey, err = H(me.sha, me.Salt, pwd); err != nil {
		return
	}
	fmt.Printf("HashedSaltAndPassword=% x\r\n", me.encryptionKey)

	var i uint = 0
	for ; i < me.E.KE.KES[0].EK.SpinCount; i++ {
		// Hn = H(iterator + Hn-1)
		// iterator is an unsigned 32-bit value that is initially set to 0x00000000 and then incremented
		if me.encryptionKey, err = H(me.sha, uintToByte(i), me.encryptionKey); err != nil {
			return
		}
	}
	fmt.Printf("FinalHashValue=% x\r\n", me.encryptionKey)
	// Hfinal = H(Hn + blockKey)
	me.blockKey = make([]byte, me.E.KE.KES[0].EK.BlockSize)
	if me.encryptionKey, err = H(me.sha, me.encryptionKey, me.blockKey); err != nil {
		return
	}
	fmt.Printf("BlockSize FinalHashValue=% x\r\n", me.encryptionKey)

	// padded by appending bytes with a value of 0x36
	me.encryptionKey = append36(me.encryptionKey, int(me.E.KE.KES[0].EK.KeyBits/8))

	//	fmt.Printf("% x\r\n", me.encryptionKey)
	//	fmt.Printf("%s\r\n", me.encryptionKey)
	return nil
}

// 2.3.4.12	Initialization Vector Generation初始化向量生成
func (me *agile) getIV() (err error) {
	if me.blockKey == nil {
		me.iv = me.keySalt
	} else {
		if me.iv, err = H(me.sha, me.keySalt, me.blockKey); err != nil {
			return
		}
		// less than the value of the blockSize attribute corresponding 相应的to the cipherAlgorithm attribute, pad 垫the array of bytes by appending 0x36 until the array is blockSize bytes.
		me.iv = append36(me.iv, int(me.E.KE.KES[0].EK.BlockSize))

	}
	return nil
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

func (me *agile) passwordVerifier() (err error) {
	return me.evPasswordVerifier(me.encryptionKey, me.sha)
}
