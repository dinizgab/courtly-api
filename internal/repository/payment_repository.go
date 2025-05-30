package repository

import (
    _ "embed"
	"context"

	"github.com/dinizgab/booking-mvp/internal/database"
	"github.com/dinizgab/booking-mvp/internal/entity"
)

//go:embed sql/payment/create_pix_subaccount.sql
var createPixSubaccountQuery string

type PaymentRepository interface {
    CreateSubaccount(ctx context.Context, subaccount entity.Subaccount) error
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
