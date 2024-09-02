package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Event struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

func work() error {
	var data = `
	{"type":"click","x":100,"y":200}
	{"type":"move","x":101,"y":202}
	`

	dec := json.NewDecoder(strings.NewReader(data))

	for {
		var e Event
		err := dec.Decode(&e)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(e)
	}

	return nil
}

func main() {
	if err := work(); err != nil {
		fmt.Println("ERROR:", err)
	}
}
