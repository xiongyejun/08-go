package main

import (
	"fmt"
	"io"
	//	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	url := os.Args[1]

	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	//	b, err := ioutil.ReadAll(resp.Body)
	_, err1 := io.Copy(os.Stdout, resp.Body)
	resp.Body.Close()
	fmt.Println(resp.Status)

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}
	//	fmt.Printf("%s", b)
}
