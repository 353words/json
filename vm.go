package main

import (
	"encoding/json"
	"fmt"
)

type StartVM struct {
	Image string
	Count int
}

func vm() error {
	data := []byte(`{"image": "debian:bookworm-slim"}`)

	var req StartVM
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}
	fmt.Println(req)
	return nil
}

func main() {
	vm()
}
