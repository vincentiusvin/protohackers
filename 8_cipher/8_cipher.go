package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"protohackers/8_cipher/cipher"
	"strconv"
	"strings"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at " + addr)

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	r := bufio.NewReader(c)
	ciphB, err := r.ReadBytes(0)
	if err != nil {
		log.Println(err)
		return
	}

	ciph := cipher.ParseCipher(ciphB)
	dec := cipher.ApplyCipherDecode(ciph, r)
	r_decoded := bufio.NewReader(dec)

	for {
		decoded, err := r_decoded.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Decoded: %v", decoded)

		decodedResult := string(decoded)

		maxNum, maxRes := 0, ""
		for _, s := range strings.Split(decodedResult, ",") {
			numRaw, _, found := strings.Cut(s, "x")
			if !found {
				panic(fmt.Errorf("failed to find number on decoded string"))
			}
			num, err := strconv.Atoi(numRaw)
			if err != nil {
				panic(fmt.Errorf("failed to parse decoded string"))
			}

			if num <= maxNum {
				continue
			}

			maxNum = num
			maxRes = s
			log.Printf("updating max toy %v at %v\n", maxRes, maxNum)
		}

		log.Printf("returning max toy %v at %v\n", maxRes, maxNum)
		output := []byte(maxRes + "\n")
		c.Write(ciph.Encode(output))
	}
}
