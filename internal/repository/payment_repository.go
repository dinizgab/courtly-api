package repository

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/gateway/openpix"
)

var (
	//go:embed sql/payment/create_pix_subaccount.sql
	createPixSubaccountQuery string
	//go:embed sql/payment/get_subaccount_pix_key_by_company_id.sql
	getSubaccountPixKeyByCompanyIDQuery string
	//go:embed sql/payment/create_charge.sql
	createChargeQuery string
)

type PaymentRepository interface {
	CreateSubaccount(ctx context.Context, subaccount entity.Subaccount) error
	GetSubaccountPixKeyByCompanyID(ctx context.Context, companyId string) (string, error)
	CreateCharge(ctx context.Context, companyId string, charge openpix.Charge) error
}

type paymentRepositoryImpl struct {
	db database.Database
}

func NewPaymentRepository(db database.Database) PaymentRepository {
	return &paymentRepositoryImpl{
		db: db,
	}
}

func (r *paymentRepositoryImpl) CreateSubaccount(ctx context.Context, subaccount entity.Subaccount) error {
	_, err := r.db.Exec(ctx, createPixSubaccountQuery, subaccount.CompanyID, subaccount.PixKey)
	if err != nil {
        return fmt.Errorf("paymentRepositoryImpl.CreateSubaccount - failed to create subaccount: %w", err)
	}

	return nil
}

func (r *paymentRepositoryImpl) GetSubaccountPixKeyByCompanyID(ctx context.Context, companyId string) (string, error) {
	var pixKey string
	err := r.db.QueryRow(ctx, getSubaccountPixKeyByCompanyIDQuery, companyId).Scan(&pixKey)
	if err != nil {
        return "", fmt.Errorf("paymentRepositoryImpl.GetSubaccountPixKeyByCompanyID - failed to get subaccount pix key: %w", err)
	}

	return pixKey, nil
}

func (r *paymentRepositoryImpl) CreateCharge(ctx context.Context, companyId string, charge openpix.Charge) error {
	bookingId := strings.Replace(charge.CorrelationID, "booking-", "", 1)
	_, err := r.db.Exec(
		ctx,
		createChargeQuery,
		companyId,
		bookingId,
		charge.CorrelationID,
		charge.PaymentLinkID,
		charge.PaymentLinkURL,
		charge.Brcode,
		charge.Value,
		charge.ExpiresDate,
	)
	if err != nil {
        return fmt.Errorf("paymentRepositoryImpl.CreateCharge - failed to create charge: %w", err)
	}

	return nil
}
