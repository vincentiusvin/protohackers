package infra

type Plate struct {
	Plate     string
	Timestamp uint32
}

type IAmACamera struct {
	Road  uint16
	Mile  uint16
	Limit uint16
}

type IAmADispatcher struct {
	Roads []uint16
}

type WantHeartbeat struct {
	Interval uint32
}

type Ticket struct {
	Plate      string
	Road       uint16
	Mile1      uint16
	Timestamp1 uint32
	Mile2      uint16
	Timestamp2 uint32
	Speed      uint16
}
