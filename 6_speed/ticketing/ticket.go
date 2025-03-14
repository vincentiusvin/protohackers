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

func (pl *Plate) String() string {
	return fmt.Sprintf("{%v %v %v %v %v}", pl.Mile, pl.Plate, pl.Road, pl.Ticketed, pl.Timestamp)
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

	for i, pl := range recs {
		last := i - 1
		if pl.Ticketed || last < 0 || last >= len(recs) {
			continue
		}
		lastPl := recs[last]

		deltaT := float64(pl.Timestamp - lastPl.Timestamp)
		deltaMile := float64(pl.Mile - lastPl.Mile)

		mph := math.Round(3600 * deltaMile / deltaT)
		limit := float64(rd.limit)

		if mph > limit {
			rd.dispatchers[0] <- infra.Ticket{
				Plate:      pl.Plate,
				Road:       rd.num,
				Mile1:      lastPl.Mile,
				Timestamp1: lastPl.Timestamp,
				Mile2:      pl.Mile,
				Timestamp2: pl.Timestamp,
				Speed:      uint16(100 * mph),
			}
			pl.Ticketed = true
			log.Printf("Ticketed %v. Speed: %v > %v\n", plate.Plate, mph, limit)
		}
	}
}
