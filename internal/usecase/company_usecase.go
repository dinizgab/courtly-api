package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type (
	CompanyUsecase interface {
		Login(ctx context.Context, email, password string) error
		Create(ctx context.Context, company entity.Company) error
		FindByID(ctx context.Context, id string) (entity.Company, error)
		FindBySlug(ctx context.Context, slug string) (entity.Company, error)
		Update(ctx context.Context, company entity.Company) error
		Delete(ctx context.Context, id string) error
	}

	companyUsecaseImpl struct {
		companyRepository repository.CompanyRepository
	}
)

func NewCompanyUsecase(companyRepository repository.CompanyRepository) CompanyUsecase {
	return &companyUsecaseImpl{
		companyRepository: companyRepository,
	}
}

func (u *companyUsecaseImpl) Create(ctx context.Context, company entity.Company) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(company.Password), 14)
	if err != nil {
		return fmt.Errorf("CompanyUsecase.Create - failed to hash password: %w", err)
	}
	company.Password = string(hash)
	company.Slug = strings.ToLower(strings.ReplaceAll(company.Name, " ", "-"))

	err = u.companyRepository.Create(ctx, company)
	if err != nil {
		return err
	}

	return nil
}

func (u *companyUsecaseImpl) Login(ctx context.Context, email, password string) error {
	company, err := u.companyRepository.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(password))
	if err != nil {
		return entity.ErrInvalidCredentials
	}

	return nil
}

func (u *companyUsecaseImpl) FindByID(ctx context.Context, id string) (entity.Company, error) {
	company, err := u.companyRepository.FindByID(ctx, id)
	if err != nil {
		return entity.Company{}, err
	}

	return company, nil
}

func (u *companyUsecaseImpl) FindBySlug(ctx context.Context, slug string) (entity.Company, error) {

	company, err := u.companyRepository.FindBySlug(ctx, slug)
	if err != nil {
		return entity.Company{}, err
	}

	return company, nil
}

func (u *companyUsecaseImpl) Update(ctx context.Context, company entity.Company) error {
	err := u.companyRepository.Update(ctx, &company)
	if err != nil {
		return err
	}

	return nil
}

func (u *companyUsecaseImpl) Delete(ctx context.Context, id string) error {
	err := u.companyRepository.Delete(ctx, id)
	if err != nil {
		return nil
	}

	return nil
}
