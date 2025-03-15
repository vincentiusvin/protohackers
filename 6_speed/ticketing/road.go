package ticketing

import (
	"log"
	"math"
	"math/rand"
	"protohackers/6_speed/infra"
)

// handles most of the traffic logic
type road struct {
	num         uint16
	dispatchers []chan *infra.Ticket
	cars        map[string]*car
	limit       uint16
	sent        map[*ticket]bool
}

func makeRoad(roadNum uint16) *road {
	return &road{
		num:         roadNum,
		dispatchers: make([]chan *infra.Ticket, 0),
		cars:        make(map[string]*car),
		sent:        make(map[*ticket]bool),
	}
}

func (rd *road) updateLimit(limit uint16) {
	rd.limit = limit
	log.Printf("Road %v got speed limit updated to: %v", rd.num, limit)
}

func (rd *road) getCar(plate string) *car {
	if _, ok := rd.cars[plate]; !ok {
		rd.cars[plate] = makeCar(plate)
	}
	return rd.cars[plate]
}

func (rd *road) addPlate(plate *Plate) {
	car := rd.getCar(plate.Plate)
	car.addPlate(plate)
	rd.processTicket()
}

func (rd *road) addDispatcher(ch chan *infra.Ticket) {
	rd.dispatchers = append(rd.dispatchers, ch)
	log.Printf("Dispatcher registered on road %v", rd.num)
	rd.processTicket()
}

func (rd *road) processTicket() {
	// set nulls to false
	for _, c := range rd.cars {
		for _, t := range c.getViolations(float64(rd.limit)) {
			if _, ok := rd.sent[t]; !ok {
				rd.sent[t] = false
			}
		}
	}

	if len(rd.dispatchers) == 0 {
		return
	}

	for ticket, sent := range rd.sent {
		if sent {
			continue
		}

		mph := ticket.speed()

		randDisp := rand.Int() % len(rd.dispatchers)
		rd.dispatchers[randDisp] <- &infra.Ticket{
			Plate:      ticket.pl1.Plate,
			Road:       ticket.pl1.Road,
			Mile1:      ticket.pl1.Mile,
			Timestamp1: ticket.pl1.Timestamp,
			Mile2:      ticket.pl2.Mile,
			Timestamp2: ticket.pl2.Timestamp,
			Speed:      uint16(math.Round(mph * 100)),
		}
		rd.sent[ticket] = true

		log.Printf("Ticketing %v. Speed: %v > %v\n", ticket.pl1.Plate, mph, rd.limit)
	}
}
