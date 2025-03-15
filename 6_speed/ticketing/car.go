package ticketing

import (
	"log"
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
