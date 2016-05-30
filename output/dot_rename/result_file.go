package GOPACKAGE

import "fmt"

var (
	resultA = map[int]string{
		1: "hello",
	}
)

const (
	resultX = 123
)

func resultAdd() {
}

type resultStruct struct {
	Val int64
}

func (s resultStruct) hello() {
	resultAdd()
	fmt.Println(resultX, resultA)
}
