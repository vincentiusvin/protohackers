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
		return nil, fmt.Errorf("cannot parse session num to int: %w", err)
	}
	if sessionNum < 0 {
		return nil, fmt.Errorf("session number is negative: %v", sessionNum)
	}

	return &Connect{
		Session: uint(sessionNum),
	}, nil
}

type Data struct {
	Session uint
	Pos     uint
	Data    string
}

func ParseData(s string) (*Data, error) {
	splits := strings.Split(s, "/")
	if splits[1] != "data" {
		return nil, fmt.Errorf("not a data request")
	}
	sessionNum, err := strconv.Atoi(splits[2])
	if err != nil {
		return nil, fmt.Errorf("cannot parse session num to int: %w", err)
	}
	if sessionNum < 0 {
		return nil, fmt.Errorf("session number is negative: %v", sessionNum)
	}

	posNum, err := strconv.Atoi(splits[3])
	if err != nil {
		return nil, fmt.Errorf("cannot parse position num to int: %w", err)
	}
	if posNum < 0 {
		return nil, fmt.Errorf("position number is negative: %v", sessionNum)
	}

	data := splits[4]

	return &Data{
		Session: uint(sessionNum),
		Pos:     uint(posNum),
		Data:    data,
	}, nil
}

type Ack struct {
	Session uint
	Length  uint
}

func ParseAck(s string) (*Ack, error) {
	splits := strings.Split(s, "/")
	if splits[1] != "ack" {
		return nil, fmt.Errorf("not an ack request")
	}
	sessionNum, err := strconv.Atoi(splits[2])
	if err != nil {
		return nil, fmt.Errorf("cannot parse session num to int: %w", err)
	}
	if sessionNum < 0 {
		return nil, fmt.Errorf("session number is negative: %v", sessionNum)
	}

	lenNum, err := strconv.Atoi(splits[3])
	if err != nil {
		return nil, fmt.Errorf("cannot parse length num to int: %w", err)
	}
	if lenNum < 0 {
		return nil, fmt.Errorf("length number is negative: %v", sessionNum)
	}

	return &Ack{
		Session: uint(sessionNum),
		Length:  uint(lenNum),
	}, nil
}
