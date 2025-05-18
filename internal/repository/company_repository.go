package repository

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/jackc/pgx/v5"
)

type (
	CompanyRepository interface {
		Create(ctx context.Context, company *entity.Company) error
		FindByID(ctx context.Context, id string) (*entity.Company, error)
		FindBySlug(ctx context.Context, slug string) (*entity.Company, error)
		Update(ctx context.Context, company *entity.Company) error
		Delete(ctx context.Context, id string) error
	}

	companyRepositoryImpl struct {
		db *pgx.Conn
	}
)

var (
	//go:embed sql/company/create_company.sql
	createCompanyQuery string
	//go:embed sql/company/find_company_by_id.sql
	findCompanyByIDQuery string
	//go:embed sql/company/find_company_by_slug.sql
	findCompanyBySlugQuery string
	//go:embed sql/company/update_company.sql
	updateCompanyQuery string
	//go:embed sql/company/delete_company.sql
	deleteCompanyQuery string
)

func NewCompanyRepository(db *pgx.Conn) CompanyRepository {
	return &companyRepositoryImpl{
		db: db,
	}
}

func (r *companyRepositoryImpl) Create(ctx context.Context, company *entity.Company) error {
	_, err := r.db.Exec(ctx, createCompanyQuery, company.Name, company.Address, company.Phone, company.Email, company.Slug)
	if err != nil {
		return fmt.Errorf("CompanyRepository.Create: %w", err)
	}

	return nil
}

func (r *companyRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Company, error) {
    var company entity.Company
    err := r.db.QueryRow(ctx, findCompanyByIDQuery, id).Scan(
        &company.ID,
        &company.Name,
        &company.Address,
        &company.Phone,
        &company.Email,
        &company.Slug,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("CompanyRepository.FindByID: company not found")
        }

        return nil, fmt.Errorf("CompanyRepository.FindByID: %w", err)
    }

    return &company, nil
}

func (r *companyRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entity.Company, error) {
    var company entity.Company
    err := r.db.QueryRow(ctx, findCompanyBySlugQuery, slug).Scan(
        &company.ID,
        &company.Name,
        &company.Address,
        &company.Phone,
        &company.Email,
        &company.Slug,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("CompanyRepository.FindByID: company not found")
        }

        return nil, fmt.Errorf("CompanyRepository.FindByID: %w", err)
    }

    return &company, nil
}

func (r *companyRepositoryImpl) Update(ctx context.Context, company *entity.Company) error {
    // TODO - check update logic
    _, err := r.db.Exec(ctx, updateCompanyQuery, company.Name, company.Address, company.Phone, company.Email, company.Slug, company.ID)
    if err != nil {
        return fmt.Errorf("CompanyRepository.Update: %w", err)
    }

	return nil
}

func (r *companyRepositoryImpl) Delete(ctx context.Context, id string) error {
    _, err := r.db.Exec(ctx, deleteCompanyQuery, id)
    if err != nil {
        return fmt.Errorf("CompanyRepository.Delete: %w", err)
    }

	return nil
}
