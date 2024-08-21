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

type StartVMPtr struct {
	Image string
	Count *int
}

func ptr() error {
	data := []byte(`{"image": "debian:bookworm-slim"}`)

	var req StartVMPtr
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	if req.Count == nil { // User didn't send "count", use default value
		c := 1
		req.Count = &c
	}

	if *req.Count < 1 {
		return fmt.Errorf("bad count: %d", *req.Count)
	}

	fmt.Printf("%+v\n", *req.Count)
	return nil
}

func mmap() error {
	data := []byte(`{"image": "debian:bookworm-slim", "count": 0}`)

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if _, ok := m["count"]; !ok { // User didn't send "count", use default value
		m["count"] = 1
	}

	image, ok := m["image"].(string)
	if !ok || image == "" {
		return fmt.Errorf("bad image: %#v", m["image"])
	}

	count, ok := m["count"].(float64)
	if !ok {
		return fmt.Errorf("bad count: %#v", m["count"])
	}

	if count < 1 {
		return fmt.Errorf("bad count: %f", count)
	}

	req := StartVM{
		Image: image,
		Count: int(count),
	}

	fmt.Println(req)
	return nil
}

func deflt() error {
	data := []byte(`{"image": "debian:bookworm-slim", "count": 0}`)

	req := StartVM{
		Count: 1,
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	if req.Count < 1 {
		return fmt.Errorf("bad count: %d", req.Count)
	}

	fmt.Printf("%+v\n", req)

	return nil
}

func main() {
	// fn := vm
	// fn := ptr
	// fn := mmap
	fn := deflt
	if err := fn(); err != nil {
		fmt.Println("ERROR:", err)
	}
}
