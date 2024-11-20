package main

import (
	"encoding/json"
	"os"
)

type VM struct {
	ID     string `json:"id,omitempty"`
	CPUs   int    `json:"cpus,omitempty"`
	Memory int    `json:"memory,omitempty"`
}

func tags() {
	vm := VM{
		ID:     "",
		CPUs:   4,
		Memory: 32,
	}
	json.NewEncoder(os.Stdout).Encode(vm)
}

type User struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	Token string `json:"-"`
}

func minus() {
	u := User{
		Login: "elliot",
		ID:    1000,
		Token: "f31469cc3a3f4ed685a6aa6bcf896a26",
	}
	json.NewEncoder(os.Stdout).Encode(u)
}

func main() {
	tags()

}
