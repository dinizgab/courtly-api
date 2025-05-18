package usecase

import (
	"context"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/repository"
)

type (
	BookingUsecase interface {
		Create(ctx context.Context, booking entity.Booking) error
		FindByID(ctx context.Context, id string) (entity.Booking, error)
		Update(ctx context.Context, booking entity.Booking) error
		Delete(ctx context.Context, id string) error
	}

	bookingUsecaseImpl struct {
		bookingRepository repository.BookingRepository
	}
)

func NewBookingUsecase(bookingRepository repository.BookingRepository) BookingUsecase {
	return &bookingUsecaseImpl{
		bookingRepository: bookingRepository,
	}
}

func (u *bookingUsecaseImpl) Create(ctx context.Context, booking entity.Booking) error {
    // TODO - Add the verification code into the booking
    // TODO - Send email to user after create booking
	err := u.bookingRepository.Create(ctx, booking)
	if err != nil {
		return err
	}

	return nil
}

func (u *bookingUsecaseImpl) FindByID(ctx context.Context, id string) (entity.Booking, error) {
	booking, err := u.bookingRepository.FindByID(ctx, id)
	if err != nil {
		return entity.Booking{}, err
	}

	return booking, nil
}

func (u *bookingUsecaseImpl) Update(ctx context.Context, booking entity.Booking) error {
	err := u.bookingRepository.Update(ctx, booking)
	if err != nil {
		return err
	}

	return nil
}

func (u *bookingUsecaseImpl) Delete(ctx context.Context, id string) error {
	err := u.bookingRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
