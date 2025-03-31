package pest

import (
	"fmt"
	"log"
	"protohackers/11_pest/types"
	"sync"
)

type Controller interface {
	AddSiteVisit(sv types.SiteVisit) error
}

type SiteFactory = func(site uint32) (Site, error)

type VisitData struct {
	species string
	site    uint32
	count   uint32
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

	for _, svEntry := range sv.Populations {
		vd := VisitData{
			species: svEntry.Species,
			site:    sv.Site,
			count:   svEntry.Count,
		}
		c.syncIndividualData(vd)
	}
	return nil
}

func (c *CController) syncIndividualData(visited VisitData) bool {
	s, err := c.getSite(visited.site)
	if err != nil {
		return false
	}
	pops, err := s.GetPops()
	if err != nil {
		return false
	}

	var pop types.TargetPopulationsEntry
	var popFound bool
	for _, targetPop := range pops.Populations {
		if targetPop.Species != visited.species {
			continue
		}
		popFound = true
		pop = targetPop
	}

	if !popFound {
		return false
	}

	pol := types.CreatePolicy{
		Species: visited.species,
	}

	var actionLog string

	if visited.count < pop.Min {
		pol.Action = types.PolicyConserve
		actionLog = fmt.Sprintf("conserve (%v < %v)", visited.count, pop.Min)
	} else if visited.count > pop.Max {
		pol.Action = types.PolicyCull
		actionLog = fmt.Sprintf("cull (%v > %v)", visited.count, pop.Max)
	} else {
		pol.Action = types.PolicyNothing
		actionLog = fmt.Sprintf("nothing (%v <= %v <= %v)", pop.Min, visited.count, pop.Max)
	}

	log.Printf("%v changing policy for %v to be %v\n", s.GetSite(), pol.Species, actionLog)
	s.UpdatePolicy(pol)

	return true
}

// this needs to be locked per site.
func (c *CController) getSite(site uint32) (Site, error) {
	if c.sites[site] == nil {
		var err error
		c.sites[site], err = c.siteFactory(site)
		if err != nil {
			return nil, err
		}
	}

	return c.sites[site], nil
}
