package ticketing

import (
	"log"
	"math"
	"math/rand/v2"
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
}

func (rd *road) addTicket(t *infra.Ticket) {
	rd.tickets = append(rd.tickets, t)
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

	for _, c := range rd.tickets {
		randDisp := rand.Int() % len(rd.dispatchers)
		rd.dispatchers[randDisp] <- c
	}

	rd.tickets = make([]*infra.Ticket, 0)
}

func (rd *road) issueTickets(plate string) {
	recs := rd.getPlateRecords(plate)

	violatedDays := make(map[int]bool)

	for i, pl := range recs {
		last := i - 1
		if last < 0 || last >= len(recs) {
			continue
		}
		lastPl := recs[last]

		deltaT := float64(pl.Timestamp - lastPl.Timestamp)
		deltaMile := float64(pl.Mile - lastPl.Mile)

		mph := math.Round(3600 * deltaMile / deltaT)
		limit := float64(rd.limit)

		if mph <= limit {
			continue
		}

		startDay := int(math.Floor(float64(lastPl.Timestamp) / 86400))
		endDay := int(math.Floor(float64(pl.Timestamp) / 86400))

		log.Printf("Ticketing %v. Speed: %v > %v. Days: %v-%v\n", pl.Plate, mph, limit, startDay, endDay)

		// we still need to process previous tickets
		// just prevent it from sending the ticket over
		// todo: find a better way to do this
		for currDay := startDay; currDay <= endDay; currDay += 1 {
			if violatedDays[currDay] {
				log.Printf("Skipped ticketing for day %v due to previous ticket", currDay)
				continue
			}
			violatedDays[currDay] = true

			if pl.Ticketed {
				continue
			}
			log.Printf("Ticketing for day %v", currDay)
			rd.addTicket(
				&infra.Ticket{
					Plate:      pl.Plate,
					Road:       rd.num,
					Mile1:      lastPl.Mile,
					Timestamp1: lastPl.Timestamp,
					Mile2:      pl.Mile,
					Timestamp2: pl.Timestamp,
					Speed:      uint16(100 * mph),
				},
			)
		}

		pl.Ticketed = true
	}
}
