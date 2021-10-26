package sender

import (
	"github.com/zinchenkom/loy-point-api/internal/model"
)

type EventSender interface {
	Send(subdomain *model.SubdomainEvent) error
}
