package fileHeader

import (
	"fmt"
	"testing"
)

func Test_fh(t *testing.T) {
	for i := range fh {
		fmt.Println(i, fh[i])
	}
	what()

	fmt.Println(IsZIP([]byte{0x50, 0x4B, 0x03}))
}

func what() {
	if ext, err := GetExt([]byte{0x17, 0x49}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ext)
	}
}
