package ticketing

import (
	"math"
)

// this is a prospective ticket
// determine whether it is applicable or not by checking the speed
type ticket struct {
	pl1 *Plate // pl1 should be lesser or equal to pl2
	pl2 *Plate
}

func makeTicket(pl1 *Plate, pl2 *Plate) *ticket {
	if pl1 == nil || pl2 == nil {
		panic("plate cannot be nil")
	}
	if pl1.Timestamp > pl2.Timestamp {
		panic("plate 1 needs to happen before plate 2")
	}

	return &ticket{
		pl1: pl1,
		pl2: pl2,
	}
}

func (t *ticket) speed() float64 {
	deltaT := float64(t.pl2.Timestamp - t.pl1.Timestamp)
	deltaMile := float64(t.pl2.Mile - t.pl1.Mile)
	mph := math.Round(3600 * deltaMile / deltaT)
	return mph
}

func (t *ticket) days() []int {
	ret := make([]int, 0)
	startDay := math.Floor(float64(t.pl1.Timestamp) / 86400)
	endDay := math.Floor(float64(t.pl2.Timestamp) / 86400)

	for curr := startDay; curr <= endDay; curr++ {
		ret = append(ret, int(curr))
	}

	return ret
}
