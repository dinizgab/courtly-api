package usecase

import (
	"context"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/repository"
)

type (
	courtUseCaseImpl struct {
		courtRepository repository.CourtRepository
	}

	CourtUseCase interface {
		Create(ctx context.Context, court entity.Court) error
		FindByID(ctx context.Context, id string) (entity.Court, error)
		ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error)
		Update(ctx context.Context, id string, court entity.Court) error
		Delete(ctx context.Context, id string) error
	}
)

func NewCourtUseCase(courtRepository repository.CourtRepository) CourtUseCase {
	return &courtUseCaseImpl{
		courtRepository: courtRepository,
	}
}

func (u *courtUseCaseImpl) Create(ctx context.Context, court entity.Court) error {
	err := u.courtRepository.Create(ctx, &court)
	if err != nil {
		return err
	}

	return nil
}

func (u *courtUseCaseImpl) FindByID(ctx context.Context, id string) (entity.Court, error) {
	court, err := u.courtRepository.FindByID(ctx, id)
	if err != nil {
		return entity.Court{}, err
	}

	return court, nil
}

func (u *courtUseCaseImpl) ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error) {
	courts, err := u.courtRepository.ListByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return courts, nil
}

func (u *courtUseCaseImpl) Update(ctx context.Context, id string, court entity.Court) error {
	err := u.courtRepository.Update(ctx, &court)
	if err != nil {
		return err
	}

	return nil
}

func (u *courtUseCaseImpl) Delete(ctx context.Context, id string) error {
	err := u.courtRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
