package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	urlTemplate = "https://api.stocktwits.com/api/2/streams/symbol/%s.json"
)

func relatedStocks(symbol string) (map[string]int, error) {
	url := fmt.Sprintf(urlTemplate, symbol)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var reply struct { // anonymous struct
		Messages []struct {
			Symbols []struct {
				Symbol string
			}
		}
	}

	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, err
	}

	related := make(map[string]int)
	for _, m := range reply.Messages {
		for _, s := range m.Symbols {
			if s.Symbol != "AAPL" {
				related[s.Symbol]++
			}
		}
	}

	return related, nil
}

func main() {
	if len(os.Args) != 2 {
		// FIXME: use the flag package
		log.Fatal("error: wrong number of arguments")
	}

	symbol := os.Args[1]
	counts, err := relatedStocks(symbol)
	if err != nil {
		log.Fatal(err)
	}

	for sym, n := range counts {
		fmt.Printf("%10s %d\n", sym, n)
	}
}
