package ticketing

import (
	"math/rand/v2"
	"protohackers/6_speed/infra"
)

type Road struct {
	num         uint16
	dispatchers []chan *infra.Ticket
	tickets     []*infra.Ticket
	plates      map[string][]*Plate // map of plate number -> plate
	limit       uint16
}

func MakeRoad(roadNum uint16) *Road {
	return &Road{
		num:         roadNum,
		dispatchers: make([]chan *infra.Ticket, 0),
		plates:      make(map[string][]*Plate),
	}
}

func (rd *Road) addTicket(t *infra.Ticket) {
	rd.tickets = append(rd.tickets, t)
	rd.processTicket()
}

func (rd *Road) processTicket() {
	if len(rd.dispatchers) == 0 {
		return
	}

	for _, c := range rd.tickets {
		randDisp := rand.Int() % len(rd.dispatchers)
		rd.dispatchers[randDisp] <- c
	}

	rd.tickets = make([]*infra.Ticket, 0)
}

func (rd *Road) getPlateRecords(plate string) []*Plate {
	if _, ok := rd.plates[plate]; !ok {
		rd.plates[plate] = make([]*Plate, 0)
	}
	return rd.plates[plate]
}
