package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Event struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

func work() error {
	events := []Event{
		{"click", 100, 200},
		{"move", 101, 202},
	}

	enc := json.NewEncoder(os.Stdout)

	for _, e := range events {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := work(); err != nil {
		fmt.Println("ERROR:", err)
	}
}
