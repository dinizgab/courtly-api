package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/gateway/openpix"
	"github.com/dinizgab/booking-mvp/internal/ports"
	"github.com/dinizgab/booking-mvp/internal/repository"
	"github.com/dinizgab/booking-mvp/internal/services/notification"
)

type PaymentUsecase interface {
	CreateSubaccount(ctx context.Context, company entity.Company) error
	CreateCharge(ctx context.Context, companyId string, booking entity.Booking) error
	ConfirmPayment(ctx context.Context, charge openpix.Charge) error
	GetBookingPaymentStatusByID(ctx context.Context, id string) (string, error)
	GetCompanyBalance(ctx context.Context, id string) (int64, error)
}

type pixGatewayUsecaseImpl struct {
	pixClient           openpix.OpenPixClient
    bookingRepo         ports.BookingSummaryReader
	repo                repository.PaymentRepository
	notificationService notification.Sender
}

func NewPixGatewayService(
	pixClient openpix.OpenPixClient,
	bookingRepo ports.BookingSummaryReader,
	repo repository.PaymentRepository,
	notificationService notification.Sender,
) PaymentUsecase {
	return &pixGatewayUsecaseImpl{
		pixClient:           pixClient,
		bookingRepo:         bookingRepo,
		repo:                repo,
		notificationService: notificationService,
	}
}

func (uc *pixGatewayUsecaseImpl) CreateSubaccount(ctx context.Context, company entity.Company) error {
	subaccount := openpix.Subaccount{
		Name:   company.Slug,
		PixKey: company.Email,
	}

	subaccount, err := uc.pixClient.CreateSubaccount(ctx, subaccount)
	if err != nil {
		return err
	}

	subaccountEntity := entity.Subaccount{
		CompanyID: company.ID,
		PixKey:    subaccount.PixKey,
	}

	err = uc.repo.CreateSubaccount(ctx, subaccountEntity)
	if err != nil {
		return err
	}

	return nil
}

func (uc *pixGatewayUsecaseImpl) CreateCharge(ctx context.Context, companyId string, booking entity.Booking) error {
	subaccountPixKey, err := uc.repo.GetSubaccountPixKeyByCompanyID(ctx, companyId)
	if err != nil {
		return err
	}

	charge, err := uc.pixClient.CreateCharge(ctx, subaccountPixKey, booking)
	if err != nil {
		return err
	}

	err = uc.repo.CreateCharge(ctx, companyId, charge)
	if err != nil {
		return err
	}

	return nil
}

func (uc *pixGatewayUsecaseImpl) ConfirmPayment(ctx context.Context, charge openpix.Charge) error {
	err := uc.repo.ConfirmPayment(ctx, charge)
	if err != nil {
		return err
	}

	bookingId := strings.TrimPrefix(charge.CorrelationID, "booking-")
	booking, err := uc.bookingRepo.GetBookingSummary(ctx, bookingId)
	if err != nil {
		return err
	}

    loc, _ := time.LoadLocation("America/Sao_Paulo")
	bookingEmailInfo := entity.BookingConfirmationInfo{
		GuestName:        booking.GuestName,
		GuestPhone:       booking.GuestPhone,
		GuestEmail:       booking.GuestEmail,
		CourtName:        booking.Court.Name,
		CourtAddress:     booking.Court.Company.Address,
		BookingDate:      booking.StartTime.In(loc).Format("02-01-2006"),
		BookingInterval:  fmt.Sprintf("%s - %s", booking.StartTime.Format("15:04"), booking.EndTime.Format("15:04")),
		TotalPrice:       booking.TotalPrice,
		VerificationCode: booking.VerificationCode,
	}

	err = uc.notificationService.Send(ctx, bookingEmailInfo)
	if err != nil {
		return err
	}

	return nil
}

func (uc *pixGatewayUsecaseImpl) GetBookingPaymentStatusByID(ctx context.Context, id string) (string, error) {
	status, err := uc.repo.GetBookingPaymentStatusByID(ctx, id)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (uc *pixGatewayUsecaseImpl) GetCompanyBalance(ctx context.Context, id string) (int64, error) {
	pixKey, err := uc.repo.GetSubaccountPixKeyByCompanyID(ctx, id)
	if err != nil {
		return 0, err
	}

	total, err := uc.pixClient.GetCompanyBalance(ctx, pixKey)
	if err != nil {
		return 0, err
	}

	return total, nil
}

