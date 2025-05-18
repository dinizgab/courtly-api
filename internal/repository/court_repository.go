package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/jackc/pgx/v5"
)

type (
	CourtRepository interface {
		Create(ctx context.Context, c *entity.Court) error
		FindByID(ctx context.Context, id string) (*entity.Court, error)
		ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error)
		Update(ctx context.Context, c *entity.Court) error
		Delete(ctx context.Context, id string) error
	}

	courtRepositoryImpl struct {
		db *pgx.Conn
	}
)

var (
	//go:embed sql/court/create_court.sql
	createCourtQuery string
	//go:embed sql/court/find_court_by_id.sql
	findCourtByIDQuery string
	//go:embed sql/court/list_court_by_company.sql
	listCourtByCompanyQuery string
    //go:embed sql/court/delete_court.sql
    deleteCourtQuery string
)

func NewCourtRepository(db *pgx.Conn) CourtRepository {
	return &courtRepositoryImpl{
		db: db,
	}
}

func (r *courtRepositoryImpl) Create(ctx context.Context, c *entity.Court) error {
	_, err := r.db.Exec(ctx, createCourtQuery, c.CompanyId, c.Name, c.SportType, c.HourlyPrice, c.IsActive)
	if err != nil {
		return fmt.Errorf("CourtRepository.Create: %w", err)
	}

	return nil
}

func (r *courtRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Court, error) {
	var court entity.Court
	err := r.db.QueryRow(ctx, findCourtByIDQuery, id).Scan(
		&court.ID,
		&court.CompanyId,
		&court.Name,
		&court.IsActive,
		&court.SportType,
		&court.HourlyPrice,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("CourtRepository.FindByID: court not found: %w", err)
		}

		return nil, fmt.Errorf("CourtRepository.FindByID: %w", err)
	}

	return &court, nil
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
			&court.IsActive,
			&court.SportType,
			&court.HourlyPrice,
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
