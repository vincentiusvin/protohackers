package main

import (
	"protohackers/6_speed/infra"
	"testing"
)

func TestLogic(t *testing.T) {
	c := make(chan any)
	send := func(b []byte) error {
		return nil
	}
	go handleConnectionLogic(c, send)

	c <- &infra.IAmACamera{
		Road:  123,
		Mile:  8,
		Limit: 60,
	}

	c <- &infra.IAmADispatcher{
		Roads: []uint16{123},
	}

	c <- &infra.Plate{
		Plate:     "UN1X",
		Timestamp: 0,
	}

	close(c)
}
