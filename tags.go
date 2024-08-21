package main

import (
	"encoding/json"
	"fmt"
)

func tags() error {
	data := []byte(`
	{
		"name": "elliot",
		"uid": 1000
	}`)

	type User struct {
		Login string `json:"name"`
		UID   int
	}

	var u User
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	fmt.Printf("%+v\n", u) // {Login:elliot UID:1000}
	return nil
}

func main() {
	tags()
}
