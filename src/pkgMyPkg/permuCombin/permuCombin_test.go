package permuCombin

import (
	"fmt"
	"testing"
)

func TestPermu(t *testing.T) {
	src := []byte("012345678")
	ch := make(chan []byte, 10)

	var selectCount uint = 5
	go Permu(src, selectCount, ch)

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
