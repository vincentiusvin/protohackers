package pest

import (
	"fmt"
	"net"
	"protohackers/11_pest/infra"
	"protohackers/11_pest/types"
)

type Site interface {
	Connect() error
	Close()
	GetPops() (types.TargetPopulations, error)
	CreatePolicy(types.CreatePolicy) (types.PolicyResult, error)
	DeletePolicy(types.DeletePolicy) (types.OK, error)
}

type CSite struct {
	site uint32
	c    net.Conn

	// cached value
	targetPop   types.TargetPopulations
	targetPopOK bool

	// authority -> us
	helloChan        chan types.Hello
	okChan           chan types.OK
	targetPopChan    chan types.TargetPopulations
	policyResultChan chan types.PolicyResult
}

func NewSite(site uint32) Site {
	s := &CSite{
		site:             site,
		helloChan:        make(chan types.Hello),
		okChan:           make(chan types.OK),
		targetPopChan:    make(chan types.TargetPopulations),
		policyResultChan: make(chan types.PolicyResult),
	}

	return s
}

func (s *CSite) Connect() error {
	conn, err := net.Dial("tcp", "pestcontrol.protohackers.com:20547")
	if err != nil {
		return err
	}

	s.c = conn
	go s.processIncoming()
	s.handshake()
	return nil
}

func (s *CSite) processIncoming() {
	defer s.Close()

	var curr []byte
	for {
		b := make([]byte, 1024)
		_, err := s.c.Read(b)
		curr = append(curr, b...)
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

	helloReply := <-s.helloChan
	if helloReply.Protocol != "pestcontrol" || helloReply.Version != 1 {
		return fmt.Errorf("got invalid handshake reply %v", helloReply)
	}

	return nil
}

func (s *CSite) GetPops() (ret types.TargetPopulations, err error) {
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

	return s.targetPop, nil
}

func (s *CSite) CreatePolicy(pol types.CreatePolicy) (ret types.PolicyResult, err error) {
	polB := infra.Encode(pol)
	_, err = s.c.Write(polB)
	if err != nil {
		return
	}

	ret = <-s.policyResultChan
	return
}

func (s *CSite) DeletePolicy(pol types.DeletePolicy) (ret types.OK, err error) {
	polB := infra.Encode(pol)
	_, err = s.c.Write(polB)
	if err != nil {
		return
	}

	ret = <-s.okChan
	return
}

func (s *CSite) Close() {
	s.c.Close()
	close(s.helloChan)
	close(s.okChan)
	close(s.targetPopChan)
	close(s.policyResultChan)
}
