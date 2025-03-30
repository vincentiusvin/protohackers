package pest

import (
	"protohackers/11_pest/types"
)

type Controller struct {
	sites map[uint32]Site
}

func NewController() *Controller {
	return &Controller{
		sites: make(map[uint32]Site),
	}
}

func (c *Controller) AddSiteVisit(sv types.SiteVisit) {
	s, err := c.getSite(sv.Site)
	count, err := s.GetPops()

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
					Action:  types.PolicyConserve,
				})
			}
		}
	}
}

func (c *Controller) getSite(site uint32) (Site, error) {
	if c.sites[site] == nil {
		ns := NewSite(site)
		err := ns.Connect()
		if err != nil {
			return nil, err
		}
		c.sites[site] = ns
	}

	return c.sites[site], nil
}
