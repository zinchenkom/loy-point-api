package loyalty

type Point struct {
	ID uint64
	Name        string
	Description string
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type PointEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Point
}
