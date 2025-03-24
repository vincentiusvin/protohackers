package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"protohackers/10_git/git"
	"strconv"
	"strings"
)

func handleIO(rw *bufio.ReadWriter, id string, vc *git.VersionControl) {
	log := func(format string, args ...any) {
		prefix := fmt.Sprintf("[%v] ", id)
		log.Printf(prefix+format, args...)
	}

	reply := func(format string, args ...any) error {
		s := fmt.Sprintf(format+"\n", args...)
		_, err := rw.WriteString(s)
		if err != nil {
			return err
		}
		rw.Flush()
		return nil
	}

	for {
		reply("READY")
		line, err := rw.ReadString('\n')
		if err != nil {
			log("%v", err)
			return
		}

		line = strings.TrimSpace(line)

		spls := strings.Split(line, " ")
		if len(spls) == 0 {
			reply("ERR illegal method")
		}

		cmd := spls[0]

		if cmd == "HELP" {
			reply("OK usage: HELP|GET|PUT|LIST")
		} else if cmd == "GET" {
			file, revision, err := parseGet(spls[1:])

			if err != nil {
				reply("ERR %v", err)
				continue
			}

			log("GET %v %v", file, revision)

			ret, err := vc.GetFile(file, revision)
			if err != nil {
				reply("ERR %v", err)
				continue
			}

			reply("OK %v", len(ret))

			_, err = rw.Write(ret)
			if err != nil {
				log("err: %v", err)
				return
			}
			rw.Flush()

		} else if cmd == "PUT" {
			file, fileLen, err := parsePut(spls[1:])

			if err != nil {
				reply("ERR %v", err)
				continue
			}

			log("PUT %v %v", file, fileLen)

			data := make([]byte, fileLen)
			_, err = io.ReadFull(rw, data)
			if err != nil {
				return
			}

			rev, err := vc.PutFile(file, data)
			if err != nil {
				reply("ERR %v", err)
				continue
			}

			reply("OK r%v", rev)
		} else if cmd == "LIST" {
			dir, err := parseList(spls[1:])

			if err != nil {
				reply("ERR %v", err)
				continue
			}

			log("LIST %v", dir)

			items, err := vc.ListFile(dir)
			if err != nil {
				reply("ERR %v", err)
				continue
			}

			reply("OK %v", len(items))
			for _, fli := range items {
				reply("%v %v", fli.Name, fli.Info)
			}
		} else {
			reply("ERR illegal method")
		}
	}
}

func parseGet(spls []string) (file string, revision int, err error) {
	l := len(spls)
	if l > 2 || l < 1 {
		err = fmt.Errorf("usage: GET file [revision]")
		return
	}

	file = spls[0]

	if l == 2 {
		rev_raw := spls[1]
		aft, _ := strings.CutPrefix(rev_raw, "r") // optional r prefix
		revision, err = strconv.Atoi(aft)
	}

	return
}

func parsePut(spls []string) (file string, dataLen int, err error) {
	l := len(spls)
	if l != 2 {
		err = fmt.Errorf("usage: PUT file length newline data")
		return
	}

	file = spls[0]
	dataLenRaw := spls[1]
	dataLen, err = strconv.Atoi(dataLenRaw)

	return
}

func parseList(spls []string) (dir string, err error) {
	l := len(spls)
	if l != 1 {
		err = fmt.Errorf("usage: LIST dir")
		return
	}
	dir = spls[0]
	return
}
