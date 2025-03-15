package ticketing

import (
	"log"
	"slices"
)

type car struct {
	plateString string
	plates      []*Plate
	violations  map[int]*ticket
}

func makeCar(plate string) *car {
	return &car{
		plateString: plate,
		plates:      make([]*Plate, 0),
		violations:  make(map[int]*ticket),
	}
}

func (c *car) addPlate(pl *Plate) {
	c.plates = append(c.plates, pl)

	log.Printf("Plate %v found on road %v mile %v time %v", pl.Plate, pl.Road, pl.Mile, pl.Timestamp)

	slices.SortFunc(c.plates, func(a *Plate, b *Plate) int {
		return int(a.Timestamp) - int(b.Timestamp)
	})
}

// calculate traffic violations
func (c *car) getViolations(limit float64) []*ticket {
	c.registerViolations(limit)

	ret := make([]*ticket, 0)
	added := make(map[*ticket]bool)

	for _, v := range c.violations {
		if added[v] {
			continue
		}
		ret = append(ret, v)
		added[v] = true
	}

	return ret
}

// calculates traffic violations and registers it to the car
func (c *car) registerViolations(limit float64) {
	tickets := func() []*ticket {
		tickets := make([]*ticket, 0)
		for i, pl := range c.plates {
			last := i - 1
			if last < 0 || last >= len(c.plates) {
				continue
			}
			lastPl := c.plates[last]

			tick := makeTicket(lastPl, pl)
			mph := tick.speed()

			if mph <= limit {
				continue
			}

			tickets = append(tickets, tick)
		}
		return tickets
	}()

	for _, t := range tickets {
		collide := false
		for _, d := range t.days() {
			if c.violations[d] != nil {
				collide = true
				break
			}
		}

		if collide {
			continue
		}

		for _, d := range t.days() {
			c.violations[d] = t
		}
	}
}
