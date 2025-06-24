package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/repository"
)

type (
	BookingUsecase interface {
		Create(ctx context.Context, booking entity.Booking) (string, error)
		FindByID(ctx context.Context, id string) (entity.Booking, error)
		FindByIDShowcase(ctx context.Context, id string) (entity.Booking, error)
		ListByCompanyID(ctx context.Context, companyId string, filter entity.BookingFilter) ([]entity.Booking, error)
		ConfirmBooking(ctx context.Context, companyId string, bookingId string, verificationCode string) error
		CancelBooking(ctx context.Context, bookingId string, cancelToken string) error
		Update(ctx context.Context, booking entity.Booking) error
		Delete(ctx context.Context, id string) error
	}

	bookingUsecaseImpl struct {
		bookingRepository repository.BookingRepository
		paymentUsecase    PaymentUsecase
		companyUsecase    CompanyUsecase
		courtUsecase      CourtUseCase
	}
)

func NewBookingUsecase(
	bookingRepository repository.BookingRepository,
	paymentUsecase PaymentUsecase,
	companyUsecase CompanyUsecase,
	courtUsecase CourtUseCase,
) BookingUsecase {
	return &bookingUsecaseImpl{
		bookingRepository: bookingRepository,
		paymentUsecase:    paymentUsecase,
		companyUsecase:    companyUsecase,
		courtUsecase:      courtUsecase,
	}
}

func (u *bookingUsecaseImpl) Create(ctx context.Context, booking entity.Booking) (string, error) {
	booking.Status = entity.StatusPending
	booking.VerificationCode = entity.GenerateVerificationCode()

	court, err := u.courtUsecase.FindByID(ctx, booking.CourtId)
	if err != nil {
		return "", err
	}
	booking.TotalPrice = court.HourlyPrice * booking.DurationInHours()

	id, err := u.bookingRepository.Create(ctx, booking)
	if err != nil {
		return "", err
	}

	booking.ID = id

	err = u.paymentUsecase.CreateCharge(ctx, court.CompanyId, booking)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (u *bookingUsecaseImpl) ListByCompanyID(ctx context.Context, companyId string, filter entity.BookingFilter) ([]entity.Booking, error) {
	bookings, err := u.bookingRepository.ListByCompanyID(ctx, companyId, filter)
	if err != nil {
		return bookings, err
	}

	return bookings, nil
}

func (u *bookingUsecaseImpl) FindByID(ctx context.Context, id string) (entity.Booking, error) {
	booking, err := u.bookingRepository.FindByID(ctx, id)
	if err != nil {
		return entity.Booking{}, err
	}

	return booking, nil
}

func (u *bookingUsecaseImpl) FindByIDShowcase(ctx context.Context, id string) (entity.Booking, error) {
	booking, err := u.bookingRepository.FindByIDShowcase(ctx, id)
	if err != nil {
		return entity.Booking{}, err
	}

	return booking, nil
}

func (u *bookingUsecaseImpl) ConfirmBooking(ctx context.Context, companyId string, bookingId string, verificationCode string) error {
	booking, err := u.bookingRepository.FindByID(ctx, bookingId)
	if err != nil {
		return err
	}

	if booking.VerificationCode != verificationCode {
		return entity.ErrInvalidVerificationCode
	}

	if booking.Status == entity.StatusConfirmed {
		return entity.ErrBookingAlreadyConfirmed
	}

	err = u.bookingRepository.ConfirmBooking(ctx, companyId, bookingId)
	if err != nil {
		return err
	}

	return nil
}

func (u *bookingUsecaseImpl) CancelBooking(ctx context.Context, bookingId string, cancelToken string) error {
	booking, err := u.bookingRepository.GetCancelTokenInfo(ctx, bookingId)
	if err != nil {
		return err
	}

	if time.Now().After(booking.CancelTokenHashExpiresAt) {
        return fmt.Errorf("BookingUsecase.CancelBooking - token expired")
	}

	if entity.HashCancelToken(cancelToken) != booking.CancelTokenHash {
        return fmt.Errorf("BookingUsecase.CancelBooking - invalid cancel token")
	}

	err = u.paymentUsecase.RefundCharge(ctx, bookingId)
    if err != nil {
        return err
    }

    err = u.bookingRepository.CancelBooking(ctx, bookingId)
    if err != nil {
        return err
    }

	return nil
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
