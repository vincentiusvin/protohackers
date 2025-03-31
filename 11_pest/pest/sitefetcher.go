package pest

import "sync"

type SiteFetcher interface {
	GetSite() (Site, error)
}

func NewSiteFetcher(site uint32, fact SiteFactory) SiteFetcher {
	return &siteFetcher{
		sitenum: site,
		factory: fact,
	}
}

type siteFetcher struct {
	mu      sync.Mutex
	site    Site
	sitenum uint32
	factory SiteFactory
}

func (s *siteFetcher) GetSite() (Site, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.site == nil {
		newSite, err := s.factory(s.sitenum)
		if err != nil {
			return nil, err
		}
		s.site = newSite
	}

	return s.site, nil
}
