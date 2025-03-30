package infra

import (
	"reflect"
	"testing"
)

func TestU32(t *testing.T) {
	p := parseUint32([]byte{0xa6, 0xa9, 0xb5, 0x67})
	exp := uint32(2796139879)
	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestString(t *testing.T) {
	p := parseString([]byte{0x00, 0x00, 0x00, 0x03, 0x66, 0x6f, 0x6f})
	exp := "foo"
	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestArray(t *testing.T) {
	parser := parseArray(parseUint32)

	p := parser([]byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x0b})
	exp := []uint32{10, 11}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestHello(t *testing.T) {
	p := parseHello([]byte{
		0x50,
		0x00, 0x00, 0x00, 0x19,
		0x00, 0x00, 0x00, 0x0b,
		0x70, 0x65, 0x73, 0x74,
		0x63, 0x6f, 0x6e, 0x74,
		0x72, 0x6f, 0x6c,
		0x00, 0x00, 0x00, 0x01,
		0xce,
	})

	exp := Hello{
		Protocol: "pestcontrol",
		Version:  1,
	}

	if !p.Ok {
		t.Fatalf("failed to parse")
	}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestError(t *testing.T) {
	p := parseError([]byte{
		0x51,
		0x00, 0x00, 0x00, 0x0d,
		0x00, 0x00, 0x00, 0x03,
		0x62, 0x61, 0x64,
		0x78,
	})

	exp := Error{
		Message: "bad",
	}

	if !p.Ok {
		t.Fatalf("failed to parse")
	}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestOK(t *testing.T) {
	p := parseOk([]byte{
		0x52,
		0x00, 0x00, 0x00, 0x06,
		0xa8,
	})

	exp := OK{}

	if !p.Ok {
		t.Fatalf("failed to parse")
	}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestDialAuthority(t *testing.T) {
	p := parseDialAuthority([]byte{
		0x53,
		0x00, 0x00, 0x00, 0x0a,
		0x00, 0x00, 0x30, 0x39,
		0x3a,
	})

	exp := DialAuthority{
		Site: 12345,
	}

	if !p.Ok {
		t.Fatalf("failed to parse")
	}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestTargetPopulations(t *testing.T) {
	p := parseTargetPopulations([]byte{
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
	})

	exp := TargetPopulations{
		Site: 12345,
		Populations: []TargetPopulationsEntry{
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
	}

	if !p.Ok {
		t.Fatalf("failed to parse")
	}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestCreatePolicy(t *testing.T) {
	p := parseCreatePolicy([]byte{
		0x55,
		0x00, 0x00, 0x00, 0x0e,
		0x00, 0x00, 0x00, 0x03,
		0x64, 0x6f, 0x67,
		0xa0,
		0xc0,
	})

	exp := CreatePolicy{
		Species: "dog",
		Action:  PolicyConserve,
	}

	if !p.Ok {
		t.Fatalf("failed to parse")
	}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestDeletePolicy(t *testing.T) {
	p := parseDeletePolicy([]byte{
		0x56,
		0x00, 0x00, 0x00, 0x0a,
		0x00, 0x00, 0x00, 0x7b,
		0x25,
	})

	exp := DeletePolicy{
		Policy: 123,
	}

	if !p.Ok {
		t.Fatalf("failed to parse")
	}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestPolicyResult(t *testing.T) {
	p := parsePolicyResult([]byte{
		0x57,
		0x00, 0x00, 0x00, 0x0a,
		0x00, 0x00, 0x00, 0x7b,
		0x24,
	})

	exp := PolicyResult{
		Policy: 123,
	}

	if !p.Ok {
		t.Fatalf("failed to parse")
	}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestSiteVisit(t *testing.T) {
	p := parseSiteVisit([]byte{
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
	})

	exp := SiteVisit{
		Site: 12345,
		Populations: []SiteVisitEntry{
			{
				Species: "dog",
				Count:   1,
			},
			{
				Species: "rat",
				Count:   5,
			},
		},
	}

	if !p.Ok {
		t.Fatalf("failed to parse")
	}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}
