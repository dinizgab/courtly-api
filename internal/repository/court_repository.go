package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/jackc/pgx/v5"
)

const CourtPhotosName = "court_photos"

type (
	CourtRepository interface {
		Create(ctx context.Context, c *entity.Court) error
		InsertPhotos(ctx context.Context, c []entity.CourtPhoto) error
		FindByID(ctx context.Context, id string) (entity.Court, error)
		ListBookingsByID(ctx context.Context, id string) ([]entity.Booking, error)
		ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error)
		ListCompanyCourtsShowcase(ctx context.Context, companyID string) ([]entity.Court, error)
		ListAvailableBookingSlots(ctx context.Context, id string, date string) ([]entity.Booking, error)
		Update(ctx context.Context, id string, c entity.Court) error
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
	//go:embed sql/court/list_company_courts_in_showcase.sql
	listCompanyCourtsShowcaseQuery string
	//go:embed sql/court/list_available_booking_slots.sql
	listAvailableBookingSlotsQuery string
	//go:embed sql/court/update_court.sql
	updateCourtQuery string
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
		c.ID,
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

func (r *courtRepositoryImpl) InsertPhotos(ctx context.Context, photos []entity.CourtPhoto) error {
	if len(photos) == 0 {
		return nil
	}

	columns := []string{"id", "court_id", "path", "position", "is_cover"}

	rows := make([][]interface{}, len(photos))
	for i, p := range photos {
		rows[i] = []any{
			p.ID,
			p.CourtId,
			p.Path,
			p.Position,
			p.IsCover,
		}
	}

	_, err := r.db.CopyFrom(ctx, CourtPhotosName, columns, rows)
	if err != nil {
		return fmt.Errorf("CourtRepository.InsertPhotos - error inserting photos in the database: %w", err)
	}

	return nil
}

func (r *courtRepositoryImpl) FindByID(ctx context.Context, id string) (entity.Court, error) {
	var court entity.Court
    // TODO - Add multiple photos to the court entity
    // An array instead of a single photo
    var courtPhoto entity.CourtPhoto
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
        &courtPhoto.ID,
        &courtPhoto.Path,
	)

    court.Photos = []entity.CourtPhoto{courtPhoto}

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

	bookings := make([]entity.Booking, 0)
	for rows.Next() {
		var booking entity.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.CourtId,
			&booking.StartTime,
			&booking.EndTime,
			&booking.CreatedAt,
			&booking.Status,
			&booking.GuestName,
			&booking.GuestEmail,
			&booking.GuestPhone,
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

	courts := make([]entity.Court, 0)
	for rows.Next() {
		var court entity.Court
		err := rows.Scan(
			&court.ID,
			&court.CompanyId,
			&court.Name,
			&court.SportType,
			&court.HourlyPrice,
			&court.IsActive,
			&court.BookingsToday,
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

func (r *courtRepositoryImpl) ListCompanyCourtsShowcase(ctx context.Context, companyID string) ([]entity.Court, error) {
	rows, err := r.db.Query(ctx, listCompanyCourtsShowcaseQuery, companyID)
	if err != nil {
		return nil, fmt.Errorf("CourtRepository.ListCompanyCourtsShowcase: %w", err)
	}
	defer rows.Close()

	courts := make([]entity.Court, 0)
	for rows.Next() {
		var court entity.Court
        var courtPhoto entity.CourtPhoto
		err := rows.Scan(
			&court.ID,
			&court.Name,
			&court.Description,
			&court.SportType,
			&court.HourlyPrice,
			&court.IsActive,
			&court.OpeningTime,
			&court.ClosingTime,
			&court.Capacity,
            &courtPhoto.ID,
            &courtPhoto.Path,
		)
		if err != nil {
			return nil, fmt.Errorf("CourtRepository.ListCompanyCourtsShowcase: %w", err)
		}
        court.Photos = []entity.CourtPhoto{courtPhoto}

		courts = append(courts, court)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CourtRepository.ListCompanyCourtsShowcase: %w", err)
	}

	return courts, nil
}

func (r *courtRepositoryImpl) ListAvailableBookingSlots(ctx context.Context, id string, date string) ([]entity.Booking, error) {
	rows, err := r.db.Query(ctx, listAvailableBookingSlotsQuery, id, date)
	if err != nil {
		return nil, fmt.Errorf("CourtRepository.ListAvailableBookingSlots: %w", err)
	}
	defer rows.Close()

	bookings := make([]entity.Booking, 0)
	for rows.Next() {
		var booking entity.Booking
		err := rows.Scan(
			&booking.StartTime,
			&booking.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("CourtRepository.ListAvailableBookingSlots: %w", err)
		}
		bookings = append(bookings, booking)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CourtRepository.ListAvailableBookingSlots: %w", err)
	}

	return bookings, nil
}

func (r *courtRepositoryImpl) Update(ctx context.Context, id string, c entity.Court) error {
	_, err := r.db.Exec(
		ctx,
		updateCourtQuery,
		c.Name,
		c.Description,
		c.SportType,
		c.HourlyPrice,
		c.IsActive,
		c.OpeningTime,
		c.ClosingTime,
		c.Capacity,
		id,
	)
	if err != nil {
		return fmt.Errorf("CourtRepository.Update: %w", err)
	}

	return nil
}

func (r *courtRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, deleteCourtQuery, id)
	if err != nil {
		return fmt.Errorf("CourtRepository.Delete: %w", err)
	}

	return nil
}
