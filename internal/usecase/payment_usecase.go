package usecase

import (
	"context"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/gateway/openpix"
	"github.com/dinizgab/booking-mvp/internal/repository"
)

type PaymentUsecase interface {
	CreateSubaccount(ctx context.Context, company entity.Company) error
	CreateCharge(ctx context.Context, companyId string, booking entity.Booking) error
    ConfirmPayment(ctx context.Context, charge openpix.Charge) error
    GetBookingPaymentStatusByID(ctx context.Context, id string) (string, error)
	GetTotalToWithdraw(ctx context.Context, slug string) (float64, error)
}

type pixGatewayUsecaseImpl struct {
	pixClient openpix.OpenPixClient
	repo      repository.PaymentRepository
}

func NewPixGatewayService(pixClient openpix.OpenPixClient, repo repository.PaymentRepository) PaymentUsecase {
	return &pixGatewayUsecaseImpl{
		pixClient: pixClient,
		repo:      repo,
	}
}

func (uc *pixGatewayUsecaseImpl) CreateSubaccount(ctx context.Context, company entity.Company) error {
	subaccount, err := uc.pixClient.CreateSubaccount(ctx, openpix.CreateSubAccountRequest{
		Name:   company.Slug,
		PixKey: company.Email,
	})
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

    return nil
}

func (uc *pixGatewayUsecaseImpl) GetBookingPaymentStatusByID(ctx context.Context, id string) (string, error) {
    status, err := uc.repo.GetBookingPaymentStatusByID(ctx, id)
    if err != nil {
        return "", err
    }

    return status, nil
}

func (uc *pixGatewayUsecaseImpl) GetTotalToWithdraw(ctx context.Context, slug string) (float64, error) {
	return 0, nil
}