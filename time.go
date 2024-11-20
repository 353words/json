package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Log struct {
	Time    time.Time
	Level   string
	Message string
}

func main() {
	l := Log{
		Time:    time.Now().UTC(),
		Level:   "ERROR",
		Message: "divide by cucumber error",
	}

	if err := json.NewEncoder(os.Stdout).Encode(l); err != nil {
		fmt.Println("ERROR:", err)
	}
}
