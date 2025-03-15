package ticketing

import (
	"log"
	"protohackers/6_speed/infra"
)

type road struct {
	num         uint16
	dispatchers []chan *infra.Ticket
	cars        map[string]*car
	limit       uint16
}

func makeRoad(roadNum uint16) *road {
	return &road{
		num:         roadNum,
		dispatchers: make([]chan *infra.Ticket, 0),
		cars:        make(map[string]*car),
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
	if len(rd.dispatchers) == 0 {
		return
	}
}
