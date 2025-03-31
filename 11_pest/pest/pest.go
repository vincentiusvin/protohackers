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
	sites       map[uint32]SiteFetcher
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
		sites:       make(map[uint32]SiteFetcher),
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

	syncIndividualData := func(visited visitSync) bool {
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

		if visited.count < pop.Min {
			pol.Action = types.PolicyConserve
			s.UpdatePolicy(pol)
		} else if visited.count > pop.Max {
			pol.Action = types.PolicyCull
			s.UpdatePolicy(pol)
		}

		return true
	}

	log.Println("synchronizing data")
	var wg sync.WaitGroup
	for _, visited := range c.visitData {
		if visited.synced {
			continue
		}
		wg.Add(1)
		copied := *visited
		go func() {
			defer wg.Done()
			res := syncIndividualData(copied)
			if res {
				visited.synced = true
			}
		}()
	}

	wg.Wait()

	log.Println("data synced")
}

// this needs to be locked per site.
func (c *CController) getSite(site uint32) (Site, error) {
	c.mu.Lock()
	if c.sites[site] == nil {
		c.sites[site] = NewSiteFetcher(site, c.siteFactory)
	}
	c.mu.Unlock()

	return c.sites[site].GetSite()
}
