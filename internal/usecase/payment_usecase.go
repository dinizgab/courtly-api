package usecase

import (
	"context"
	"fmt"

	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/gateway/openpix"
	"github.com/dinizgab/booking-mvp/internal/repository"
)

type PaymentUsecase interface {
	CreateSubaccount(ctx context.Context, company entity.Company) error
	CreateCharge(ctx context.Context, companyId string, booking entity.Booking) error
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

    fmt.Printf("Charge created successfully: %+v\n", charge)

	return nil
}
