package infra_test

import (
	"protohackers/11_pest/infra"
	"protohackers/11_pest/types"
	"reflect"
	"testing"
)

type parseCase struct {
	inB    []byte
	expVal any
	expOk  bool
}

var helloCases = []parseCase{
	{
		inB: []byte{
			0x50,
			0x00, 0x00, 0x00, 0x19,
			0x00, 0x00, 0x00, 0x0b,
			0x70, 0x65, 0x73, 0x74,
			0x63, 0x6f, 0x6e, 0x74,
			0x72, 0x6f, 0x6c,
			0x00, 0x00, 0x00, 0x01,
			0xce,
		},
		expVal: types.Hello{
			Protocol: "pestcontrol",
			Version:  1,
		},
		expOk: true,
	},
}

var errorCases = []parseCase{
	{
		inB: []byte{
			0x51,
			0x00, 0x00, 0x00, 0x0d,
			0x00, 0x00, 0x00, 0x03,
			0x62, 0x61, 0x64,
			0x78,
		},
		expVal: types.Error{
			Message: "bad",
		},
		expOk: true,
	},
}

var okCases = []parseCase{
	{
		inB: []byte{
			0x52,
			0x00, 0x00, 0x00, 0x06,
			0xa8,
		},
		expVal: types.OK{},
		expOk:  true,
	},
}

var dialAuthorityCases = []parseCase{
	{
		inB: []byte{
			0x53,
			0x00, 0x00, 0x00, 0x0a,
			0x00, 0x00, 0x30, 0x39,
			0x3a,
		},
		expVal: types.DialAuthority{
			Site: 12345,
		},
		expOk: true,
	},
}

var targetPopCases = []parseCase{
	{
		inB: []byte{
			0x54,
			0x00, 0x00, 0x00, 0x2c,
			0x00, 0x00, 0x30, 0x39,
			0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x00, 0x03,
			0x64, 0x6f, 0x67,
			0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x03,
			0x00, 0x00, 0x00, 0x03,
			0x72, 0x61, 0x74,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x0a,
			0x80,
		},
		expVal: types.TargetPopulations{
			Site: 12345,
			Populations: []types.TargetPopulationsEntry{
				{
					Species: "dog",
					Min:     1,
					Max:     3,
				},
				{
					Species: "rat",
					Min:     0,
					Max:     10,
				},
			},
		},
		expOk: true,
	},
}

var createPolicyCases = []parseCase{
	{
		inB: []byte{
			0x55,
			0x00, 0x00, 0x00, 0x0e,
			0x00, 0x00, 0x00, 0x03,
			0x64, 0x6f, 0x67,
			0xa0,
			0xc0,
		},
		expVal: types.CreatePolicy{
			Species: "dog",
			Action:  types.PolicyConserve,
		},
		expOk: true,
	},
}

var deletePolicyCases = []parseCase{
	{
		inB: []byte{
			0x56,
			0x00, 0x00, 0x00, 0x0a,
			0x00, 0x00, 0x00, 0x7b,
			0x25,
		},
		expVal: types.DeletePolicy{
			Policy: 123,
		},
		expOk: true,
	},
}

var policyResultCases = []parseCase{
	{
		inB: []byte{
			0x57,
			0x00, 0x00, 0x00, 0x0a,
			0x00, 0x00, 0x00, 0x7b,
			0x24,
		},
		expVal: types.PolicyResult{
			Policy: 123,
		},
		expOk: true,
	},
}

var siteVisitCases = []parseCase{
	{
		inB: []byte{
			0x58,
			0x00, 0x00, 0x00, 0x24,
			0x00, 0x00, 0x30, 0x39,
			0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x00, 0x03,
			0x64, 0x6f, 0x67,
			0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x03,
			0x72, 0x61, 0x74,
			0x00, 0x00, 0x00, 0x05,
			0x8c,
		},
		expVal: types.SiteVisit{
			Site: 12345,
			Populations: []types.SiteVisitEntry{
				{
					Species: "dog",
					Count:   1,
				},
				{
					Species: "rat",
					Count:   5,
				},
			},
		},
		expOk: true,
	},
}

func TestParser(t *testing.T) {
	runParseCase := func(t *testing.T, c parseCase) {
		res := infra.Parse(c.inB)
		ok := res.Error == nil
		if c.expOk != ok {
			t.Fatalf("parse status wrong. exp %v got %v", c.expOk, res.Error)
		}
		if !reflect.DeepEqual(res.Value, c.expVal) {
			t.Fatalf("wrong parse result. exp %v got %v", c.expVal, res.Value)
		}
	}

	t.Run("Hello", func(t *testing.T) {
		for _, c := range helloCases {
			runParseCase(t, c)
		}
	})

	t.Run("Error", func(t *testing.T) {
		for _, c := range errorCases {
			runParseCase(t, c)
		}
	})

	t.Run("OK", func(t *testing.T) {
		for _, c := range okCases {
			runParseCase(t, c)
		}
	})

	t.Run("DialAuthority", func(t *testing.T) {
		for _, c := range dialAuthorityCases {
			runParseCase(t, c)
		}
	})

	t.Run("TargetPopulations", func(t *testing.T) {
		for _, c := range targetPopCases {
			runParseCase(t, c)
		}
	})

	t.Run("CreatePolicy", func(t *testing.T) {
		for _, c := range createPolicyCases {
			runParseCase(t, c)
		}
	})

	t.Run("DeletePolicy", func(t *testing.T) {
		for _, c := range deletePolicyCases {
			runParseCase(t, c)
		}
	})

	t.Run("PolicyResult", func(t *testing.T) {
		for _, c := range policyResultCases {
			runParseCase(t, c)
		}
	})

	t.Run("SiteVisit", func(t *testing.T) {
		for _, c := range siteVisitCases {
			runParseCase(t, c)
		}
	})
}

func TestEncoder(t *testing.T) {
	runEncodeCase := func(t *testing.T, c parseCase) {
		res := infra.Encode(c.expVal)
		if !reflect.DeepEqual(res, c.inB) {
			t.Fatalf("wrong parse result. exp %v got %v", c.inB, res)
		}
	}

	t.Run("Hello", func(t *testing.T) {
		for _, c := range helloCases {
			runEncodeCase(t, c)
		}
	})

	t.Run("Error", func(t *testing.T) {
		for _, c := range errorCases {
			runEncodeCase(t, c)
		}
	})

	t.Run("OK", func(t *testing.T) {
		for _, c := range okCases {
			runEncodeCase(t, c)
		}
	})

	t.Run("DialAuthority", func(t *testing.T) {
		for _, c := range dialAuthorityCases {
			runEncodeCase(t, c)
		}
	})

	t.Run("TargetPopulations", func(t *testing.T) {
		for _, c := range targetPopCases {
			runEncodeCase(t, c)
		}
	})

	t.Run("CreatePolicy", func(t *testing.T) {
		for _, c := range createPolicyCases {
			runEncodeCase(t, c)
		}
	})

	t.Run("DeletePolicy", func(t *testing.T) {
		for _, c := range deletePolicyCases {
			runEncodeCase(t, c)
		}
	})

	t.Run("PolicyResult", func(t *testing.T) {
		for _, c := range policyResultCases {
			runEncodeCase(t, c)
		}
	})

	t.Run("SiteVisit", func(t *testing.T) {
		for _, c := range siteVisitCases {
			runEncodeCase(t, c)
		}
	})
}
