package ticketing

import (
	"fmt"
	"log"
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

func (pl *Plate) String() string {
	return fmt.Sprintf("{%v %v %v %v %v}", pl.Mile, pl.Plate, pl.Road, pl.Ticketed, pl.Timestamp)
}

type Controller struct {
	roads map[uint16]*road
	mu    sync.Mutex
}

func MakeController() *Controller {
	return &Controller{
		roads: make(map[uint16]*road),
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

	rd.issueTickets(plate.Plate)
}

func (g *Controller) getRoad(roadNum uint16) *road {
	if _, ok := g.roads[roadNum]; !ok {
		g.roads[roadNum] = makeRoad(roadNum)
	}
	return g.roads[roadNum]
}
