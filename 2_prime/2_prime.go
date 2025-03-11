package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net"
)

type PrimeRequest struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

type PrimeResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server listening at " + addr)

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go prime(c)
	}
}

func prime(c net.Conn) {
	defer c.Close()

	sc := bufio.NewScanner(c)
	d := json.NewEncoder(c)

	for sc.Scan() {
		s := sc.Text()
		pr, err := parsePrimeRequest(s)
		var resp PrimeResponse

		if err != nil {
			fmt.Println(err)
			d.Encode(resp) // a response struct with zero values is considered malformed
			break
		}

		resp.Method = "isPrime"
		resp.Prime = isPrime(pr.Number)
		d.Encode(resp)
	}
	fmt.Println("EOF")
}

func isPrime(number int) bool {
	number = int(math.Abs(float64(number)))

	if number <= 1 {
		return false
	}
	for i := 2; i < number; i++ { // should use sqrt but whatever
		if (number % i) == 0 {
			return false
		}
	}
	return true
}

func parsePrimeRequest(s string) (*PrimeRequest, error) {
	// to differentiate zero values and non existant fields
	type PrimeRequestNullable struct {
		Method *string `json:"method"`
		Number *int    `json:"number"`
	}

	b := new(bytes.Buffer)
	b.WriteString(s)

	d := json.NewDecoder(b)

	var pr PrimeRequestNullable
	err := d.Decode(&pr)
	if err != nil {
		return nil, err
	}

	if pr.Method == nil {
		return nil, fmt.Errorf("no method")
	}

	if pr.Number == nil {
		return nil, fmt.Errorf("no number")
	}

	if *pr.Method != "isPrime" {
		return nil, fmt.Errorf("wrong method")
	}

	return &PrimeRequest{
		Method: *pr.Method,
		Number: *pr.Number,
	}, nil
}
