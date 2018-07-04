package offCrypto

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestDecrypt(t *testing.T) {
	var b []byte
	var err error

	b, err = ioutil.ReadFile(`C:\Users\Administrator\Desktop\ecma376rc4\EncryptionInfo`)
	if err != nil {
		t.Error(err)
	}

	var ie IEncryptedType
	ie, err = newECMA376(b)
	if err != nil {
		t.Error(err)
	}

	ie.initData()
	err = ie.CheckPassword(asc2Unicode([]byte("green")))

	if err != nil {
		println(err.Error())
	} else {
		println("ok")
	}

	var bEncryptedPackage []byte
	bEncryptedPackage, err = ioutil.ReadFile(`C:\Users\Administrator\Desktop\ecma376rc4\EncryptedPackage`)
	if err != nil {
		t.Error(err)
	}
	b, err = ie.Decrypt(bEncryptedPackage)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("len b = %d\r\n", len(b))
	err = ioutil.WriteFile(`C:\Users\Administrator\Desktop\ecma376rc4\Decrypt.xlsx`, b, 0666)
	if err != nil {
		t.Error(err)
	}
}

func asc2Unicode(b []byte) []byte {
	var bb []byte = make([]byte, len(b)*2)
	for i := range b {
		bb[i*2] = b[i]
		bb[i*2+1] = 0
	}
	return bb
}
