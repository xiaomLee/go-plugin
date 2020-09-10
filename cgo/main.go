package main

//#include "hello.h"
import "C"

func main() {
	C.SayHello(C.CString("hello"))
	C.SaySomething(C.CString("my name is tom"))
	C.SayBye(C.CString("bye"))
}
