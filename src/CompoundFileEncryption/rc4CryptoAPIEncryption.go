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

type rc4CryptoAPI struct {
	rc4
}

func (me *rc4CryptoAPI) CheckPassword(password string) (err error) {
	fmt.Println("password=", password)

	// 生成加密密钥
	if err = me.getEncryptionKey(string2Unicode(password)); err != nil {
		return
	}

	// 验证密码
	return me.passwordVerifier()
}

func (me *rc4CryptoAPI) getEncryptionKey(pwd []byte) (err error) {
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
func (me *rc4CryptoAPI) deriveKey() (err error) {
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
func (me *rc4CryptoAPI) passwordVerifier() (err error) {
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
