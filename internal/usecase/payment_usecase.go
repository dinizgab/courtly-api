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

const (
    bookingConfirmationEmailSubject = "Confirmação de reserva"
    refundEmailSubject = "Confirmação de solicitação de reembolso"

    bookingConfirmationTemplateName = "booking_confirmation.html"
    refundTemplateName = "refund_request_confirmation.html"
)

type PaymentUsecase interface {
	CreateSubaccount(ctx context.Context, company entity.Company) error
	CreateCharge(ctx context.Context, companyId string, booking entity.Booking) error
	ConfirmPayment(ctx context.Context, charge openpix.Charge) error
	GetBookingPaymentStatusByID(ctx context.Context, id string) (string, error)
	GetCompanyBalance(ctx context.Context, id string) (int64, error)
	CreateWithdrawRequest(ctx context.Context, companyId string) error
	ExpirePayment(ctx context.Context, charge openpix.Charge) error
	GetBookingChargeInformation(ctx context.Context, id string) (entity.Payment, error)
	RefundCharge(ctx context.Context, bookingId string) error
}

type pixGatewayUsecaseImpl struct {
	pixClient           openpix.OpenPixClient
	summaryReader       ports.BookingSummaryReader
	tokenWriter         ports.BookingCancelTokenWriter
	repo                repository.PaymentRepository
	notificationService notification.Sender
}

func NewPixGatewayService(
	pixClient openpix.OpenPixClient,
	summaryReader ports.BookingSummaryReader,
	tokenWriter ports.BookingCancelTokenWriter,
	repo repository.PaymentRepository,
	notificationService notification.Sender,
) PaymentUsecase {
	return &pixGatewayUsecaseImpl{
		pixClient:           pixClient,
		summaryReader:       summaryReader,
		tokenWriter:         tokenWriter,
		repo:                repo,
		notificationService: notificationService,
	}
}

func (uc *pixGatewayUsecaseImpl) CreateSubaccount(ctx context.Context, company entity.Company) error {
	subaccount := openpix.Subaccount{
		Name:   company.Slug,
		PixKey: company.PixKey,
	}

	subaccount, err := uc.pixClient.CreateSubaccount(ctx, subaccount)
	if err != nil {
		return err
	}

	subaccountEntity := entity.Subaccount{
		CompanyID: company.ID,
		PixKey:    subaccount.PixKey,
        PixKeyType: company.PixKeyType,
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
	booking, err := uc.summaryReader.GetBookingSummary(ctx, bookingId)
	if err != nil {
		return err
	}

	token, err := entity.GenerateCancelToken()
	if err != nil {
		return err
	}
	tokenHash := entity.HashCancelToken(token)
	if err := uc.tokenWriter.SetCancelTokenHash(ctx, bookingId, tokenHash); err != nil {
		return err
	}

    loc := time.FixedZone("BRT", -3*3600)
	bookingEmailInfo := entity.BookingConfirmationInfo{
		ID:               bookingId,
		GuestName:        booking.GuestName,
		GuestPhone:       booking.GuestPhone,
		GuestEmail:       booking.GuestEmail,
		CourtName:        booking.Court.Name,
		CourtAddress:     booking.Court.Company.Address,
		BookingDate:      booking.StartTime.In(loc).Format("02-01-2006"),
		BookingInterval:  fmt.Sprintf("%s - %s", booking.StartTime.In(loc).Format("15:04"), booking.EndTime.In(loc).Format("15:04")),
		TotalPrice:       fmt.Sprintf("%.2f", float64(booking.TotalPrice)/100),
		VerificationCode: booking.VerificationCode,
		CancelToken:      token,
	}

	err = uc.notificationService.Send(ctx, bookingConfirmationTemplateName, bookingConfirmationEmailSubject, bookingEmailInfo, booking.GuestEmail)
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

func (uc *pixGatewayUsecaseImpl) CreateWithdrawRequest(ctx context.Context, companyId string) error {
	pixKey, err := uc.repo.GetSubaccountPixKeyByCompanyID(ctx, companyId)
	if err != nil {
		return err
	}

	withdraw, err := uc.pixClient.WithdrawSubaccount(ctx, pixKey)
	if err != nil {
		return err
	}

	err = uc.repo.CreateWithdrawRequest(ctx, companyId, withdraw)
	if err != nil {
		return fmt.Errorf("failed to create withdraw request: %w", err)
	}

	return nil
}

func (uc *pixGatewayUsecaseImpl) ExpirePayment(ctx context.Context, charge openpix.Charge) error {
	err := uc.repo.ExpirePayment(ctx, charge)
	if err != nil {
		return err
	}

	return nil
}

func (uc *pixGatewayUsecaseImpl) GetBookingChargeInformation(ctx context.Context, id string) (entity.Payment, error) {
	payment, err := uc.repo.GetBookingChargeInformation(ctx, id)
	if err != nil {
		return entity.Payment{}, err
	}

	return payment, nil
}

func (uc *pixGatewayUsecaseImpl) RefundCharge(ctx context.Context, bookingId string) error {
	payment, err := uc.repo.GetPaymentByBookingID(ctx, bookingId)
	if err != nil {
		return err
	}

	refund, err := uc.pixClient.RefundCharge(ctx, payment)
	if err != nil {
		return err
	}

	err = uc.repo.SaveRefundRequest(ctx, bookingId, refund)
	if err != nil {
		return err
	}

    booking, err := uc.summaryReader.GetBookingSummary(ctx, bookingId)
    if err != nil {
        return err
    }

    loc := time.FixedZone("BRT", -3*3600)
	bookingEmailInfo := entity.BookingConfirmationInfo{
		ID:               bookingId,
		GuestName:        booking.GuestName,
		GuestPhone:       booking.GuestPhone,
		GuestEmail:       booking.GuestEmail,
		CourtName:        booking.Court.Name,
		CourtAddress:     booking.Court.Company.Address,
		BookingDate:      booking.StartTime.In(loc).Format("02-01-2006"),
		BookingInterval:  fmt.Sprintf("%s - %s", booking.StartTime.In(loc).Format("15:04"), booking.EndTime.In(loc).Format("15:04")),
		TotalPrice:       fmt.Sprintf("%.2f", float64(booking.TotalPrice)/100),
		VerificationCode: booking.VerificationCode,
	}

	err = uc.notificationService.Send(ctx, refundTemplateName, refundEmailSubject, bookingEmailInfo, booking.GuestEmail)
	if err != nil {
		return err
	}

	return nil
}
