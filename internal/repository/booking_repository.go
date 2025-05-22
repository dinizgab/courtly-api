package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/jackc/pgx/v5"
)

type (
	BookingRepository interface {
		Create(ctx context.Context, booking entity.Booking) (string, error) 
		FindByID(ctx context.Context, id string) (entity.Booking, error)
        FindByIDShowcase(ctx context.Context, id string) (entity.Booking, error)
		ListByCompanyID(ctx context.Context, companyId string) ([]entity.Booking, error)
        ConfirmBooking(ctx context.Context, companyId string, bookingId string) error
		Update(ctx context.Context, booking entity.Booking) error
		Delete(ctx context.Context, id string) error
	}

	bookingRepositoryImpl struct {
		db database.Database
	}
)

var (
	//go:embed sql/booking/create_booking.sql
	createBookingQuery string
	//go:embed sql/booking/list_bookings_by_company_id.sql
	listBookingsByCompanyIDQuery string
	//go:embed sql/booking/find_booking_by_id.sql
	findBookingByIDQuery string
    //go:embed sql/booking/find_booking_by_id_in_showcase.sql
    findBookingByIDShowcaseQuery string
    //go:embed sql/booking/confirm_booking.sql
    confirmBookingQuery string
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

func (r *bookingRepositoryImpl) Create(ctx context.Context, booking entity.Booking) (string, error) {
	row := r.db.QueryRow(
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
        booking.TotalPrice,
        booking.Court.CompanyId,
	)

    var id string

    err := row.Scan(&id)
	if err != nil {
        return "", fmt.Errorf("BookingRepository.Create - error scanning row: %w", err)
	}

	return id, nil
}

func (r *bookingRepositoryImpl) ListByCompanyID(ctx context.Context, companyId string) ([]entity.Booking, error) {
	bookings := make([]entity.Booking, 0)

	rows, err := r.db.Query(
		ctx,
		listBookingsByCompanyIDQuery,
		companyId,
	)
	if err != nil {
		return bookings, err 
	}
	defer rows.Close()

	for rows.Next() {
		var booking entity.Booking
		var court entity.Court

		err := rows.Scan(
			&booking.ID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.CreatedAt,
			&booking.Status,
			&booking.GuestName,
			&booking.GuestPhone,
			&booking.GuestEmail,
			&court.Name,
		)
		if err != nil {
			return []entity.Booking{}, fmt.Errorf("BookingRepository.ListByCompanyID - error scanning rows: %w", err)
		}

        booking.Court = &court
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CourtRepository.ListBookingsByID - error in rows: %w", err)
	}

	return bookings, nil
}

func (r *bookingRepositoryImpl) FindByID(ctx context.Context, id string) (entity.Booking, error) {
	var booking entity.Booking
    var court entity.Court
	err := r.db.QueryRow(ctx, findBookingByIDQuery, id).Scan(
		&booking.ID,
		&booking.CourtId,
		&booking.StartTime,
		&booking.EndTime,
        &booking.CreatedAt,
		&booking.Status,
		&booking.GuestName,
		&booking.GuestPhone,
		&booking.GuestEmail,
		&booking.VerificationCode,
        &booking.TotalPrice,
        &court.Name,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Booking{}, fmt.Errorf("BookingRepository.FindByID: booking not found")
		}
		return entity.Booking{}, fmt.Errorf("BookingRepository.FindByID: %w", err)
	}

    booking.Court = &court

	return booking, nil
}

func (r *bookingRepositoryImpl) FindByIDShowcase(ctx context.Context, id string) (entity.Booking, error) {
    var booking entity.Booking
    var court entity.Court
    var company entity.Company

    err := r.db.QueryRow(ctx, findBookingByIDShowcaseQuery, id).Scan(
        &booking.StartTime,
        &booking.EndTime,
        &booking.TotalPrice,
        &court.Name,
        &company.Address,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return entity.Booking{}, fmt.Errorf("BookingRepository.FindByID: booking not found")
        }
        return entity.Booking{}, fmt.Errorf("BookingRepository.FindByID: %w", err)
    }

    court.Company = &company
    booking.Court = &court

    return booking, nil
}

func (r *bookingRepositoryImpl) ConfirmBooking(ctx context.Context, companyId string, bookingId string) error {
    _, err := r.db.Exec(ctx, confirmBookingQuery, bookingId, companyId)
    if err != nil {
        return fmt.Errorf("BookingRepository.ConfirmBooking: %w", err)
    }

    return nil
}

func (r *bookingRepositoryImpl) Update(ctx context.Context, booking entity.Booking) error {
	return nil
}

func (r *bookingRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, deleteBookingQuery, id)
	if err != nil {
		return fmt.Errorf("BookingRepository.Delete: %w", err)
	}

	return nil
}
