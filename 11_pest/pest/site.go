package pest

import (
	"fmt"
	"io"
	"log"
	"net"
	"protohackers/11_pest/infra"
	"protohackers/11_pest/types"
	"sync"
)

type Site interface {
	GetSite() uint32
	GetPops() (types.TargetPopulations, error)
	UpdatePolicy(types.CreatePolicy) error
}

type CSite struct {
	mu sync.Mutex

	site uint32
	c    io.ReadWriteCloser

	// cached value
	targetPop   types.TargetPopulations
	targetPopOK bool

	// authority -> us
	helloChan        chan types.Hello
	okChan           chan types.OK
	targetPopChan    chan types.TargetPopulations
	policyResultChan chan types.PolicyResult

	// need to ensure only one policy is active
	// key is species name
	policies map[string]types.PolicyResult
}

func NewSiteTCP(site uint32) (Site, error) {
	conn, err := net.Dial("tcp", "pestcontrol.protohackers.com:20547")
	if err != nil {
		return nil, err
	}
	return NewSite(site, conn), nil
}

func NewSite(site uint32, c io.ReadWriteCloser) Site {
	s := &CSite{
		c:                c,
		site:             site,
		helloChan:        make(chan types.Hello),
		okChan:           make(chan types.OK),
		targetPopChan:    make(chan types.TargetPopulations),
		policyResultChan: make(chan types.PolicyResult),
		policies:         make(map[string]types.PolicyResult),
	}

	go s.processIncoming()
	s.handshake()
	log.Printf("%v setup finished\n", s.site)

	return s
}

func (s *CSite) GetSite() uint32 {
	return s.site
}

func (s *CSite) processIncoming() {
	defer s.close()

	var curr []byte
	for {
		b := make([]byte, 1024)
		n, err := s.c.Read(b)
		curr = append(curr, b[:n]...)
		if err != nil {
			break
		}

		for {
			res := infra.Parse(curr)
			if !res.Ok {
				break
			}

			switch v := res.Value.(type) {
			case types.Hello:
				s.helloChan <- v
			case types.OK:
				s.okChan <- v
			case types.TargetPopulations:
				s.targetPopChan <- v
			case types.PolicyResult:
				s.policyResultChan <- v
			}

			curr = res.Next
		}
	}
}

func (s *CSite) handshake() error {
	helloB := infra.Encode(types.Hello{
		Protocol: "pestcontrol",
		Version:  1,
	})

	_, err := s.c.Write(helloB)
	if err != nil {
		return err
	}

	log.Printf("%v sent hello\n", s.site)

	helloReply := <-s.helloChan
	if helloReply.Protocol != "pestcontrol" || helloReply.Version != 1 {
		return fmt.Errorf("got invalid handshake reply %v", helloReply)
	}

	log.Printf("%v got hello\n", s.site)

	return nil
}

func (s *CSite) GetPops() (ret types.TargetPopulations, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.targetPopOK {
		return s.targetPop, nil
	}

	dialB := infra.Encode(types.DialAuthority{
		Site: s.site,
	})

	_, err = s.c.Write(dialB)
	if err != nil {
		return
	}

	s.targetPop = <-s.targetPopChan
	s.targetPopOK = true
	log.Printf("%v pop reading updated %v\n", s.site, s.targetPop)

	return s.targetPop, nil
}

// update policy.
// also ensures that there is only 1 policy in place
func (s *CSite) UpdatePolicy(pol types.CreatePolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	prev, ok := s.policies[pol.Species]
	if ok {
		_, err := s.deletePolicy(types.DeletePolicy(prev))
		if err != nil {
			return err
		}
	}

	ret, err := s.createPolicy(pol)
	if err != nil {
		return err
	}
	s.policies[pol.Species] = ret
	return nil
}

func (s *CSite) createPolicy(pol types.CreatePolicy) (ret types.PolicyResult, err error) {
	polB := infra.Encode(pol)
	_, err = s.c.Write(polB)
	if err != nil {
		return
	}

	ret = <-s.policyResultChan

	var action string
	if pol.Action == types.PolicyConserve {
		action = "conserve"
	} else if pol.Action == types.PolicyCull {
		action = "cull"
	}

	log.Printf("%v policy for %v: %v\n", s.site, pol.Species, action)
	return
}

func (s *CSite) deletePolicy(pol types.DeletePolicy) (ret types.OK, err error) {
	polB := infra.Encode(pol)
	_, err = s.c.Write(polB)
	if err != nil {
		return
	}

	ret = <-s.okChan
	log.Printf("%v policy %v deleted\n", s.site, pol.Policy)
	return
}

func (s *CSite) close() {
	s.c.Close()
	close(s.helloChan)
	close(s.okChan)
	close(s.targetPopChan)
	close(s.policyResultChan)
}
