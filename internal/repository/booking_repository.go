package repository

import (
	"context"
	"fmt"
    _ "embed"

	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/jackc/pgx/v5"
)

type (
	BookingRepository interface {
		Create(ctx context.Context, booking *entity.Booking) error
        FindByID(ctx context.Context, id string) (*entity.Booking, error)
        Update(ctx context.Context, booking *entity.Booking) error
        Delete(ctx context.Context, id string) error
	}

	bookingRepositoryImpl struct {
		db database.Database
	}
)

var (
	//go:embed sql/booking/create_booking.sql
	createBookingQuery string
    //go:embed sql/booking/find_booking_by_id.sql
    findBookingByIDQuery string
    //go:embed sql/booking/update_booking.sql
    updateBookingQuery string
    //go:embed sql/booking/delete_booking.sql
    deleteBookingQuery string
)

func NewBookingRepository(db database.Database) BookingRepository {
	return &bookingRepositoryImpl{
		db: db,
	}
}

func (r *bookingRepositoryImpl) Create(ctx context.Context, booking *entity.Booking) error {
	_, err := r.db.Exec(
		ctx,
		createBookingQuery,
		booking.CourtId,
		booking.StartTime,
		booking.EndTime,
		booking.GuestName,
		booking.GuestEmail,
		booking.GuestPhone,
		booking.Status,
		booking.VerificationCode,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *bookingRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Booking, error) {
    var booking entity.Booking
    err := r.db.QueryRow(ctx, findBookingByIDQuery, id).Scan(
        &booking.ID,
        &booking.CourtId,
        &booking.StartTime,
        &booking.EndTime,
        &booking.GuestName,
        &booking.GuestEmail,
        &booking.GuestPhone,
        &booking.Status,
        &booking.VerificationCode,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("BookingRepository.FindByID: booking not found")
        }
        return nil, fmt.Errorf("BookingRepository.FindByID: %w", err)
    }

    return &booking, nil
}

func (r *bookingRepositoryImpl) Update(ctx context.Context, booking *entity.Booking) error {
    return nil
}

func (r *bookingRepositoryImpl) Delete(ctx context.Context, id string) error {
    _, err := r.db.Exec(ctx, deleteBookingQuery, id)
    if err != nil {
        return fmt.Errorf("BookingRepository.Delete: %w", err)
    }

    return nil
}
