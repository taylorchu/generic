package test

type Type int
type Type2 int

func (Type Type) Type(Type) {}

func add(_ Type, _ Type2)  {}
func add2(_ Type2, _ Type) {}

type Struct struct {
	Val  Type
	Val2 Type2
}
