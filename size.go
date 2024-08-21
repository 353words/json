package main

import (
	"fmt"
	"unsafe"
)

func main() {
	i := uint8(123)
	s := "123"

	fmt.Println("int: %d\n", unsafe.Sizeof(i))
}
