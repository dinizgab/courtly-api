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
	CourtRepository interface {
		Create(ctx context.Context, c *entity.Court) error
		FindByID(ctx context.Context, id string) (entity.Court, error)
		ListBookingsByID(ctx context.Context, id string) ([]entity.Booking, error)
		ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error)
		Update(ctx context.Context, c *entity.Court) error
		Delete(ctx context.Context, id string) error
	}

	courtRepositoryImpl struct {
		db database.Database
	}
)

var (
	//go:embed sql/court/create_court.sql
	createCourtQuery string
	//go:embed sql/court/find_court_by_id.sql
	findCourtByIDQuery string
	//go:embed sql/court/list_bookings_by_id.sql
	listBookingsByIDQuery string
	//go:embed sql/court/list_court_by_company.sql
	listCourtByCompanyQuery string
	//go:embed sql/court/delete_court.sql
	deleteCourtQuery string
)

func NewCourtRepository(db database.Database) CourtRepository {
	return &courtRepositoryImpl{
		db: db,
	}
}

func (r *courtRepositoryImpl) Create(ctx context.Context, c *entity.Court) error {
	_, err := r.db.Exec(
		ctx,
		createCourtQuery,
		c.CompanyId,
		c.Name,
		c.Description,
		c.SportType,
		c.HourlyPrice,
		c.IsActive,
		c.OpeningTime,
		c.ClosingTime,
		c.Capacity,
	)
	if err != nil {
		return fmt.Errorf("CourtRepository.Create: %w", err)
	}

	return nil
}

func (r *courtRepositoryImpl) FindByID(ctx context.Context, id string) (entity.Court, error) {
	var court entity.Court
	err := r.db.QueryRow(ctx, findCourtByIDQuery, id).Scan(
		&court.ID,
		&court.CompanyId,
		&court.Name,
		&court.Description,
		&court.SportType,
		&court.HourlyPrice,
		&court.IsActive,
		&court.OpeningTime,
		&court.ClosingTime,
		&court.Capacity,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Court{}, fmt.Errorf("CourtRepository.FindByID: court not found: %w", err)
		}

		return entity.Court{}, fmt.Errorf("CourtRepository.FindByID: %w", err)
	}

	return court, nil
}

func (r *courtRepositoryImpl) ListBookingsByID(ctx context.Context, id string) ([]entity.Booking, error) {
	rows, err := r.db.Query(ctx, listBookingsByIDQuery, id)
	if err != nil {
		return nil, fmt.Errorf("CourtRepository.ListBookingsByID: %w", err)
	}

	defer rows.Close()
	var bookings []entity.Booking
	for rows.Next() {
		var booking entity.Booking
		err := rows.Scan(
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
			return nil, fmt.Errorf("CourtRepository.ListBookingsByID: %w", err)
		}
		bookings = append(bookings, booking)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CourtRepository.ListBookingsByID: %w", err)
	}

	return bookings, nil
}

func (r *courtRepositoryImpl) ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error) {
	rows, err := r.db.Query(ctx, listCourtByCompanyQuery, companyID)
	if err != nil {
		return nil, fmt.Errorf("CourtRepository.ListByCompany: %w", err)
	}
	defer rows.Close()

	var courts []entity.Court
	for rows.Next() {
		var court entity.Court
		err := rows.Scan(
			&court.ID,
			&court.CompanyId,
			&court.Name,
			&court.SportType,
			&court.HourlyPrice,
			&court.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("CourtRepository.ListByCompany: %w", err)
		}
		courts = append(courts, court)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CourtRepository.ListByCompany: %w", err)
	}

	return courts, nil
}

func (r *courtRepositoryImpl) Update(ctx context.Context, c *entity.Court) error {
	return nil
}

func (r *courtRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, deleteCourtQuery, id)
	if err != nil {
		return fmt.Errorf("CourtRepository.Delete: %w", err)
	}

	return nil
}
