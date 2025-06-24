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
	//go:embed sql/payment/expire_payment.sql
	expirePaymentQuery string
    //go:embed sql/payment/get_booking_charge_information_by_booking_id.sql
    getBookingChargeInformationByBookingId string
    //go:embed sql/payment/get_payment_by_booking_id.sql
    getPaymentByBoookingIdQuery string
    //go:embed sql/payment/save_refund_request.sql
    saveRefundRequestQuery string
)

type PaymentRepository interface {
	CreateSubaccount(ctx context.Context, subaccount entity.Subaccount) error
	GetSubaccountPixKeyByCompanyID(ctx context.Context, companyId string) (string, error)
	CreateCharge(ctx context.Context, companyId string, charge openpix.Charge) error
    ConfirmPayment(ctx context.Context, charge openpix.Charge) error
    GetBookingPaymentStatusByID(ctx context.Context, id string) (string, error)
    CreateWithdrawRequest(ctx context.Context, companyId string, withdraw openpix.Withdraw) error
	ExpirePayment(ctx context.Context, charge openpix.Charge) error
    GetBookingChargeInformation(ctx context.Context, id string) (entity.Payment, error)
    GetPaymentByBookingID(ctx context.Context, id string) (entity.Payment, error)
    SaveRefundRequest(ctx context.Context, bookingId string, refund openpix.Refund) error
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
        charge.QrCodeImage,
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

func (r *paymentRepositoryImpl) GetBookingChargeInformation(ctx context.Context, id string) (entity.Payment, error) {
    var payment entity.Payment
    err := r.db.QueryRow(ctx, getBookingChargeInformationByBookingId, id).Scan(
        &payment.BrCode,
        &payment.QrCodeImage,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return entity.Payment{}, fmt.Errorf("paymentRepositoryImpl.GetPaymentBookingPaymentInformation - booking payment not found: %w", err)
        }
        return entity.Payment{}, fmt.Errorf("paymentRepositoryImpl.GetPaymentBookingPaymentInformation - failed to get booking payment information: %w", err)
    }

    return payment, nil
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

func (r *paymentRepositoryImpl) ExpirePayment(ctx context.Context, charge openpix.Charge) error {
	_, err := r.db.Exec(
		ctx,
		expirePaymentQuery,
		charge.CorrelationID,
	)
	if err != nil {
		return fmt.Errorf("paymentRepositoryImpl.ExpirePayment - failed to expire payment: %w", err)
	}

	return nil
}

func (r *paymentRepositoryImpl) GetPaymentByBookingID(ctx context.Context, id string) (entity.Payment, error) {
    var payment entity.Payment
    err := r.db.QueryRow(ctx, getPaymentByBoookingIdQuery, id).Scan(
        &payment.ID,
        &payment.CorrelationID,
        &payment.BookingID,
        &payment.PaidAt,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return entity.Payment{}, fmt.Errorf("paymentRepositoryImpl.GetPaymentByBookingID - payment not found: %w", err)
        }
        return entity.Payment{}, fmt.Errorf("paymentRepositoryImpl.GetPaymentByBookingID - failed to get payment by booking ID: %w", err)
    }

    return payment, nil
}

func (r *paymentRepositoryImpl) SaveRefundRequest(ctx context.Context, bookingId string, refund openpix.Refund) error {
    _, err := r.db.Exec(
        ctx,
        saveRefundRequestQuery,
        bookingId,
        refund.EndToEndID,
        refund.Status,
        // TODO = Check if this variable name 
        refund.RefundedAt,
    )
    if err != nil {
        return fmt.Errorf("paymentRepositoryImpl.SaveRefundRequest - failed to save refund request: %w", err)
    }

    return nil
}

