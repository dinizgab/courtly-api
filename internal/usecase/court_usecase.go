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
	return nil
}

func (u *CourtUseCase) FindByID(ctx context.Context, id string) (entity.Court, error) {
	return entity.Court{}, nil
}

func (u *CourtUseCase) ListByCompany(ctx context.Context, companyID string) ([]entity.Court, error) {

	return nil, nil
}

func (u *CourtUseCase) Update(ctx context.Context, id string, court entity.Court) error {
	return nil
}
