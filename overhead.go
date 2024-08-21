package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func add(a, b int) int {
	return a + b
}

const ADD = "add"

func main() {
	// Regular call
	val := add(27, 15)
	fmt.Println(val)

	// RPC
	// client
	request := map[string]any{
		"func": ADD,
		"args": []int{27, 15},
	}

	var network bytes.Buffer

	// client
	cenc := json.NewEncoder(&network)
	if err := cenc.Encode(request); err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	// server
	var req struct {
		Func string
		Args []int
	}
	dec := json.NewDecoder(&network)
	if err := dec.Decode(&req); err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	var out int
	switch req.Func {
	case ADD:
		out = add(req.Args[0], req.Args[1])
	default:
		fmt.Println("ERROR: unknown function - ", req.Func)
		return
	}

	senc := json.NewEncoder(&network)
	if err := senc.Encode(out); err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	// client
	var result int
	if err := json.NewDecoder(&network).Decode(&result); err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Println("result:", result)
}
