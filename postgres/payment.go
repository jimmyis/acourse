package postgres

import (
	"context"
	"database/sql"

	"github.com/acoshift/acourse/payment"
	"github.com/lib/pq"
)

// NewPaymentRepository creates new payment repository
func NewPaymentRepository(db *sql.DB) (payment.Repository, error) {
	r := paymentRepository{db}
	return r, nil
}

type paymentRepository struct {
	db *sql.DB
}

func (r *paymentRepository) Store(ctx context.Context, payment payment.Payment) error {
}

func (r *paymentRepository) FindID(ctx context.Context, id string) (*payment.Payment, error) {
}

func (r *paymentRepository) List(ctx context.Context, status []payment.Status, limit, offset int64) ([]*payment.Payment, error) {
}

func (r *paymentRepository) Count(ctx context.Context, status []payment.Status) (int64, error) {
	var cnt int64
	err := r.db.QueryRowContext(ctx, `
		select count(*)
		from payments
		where status = any($1)
	`, pq.Array(status)).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
