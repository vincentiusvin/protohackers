package ticketing

import (
	"log"
	"protohackers/6_speed/infra"
	"slices"
	"sync"
	"time"
)

type Plate struct {
	Plate     string
	Timestamp uint32
	Road      uint16
	Ticketed  bool
	Mile      uint16
}

type Road struct {
	num         uint16
	dispatchers []chan infra.Ticket
	plates      map[string][]*Plate // map of plate number -> plate
	limit       uint16
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

func (g *Controller) getRoad(roadNum uint16) *Road {
	if _, ok := g.roads[roadNum]; !ok {
		g.roads[roadNum] = &Road{
			num:         roadNum,
			dispatchers: make([]chan infra.Ticket, 0),
			plates:      make(map[string][]*Plate),
		}
	}

	return g.roads[roadNum]
}

func (g *Controller) UpdateLimit(roadNum uint16, limit uint16) {
	g.mu.Lock()
	defer g.mu.Unlock()

	rd := g.getRoad(roadNum)
	rd.limit = limit
	log.Printf("Road %v got speed limit updated to: %v", roadNum, limit)
}

func (g *Controller) AddDispatcher(roads []uint16, ch chan infra.Ticket) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, roadNum := range roads {
		rd := g.getRoad(roadNum)
		rd.dispatchers = append(rd.dispatchers, ch)
		log.Printf("Dispatcher registered on road %v", roadNum)
	}
}

func (rd *Road) getPlateRecords(plate string) []*Plate {
	if _, ok := rd.plates[plate]; !ok {
		rd.plates[plate] = make([]*Plate, 0)
	}
	return rd.plates[plate]
}

func (g *Controller) AddPlates(plate *Plate) {
	g.mu.Lock()
	defer g.mu.Unlock()

	rd := g.getRoad(plate.Road)
	recs := rd.getPlateRecords(plate.Plate)

	recs = append(recs, plate)

	log.Printf("Plate %v found on road %v at %v", plate.Plate, plate.Road, plate.Timestamp)

	slices.SortFunc(recs, func(a *Plate, b *Plate) int {
		return int(a.Timestamp) - int(b.Timestamp)
	})

	for i, pl := range recs {
		last := i - 1
		if last < 0 || last >= len(recs) {
			continue
		}
		lastPl := recs[last]

		delta := lastPl.Timestamp - pl.Timestamp
		dur := time.Duration(delta)
		limit := rd.limit
		log.Println(dur, limit)

	}
}
