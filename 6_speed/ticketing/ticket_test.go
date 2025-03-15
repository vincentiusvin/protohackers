package ticketing_test

import (
	"protohackers/6_speed/infra"
	"protohackers/6_speed/ticketing"
	"reflect"
	"testing"
)

func TestTicketingBasic(t *testing.T) {
	type basicTestCases struct {
		expected *infra.Ticket
		in1      *ticketing.Plate
		in2      *ticketing.Plate
	}

	var limit uint16 = 0
	var roadNum uint16 = 10
	cases := []basicTestCases{
		{
			in1: &ticketing.Plate{
				Plate:     "UN1X",
				Road:      roadNum,
				Mile:      8,
				Timestamp: 0,
			},
			in2: &ticketing.Plate{
				Plate:     "UN1X",
				Road:      roadNum,
				Mile:      9,
				Timestamp: 45,
			},
			expected: &infra.Ticket{
				Plate:      "UN1X",
				Road:       roadNum,
				Mile1:      8,
				Timestamp1: 0,
				Mile2:      9,
				Timestamp2: 45,
				Speed:      8000,
			},
		},
		{
			in1: &ticketing.Plate{
				Plate:     "UN1X",
				Road:      roadNum,
				Mile:      9,
				Timestamp: 0,
			},
			in2: &ticketing.Plate{
				Plate:     "UN1X",
				Road:      roadNum,
				Mile:      8,
				Timestamp: 45,
			},
			expected: &infra.Ticket{
				Plate:      "UN1X",
				Road:       roadNum,
				Mile1:      9,
				Timestamp1: 0,
				Mile2:      8,
				Timestamp2: 45,
				Speed:      8000,
			},
		},
	}

	for _, ticketCase := range cases {
		c := ticketing.MakeController()

		c.UpdateLimit(roadNum, limit)
		c.AddPlates(ticketCase.in1)
		c.AddPlates(ticketCase.in2)

		// Need to buffer ticket first if no dispatcher.
		// So we register it late to test for that.
		outCh := make(chan *infra.Ticket, 1)
		c.AddDispatcher([]uint16{roadNum}, outCh)
		out := <-outCh

		if !reflect.DeepEqual(out, ticketCase.expected) {
			t.Fatalf("ticket different. expected %v. got %v", ticketCase.expected, out)
		}
	}
}

// 3 traffic violations but we should get 2 tickets
func TestTicketingDay1(t *testing.T) {
	c := ticketing.MakeController()

	var roadNum uint16 = 10
	var plate string = "UN1X"

	c.UpdateLimit(roadNum, 60)

	outCh := make(chan *infra.Ticket, 2)
	c.AddDispatcher([]uint16{roadNum}, outCh)

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

	c.AddPlates(&ticketing.Plate{
		Plate:     plate,
		Road:      roadNum,
		Mile:      10,
		Timestamp: 90,
	})

	<-outCh

	select {
	case <-outCh:
		t.Fatalf("got two tickets on same day. expected only 1")
	default:
	}
}

// 3 traffic violations but we should get 2 tickets
// same as above but in different orders
func TestTicketingDay2(t *testing.T) {
	c := ticketing.MakeController()

	var roadNum uint16 = 10
	var plate string = "UN1X"

	outCh := make(chan *infra.Ticket, 2)
	c.AddDispatcher([]uint16{roadNum}, outCh)

	c.UpdateLimit(roadNum, 60)

	c.AddPlates(&ticketing.Plate{
		Plate:     plate,
		Road:      roadNum,
		Mile:      9,
		Timestamp: 45,
	})

	c.AddPlates(&ticketing.Plate{
		Plate:     plate,
		Road:      roadNum,
		Mile:      10,
		Timestamp: 90,
	})

	c.AddPlates(&ticketing.Plate{
		Plate:     plate,
		Road:      roadNum,
		Mile:      8,
		Timestamp: 0,
	})

	<-outCh

	select {
	case <-outCh:
		t.Fatalf("got two tickets on same day. expected only 1")
	default:
	}
}
