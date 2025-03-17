package lrcp

import (
	"fmt"
	"strconv"
	"strings"
)

type Connect struct {
	Session uint
}

func ParseConnect(s string) (*Connect, error) {
	splits := strings.Split(s, "/")
	if splits[1] != "connect" {
		return nil, fmt.Errorf("not a connect request")
	}
	sessionNum, err := strconv.Atoi(splits[2])
	if err != nil {
		return nil, fmt.Errorf("cannot serialize session num to int: %w", err)
	}
	if sessionNum < 0 {
		return nil, fmt.Errorf("session number is negative: %v", sessionNum)
	}

	return &Connect{
		Session: uint(sessionNum),
	}, nil
}
