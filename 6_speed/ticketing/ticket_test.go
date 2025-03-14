package ticketing_test

import (
	"protohackers/6_speed/infra"
	"protohackers/6_speed/ticketing"
	"reflect"
	"testing"
)

func TestTicketing(t *testing.T) {
	c := ticketing.MakeController()

	var roadNum uint16 = 10
	var plate string = "UN1X"
	expected := infra.Ticket{
		Plate:      "UN1X",
		Road:       10,
		Mile1:      8,
		Timestamp1: 0,
		Mile2:      9,
		Timestamp2: 45,
		Speed:      8000,
	}

	outCh := make(chan infra.Ticket, 1)
	c.AddDispatcher([]uint16{roadNum}, outCh)

	c.UpdateLimit(roadNum, 60)

	c.AddPlates(&ticketing.Plate{
		Plate:     plate,
		Road:      roadNum,
		Mile:      8,
		Timestamp: 0,
	})

	c.AddPlates(&ticketing.Plate{
		Plate:     plate,
		Road:      roadNum,
		Mile:      9,
		Timestamp: 45,
	})

	out := <-outCh

	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("ticket different. expected %v. got %v", expected, out)
	}
}
