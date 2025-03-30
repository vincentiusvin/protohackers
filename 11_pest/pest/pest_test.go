package pest_test

import (
	"protohackers/11_pest/pest"
	"protohackers/11_pest/types"
	"reflect"
	"testing"
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
	}

	for _, cs := range cases {
		t.Run("visit", func(t *testing.T) {
			var sitenum uint32 = 12345
			s := newMockSite(sitenum)

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

			err := c.AddSiteVisit(sv1)
			if err != nil {
				t.Fatal(err)
			}

			outPol := <-s.policies
			expPol := types.CreatePolicy{
				Species: "kucing",
				Action:  cs.action,
			}

			if !reflect.DeepEqual(expPol, outPol) {
				t.Fatalf("wrong policy. exp %v got %v", expPol, outPol)
			}
		})
	}

}

type mockSite struct {
	site     uint32
	policies chan types.CreatePolicy
}

func newMockSite(site uint32) *mockSite {
	return &mockSite{
		site:     site,
		policies: make(chan types.CreatePolicy, 1),
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
		},
	}, nil
}

func (ms *mockSite) UpdatePolicy(pol types.CreatePolicy) error {
	ms.policies <- pol
	return nil
}
