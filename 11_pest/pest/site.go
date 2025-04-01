package pest

import (
	"errors"
	"io"
	"log"
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

func NewSite(site uint32, c io.ReadWriteCloser) (Site, error) {
	s := &CSite{
		c:                c,
		site:             site,
		helloChan:        make(chan types.Hello),
		okChan:           make(chan types.OK),
		targetPopChan:    make(chan types.TargetPopulations),
		policyResultChan: make(chan types.PolicyResult),
		policies:         make(map[string]types.PolicyResult),
	}

	err := s.run()
	if err != nil {
		s.sendError(err)
		return nil, err
	}

	return s, nil
}

func (s *CSite) run() error {
	go s.processIncoming()

	err := s.runHandshake()
	if err != nil {
		return err
	}

	return nil
}

func (s *CSite) runHandshake() error {
	helloB := infra.Encode(types.Hello{
		Protocol: "pestcontrol",
		Version:  1,
	})

	_, err := s.c.Write(helloB)
	if err != nil {
		return err
	}

	select {
	case helloReply := <-s.helloChan:
		if helloReply.Protocol != "pestcontrol" || helloReply.Version != 1 {
			return ErrInvalidHandshake
		}
	case <-s.okChan:
		return ErrInvalidHandshake
	case <-s.policyResultChan:
		return ErrInvalidHandshake
	case <-s.targetPopChan:
		return ErrInvalidHandshake
	}

	return nil
}

func (s *CSite) GetSite() uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()

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
			if res.Error != nil {
				if errors.Is(res.Error, infra.ErrNotEnough) {
					break
				}
				s.sendError(res.Error)
				return
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
			default:
				s.sendError(ErrInvalidData)
				return
			}

			curr = res.Next
		}
	}
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
	log.Printf("%v pop limits updated %v\n", s.site, s.targetPop.Populations)

	return s.targetPop, nil
}

// update policy.
// also ensures that there is only 1 policy in place
func (s *CSite) UpdatePolicy(pol types.CreatePolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	prev, ok := s.policies[pol.Species]
	if ok {
		log.Printf("%v policy for %v detected (%v)\n", s.site, pol.Species, prev.Policy)
		_, err := s.deletePolicy(types.DeletePolicy(prev), pol.Species)
		if err != nil {
			return err
		}
		delete(s.policies, pol.Species)
	} else {
		log.Printf("%v policy for %v not detected\n", s.site, pol.Species)
	}

	if pol.Action != types.PolicyNothing {
		ret, err := s.createPolicy(pol)
		if err != nil {
			return err
		}
		s.policies[pol.Species] = ret
	}

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

	log.Printf("%v created policy for %v: %v (%v)\n", s.site, pol.Species, action, ret.Policy)
	return
}

func (s *CSite) deletePolicy(pol types.DeletePolicy, species string) (ret types.OK, err error) {
	polB := infra.Encode(pol)
	_, err = s.c.Write(polB)
	if err != nil {
		return
	}

	ret = <-s.okChan
	log.Printf("%v policy for %v (%v) deleted\n", s.site, species, pol.Policy)
	return
}

func (s *CSite) close() {
	s.c.Close()
	close(s.helloChan)
	close(s.okChan)
	close(s.targetPopChan)
	close(s.policyResultChan)
}

func (c *CSite) sendError(err error) {
	log.Println("sent error auth")
	errorB := infra.Encode(types.Error{
		Message: err.Error(),
	})
	c.c.Write(errorB)
}
