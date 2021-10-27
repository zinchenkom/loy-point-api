package repo

import (
	"github.com/zinchenkom/loy-point-api/internal/model"
)

type EventRepo interface {
	Lock(n uint64) ([]loyalty.PointEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []loyalty.PointEvent) error
	Remove(eventIDs []uint64) error
}
