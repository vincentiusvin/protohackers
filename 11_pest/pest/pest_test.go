package pest_test

import (
	"log"
	"protohackers/11_pest/pest"
	"protohackers/11_pest/types"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func TestSiteVisit(t *testing.T) {
	type visitCases struct {
		inCount uint32
		action  types.Policy
	}

	cases := []visitCases{
		{
			inCount: 250,
			action:  types.PolicyCull,
		},
		{
			inCount: 5,
			action:  types.PolicyConserve,
		},
		{
			inCount: 15,
			action:  types.PolicyNothing,
		},
	}

	for _, cs := range cases {
		t.Run("visit", func(t *testing.T) {
			var sitenum uint32 = 12345
			s := newMockSite(sitenum, 0)

			factory := func(site uint32) (pest.Site, error) {
				return s, nil
			}
			c := pest.NewController(factory)

			sv1 := types.SiteVisit{
				Site: sitenum,
				Populations: []types.SiteVisitEntry{
					{
						Species: "kucing",
						Count:   cs.inCount,
					},
				},
			}

			retval := make(chan error)
			go func() {
				retval <- c.AddSiteVisit(sv1)
			}()

			outPol := <-s.policies
			expPol := types.CreatePolicy{
				Species: "kucing",
				Action:  cs.action,
			}

			if !reflect.DeepEqual(expPol, outPol) {
				t.Fatalf("wrong policy. exp %v got %v", expPol, outPol)
			}

			err := <-retval
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSubsequent(t *testing.T) {
	var sitenum uint32 = 12345
	s := newMockSite(sitenum, 0)

	factory := func(site uint32) (pest.Site, error) {
		return s, nil
	}
	c := pest.NewController(factory)

	svs := []types.SiteVisit{
		{
			Site: sitenum,
			Populations: []types.SiteVisitEntry{
				{
					Species: "kucing",
					Count:   200,
				},
			},
		},
		{
			Site: sitenum,
			Populations: []types.SiteVisitEntry{
				{
					Species: "kucing",
					Count:   199,
				},
			},
		},
		{
			Site: sitenum,
			Populations: []types.SiteVisitEntry{
				{
					Species: "kucing",
					Count:   198,
				},
			},
		},
		{
			Site: sitenum,
			Populations: []types.SiteVisitEntry{
				{
					Species: "kucing",
					Count:   197,
				},
			},
		},
		{
			Site: sitenum,
			Populations: []types.SiteVisitEntry{
				{
					Species: "kucing",
					Count:   1,
				},
			},
		},
	}

	go func() {
		for _, sv := range svs {
			c.AddSiteVisit(sv)
		}
	}()

	for i := 0; i < 4; i++ {
		outPol := <-s.policies
		expPol := types.CreatePolicy{
			Species: "kucing",
			Action:  types.PolicyCull,
		}

		if !reflect.DeepEqual(expPol, outPol) {
			t.Fatalf("wrong policy. exp %v got %v", expPol, outPol)
		}
	}
	for i := 0; i < 1; i++ {
		outPol := <-s.policies
		expPol := types.CreatePolicy{
			Species: "kucing",
			Action:  types.PolicyConserve,
		}

		if !reflect.DeepEqual(expPol, outPol) {
			t.Fatalf("wrong policy. exp %v got %v", expPol, outPol)
		}
	}
}

func TestConccurentCreation(t *testing.T) {
	var sitenum uint32 = 12345
	var called atomic.Int32
	s := newMockSite(sitenum, 0)
	factory := func(site uint32) (pest.Site, error) {
		time.Sleep(50 * time.Millisecond)
		called.Add(1)
		return s, nil
	}
	c := pest.NewController(factory)

	svs := []types.SiteVisit{
		{
			Site: sitenum,
			Populations: []types.SiteVisitEntry{
				{
					Species: "kucing",
					Count:   200,
				},
			},
		},
		{
			Site: sitenum,
			Populations: []types.SiteVisitEntry{
				{
					Species: "anjing",
					Count:   199,
				},
			},
		},
	}

	go func() {
		for range svs {
			<-s.policies
		}
	}()

	for _, sv := range svs {
		err := c.AddSiteVisit(sv)
		if err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(200 * time.Millisecond)
	log.Println(called.Load())
	if called.Load() != 1 {
		t.Fatal("expected factory function to be called once")
	}

}

type mockSite struct {
	site     uint32
	policies chan types.CreatePolicy
}

func newMockSite(site uint32, chanlen int) *mockSite {
	return &mockSite{
		site:     site,
		policies: make(chan types.CreatePolicy, chanlen),
	}
}

func (ms *mockSite) GetSite() uint32 {
	return ms.site
}

func (ms *mockSite) GetPops() (types.TargetPopulations, error) {
	return types.TargetPopulations{
		Site: ms.site,
		Populations: []types.TargetPopulationsEntry{
			{
				Species: "kucing",
				Min:     10,
				Max:     20,
			},
			{
				Species: "anjing",
				Min:     10,
				Max:     20,
			},
		},
	}, nil
}

func (ms *mockSite) UpdatePolicy(pol types.CreatePolicy) error {
	ms.policies <- pol
	return nil
}
