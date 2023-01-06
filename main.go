package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Counts map[string]int

const url = "https://google.com"

func content(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	output := &strings.Builder{}
	buff := make([]byte, 1024)

	for {
		n, _ := res.Body.Read(buff)
		if n <= 0 {
			break
		}

		output.Write(buff[:n])
	}

	time.Sleep(2 * time.Second)

	return output.String(), nil
}

func countWords(content string) Counts {
	out := Counts{}

	words := strings.Split(content, ",")

	for _, w := range words {
		out[w]++
	}

	return out
}

func syncJob() {
	c, err := content(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", countWords(c))
}

func main() {
	syncJob()
}
