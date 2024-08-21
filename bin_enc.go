package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	i := int64(1234567890)
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, i); err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Printf("%d: %x\n", i, buf.Bytes())

	buf.Reset()
	f := 1234567.89
	if err := binary.Write(&buf, binary.BigEndian, f); err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Printf("%.2f: %x\n", f, buf.Bytes())

	s := "I ♡ Go"
	fmt.Printf("%s    : %x\n", s, []byte(s))
}
