package infra

type Hello struct {
	Protocol string
	Version  uint32
}

type Error struct {
	Message string
}

type OK struct{}

type DialAuthority struct {
	Site uint32
}

type TargetPopulationsEntry struct {
	Species string
	Min     uint32
	Max     uint32
}

type TargetPopulations struct {
	Site        uint32
	Populations []TargetPopulationsEntry
}

type Policy uint8

var (
	PolicyCull     Policy = 0x90
	PolicyConserve Policy = 0xa0
)

type CreatePolicy struct {
	Species string
	Action  Policy
}

type DeletePolicy struct {
	Policy uint32
}

type PolicyResult struct {
	Policy uint32
}

type SiteVisitEntry struct {
	Species string
	Count   uint32
}

type SiteVisit struct {
	Site        uint32
	Populations []SiteVisitEntry
}
