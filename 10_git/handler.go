package main

import (
	"bufio"
	"fmt"
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
