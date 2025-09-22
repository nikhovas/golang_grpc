package main

import (
	"encoding/json"
	"fmt"
)

type A struct {
	F int `json:"f" log_status:"ignore"`
	s int `json:"s"`
}

type UserID int

func ggg(follower UserID) {

}

func main() {
	a := A{F: 5, s: 6}
	byteVal, _ := json.Marshal(&a)
	fmt.Println(string(byteVal))
	var s2 A
	_ = json.Unmarshal([]byte(`{"f": 5, "s": 6}`), &s2)
	fmt.Println(s2.F) // 5
	fmt.Println(s2.s) // 0

	ggg(5)
}
