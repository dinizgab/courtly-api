package usecase

import (
	"context"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/repository"
)

type (
	CompanyUsecase interface {
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
	err := u.companyRepository.Create(ctx, company)
	if err != nil {
		return err
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
