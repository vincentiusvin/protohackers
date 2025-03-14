package ticketing_test

import (
	"protohackers/6_speed/infra"
	"protohackers/6_speed/ticketing"
	"testing"
)

func TestTicketing(t *testing.T) {
	c := ticketing.MakeController()

	var roadNum uint16 = 10
	var plate string = "UN1X"

	out := make(chan infra.Ticket, 1)
	c.AddDispatcher([]uint16{roadNum}, out)

	c.UpdateLimit(roadNum, 10)

	c.AddPlates(&ticketing.Plate{
		Plate:     plate,
		Road:      roadNum,
		Timestamp: 123456,
		Mile:      100,
	})

	c.AddPlates(&ticketing.Plate{
		Plate:     plate,
		Road:      roadNum,
		Timestamp: 123816,
		Mile:      110,
	})

}
