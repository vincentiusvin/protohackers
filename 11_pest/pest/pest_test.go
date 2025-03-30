package pest_test

import (
	"fmt"
	"io"
	"protohackers/11_pest/infra"
	"protohackers/11_pest/pest"
	"protohackers/11_pest/types"
	"reflect"
	"testing"
)

func TestGetPops(t *testing.T) {
	var siteNum uint32 = 12345
	s, in, out, err := fixture(siteNum)
	if err != nil {
		t.Fatal(err)
	}

	popRes := make(chan types.TargetPopulations)
	go func() {
		res, err := s.GetPops()
		if err != nil {
			close(popRes)
		}
		popRes <- res
	}()

	popsOut := <-out
	expPops := types.DialAuthority{
		Site: siteNum,
	}
	if !reflect.DeepEqual(expPops, popsOut) {
		t.Fatalf("wrong out exp %v got %v", expPops, popsOut)
	}

	target := types.TargetPopulations{
		Site: siteNum,
		Populations: []types.TargetPopulationsEntry{
			{
				Species: "kucing",
				Min:     10,
				Max:     20,
			},
		},
	}

	in <- target
	result := <-popRes
	if !reflect.DeepEqual(target, result) {
		t.Fatalf("different pops. exp %v got %v", target, result)
	}
}

func TestUpdatePolicy(t *testing.T) {
	var siteNum uint32 = 12345
	s, _, out, err := fixture(siteNum)
	if err != nil {
		t.Fatal(err)
	}

	go s.UpdatePolicy(types.CreatePolicy{
		Species: "kucing",
		Action:  types.PolicyCull,
	})

	popsOut := <-out
	expPops := types.CreatePolicy{
		Species: "kucing",
		Action:  types.PolicyCull,
	}
	if !reflect.DeepEqual(expPops, popsOut) {
		t.Fatalf("wrong out exp %v got %v", expPops, popsOut)
	}
}

func fixture(site uint32) (s pest.Site, in chan any, out chan any, err error) {
	var rw io.ReadWriteCloser
	rw, in, out = createRW()

	sch := make(chan pest.Site)
	go func() {
		s := pest.NewSite(site, rw)
		sch <- s
	}()

	helloOut := <-out
	expHello := types.Hello{
		Protocol: "pestcontrol",
		Version:  1,
	}

	if !reflect.DeepEqual(expHello, helloOut) {
		err = fmt.Errorf("wrong out exp %v got %v", expHello, helloOut)
		return
	}

	in <- types.Hello{
		Protocol: "pestcontrol",
		Version:  1,
	}

	s = <-sch
	return
}

func createRW() (rw io.ReadWriteCloser, in chan any, out chan any) {
	in = make(chan any, 1)
	out = make(chan any, 1)

	inr, inw := io.Pipe()

	go func() {
		defer inw.Close()
		for {
			d := <-in
			b := infra.Encode(d)
			_, err := inw.Write(b)
			if err != nil {
				return
			}
		}
	}()

	outr, outw := io.Pipe()
	go func() {
		var curr []byte
		for {
			b := make([]byte, 1024)
			n, err := outr.Read(b)
			curr = append(curr, b[:n]...)
			if err != nil {
				break
			}

			for {
				res := infra.Parse(curr)
				if !res.Ok {
					break
				}
				out <- res.Value
				curr = res.Next
			}
		}
	}()

	rw = &mockRW{
		r: inr,
		w: outw,
	}

	return
}

type mockRW struct {
	r io.Reader
	w io.Writer
}

func (rw *mockRW) Read(b []byte) (n int, err error) {
	return rw.r.Read(b)
}

func (rw *mockRW) Write(b []byte) (n int, err error) {
	return rw.w.Write(b)
}

func (rw *mockRW) Close() error {
	return nil
}
