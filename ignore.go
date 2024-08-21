package main

import (
	"encoding/json"
	"fmt"
)

func ignore() error {
	data := []byte(`
	{
		"login": "elliot",
		"nick": "Mr. Robot"
	}`)

	type User struct {
		Login string
		UID   int
	}

	var u User
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	fmt.Printf("%+v\n", u) // {Login:elliot UID:0}
	return nil
}

func main() {
	ignore()
}
