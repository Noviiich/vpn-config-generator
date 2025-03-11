package ports

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/internal/core/domain"
)

type DevicePort interface {
	All(context.Context) ([]domain.Device, error)
	Create(context.Context, *domain.Device) error
	Delete(context.Context, string)
	Get(context.Context, string) (domain.Device, error)
}
