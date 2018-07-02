package offCrypto

import (
	"bytes"
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha1"
	"errors"
	"hash"
	"strconv"
)

type rc4CryptoAPI struct {
	ecma376RC4
}

func (me *rc4CryptoAPI) CheckPassword(passwordUnicodeByte []byte) (err error) {
	// 生成加密密钥
	if err = me.getEncryptionKey(passwordUnicodeByte); err != nil {
		return
	}

	// 验证密码
	return me.passwordVerifier()
}

// 2.3.5.1
func (me *rc4CryptoAPI) initData() (err error) {
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
	//AlgID	MUST be 0x00006801 (RC4 encryption).
	var algID uint32
	if algID, err = byteToUint32(me.b[startIndex : startIndex+4]); err != nil {
		return
	}
	if algID != 0x00006801 {
		return errors.New("AlgID MUST be 0x00006801 (RC4 encryption).")
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
	//KeySize	MUST be greater than or equal to 0x00000028 bits and less than or equal to 0x00000080 bits, in increments 增量of 8 bits. If set to 0x00000000, it MUST be interpreted 解释as 0x00000028 bits.
	var keySize uint32
	if keySize, err = byteToUint32(me.b[startIndex : startIndex+4]); err != nil {
		return
	}
	if keySize == 0 {
		keySize = 0x00000028
	}
	if keySize < 0x00000028 || keySize > 0x00000080 {
		return errors.New("KeySize MUST be greater than or equal to 0x00000028 bits and less than or equal to 0x00000080 bits.")
	}
	if (keySize-0x00000028)%8 != 0 {
		return errors.New("KeySize MUST in increments 增量of 8 bits.")
	}
	me.keySize = int(keySize)
	startIndex += 4
	//ProviderType	MUST be 0x00000001.
	var providerType uint32
	if providerType, err = byteToUint32(me.b[startIndex : startIndex+4]); err != nil {
		return
	}
	if providerType != 0x00000001 {
		return errors.New("ProviderType=" + strconv.Itoa(int(providerType)) + ", MUST be 0x00000001.")
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

// 2.3.5.2
func (me *rc4CryptoAPI) getEncryptionKey(pwd []byte) (err error) {
	if me.encryptionKey, err = H(me.sha, me.Salt, pwd); err != nil {
		return
	}
	// Hfinal = H(Hn + block)
	if me.encryptionKey, err = H(me.sha, me.encryptionKey, []byte{0, 0, 0, 0}); err != nil {
		return
	}

	var keySize int = me.keySize / 8
	if len(me.encryptionKey) > keySize {
		me.encryptionKey = me.encryptionKey[:keySize]
	} else {
		for i := me.keySize; i < len(me.encryptionKey); i++ {
			me.encryptionKey = append(me.encryptionKey, 0)
		}
	}

	return nil
}

// 密码验证
func (me *rc4CryptoAPI) passwordVerifier() (err error) {
	return passwordVerifier(me.encryptionKey, me.EncryptedVerifier, me.EncryptedVerifierHash, me.sha)
}

func rc4EncryptDecrypt(src []byte, key []byte) ([]byte, error) {
	if r, err := rc4.NewCipher(key); err != nil {
		return nil, err
	} else {
		des := make([]byte, len(src))
		r.XORKeyStream(des, src)
		return des, nil
	}
}

type officeBinRC4 struct {
	ecma376RC4
}

func (me *officeBinRC4) CheckPassword(passwordUnicodeByte []byte) (err error) {
	// 生成加密密钥
	if err = me.getEncryptionKey(passwordUnicodeByte); err != nil {
		return
	}

	// 验证密码
	return me.passwordVerifier()
}

// 2.3.6.1
func (me *officeBinRC4) initData() (err error) {
	var startIndex int = 4 // versio大小是4，前面已经判断过了
	// Salt (16 bytes):
	me.Salt = me.b[startIndex : startIndex+16]
	startIndex += 16
	// EncryptedVerifier (16 bytes):
	me.EncryptedVerifier = me.b[startIndex : startIndex+16]
	startIndex += 16
	// EncryptedVerifierHash (16 bytes):
	me.EncryptedVerifierHash = me.b[startIndex : startIndex+16]
	startIndex += 16

	me.sha = md5.New()

	return nil
}

// 2.3.6.2
func (me *officeBinRC4) getEncryptionKey(pwd []byte) (err error) {
	if me.encryptionKey, err = H(me.sha, nil, pwd); err != nil {
		return
	}
	tmp := append(me.encryptionKey[:5], me.Salt...)
	me.encryptionKey = nil
	for i := 0; i < 16; i++ {
		me.encryptionKey = append(me.encryptionKey, tmp...)
	}
	if me.encryptionKey, err = H(me.sha, nil, me.encryptionKey); err != nil {
		return
	}

	me.encryptionKey = me.encryptionKey[:5]
	// Hfinal equals H(TruncatedHash + block).
	if me.encryptionKey, err = H(me.sha, me.encryptionKey, []byte{0, 0, 0, 0}); err != nil {
		return
	}

	return nil
}

func (me *officeBinRC4) passwordVerifier() (err error) {
	return passwordVerifier(me.encryptionKey, me.EncryptedVerifier, me.EncryptedVerifierHash, me.sha)
}

func passwordVerifier(encryptionKey, EncryptedVerifier, EncryptedVerifierHash []byte, sha hash.Hash) (err error) {
	// Decrypt the encrypted verifier 解密加密验证器
	var r *rc4.Cipher
	if r, err = rc4.NewCipher(encryptionKey); err != nil {
		return err
	}

	// Decrypt the encrypted verifier 解密加密验证器
	var decryptedVerifier []byte = make([]byte, len(EncryptedVerifier))
	r.XORKeyStream(decryptedVerifier, EncryptedVerifier)

	var decryptedVerifierHash []byte = make([]byte, len(EncryptedVerifierHash))
	r.XORKeyStream(decryptedVerifierHash, EncryptedVerifierHash)

	var checkHash []byte
	// 4.	Calculate the SHA-1 hash value of the Verifier value calculated in step 2.
	if checkHash, err = H(sha, decryptedVerifier, nil); err != nil {
		return
	}

	if bytes.Compare(checkHash, decryptedVerifierHash) != 0 {
		return errors.New("密码错误。")
	}

	return nil
}
