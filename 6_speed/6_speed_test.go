package main

import (
	"protohackers/6_speed/infra"
	"testing"
)

func TestLogic(t *testing.T) {
	in := make(chan any)
	out := make(chan any)
	go handleConnectionLogic(nil, in, out)

	in <- &infra.IAmACamera{
		Road:  123,
		Mile:  8,
		Limit: 60,
	}

	in <- &infra.IAmADispatcher{
		Roads: []uint16{123},
	}

	in <- &infra.Plate{
		Plate:     "UN1X",
		Timestamp: 0,
	}

	close(in)
}
