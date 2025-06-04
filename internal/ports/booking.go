package ports

import (
	"context"

	"github.com/dinizgab/booking-mvp/internal/entity"
)

type BookingSummaryReader interface {
	GetBookingSummary(ctx context.Context, bookingId string) (entity.Booking, error)
}
