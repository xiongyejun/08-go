package main

import (
	"fmt"
	"os"
)

func main() {
	var s, sep string
	//	for i := 1; i < len(os.Args); i++ {
	for _, arg := range os.Args {
		//		s += sep + os.Args[i]
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
	fmt.Println("hello")
}

/*
package main

import "fmt"

type testInt func(int) bool

func main() {

	slice := []int{1, 2, 3, 4, 5, 7}
	fmt.Println("slice = ", slice)

	odd := filter(slice, isOdd)
	fmt.Println("odd elements in slice are: ", odd)

	event := filter(slice, isEven)
	fmt.Println("event elements in slice are :", event)

	fmt.Println(isOdd(2))

	fmt.Println("hello world")

}

func isOdd(a int) bool {
	if a%2 == 0 {
		return false
	}
	return true
}

func isEven(a int) bool {
	if a%2 == 0 {
		return true
	}
	return false
}

func filter(slice []int, f testInt) []int {
	var result []int
	for _, value := range slice {
		if f(value) {
			result = append(result, value)
		}
	}
	return result
}
*/
