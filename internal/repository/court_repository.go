package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/jackc/pgx/v5"
)

const CourtPhotosName = "court_photos"
const CourtSchedulesName = "court_schedules"

type (
	CourtRepository interface {
		Create(ctx context.Context, c *entity.Court) (string, error)
		InsertPhotos(ctx context.Context, c []entity.CourtPhoto) error
		FindByID(ctx context.Context, id string) (entity.Court, error)
		ListBookingsByID(ctx context.Context, id string) ([]entity.Booking, error)
		ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error)
		ListCompanyCourtsShowcase(ctx context.Context, companyID string) ([]entity.Court, error)
		ListAvailableBookingSlots(ctx context.Context, id string, date string) ([]entity.Booking, error)
		Update(ctx context.Context, id string, c entity.Court) error
		Delete(ctx context.Context, id string) error
		UpdateCourtStatus(ctx context.Context, id string, court entity.Court) error
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
	//go:embed sql/court/list_court_photos.sql
	listCourtPhotosQuery string
	//go:embed sql/court/list_court_schedule.sql
	listCourtScheduleQuery string
	//go:embed sql/court/update_court.sql
	updateCourtQuery string
	//go:embed sql/court/update_court_schedule.sql
	updateCourtScheduleQuery string
	//go:embed sql/court/delete_court.sql
	deleteCourtQuery string
	//go:embed sql/court/update_court_status.sql
	updateCourtStatusQuery string
)

func NewCourtRepository(db database.Database) CourtRepository {
	return &courtRepositoryImpl{
		db: db,
	}
}

