package repository

import (
	"context"
	_ "embed"

	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/entity"
)

var (
	//go:embed sql/payment/create_pix_subaccount.sql
	createPixSubaccountQuery string
    //go:embed sql/payment/get_subaccount_pix_key_by_company_id.sql
    getSubaccountPixKeyByCompanyIDQuery string
)

type PaymentRepository interface {
	CreateSubaccount(ctx context.Context, subaccount entity.Subaccount) error
	GetSubaccountPixKeyByCompanyID(ctx context.Context, companyId string) (string, error)
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
		return err
	}

	return nil
}

func (r *paymentRepositoryImpl) GetSubaccountPixKeyByCompanyID(ctx context.Context, companyId string) (string, error) {
    var pixKey string
    err := r.db.QueryRow(ctx, getSubaccountPixKeyByCompanyIDQuery, companyId).Scan(&pixKey)
    if err != nil {
        return "", err
    }

    return pixKey, nil
}
