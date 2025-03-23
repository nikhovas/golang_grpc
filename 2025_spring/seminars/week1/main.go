package main

import (
	"fmt"
	"unsafe"
)

// type A struct {
// 	F int `json:"f" asdf:"asdf" jklasd:5`
// 	s int `json:"s"`
// }

// _= json.Unmarshal(`{“f”: 5, “s”: 6}`, &s2)

type A struct{}

type B struct {
	A
	B int
}

func main() {
	var b B

	fmt.Println(unsafe.Sizeof(b))
}
