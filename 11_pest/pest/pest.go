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
	visitData map[string]visitSync
}

func NewControllerTCP() Controller {
	return NewController(NewSiteTCP)
}

func NewController(siteFactory SiteFactory) Controller {
	c := &CController{
		sites:       make(map[uint32]Site),
		siteFactory: siteFactory,
		visitData:   make(map[string]visitSync),
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
		}
		c.visitData[vs.hash()] = vs
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
		s, err := c.getSite(visited.site)
		if err != nil {
			return
		}
		pops, err := s.GetPops()
		if err != nil {
			return
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
			continue
		}

		pol := types.CreatePolicy{
			Species: visited.species,
		}

		if visited.count < pop.Min {
			pol.Action = types.PolicyConserve
		} else if visited.count > pop.Max {
			pol.Action = types.PolicyCull
		}

		s.UpdatePolicy(pol)
	}
	log.Println("data synced")
}

func (c *CController) getSite(site uint32) (Site, error) {
	if c.sites[site] == nil {
		ns, err := c.siteFactory(site)
		if err != nil {
			return nil, err
		}
		c.sites[site] = ns
	}

	return c.sites[site], nil
}
