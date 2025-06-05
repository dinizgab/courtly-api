package repository

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/entity"
	"github.com/dinizgab/booking-mvp/internal/gateway/openpix"
	"github.com/jackc/pgx/v5"
)

var (
	//go:embed sql/payment/create_pix_subaccount.sql
	createPixSubaccountQuery string
	//go:embed sql/payment/get_subaccount_pix_key_by_company_id.sql
	getSubaccountPixKeyByCompanyIDQuery string
	//go:embed sql/payment/create_charge.sql
	createChargeQuery string
    //go:embed sql/payment/confirm_booking_payment.sql
    confirmPaymentQuery string
    //go:embed sql/payment/get_booking_payment_status_by_id.sql
    getBookingPaymentStatusByIDQuery string
    //go:embed sql/payment/create_withdraw_request.sql
    createWithdrawRequestQuery string
)

type PaymentRepository interface {
	CreateSubaccount(ctx context.Context, subaccount entity.Subaccount) error
	GetSubaccountPixKeyByCompanyID(ctx context.Context, companyId string) (string, error)
	CreateCharge(ctx context.Context, companyId string, charge openpix.Charge) error
    ConfirmPayment(ctx context.Context, charge openpix.Charge) error
    GetBookingPaymentStatusByID(ctx context.Context, id string) (string, error)
    CreateWithdrawRequest(ctx context.Context, companyId string, withdraw openpix.Withdraw) error
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

func (r *paymentRepositoryImpl) ConfirmPayment(ctx context.Context, charge openpix.Charge) error {
    _, err := r.db.Exec(
        ctx,
        confirmPaymentQuery,
        charge.CorrelationID,
        charge.PaidAt,
    )
    if err != nil {
        return fmt.Errorf("paymentRepositoryImpl.ConfirmPayment - failed to confirm payment: %w", err)
    }

    return nil
}

func (r *paymentRepositoryImpl) GetBookingPaymentStatusByID(ctx context.Context, id string) (string, error) {
    var status string
    err := r.db.QueryRow(ctx, getBookingPaymentStatusByIDQuery, id).Scan(&status)
    if err != nil {
        if err == pgx.ErrNoRows {
            return "", fmt.Errorf("paymentRepositoryImpl.GetBookingPaymentStatusByID - booking payment not found: %w", err)
        }
        return "", fmt.Errorf("paymentRepositoryImpl.GetBookingPaymentStatusByID - failed to get booking payment status: %w", err)
    }

    return status, nil
}

func (r *paymentRepositoryImpl) CreateWithdrawRequest(ctx context.Context, companyId string, withdraw openpix.Withdraw) error {
    _, err := r.db.Exec(
        ctx,
        createWithdrawRequestQuery,
        companyId,
        withdraw.CorrelationId,
        withdraw.Value,
    )
    if err != nil {
        return fmt.Errorf("paymentRepositoryImpl.CreateWithdrawRequest - failed to create withdraw request: %w", err)
    }

    return nil
}
