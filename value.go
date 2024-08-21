package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Unit string

const (
	Meter = "meter"
	Inch  = "inch"
)

type Value struct {
	Unit   Unit
	Amount float64
}

func (v Value) MarshalJSON() ([]byte, error) {
	// Step 1: Convert to type known to encoding/json
	s := fmt.Sprintf("%f%s", v.Amount, v.Unit)

	// Step 2: Use json.Marshal
	return json.Marshal(s)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("value too small")
	}

	// "2.1meter"
	r := bytes.NewReader(data[1 : len(data)-1]) // trim ""
	var a float64
	var u Unit
	if _, err := fmt.Fscanf(r, "%f%s", &a, &u); err != nil {
		return err
	}

	v.Amount = a
	v.Unit = u

	return nil
}

func work() error {
	v := Value{
		Unit:   Meter,
		Amount: 2.1,
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Println(string(data)) // "2.1meter"

	var v2 Value
	if err := json.Unmarshal(data, &v2); err != nil {
		return err
	}
	fmt.Println(v2) // {meter 2.1}
	return nil
}

func main() {
	if err := work(); err != nil {
		fmt.Println("ERROR:", err)
	}
}
