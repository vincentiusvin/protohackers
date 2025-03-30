package infra_test

import (
	"protohackers/11_pest/infra"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	type parseCase struct {
		inB    []byte
		expVal any
		expOk  bool
	}

	runParseCase := func(t *testing.T, c parseCase) {
		res := infra.Parse(c.inB)
		if c.expOk != res.Ok {
			t.Fatalf("parse status wrong. exp %v got %v", c.expOk, res.Ok)
		}
		if !reflect.DeepEqual(res.Value, c.expVal) {
			t.Fatalf("wrong parse result. exp %v got %v", c.expVal, res.Value)
		}
	}

	t.Run("Hello", func(t *testing.T) {
		cases := []parseCase{
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
				expVal: infra.Hello{
					Protocol: "pestcontrol",
					Version:  1,
				},
				expOk: true,
			},
		}
		for _, c := range cases {
			runParseCase(t, c)
		}
	})

	t.Run("Error", func(t *testing.T) {
		cases := []parseCase{
			{
				inB: []byte{
					0x51,
					0x00, 0x00, 0x00, 0x0d,
					0x00, 0x00, 0x00, 0x03,
					0x62, 0x61, 0x64,
					0x78,
				},
				expVal: infra.Error{
					Message: "bad",
				},
				expOk: true,
			},
		}
		for _, c := range cases {
			runParseCase(t, c)
		}
	})

	t.Run("OK", func(t *testing.T) {
		cases := []parseCase{
			{
				inB: []byte{
					0x52,
					0x00, 0x00, 0x00, 0x06,
					0xa8,
				},
				expVal: infra.OK{},
				expOk:  true,
			},
		}
		for _, c := range cases {
			runParseCase(t, c)
		}
	})

	t.Run("DialAuthority", func(t *testing.T) {
		cases := []parseCase{
			{
				inB: []byte{
					0x53,
					0x00, 0x00, 0x00, 0x0a,
					0x00, 0x00, 0x30, 0x39,
					0x3a,
				},
				expVal: infra.DialAuthority{
					Site: 12345,
				},
				expOk: true,
			},
		}
		for _, c := range cases {
			runParseCase(t, c)
		}
	})

	t.Run("TargetPopulations", func(t *testing.T) {
		cases := []parseCase{
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
				expVal: infra.TargetPopulations{
					Site: 12345,
					Populations: []infra.TargetPopulationsEntry{
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
		for _, c := range cases {
			runParseCase(t, c)
		}
	})

	t.Run("CreatePolicy", func(t *testing.T) {
		cases := []parseCase{
			{
				inB: []byte{
					0x55,
					0x00, 0x00, 0x00, 0x0e,
					0x00, 0x00, 0x00, 0x03,
					0x64, 0x6f, 0x67,
					0xa0,
					0xc0,
				},
				expVal: infra.CreatePolicy{
					Species: "dog",
					Action:  infra.PolicyConserve,
				},
				expOk: true,
			},
		}
		for _, c := range cases {
			runParseCase(t, c)
		}
	})

	t.Run("DeletePolicy", func(t *testing.T) {
		cases := []parseCase{
			{
				inB: []byte{
					0x56,
					0x00, 0x00, 0x00, 0x0a,
					0x00, 0x00, 0x00, 0x7b,
					0x25,
				},
				expVal: infra.DeletePolicy{
					Policy: 123,
				},
				expOk: true,
			},
		}
		for _, c := range cases {
			runParseCase(t, c)
		}
	})

	t.Run("PolicyResult", func(t *testing.T) {
		cases := []parseCase{
			{
				inB: []byte{
					0x57,
					0x00, 0x00, 0x00, 0x0a,
					0x00, 0x00, 0x00, 0x7b,
					0x24,
				},
				expVal: infra.PolicyResult{
					Policy: 123,
				},
				expOk: true,
			},
		}
		for _, c := range cases {
			runParseCase(t, c)
		}
	})

	t.Run("SiteVisit", func(t *testing.T) {
		cases := []parseCase{
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
				expVal: infra.SiteVisit{
					Site: 12345,
					Populations: []infra.SiteVisitEntry{
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
		for _, c := range cases {
			runParseCase(t, c)
		}
	})

}
