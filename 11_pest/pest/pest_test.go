package pest_test

import (
	"fmt"
	"protohackers/11_pest/pest"
	"testing"
)

func TestSite(t *testing.T) {
	s := pest.NewSite(12345)
	err := s.Connect()
	if err != nil {
		t.Fatal(err)
	}

	pops, err := s.GetPops()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(pops)
}
