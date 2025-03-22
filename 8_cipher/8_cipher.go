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

	log := func(str string, args ...any) {
		baseStr := fmt.Sprintf("[%v] %v", c.RemoteAddr(), str)
		log.Printf(baseStr, args...)
	}

	r := bufio.NewReader(c)
	ciphB, err := r.ReadBytes(0)
	if err != nil {
		log("%v", err)
		return
	}

	ciph, err := cipher.ParseCipher(ciphB)
	if err != nil {
		log("%v", err)
		return
	}
	dec := cipher.ApplyCipherDecode(ciph, r)
	r_decoded := bufio.NewReader(dec)
	log("cipher: %v\n", ciph)

	for {
		decoded, err := r_decoded.ReadString('\n')
		if err != nil {
			log("%v\n", err)
			return
		}
		decoded = decoded[:len(decoded)-1]
		log("decoded: %v", decoded)

		maxNum, maxRes := 0, ""
		for _, s := range strings.Split(decoded, ",") {
			numRaw, _, found := strings.Cut(s, "x")
			if !found {
				log("%v\n", fmt.Errorf("failed to find number on decoded string"))
				log("cipher: %v\n", ciph)
				return
			}
			num, err := strconv.Atoi(numRaw)
			if err != nil {
				log("%v\n", fmt.Errorf("failed to parse decoded string"))
				log("cipher: %v\n", ciph)
				return
			}

			if num <= maxNum {
				continue
			}

			maxNum = num
			maxRes = s
			log("updating max toy %v at %v\n", maxRes, maxNum)
		}

		log("returning max toy %v at %v\n", maxRes, maxNum)
		output := []byte(maxRes + "\n")
		c.Write(ciph.Encode(output))
	}
}
