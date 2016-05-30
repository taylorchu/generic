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

func result_add() {
}

type resultStruct struct {
	Val Data
}

func (s resultStruct) hello() {
	result_add()
	fmt.Println(resultX, resultA)
}
