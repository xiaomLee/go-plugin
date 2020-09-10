package main

import "C"
import "fmt"

//export  SayBye
func SayBye(s *C.char) {
	fmt.Println(C.GoString(s))
}
