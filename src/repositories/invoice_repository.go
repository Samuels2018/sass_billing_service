package repositories

import (
	"context"
	"database/sql"
	"sass-billing-service/src/models"
	"time"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) GetByID(ctx context.Context, id int) (*models.Invoice, error) {
	query := `SELECT id, user_id, amount, description, status, payment_method, created_at, updated_at 
	FROM invoices WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var invoice models.Invoice
	err := row.Scan(
		&invoice.ID,
		&invoice.UserID,
		&invoice.Amount,
		&invoice.Description,
		&invoice.Status,
		&invoice.PaymentMethod,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (r *InvoiceRepository) GetByUserID(ctx context.Context, userID int) ([]models.Invoice, error) {
	query := `SELECT id, user_id, amount, description, status, payment_method, created_at, updated_at 
	FROM invoices WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		err := rows.Scan(
			&invoice.ID,
			&invoice.UserID,
			&invoice.Amount,
			&invoice.Description,
			&invoice.Status,
			&invoice.PaymentMethod,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (r *InvoiceRepository) Create(ctx context.Context, invoice *models.CreateInvoiceRequest) (*models.Invoice, error) {
	query := `INSERT INTO invoices (user_id, amount, description, status, payment_method, created_at, updated_at)
	VALUES ($1, $2, $3, 'pending', $4, $5, $5) 
	RETURNING id, user_id, amount, description, status, payment_method, created_at, updated_at`

	now := time.Now()
	row := r.db.QueryRowContext(ctx, query,
		invoice.UserID,
		invoice.Amount,
		invoice.Description,
		invoice.PaymentMethod,
		now,
	)

	var createdInvoice models.Invoice
	err := row.Scan(
		&createdInvoice.ID,
		&createdInvoice.UserID,
		&createdInvoice.Amount,
		&createdInvoice.Description,
		&createdInvoice.Status,
		&createdInvoice.PaymentMethod,
		&createdInvoice.CreatedAt,
		&createdInvoice.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &createdInvoice, nil
}
