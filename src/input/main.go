package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	//	input := bufio.NewScanner(os.Stdin)

	//	for input.Scan() {
	//		counts[input.Text()]++
	//	}

	//	for line, n := range counts {
	//		if n > 1 {
	//			fmt.Printf("%d\t%s", n, line)
	//		}
	//	}

	files := os.Args[1:]
	if len(files) == 0 {
		//		countLines(os.Stdin, counts)

	} else {
		for _, arg := range files {
			//			f, err := os.Open(arg)
			data, err := ioutil.ReadFile(arg) //一次读取
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			//			countLines(f, counts)
			//			f.Close()

			for _, line := range strings.Split(string(data), "\n") {
				counts[line]++
			}
		}

		for line, n := range counts {
			if n > 1 {
				fmt.Printf("%d\t %s\n", n, line)
			}
		}
	}
}

//逐行读取文件内容
func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
}
