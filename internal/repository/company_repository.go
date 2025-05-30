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
	CompanyRepository interface {
		Create(ctx context.Context, company entity.Company) (entity.Company, error)
		FindByID(ctx context.Context, id string) (entity.Company, error)
        FindByEmail(ctx context.Context, email string) (entity.Company, error)
		FindBySlug(ctx context.Context, slug string) (entity.Company, error)
        GetDashboardInfo(ctx context.Context, companyId string) (entity.CompanyDashboard, error)
		Update(ctx context.Context, id string, company entity.Company) error
		Delete(ctx context.Context, id string) error
	}

	companyRepositoryImpl struct {
		db database.Database
	}
)

var (
	//go:embed sql/company/create_company.sql
	createCompanyQuery string
	//go:embed sql/company/find_company_by_id.sql
	findCompanyByIDQuery string
    //go:embed sql/company/find_company_by_email.sql
    findCompanyByEmailQuery string
	//go:embed sql/company/find_company_by_slug.sql
	findCompanyBySlugQuery string
    //go:embed sql/company/get_dashboard_info.sql
    getDashboardInfoQuery string
	//go:embed sql/company/update_company.sql
	updateCompanyQuery string
	//go:embed sql/company/delete_company.sql
	deleteCompanyQuery string
)

func NewCompanyRepository(db database.Database) CompanyRepository {
	return &companyRepositoryImpl{
		db: db,
	}
}

func (r *companyRepositoryImpl) Create(ctx context.Context, company entity.Company) (entity.Company, error) {
	err := r.db.QueryRow(ctx,
		createCompanyQuery,
		company.Name,
		company.Address,
		company.Phone,
		company.Email,
		company.Password,
		company.CNPJ,
		company.Slug,
	).Scan(&company.ID)
	if err != nil {
		return company, fmt.Errorf("CompanyRepository.Create - error creating company: %w", err)
	}

	return company, nil
}

func (r *companyRepositoryImpl) FindByID(ctx context.Context, id string) (entity.Company, error) {
	var company entity.Company
	err := r.db.QueryRow(ctx, findCompanyByIDQuery, id).Scan(
		&company.ID,
		&company.Name,
		&company.Address,
		&company.Phone,
		&company.Email,
        &company.CNPJ,
		&company.Slug,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Company{}, fmt.Errorf("CompanyRepository.FindByID: company not found")
		}

		return entity.Company{}, fmt.Errorf("CompanyRepository.FindByID: %w", err)
	}

	return company, nil
}

func (r *companyRepositoryImpl) FindByEmail(ctx context.Context, email string) (entity.Company, error) {
    var company entity.Company
    err := r.db.QueryRow(ctx, findCompanyByEmailQuery, email).Scan(
        &company.ID,
        &company.Email,
        &company.Password,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return entity.Company{}, err
        }

        return entity.Company{}, fmt.Errorf("CompanyRepository.FindByEmail: %w", err)
    }

    return company, nil
}

func (r *companyRepositoryImpl) FindBySlug(ctx context.Context, slug string) (entity.Company, error) {
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
			return entity.Company{}, fmt.Errorf("CompanyRepository.FindByID: company not found")
		}

		return entity.Company{}, fmt.Errorf("CompanyRepository.FindByID: %w", err)
	}

	return company, nil
}

func (r *companyRepositoryImpl) GetDashboardInfo(ctx context.Context, companyId string) (entity.CompanyDashboard, error) {
    var dashboard entity.CompanyDashboard
    err := r.db.QueryRow(ctx, getDashboardInfoQuery, companyId).Scan(
        &dashboard.TotalEarnings,
        &dashboard.TotalBookedHours,
        &dashboard.TotalBookings,
        &dashboard.TotalClients,
    )
    if err != nil {
        return entity.CompanyDashboard{}, fmt.Errorf("CompanyRepository.GetDashboardInfo: %w", err)
    }

    return dashboard, nil
}

func (r *companyRepositoryImpl) Update(ctx context.Context, id string, company entity.Company) error {
	_, err := r.db.Exec(
        ctx,
        updateCompanyQuery,
        company.Name,
        company.Address,
        company.Phone,
        company.Email,
        company.CNPJ,
        company.Slug,
        id,
    )
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
