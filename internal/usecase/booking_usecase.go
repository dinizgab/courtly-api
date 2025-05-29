package usecase

import (
	"context"
	"fmt"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/repository"
	"github.com/dinizgab/booking-mvp/internal/services/notification"
)

type (
	BookingUsecase interface {
		Create(ctx context.Context, booking entity.Booking) (string, error)
		FindByID(ctx context.Context, id string) (entity.Booking, error)
		FindByIDShowcase(ctx context.Context, id string) (entity.Booking, error)
		ListByCompanyID(ctx context.Context, companyId string) ([]entity.Booking, error)
		ConfirmBooking(ctx context.Context, companyId string, bookingId string, verificationCode string) error
		Update(ctx context.Context, booking entity.Booking) error
		Delete(ctx context.Context, id string) error
	}

	bookingUsecaseImpl struct {
		bookingRepository   repository.BookingRepository
		companyUsecase      CompanyUsecase
        courtUsecase        CourtUseCase
		notificationService notification.Sender
	}
)

func NewBookingUsecase(
	bookingRepository repository.BookingRepository,
	companyUsecase CompanyUsecase,
    courtUsecase CourtUseCase,
	notificationService notification.Sender,
) BookingUsecase {
	return &bookingUsecaseImpl{
		bookingRepository:   bookingRepository,
		companyUsecase:      companyUsecase,
        courtUsecase:        courtUsecase,
		notificationService: notificationService,
	}
}

func (u *bookingUsecaseImpl) Create(ctx context.Context, booking entity.Booking) (string, error) {
	booking.Status = entity.StatusPending
	booking.VerificationCode = entity.GenerateVerificationCode()

    court, err := u.courtUsecase.FindByID(ctx, booking.CourtId)
    if err != nil {
        return "", err
    }

	company, err := u.companyUsecase.FindByID(ctx, court.CompanyId)
	if err != nil {
		return "", err
	}

    booking.TotalPrice = court.HourlyPrice * booking.DurationInHours()

	id, err := u.bookingRepository.Create(ctx, booking)
	if err != nil {
		return "", err
	}

	bookingEmailInfo := entity.BookingConfirmationDTO{
		GuestName:        booking.GuestName,
		GuestPhone:       booking.GuestPhone,
		GuestEmail:       booking.GuestEmail,
		CourtName:        court.Name,
		CourtAddress:     company.Address,
        BookingDate:      booking.StartTime.Format("02-12-2006"),
		BookingInterval:  fmt.Sprintf("%s - %s", booking.StartTime.Format("15:04"), booking.EndTime.Format("15:04")),
		TotalPrice:       booking.TotalPrice,
		VerificationCode: booking.VerificationCode,
	}
   
    errCh := make(chan error, 1)
    go func() {
        errCh <- u.notificationService.Send(ctx, bookingEmailInfo)
    }()

    if err = <- errCh; err != nil {
        return "", err
    }

	return id, nil
}

func (u *bookingUsecaseImpl) ListByCompanyID(ctx context.Context, companyId string) ([]entity.Booking, error) {
	bookings, err := u.bookingRepository.ListByCompanyID(ctx, companyId)
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
