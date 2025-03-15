package ticketing

import (
	"log"
	"math"
	"protohackers/6_speed/infra"
	"slices"
)

type car struct {
	plateString string
	tickets     []*ticket
	plates      []*Plate
}

func makeCar(plate string) *car {
	return &car{
		plateString: plate,
		tickets:     make([]*ticket, 0),
		plates:      make([]*Plate, 0),
	}
}

func (c *car) addPlate(pl *Plate) {
	c.plates = append(c.plates, pl)

	log.Printf("Plate %v found on road %v at %v", pl.Plate, pl.Road, pl.Timestamp)

	slices.SortFunc(c.plates, func(a *Plate, b *Plate) int {
		return int(a.Timestamp) - int(b.Timestamp)
	})
}

func (c *car) getTickets() {

}

func (c *car) issueTickets(limit float64) {
	for i, pl := range c.plates {
		last := i - 1
		if last < 0 || last >= len(c.plates) {
			continue
		}
		lastPl := c.plates[last]

		tick := makeTicket(pl, lastPl)
		mph := tick.speed()

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
