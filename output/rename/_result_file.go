package GOPACKAGE

import "fmt"

var (
	_result_A = map[int]string{
		1: "hello",
	}
)

const (
	_result_X = 123
)

func _result_add() {
}

type _result_Struct struct {
	Val int64
}

func (s _result_Struct) hello() {
	_result_add()
	fmt.Println(_result_X, _result_A)
}
