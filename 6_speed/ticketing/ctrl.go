package ticketing

import (
	"fmt"
	"protohackers/6_speed/infra"
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
	rd.updateLimit(limit)
}

func (g *Controller) AddDispatcher(roads []uint16, ch chan *infra.Ticket) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, roadNum := range roads {
		rd := g.getRoad(roadNum)
		rd.addDispatcher(ch)
	}

}

func (g *Controller) AddPlates(plate *Plate) {
	g.mu.Lock()
	defer g.mu.Unlock()

	rd := g.getRoad(plate.Road)
	rd.addPlate(plate)
}

func (g *Controller) getRoad(roadNum uint16) *road {
	if _, ok := g.roads[roadNum]; !ok {
		g.roads[roadNum] = makeRoad(roadNum)
	}
	return g.roads[roadNum]
}
