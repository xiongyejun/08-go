package main

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"hash"
	"pkgMyPkg/aesECB"
	"strconv"
)

type rc4 struct {
	b []byte

	sha           hash.Hash
	encryptionKey []byte
	keySize       int
	EncryptionVerifier
}

func (me *rc4) CheckPassword(password string) (err error) {
	fmt.Println("password=", password)

	// 生成加密密钥
	if err = me.getEncryptionKey(string2Unicode(password)); err != nil {
		return
	}

	// 验证密码
	return me.passwordVerifier()
}

// 2.3.4.5
func (me *rc4) initData() (err error) {
	var startIndex int = 4 // versio大小是4，前面已经判断过来
	//	var encryptionHeaderFlags uint32
	//	if encryptionHeaderFlags, err = byteToUint32(me.b[startIndex : startIndex+4]); err != nil {
	//		return
	//	}
	startIndex += 4

	var encryptionHeaderSize uint32
	if encryptionHeaderSize, err = byteToUint32(me.b[startIndex : startIndex+4]); err != nil {
		return
	}
	startIndex += 4
	// Flags	The fCryptoAPI and fAES bits MUST be set. The fDocProps bit MUST be 0.
	// 这个可以跳过
	startIndex += 4
	//SizeExtra	This value MUST be 0x00000000.
	startIndex += 4
	//AlgID	This value MUST be 0x0000660E (AES-128), 0x0000660F (AES-192), or 0x00006610 (AES-256).
	var algID uint32
	if algID, err = byteToUint32(me.b[startIndex : startIndex+4]); err != nil {
		return
	}
	if algID != 0x0000660E && algID != 0x0000660E && algID != 0x0000660E {
		return errors.New("AlgID	This value MUST be 0x0000660E (AES-128), 0x0000660F (AES-192), or 0x00006610 (AES-256).")
	}
	startIndex += 4
	//AlgIDHash	This value MUST be 0x00008004 (SHA-1).
	var algIDHash uint32
	if algIDHash, err = byteToUint32(me.b[startIndex : startIndex+4]); err != nil {
		return
	}
	if algIDHash != 0x00008004 {
		return errors.New("AlgIDHash=" + strconv.Itoa(int(algIDHash)) + ", This value MUST be 0x00008004 (SHA-1).")
	}
	me.sha = sha1.New()
	startIndex += 4
	//KeySize	This value MUST be 0x00000080 (AES-128), 0x000000C0 (AES-192), or 0x00000100 (AES-256).
	var keySize uint32
	if keySize, err = byteToUint32(me.b[startIndex : startIndex+4]); err != nil {
		return
	}
	me.keySize = int(keySize)
	startIndex += 4
	//ProviderType	This value SHOULD be 0x00000018 (AES).
	var providerType uint32
	if providerType, err = byteToUint32(me.b[startIndex : startIndex+4]); err != nil {
		return
	}
	if providerType != 0x00000018 {
		return errors.New("ProviderType=" + strconv.Itoa(int(providerType)) + ", This value SHOULD be 0x00000018 (AES).")
	}
	startIndex += 4
	//Reserved1	This value is undefined and MUST be ignored.
	startIndex += 4
	//Reserved2	This value MUST be 0x00000000 and MUST be ignored.
	startIndex += 4
	//CSPName	This value SHOULD<11> be set to either "Microsoft Enhanced RSA and AES Cryptographic Provider" or "Microsoft Enhanced RSA and AES Cryptographic Provider (Prototype)" as a null-terminated Unicode string.
	encryptionHeaderSize -= 4 * 8 // 剩余的长度是cspName
	//	cspName := me.b[startIndex : startIndex+int(encryptionHeaderSize)]
	startIndex += int(encryptionHeaderSize)

	// fmt.Printf("encryptionHeaderFlags=%d, algID=%d, keySize=%d, cspName=% x\r\n", encryptionHeaderFlags, algID, keySize, cspName)

	// 接下来是EncryptionVerifier
	if startIndex, err = me.EncryptionVerifier.getEncryptionVerifier(me.b, startIndex, providerType); err != nil {
		return
	}

	return nil
}

// 2.3.4.7
func (me *rc4) getEncryptionKey(pwd []byte) (err error) {
	if me.encryptionKey, err = H(me.sha, me.Salt, pwd); err != nil {
		return
	}
	var i uint = 0
	for ; i < 50000; i++ {
		// Hn = H(iterator + Hn-1)
		// iterator is an unsigned 32-bit value that is initially set to 0x00000000 and then incremented
		if me.encryptionKey, err = H(me.sha, uintToByte(i), me.encryptionKey); err != nil {
			return
		}
	}
	// Hfinal = H(Hn + block)
	if me.encryptionKey, err = H(me.sha, me.encryptionKey, []byte{0, 0, 0, 0}); err != nil {
		return
	}

	me.deriveKey()
	// Should handle the case of longer key lengths as shown in 2.3.4.9
	// Grab the key length bytes of the final hash as the encrypytion key
	var keySize int = me.keySize / 8
	if len(me.encryptionKey) > keySize {
		me.encryptionKey = me.encryptionKey[:keySize]
	} else {
		for i := keySize; i < len(me.encryptionKey); i++ {
			me.encryptionKey = append(me.encryptionKey, 0)
		}
	}

	return nil
}

// 2.3.4.7
func (me *rc4) deriveKey() (err error) {
	n := len(me.encryptionKey)
	if n > 64 {
		n = 64
	}

	key := make([]byte, 64)
	for i := 0; i < n; i++ {
		key[i] = me.encryptionKey[i] ^ 0x36
	}
	for i := n; i < 64; i++ {
		key[i] = 0x36
	}

	me.sha.Reset()
	if _, err = me.sha.Write(key); err != nil {
		return err
	}
	x1 := me.sha.Sum(nil)

	if int(me.VerifierHashSize) > me.keySize/8 {
		me.encryptionKey = x1
		return
	}

	n = len(me.encryptionKey)
	if n > 64 {
		n = 64
	}
	for i := range me.encryptionKey {
		key[i] = me.encryptionKey[i] ^ 0x5C
	}
	for i := len(me.encryptionKey); i < 64; i++ {
		key[i] = 0x5C
	}

	me.sha.Reset()
	if _, err = me.sha.Write(key); err != nil {
		return err
	}
	x2 := me.sha.Sum(nil)

	me.encryptionKey = append(x1, x2...)

	return nil
}

// 密码验证
func (me *rc4) passwordVerifier() (err error) {
	ecb := aesECB.NewAesTool(me.encryptionKey, me.keySize/8)

	// Decrypt the encrypted verifier 解密加密验证器
	var decryptedVerifier []byte
	if decryptedVerifier, err = ecb.Decrypt(me.EncryptedVerifier); err != nil {
		return errors.New("decryptedVerifier:" + err.Error())
	}
	decryptedVerifier = decryptedVerifier[:16]

	var decryptedVerifierHash []byte
	if decryptedVerifierHash, err = ecb.Decrypt(me.EncryptedVerifierHash); err != nil {
		return errors.New("decryptedVerifierHash:" + err.Error())
	}
	decryptedVerifierHash = decryptedVerifierHash[:20]
	// Hash the decrypted verifier (2.3.4.9)
	me.sha.Reset()
	if _, err = me.sha.Write(decryptedVerifier); err != nil {
		return errors.New("sha.Write:" + err.Error())
	}
	checkHash := me.sha.Sum(nil)

	if bytes.Compare(checkHash, decryptedVerifierHash) != 0 {
		return errors.New("密码错误。")
	}

	return nil
}
