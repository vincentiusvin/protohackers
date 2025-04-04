package lrcp

import (
	"fmt"
	"strconv"
	"strings"
)

type LRCPPackets interface {
	Encode() string
	GetSession() uint
}

func Parse(s string) (LRCPPackets, error) {
	conn, err := parseConnect(s)
	if err == nil {
		return conn, nil
	}

	data, err := parseData(s)
	if err == nil {
		return data, nil
	}

	ack, err := parseAck(s)
	if err == nil {
		return ack, nil
	}

	cls, err := parseClose(s)
	if err == nil {
		return cls, nil
	}

	return nil, fmt.Errorf("failed to parse %v", s)
}

type Connect struct {
	Session uint
}

func (c *Connect) Encode() string {
	return fmt.Sprintf("/connect/%v/", c.Session)
}

func (c *Connect) GetSession() uint {
	return c.Session
}

func parseConnect(s string) (*Connect, error) {
	curr, found := strings.CutPrefix(s, "/connect/")
	if !found {
		return nil, fmt.Errorf("not a connect request")
	}
	curr, found = strings.CutSuffix(curr, "/")
	if !found {
		return nil, fmt.Errorf("not a connect request")
	}
	sessionNum, err := strconv.Atoi(curr)
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

func (c *Data) Encode() string {
	escaped := c.Data
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "/", "\\/")

	return fmt.Sprintf("/data/%v/%v/%v/", c.Session, c.Pos, escaped)
}

func (c *Data) GetSession() uint {
	return c.Session
}

func parseData(s string) (*Data, error) {
	curr, found := strings.CutPrefix(s, "/data/")
	if !found {
		return nil, fmt.Errorf("not an ack request")
	}
	curr, found = strings.CutSuffix(curr, "/")
	if !found {
		return nil, fmt.Errorf("not an ack request")
	}

	sesRaw, lenDataRaw, found := strings.Cut(curr, "/")
	if !found {
		return nil, fmt.Errorf("not an ack request")
	}
	posRaw, dataRaw, found := strings.Cut(lenDataRaw, "/")
	if !found {
		return nil, fmt.Errorf("not an ack request")
	}

	sessionNum, err := strconv.Atoi(sesRaw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse session num to int: %w", err)
	}
	if sessionNum < 0 {
		return nil, fmt.Errorf("session number is negative: %v", sessionNum)
	}

	posNum, err := strconv.Atoi(posRaw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse position num to int: %w", err)
	}
	if posNum < 0 {
		return nil, fmt.Errorf("position number is negative: %v", sessionNum)
	}

	data := []rune{}
	escaped := false
	for _, c := range dataRaw {
		switch c {
		case '\\':
			if escaped {
				data = append(data, c)
				escaped = false
			} else {
				escaped = true
			}
		case '/':
			if escaped {
				data = append(data, c)
				escaped = false
			} else {
				return nil, fmt.Errorf("data has too many segments")
			}
		default:
			if escaped {
				return nil, fmt.Errorf("illegal escaped character")
			} else {
				data = append(data, c)
			}
		}
	}

	return &Data{
		Session: uint(sessionNum),
		Pos:     uint(posNum),
		Data:    string(data),
	}, nil
}

type Ack struct {
	Session uint
	Length  uint
}

func (c *Ack) Encode() string {
	return fmt.Sprintf("/ack/%v/%v/", c.Session, c.Length)
}

func (c *Ack) GetSession() uint {
	return c.Session
}

func parseAck(s string) (*Ack, error) {
	curr, found := strings.CutPrefix(s, "/ack/")
	if !found {
		return nil, fmt.Errorf("not an ack request")
	}
	curr, found = strings.CutSuffix(curr, "/")
	if !found {
		return nil, fmt.Errorf("not an ack request")
	}

	sesRaw, lenRaw, found := strings.Cut(curr, "/")
	if !found {
		return nil, fmt.Errorf("not an ack request")
	}

	sessionNum, err := strconv.Atoi(sesRaw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse session num to int: %w", err)
	}
	if sessionNum < 0 {
		return nil, fmt.Errorf("session number is negative: %v", sessionNum)
	}

	lenNum, err := strconv.Atoi(lenRaw)
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

type Close struct {
	Session uint
}

func (c *Close) Encode() string {
	return fmt.Sprintf("/close/%v/", c.Session)
}

func (c *Close) GetSession() uint {
	return c.Session
}

func parseClose(s string) (*Close, error) {
	curr, found := strings.CutPrefix(s, "/close/")
	if !found {
		return nil, fmt.Errorf("not an close request")
	}
	curr, found = strings.CutSuffix(curr, "/")
	if !found {
		return nil, fmt.Errorf("not an close request")
	}
	sessionNum, err := strconv.Atoi(curr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse session num to int: %w", err)
	}
	if sessionNum < 0 {
		return nil, fmt.Errorf("session number is negative: %v", sessionNum)
	}

	return &Close{
		Session: uint(sessionNum),
	}, nil
}
