package offCrypto

import (
	"encoding/binary"
	"errors"
)

func NewIEncrypted(b []byte, bECMA376 bool) (iEncryptedType IEncryptedType, err error) {
	if bECMA376 {
		return newECMA376(b)
	} else {
		return newOffBin(b)
	}
}

func newECMA376(b []byte) (iEncryptedType IEncryptedType, err error) {
	p := new(version)
	var startIndex int
	if startIndex, err = readVersion(p, b, startIndex); err != nil {
		return
	}

	if p.vMajor == 0x0004 &&
		p.vMinor == 0x0004 {
		// Agile敏捷 Encryption
		println("ECMA-376 Agile Encryption")
		agl := &agile{}
		agl.b = b
		if err = agl.initData(); err != nil {
			return nil, err
		}
		return agl, nil

	} else if (p.vMajor == 0x0002 ||
		p.vMajor == 0x0003 ||
		p.vMajor == 0x0004) &&
		p.vMinor == 0x0002 {
		// Standard Encryption
		println("ECMA-376 rc4 Encryption")
		r := &ecma376RC4{}
		r.b = b
		if err = r.initData(); err != nil {
			return nil, err
		}
		return r, nil

	} else if (p.vMajor == 0x0003 ||
		p.vMajor == 0x0004) &&
		p.vMinor == 0x0003 {
		// Extensible Encryption
		println("Extensible Encryption")
		return nil, errors.New("未实现的加密类型。")
	} else {
		return nil, errors.New("未知加密类型。")
	}

	return nil, nil
}

func newOffBin(b []byte) (iEncryptedType IEncryptedType, err error) {
	p := new(version)
	var startIndex int
	if startIndex, err = readVersion(p, b, startIndex); err != nil {
		return
	}

	if p.vMajor == 0x0001 &&
		p.vMinor == 0x0001 {
		println("OffBinary rc4 Encryption")
		r := &officeBinRC4{}
		r.b = b
		if err = r.initData(); err != nil {
			return nil, err
		}
		return r, nil

	} else if (p.vMajor == 0x0002 ||
		p.vMajor == 0x0003 ||
		p.vMajor == 0x0004) &&
		p.vMinor == 0x0002 {
		// Standard Encryption
		println("OffBinary rc4 CryptoAPI Encryption")
		r := &rc4CryptoAPI{}
		r.b = b
		if err = r.initData(); err != nil {
			return nil, err
		}
		return r, nil

	} else {
		return nil, errors.New("未知加密类型。")
	}

	return nil, nil
}

type version struct {
	vMajor uint16
	vMinor uint16
}

// 读取Version结构
func readVersion(p *version, b []byte, startIndex int) (endIndex int, err error) {
	if p.vMajor, err = byteToUint16(b[startIndex : startIndex+binary.Size(p.vMajor)]); err != nil {
		return
	}
	startIndex += binary.Size(p.vMajor)

	if p.vMinor, err = byteToUint16(b[startIndex : startIndex+binary.Size(p.vMinor)]); err != nil {
		return
	}
	startIndex += binary.Size(p.vMinor)

	return startIndex, nil
}
