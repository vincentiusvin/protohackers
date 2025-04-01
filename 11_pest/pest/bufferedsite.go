package pest

import (
	"io"
	"net"
	"protohackers/11_pest/types"
)

type BufferedSite struct {
	site Site

	// queue, group by species
	store chan types.CreatePolicy
}

func NewBufferedSiteTCP(site uint32) (Site, error) {
	conn, err := net.Dial("tcp", "pestcontrol.protohackers.com:20547")
	if err != nil {
		return nil, err
	}
	s, err := NewBufferedSite(site, conn)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func NewBufferedSite(site uint32, c io.ReadWriteCloser) (Site, error) {
	s, err := NewSite(site, c)
	if err != nil {
		return nil, err
	}

	ret := &BufferedSite{
		site:  s,
		store: make(chan types.CreatePolicy, 100),
	}

	go ret.run()

	return ret, nil
}

func (bs *BufferedSite) GetSite() uint32 {
	return bs.site.GetSite()
}

func (bs *BufferedSite) GetPops() (types.TargetPopulations, error) {
	return bs.site.GetPops()
}

func (bs *BufferedSite) UpdatePolicy(pol types.CreatePolicy) error {
	bs.store <- pol
	return nil
}

func (bs *BufferedSite) run() {
	for v := range bs.store {
		bs.site.UpdatePolicy(v)
	}
}
