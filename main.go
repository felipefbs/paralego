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

func contentAsync(url string, output chan<- string) error {
	c, err := content(url)
	if err != nil {
		return err
	}

	output <- c

	return nil
}

func countAsync(input <-chan string, output chan<- Counts) {
	for c := range input {
		output <- countWords(c)
	}

	close(output)
}

func asyncJob() {
	size := 1000

	q := make(chan string, size)
	o := make(chan Counts, size)

	for i := 0; i < size*2; i++ {
		go contentAsync(url, q)
	}

	go countAsync(q, o)

	for i := 0; i < size*2; i++ {
		fmt.Println(len(<-o))
	}
}

func main() {
	asyncJob()
}
