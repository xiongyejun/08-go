package offCrypto

import (
	"io/ioutil"
	"testing"
)

func TestRc4(t *testing.T) {
	var b []byte
	var err error

	b, err = ioutil.ReadFile(`C:\Users\Administrator\Desktop\03encryptionRC4API\Workbook`)
	if err != nil {
		t.Error(err)
	}

	var ie IEncryptedType
	ie, err = newOffBin(b)
	if err != nil {
		t.Error(err)
	}

	ie.initData()
	err = ie.CheckPassword([]byte{49, 0})

	if err != nil {
		println(err.Error())
	} else {
		println("ok")
	}
}
