package ticketing

import (
	"fmt"
	"log"
	"math"
	"protohackers/6_speed/infra"
	"slices"
	"sync"
)

type Plate struct {
	Plate     string
	Timestamp uint32
	Road      uint16
	Ticketed  bool
	Mile      uint16
}

type Controller struct {
	roads map[uint16]*Road
	mu    sync.Mutex
}

func MakeController() *Controller {
	return &Controller{
		roads: make(map[uint16]*Road),
	}
}

func (g *Controller) UpdateLimit(roadNum uint16, limit uint16) {
	g.mu.Lock()
	defer g.mu.Unlock()

	rd := g.getRoad(roadNum)
	rd.limit = limit
	log.Printf("Road %v got speed limit updated to: %v", roadNum, limit)
}

func (g *Controller) AddDispatcher(roads []uint16, ch chan *infra.Ticket) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, roadNum := range roads {
		rd := g.getRoad(roadNum)
		rd.dispatchers = append(rd.dispatchers, ch)
		rd.processTicket()
		log.Printf("Dispatcher registered on road %v", roadNum)
	}

}

func (g *Controller) AddPlates(plate *Plate) {
	g.mu.Lock()
	defer g.mu.Unlock()

	rd := g.getRoad(plate.Road)

	recs := rd.getPlateRecords(plate.Plate)
	recs = append(recs, plate)
	rd.plates[plate.Plate] = recs

	log.Printf("Plate %v found on road %v at %v", plate.Plate, plate.Road, plate.Timestamp)

	slices.SortFunc(recs, func(a *Plate, b *Plate) int {
		return int(a.Timestamp) - int(b.Timestamp)
	})

	g.issueTickets(plate.Road, plate.Plate)
}

func (pl *Plate) String() string {
	return fmt.Sprintf("{%v %v %v %v %v}", pl.Mile, pl.Plate, pl.Road, pl.Ticketed, pl.Timestamp)
}

func (g *Controller) issueTickets(road uint16, plate string) {
	rd := g.getRoad(road)
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

func (g *Controller) getRoad(roadNum uint16) *Road {
	if _, ok := g.roads[roadNum]; !ok {
		g.roads[roadNum] = MakeRoad(roadNum)
	}
	return g.roads[roadNum]
}
