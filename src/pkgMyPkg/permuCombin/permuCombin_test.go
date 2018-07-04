package permuCombin

import (
	"fmt"
	"testing"
)

func TestPermu(t *testing.T) {
	src := make([]string, 0)
	src = append(src, "zhang")
	src = append(src, "jia")
	src = append(src, "0")
	src = append(src, "1")
	src = append(src, "2")

	ch := make(chan []byte, 10)

	var selectCount uint = 6
	go PermuString(src, selectCount, ch)

	for i := 0; i < getCount(len(src), int(selectCount)); i++ {
		fmt.Printf("%d\t%s\r\n", i, <-ch)
	}
}

func getCount(m, n int) (icount int) {
	icount = 1
	for i := 0; i < n; i++ {
		icount *= m
	}
	return
}
