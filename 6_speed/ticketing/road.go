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
	tickets     []*infra.Ticket
	plates      map[string][]*Plate // map of plate number -> plate
	limit       uint16
}

func makeRoad(roadNum uint16) *road {
	return &road{
		num:         roadNum,
		dispatchers: make([]chan *infra.Ticket, 0),
		plates:      make(map[string][]*Plate),
	}
}

func (rd *road) addTicket(t *infra.Ticket) {
	rd.tickets = append(rd.tickets, t)
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

func (rd *road) getPlateRecords(plate string) []*Plate {
	if _, ok := rd.plates[plate]; !ok {
		rd.plates[plate] = make([]*Plate, 0)
	}
	return rd.plates[plate]
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
