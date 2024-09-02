package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Event struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

func queryEvents() chan Event {
	events := []Event{
		{"click", 100, 200},
		{"move", 101, 202},
		{"move", 102, 203},
		{"move", 103, 204},
		{"move", 104, 204},
		{"click", 104, 204},
	}

	ch := make(chan Event)

	go func() {
		for _, evt := range events {
			ch <- evt
			time.Sleep(100 * time.Millisecond)
		}
		close(ch)
	}()

	return ch
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	ctrl := http.NewResponseController(w)

	enc := json.NewEncoder(w)
	for evt := range queryEvents() {
		if err := enc.Encode(evt); err != nil {
			// Can't set error
			slog.Error("JSON encode", "error", err)
			return
		}

		if err := ctrl.Flush(); err != nil {
			slog.Error("flush", "error", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/events", eventsHandler)

	addr := ":8080"
	slog.Info("server starting", "address", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
