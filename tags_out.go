package main

import (
	"encoding/json"
	"os"
)

type VM struct {
	ID     string `json:"id"`
	CPUs   int    `json:"cpus"`
	Memory int    `json:"memory"`
}

func tags() {
	vm := VM{
		ID:     "b70229443f8d489bbc733f13a9268f63",
		CPUs:   4,
		Memory: 32,
	}
	json.NewEncoder(os.Stdout).Encode(vm)
}

func main() {
	tags()

}
