package main

import (
	"encoding/json"
	"os"
)

type VM struct {
	ID     string
	CPUs   int
	Memory int
}

func s() {
	vm := VM{
		ID:     "b70229443f8d489bbc733f13a9268f63",
		CPUs:   4,
		Memory: 32,
	}

	json.NewEncoder(os.Stdout).Encode(vm)
}

func main() {
	vm := map[string]any{
		"id":     "b70229443f8d489bbc733f13a9268f63",
		"cpus":   4,
		"memory": 32,
	}

	json.NewEncoder(os.Stdout).Encode(vm)
	s()
}
