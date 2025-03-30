package pest

import (
	"protohackers/11_pest/types"
	"sync"
)

type Controller interface {
	AddSiteVisit(sv types.SiteVisit) error
}

type SiteFactory = func(site uint32) (Site, error)

type CController struct {
	sites       map[uint32]Site
	siteFactory SiteFactory
	mu          sync.Mutex
}

func NewControllerTCP() Controller {
	return NewController(NewSiteTCP)
}

func NewController(siteFactory SiteFactory) Controller {
	return &CController{
		sites:       make(map[uint32]Site),
		siteFactory: siteFactory,
	}
}

func (c *CController) AddSiteVisit(sv types.SiteVisit) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	s, err := c.getSite(sv.Site)
	if err != nil {
		return err
	}
	count, err := s.GetPops()
	if err != nil {
		return err
	}

	// change to map if too slow
	for _, svEntry := range sv.Populations {
		for _, countEntry := range count.Populations {
			if svEntry.Species != countEntry.Species {
				continue
			}
			if svEntry.Count < countEntry.Min {
				s.UpdatePolicy(types.CreatePolicy{
					Species: svEntry.Species,
					Action:  types.PolicyConserve,
				})
			} else if svEntry.Count > countEntry.Max {
				s.UpdatePolicy(types.CreatePolicy{
					Species: svEntry.Species,
					Action:  types.PolicyCull,
				})
			}
		}
	}

	return nil
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
