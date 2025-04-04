package main

import (
	"context"
	"protohackers/6_speed/infra"
	"protohackers/6_speed/ticketing"
	"testing"
)

func TestLogic(t *testing.T) {
	in := make(chan any)
	defer close(in)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := make(chan infra.Encode)

	ctrl := ticketing.MakeController()

	go handleConnectionLogic(ctx, in, out, ctrl)

	in <- &infra.IAmACamera{
		Road:  123,
		Mile:  8,
		Limit: 60,
	}

	in <- &infra.IAmADispatcher{
		Roads: []uint16{123},
	}
	outerr := <-out
	switch outerr.(type) {
	case *infra.SpeedError:
	default:
		t.Fatalf("expected error when sent iamdispatcher for camera")
	}

	in <- &infra.Plate{
		Plate:     "UN1X",
		Timestamp: 0,
	}
}
