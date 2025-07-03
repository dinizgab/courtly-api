package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/dinizgab/booking-mvp/internal/auth"
	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type (
	CompanyUsecase interface {
		Login(ctx context.Context, email, password string) (string, error)
		Create(ctx context.Context, company entity.Company) (string, error)
		GetDashboardInfo(ctx context.Context, companyId string) (entity.CompanyDashboard, error)
		FindByID(ctx context.Context, id string) (entity.Company, error)
		Update(ctx context.Context, id string, company entity.Company) error
		Delete(ctx context.Context, id string) error
        FindByIDShowcase(ctx context.Context, id string) (entity.Company, error)
	}

	companyUsecaseImpl struct {
		companyRepository repository.CompanyRepository
		authService       auth.AuthService
		paymentUsecase    PaymentUsecase
	}
)

func NewCompanyUsecase(
	companyRepository repository.CompanyRepository,
	authService auth.AuthService,
	paymentUsecase PaymentUsecase,
) CompanyUsecase {
	return &companyUsecaseImpl{
		companyRepository: companyRepository,
		authService:       authService,
		paymentUsecase:    paymentUsecase,
	}
}

// TODO - Make this function atomic (create company and subaccount in one transaction)
func (u *companyUsecaseImpl) Create(ctx context.Context, company entity.Company) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(company.Password), 14)
	if err != nil {
		return "", fmt.Errorf("CompanyUsecase.Create - failed to hash password: %w", err)
	}

	company.Password = string(hash)
	company.Slug = strings.ToLower(strings.ReplaceAll(company.Name, " ", "-"))

	company, err = u.companyRepository.Create(ctx, company)
	if err != nil {
		return "", err
	}

	err = u.paymentUsecase.CreateSubaccount(ctx, company)
	if err != nil {
        _ = u.companyRepository.Delete(ctx, company.ID)

		return "", err
	}

	token, err := u.authService.GenerateToken(company.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *companyUsecaseImpl) Login(ctx context.Context, email, password string) (string, error) {
	company, err := u.companyRepository.FindByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", entity.ErrInvalidCredentials
		}

		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(password))
	if err != nil {
		return "", entity.ErrInvalidCredentials
	}

	token, err := u.authService.GenerateToken(company.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *companyUsecaseImpl) GetDashboardInfo(ctx context.Context, companyId string) (entity.CompanyDashboard, error) {
	dashboard, err := u.companyRepository.GetDashboardInfo(ctx, companyId)
	if err != nil {
		return entity.CompanyDashboard{}, err
	}

	return dashboard, nil
}

func (u *companyUsecaseImpl) FindByID(ctx context.Context, id string) (entity.Company, error) {
	company, err := u.companyRepository.FindByID(ctx, id)
	if err != nil {
		return entity.Company{}, err
	}

	return company, nil
}

func (u *companyUsecaseImpl) Update(ctx context.Context, id string, company entity.Company) error {
	company.Slug = strings.ToLower(strings.ReplaceAll(company.Name, " ", "-"))

	err := u.companyRepository.Update(ctx, id, company)
	if err != nil {
		return err
	}

	return nil
}

func (u *companyUsecaseImpl) Delete(ctx context.Context, id string) error {
	err := u.companyRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *companyUsecaseImpl) FindByIDShowcase(ctx context.Context, id string) (entity.Company, error) {
    company, err := u.companyRepository.FindByIDShowcase(ctx, id)
    if err != nil {
        return entity.Company{}, err
    }

    return company, nil
}

