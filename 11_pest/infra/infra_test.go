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
	expErr error
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
		expErr: nil,
	},
	{
		inB: []byte{
			80, 0, 0, 0, 25, 0, 0, 0, 111, 112, 101, 115, 116, 99, 111, 110, 116, 114, 111, 108, 0, 0, 0, 1, 106,
		},
		expVal: nil,
		expErr: infra.ErrTooLong,
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
		expErr: nil,
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
		expErr: nil,
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
		expErr: nil,
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
		expErr: nil,
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
		expErr: nil,
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
		expErr: nil,
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
		expErr: nil,
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
		expErr: nil,
	},
	{
		inB: []byte{
			88,
			0, 0, 0, 60,
			0, 0, 48, 57,
			0, 0, 0, 102, 0, 0, 0, 15, 108, 111, 110, 103, 45, 116, 97, 105, 108, 101, 100, 32, 114, 97, 116, 0, 0, 0, 50, 0, 0, 0, 15, 98, 105, 103, 45, 119, 105, 110, 103, 101, 100, 32, 98, 105, 114, 100, 0, 0, 0, 3, 245,
		},
		expVal: nil,
		expErr: infra.ErrTooLong,
	},
}

func TestParser(t *testing.T) {
	runParseCase := func(t *testing.T, c parseCase) {
		res := infra.Parse(c.inB)
		if c.expErr != res.Error {
			t.Fatalf("parse status wrong. exp %v got %v", c.expErr, res.Error)
		}
		if c.expErr == nil {
			if !reflect.DeepEqual(res.Value, c.expVal) {
				t.Fatalf("wrong parse result. exp %v got %v", c.expVal, res.Value)
			}
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
		if c.expErr != nil {
			return
		}
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
