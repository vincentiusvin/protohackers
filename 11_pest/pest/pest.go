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

type visitSync struct {
	species string
	site    uint32
	count   uint32
	synced  bool
}

func (v visitSync) hash() string {
	return fmt.Sprintf("%v-%v", v.site, v.species)
}

type CController struct {
	mu          sync.Mutex
	sites       map[uint32]Site
	siteFactory SiteFactory

	// synchronize with authority server
	syncCh chan struct{}
	// based on hash() of visitSync
	visitData map[string]*visitSync
}

func NewControllerTCP() Controller {
	return NewController(NewSiteTCP)
}

func NewController(siteFactory SiteFactory) Controller {
	c := &CController{
		sites:       make(map[uint32]Site),
		siteFactory: siteFactory,
		visitData:   make(map[string]*visitSync),
		syncCh:      make(chan struct{}, 1), // buffered to represent needing to sync
	}
	go c.runSynchronize()
	return c
}

func (c *CController) AddSiteVisit(sv types.SiteVisit) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, svEntry := range sv.Populations {
		vs := visitSync{
			species: svEntry.Species,
			site:    sv.Site,
			count:   svEntry.Count,
			synced:  false,
		}
		c.visitData[vs.hash()] = &vs
	}
	log.Println("updating site visit data")

	select {
	case c.syncCh <- struct{}{}:
	default:
	}
	return nil
}

func (c *CController) runSynchronize() {
	for range c.syncCh {
		c.synchronize()
	}
}

func (c *CController) synchronize() {
	c.mu.Lock()
	defer c.mu.Unlock()

	log.Println("synchronizing data")
	for _, visited := range c.visitData {
		if visited.synced {
			continue
		}
		copied := *visited
		res := c.syncIndividualData(copied)
		if res {
			visited.synced = true
		}
	}
	log.Println("data synced")
}

func (c *CController) syncIndividualData(visited visitSync) bool {
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
	var changed bool

	if visited.count < pop.Min {
		pol.Action = types.PolicyConserve
		actionLog = fmt.Sprintf("conserve (%v < %v)", visited.count, pop.Min)
		changed = true
	} else if visited.count > pop.Max {
		pol.Action = types.PolicyCull
		actionLog = fmt.Sprintf("cull (%v > %v)", visited.count, pop.Max)
		changed = true
	}

	if changed {
		log.Printf("%v changing policy for %v to be %v\n", s.GetSite(), pol.Species, actionLog)
		s.UpdatePolicy(pol)
	}

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
