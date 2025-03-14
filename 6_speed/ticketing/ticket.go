package ticketing

import (
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
}

type Road struct {
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
}

func (g *Controller) AddDispatcher(roads []uint16, ch chan infra.Ticket) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, roadNum := range roads {
		rd := g.getRoad(roadNum)
		rd.dispatchers = append(rd.dispatchers, ch)
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

	slices.SortFunc(recs, func(a *Plate, b *Plate) int {
		return int(a.Timestamp) - int(b.Timestamp)
	})

	for i, pl := range recs {
		last := i - 1
		if last < 0 && last >= len(recs) {
			continue
		}
		lastPl := recs[last]

		delta := lastPl.Timestamp - pl.Timestamp
		dur := time.Duration(delta)
		limit := rd.limit

	}
}
