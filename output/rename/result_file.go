package GOPACKAGE

import "fmt"

var (
	result_A = map[int]string{
		1: "hello",
	}
)

const (
	result_X = 123
)

func result_add() {
}

type result_Struct struct {
	Val int64
}

func (s result_Struct) hello() {
	result_add()
	fmt.Println(result_X, result_A)
}
