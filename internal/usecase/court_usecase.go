package usecase

import (
	"context"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/repository"
)

type (
	CourtUseCase struct {
		courtRepository repository.CourtRepository
	}

	CourtUseCaseInterface interface {
		Create(ctx context.Context, court entity.Court) error
		FindByID(ctx context.Context, id string) (entity.Court, error)
		ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error)
		Update(ctx context.Context, id string, court entity.Court) error
		Delete(ctx context.Context, id string) error
	}
)

func NewCourtUseCase(courtRepository repository.CourtRepository) CourtUseCaseInterface {
	return &CourtUseCase{
		courtRepository: courtRepository,
	}
}

func (u *CourtUseCase) Create(ctx context.Context, court entity.Court) error {
	err := u.courtRepository.Create(ctx, &court)
	if err != nil {
		return err
	}

	return nil
}

func (u *CourtUseCase) FindByID(ctx context.Context, id string) (entity.Court, error) {
	court, err := u.courtRepository.FindByID(ctx, id)
	if err != nil {
		return entity.Court{}, err
	}

	return court, nil
}

func (u *CourtUseCase) ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error) {
	courts, err := u.courtRepository.ListByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return courts, nil
}

func (u *CourtUseCase) Update(ctx context.Context, id string, court entity.Court) error {
	err := u.courtRepository.Update(ctx, &court)
	if err != nil {
		return err
	}

	return nil
}

func (u *CourtUseCase) Delete(ctx context.Context, id string) error {
	err := u.courtRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