func (r *courtRepositoryImpl) Create(ctx context.Context, c *entity.Court) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("CourtRepository.Create: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Printf("CourtRepository.Create: could not rollback transaction: %v\n", rollbackErr)
			}
		}
	}()

	var id string
	err = tx.QueryRow(
		ctx,
		createCourtQuery,
		c.CompanyId,
		c.Name,
		c.Description,
		c.SportType,
		c.HourlyPrice,
		c.IsActive,
		c.Capacity,
	).Scan(
		&id,
	)
	if err != nil {
		return "", fmt.Errorf("CourtRepository.Create: %w", err)
	}

	columns := []string{"court_id", "day_of_week", "is_open", "opening_time", "closing_time"}
	rows := make([][]interface{}, 0, len(c.CourtSchedule))
	for _, s := range c.CourtSchedule {
		rows = append(rows, []interface{}{
			id,
			s.Weekday,
			s.IsOpen,
			s.OpeningTime,
			s.ClosingTime,
		})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{CourtSchedulesName},
		columns,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return "", fmt.Errorf("CourtRepository.Create: copy court_schedules: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("CourtRepository.Create: commit tx: %w", err)
	}

	return id, nil
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
	var courtPhotos []entity.CourtPhoto
	var courtSchedule []entity.CourtSchedule

	err := r.db.QueryRow(ctx, findCourtByIDQuery, id).Scan(
		&court.ID,
		&court.CompanyId,
		&court.Name,
		&court.Description,
		&court.SportType,
		&court.HourlyPrice,
		&court.IsActive,
		&court.Capacity,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Court{}, fmt.Errorf("CourtRepository.FindByID: court not found: %w", err)
		}
		return entity.Court{}, fmt.Errorf("CourtRepository.FindByID: %w", err)
	}

	rows, err := r.db.Query(ctx, listCourtScheduleQuery, id)
	if err != nil {
		return entity.Court{}, fmt.Errorf("CourtRepository.FindByID: error querying court: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var daySchedule entity.CourtSchedule
		err := rows.Scan(
			&daySchedule.ID,
			&daySchedule.IsOpen,
			&daySchedule.Weekday,
			&daySchedule.OpeningTime,
			&daySchedule.ClosingTime,
		)
		if err != nil {
			return entity.Court{}, fmt.Errorf("CourtRepository.FindByID: error scanning court: %w", err)
		}

		courtSchedule = append(courtSchedule, daySchedule)
	}

	rows, err = r.db.Query(ctx, listCourtPhotosQuery, id)
	if err != nil {
		return entity.Court{}, fmt.Errorf("CourtRepository.FindByID: error querying court photos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var courtPhoto entity.CourtPhoto
		err := rows.Scan(
			&courtPhoto.ID,
			&courtPhoto.Path,
		)
		if err != nil {
			return entity.Court{}, fmt.Errorf("CourtRepository.FindByID: error scanning court photo: %w", err)
		}

		courtPhotos = append(courtPhotos, courtPhoto)
	}

	court.Photos = courtPhotos
	court.CourtSchedule = courtSchedule

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

	courtsMap := make(map[string]*entity.Court)
	for rows.Next() {
		var (
			courtID                  string
			name                     string
			description              string
			sportType                string
			hourlyPrice              int64
			isActive                 bool
			capacity                 int
			openingTime, closingTime sql.NullString
			photoID, photoPath       sql.NullString
		)

		if err := rows.Scan(
			&courtID,
			&name,
			&description,
			&sportType,
			&hourlyPrice,
			&isActive,
			&capacity,
			&openingTime,
			&closingTime,
			&photoID,
			&photoPath,
		); err != nil {
			return nil, fmt.Errorf("CourtRepository.ListCompanyCourtsShowcase: %w", err)
		}

		court, exists := courtsMap[courtID]
		if !exists {
			court = &entity.Court{
				ID:          courtID,
				Name:        name,
				Description: description,
				SportType:   sportType,
				HourlyPrice: hourlyPrice,
				IsActive:    isActive,
				Capacity:    capacity,
			}

			if openingTime.Valid && closingTime.Valid {
				// TODO - Fix time parsing
				parsedOpening, err1 := time.Parse("15:04:05.000000", openingTime.String)
				parsedClosing, err2 := time.Parse("15:04:05.000000", closingTime.String)
				if err1 != nil || err2 != nil {
					return nil, fmt.Errorf("CourtRepository.ListCompanyCourtsShowcase: error parsing time: %w", err1)
				}

				court.CourtSchedule = []entity.CourtSchedule{
					{
						OpeningTime: parsedOpening,
						ClosingTime: parsedClosing,
					},
				}
			}
			courtsMap[courtID] = court
		}

		if photoID.Valid && photoPath.Valid {
			court.Photos = append(court.Photos, entity.CourtPhoto{
				ID:      photoID.String,
				CourtId: courtID,
				Path:    photoPath.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CourtRepository.ListCompanyCourtsShowcase: %w", err)
	}

	courts := make([]entity.Court, 0, len(courtsMap))
	for _, c := range courtsMap {
		courts = append(courts, *c)
	}
	sort.SliceStable(courts, func(i, j int) bool { return courts[i].ID < courts[j].ID })

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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("CourtRepository.Update (begin tx): %w", err)
	}

	defer func() {
		if tx != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = r.db.Exec(
		ctx,
		updateCourtQuery,
		c.Name,
		c.Description,
		c.SportType,
		c.HourlyPrice,
		c.IsActive,
		c.Capacity,
		id,
	)
	if err != nil {
		return fmt.Errorf("CourtRepository.Update: %w", err)
	}

	if len(c.CourtSchedule) > 0 {
		valueStrings := make([]string, 0, len(c.CourtSchedule))
		args := make([]interface{}, 0, len(c.CourtSchedule)*4)

		for i, s := range c.CourtSchedule {
			startIdx := i*4 + 1
			valueStrings = append(valueStrings, fmt.Sprintf("($%d::uuid, $%d::boolean, $%d::time, $%d::time)", startIdx, startIdx+1, startIdx+2, startIdx+3))
			args = append(args, s.ID, s.IsOpen, s.OpeningTime, s.ClosingTime)
		}

		query := fmt.Sprintf(updateCourtScheduleQuery, strings.Join(valueStrings, ","))

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("CourtRepository.Update (schedules): %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("CourtRepository.Update (commit): %w", err)
	}
	// TODO - Switch this, this is bad practice
	tx = nil

	return nil
}

func (r *courtRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, deleteCourtQuery, id)
	if err != nil {
		return fmt.Errorf("CourtRepository.Delete: %w", err)
	}

	return nil
}

func (r *courtRepositoryImpl) UpdateCourtStatus(ctx context.Context, id string, court entity.Court) error {
	_, err := r.db.Exec(ctx, updateCourtStatusQuery, id, court.IsActive)
	if err != nil {
		return fmt.Errorf("CourtRepository.UpdateCourtStatus: %w", err)
	}

	return nil
}
