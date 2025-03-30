package pest

import (
	"net"
	"protohackers/11_pest/infra"
	"protohackers/11_pest/types"
)

type ISite interface {
	Connect() error
	Close()
}

type Site struct {
	site uint32
	c    net.Conn
}

var x *Site = nil
var _ ISite = x

func NewSite(site uint32) *Site {
	s := &Site{
		site: site,
	}

	return s
}

func (s *Site) Connect() error {
	conn, err := net.Dial("tcp", "pestcontrol.protohackers.com:20547")
	if err != nil {
		return err
	}

	s.c = conn

	dialB := infra.Encode(types.DialAuthority{
		Site: s.site,
	})

	_, err = s.c.Write(dialB)
	if err != nil {
		return err
	}

	return nil
}

func (s *Site) Close() {
	s.c.Close()
}
