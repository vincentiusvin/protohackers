package pest

import (
	"fmt"
	"log"
	"protohackers/11_pest/types"
	"sync"
)

var (
	ErrInvalidSiteVisit = fmt.Errorf("invalid site visit. found conflicting keys")
	ErrInvalidHandshake = fmt.Errorf("invalid handshake")
)

type Controller interface {
	AddSiteVisit(sv types.SiteVisit) error
}

type SiteFactory = func(site uint32) (Site, error)

type VisitData struct {
	species string
	site    uint32
	count   uint32
	min     uint32
	max     uint32
}

type CController struct {
	mu          sync.Mutex
	sites       map[uint32]Site
	siteFactory SiteFactory
}

func NewControllerTCP() Controller {
	return NewController(NewBufferedSiteTCP)
}

func NewController(siteFactory SiteFactory) Controller {
	c := &CController{
		sites:       make(map[uint32]Site),
		siteFactory: siteFactory,
	}
	return c
}

func (c *CController) AddSiteVisit(sv types.SiteVisit) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	site := c.getSite(sv.Site)

	pops, err := site.GetPops()
	if err != nil {
		return err
	}

	visitLkp := make(map[string]types.SiteVisitEntry)

	for _, visitEntry := range sv.Populations {
		prev, ok := visitLkp[visitEntry.Species]
		if ok && prev.Count != visitEntry.Count {
			return ErrInvalidSiteVisit
		}
		visitLkp[visitEntry.Species] = visitEntry
	}

	for _, popLimit := range pops.Populations {
		vd := VisitData{
			species: popLimit.Species,
			site:    sv.Site,
			min:     popLimit.Min,
			max:     popLimit.Max,
		}

		visit, ok := visitLkp[popLimit.Species]
		if ok {
			vd.count = visit.Count
		}

		err := c.updatePolicy(vd, site)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (c *CController) updatePolicy(visited VisitData, site Site) error {
	pol := types.CreatePolicy{
		Species: visited.species,
	}

	var actionLog string

	if visited.count < visited.min {
		pol.Action = types.PolicyConserve
		actionLog = fmt.Sprintf("conserve (%v < %v)", visited.count, visited.min)
	} else if visited.count > visited.max {
		pol.Action = types.PolicyCull
		actionLog = fmt.Sprintf("cull (%v > %v)", visited.count, visited.max)
	} else {
		pol.Action = types.PolicyNothing
		actionLog = fmt.Sprintf("nothing (%v <= %v <= %v)", visited.min, visited.count, visited.max)
	}

	log.Printf("%v changing policy for %v to be %v\n", visited.site, pol.Species, actionLog)
	err := site.UpdatePolicy(pol)
	if err != nil {
		return err
	}

	return nil
}

// this needs to be locked per site.
func (c *CController) getSite(site uint32) Site {
	if c.sites[site] == nil {
		for {
			var err error
			c.sites[site], err = c.siteFactory(site)
			if err != nil {
				continue
			}
			break
		}
	}

	return c.sites[site]
}
