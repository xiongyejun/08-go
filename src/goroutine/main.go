package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs  slapsed\n", time.Since(start).Seconds())

}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)

	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body) //ioutil.Discard：可以把这个变量看作一个垃圾桶，可以向里面写一些不需要的数据
	resp.Body.Close()

	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

//func say(s string) {
//	for i := 0; i < 5; i++ {
//		runtime.Gosched()
//		fmt.Println(s)
//	}
//}

//func sum(a []int, c chan int) {
//	total := 0
//	for _, v := range a {
//		total += v
//	}
//	c <- total
//}

//func main() {
//	/*go say("world")
//	say("hello")
//	*/
//	a := []int {7,2,8,-9,4,0}

//	c := make(chan int)
//	go sum(a[:len(a)/2], c)
//	go sum(a[len(a)/2:], c)
//	x, y := <-c, <-c

//	fmt.Println(x, y, x + y)
//}
