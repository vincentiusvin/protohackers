package pest

import (
	"io"
	"net"
	"protohackers/11_pest/types"
	"sync"
)

type BufferedSite struct {
	site Site
	mu   sync.Mutex

	// queue, group by species
	queue map[string]chan types.CreatePolicy
}

func NewBufferedSiteTCP(site uint32) (Site, error) {
	conn, err := net.Dial("tcp", "pestcontrol.protohackers.com:20547")
	if err != nil {
		return nil, err
	}
	return NewBufferedSite(site, conn), nil
}

func NewBufferedSite(site uint32, c io.ReadWriteCloser) Site {
	return &BufferedSite{
		site:  NewSite(site, c),
		queue: make(map[string]chan types.CreatePolicy),
	}
}

func (bs *BufferedSite) GetSite() uint32 {
	return bs.site.GetSite()
}

func (bs *BufferedSite) GetPops() (types.TargetPopulations, error) {
	return bs.site.GetPops()
}

func (bs *BufferedSite) UpdatePolicy(pol types.CreatePolicy) error {
	worker := bs.getWorker(pol.Species)
	worker <- pol
	return nil
}

func (bs *BufferedSite) getWorker(species string) chan types.CreatePolicy {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if _, ok := bs.queue[species]; !ok {
		ch := make(chan types.CreatePolicy)
		bs.queue[species] = ch

		go func() {
			for pol := range ch {
				bs.site.UpdatePolicy(pol)
			}
		}()
	}

	return bs.queue[species]
}
