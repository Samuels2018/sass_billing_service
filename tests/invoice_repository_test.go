package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"sass-billing-service/src/models"
	"sass-billing-service/src/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewInvoiceRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repositories.NewInvoiceRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestGetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)

		// Mock data
		expectedID := 1
		expectedInvoice := &models.Invoice{
			ID:            expectedID,
			UserID:        123,
			Amount:        100.50,
			Description:   "Test invoice",
			Status:        "pending",
			PaymentMethod: "credit_card",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Set up expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "description", "status", "payment_method", "created_at", "updated_at"}).
			AddRow(
				expectedInvoice.ID,
				expectedInvoice.UserID,
				expectedInvoice.Amount,
				expectedInvoice.Description,
				expectedInvoice.Status,
				expectedInvoice.PaymentMethod,
				expectedInvoice.CreatedAt,
				expectedInvoice.UpdatedAt,
			)

		mock.ExpectQuery(`SELECT id, user_id, amount, description, status, payment_method, created_at, updated_at 
			FROM invoices WHERE id = \$1`).
			WithArgs(expectedID).
			WillReturnRows(rows)

		// Execute
		ctx := context.Background()
		result, err := repo.GetByID(ctx, expectedID)

		// Validate
		assert.NoError(t, err)
		assert.Equal(t, expectedInvoice, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)
		expectedID := 999

		mock.ExpectQuery(`SELECT id, user_id, amount, description, status, payment_method, created_at, updated_at 
			FROM invoices WHERE id = \$1`).
			WithArgs(expectedID).
			WillReturnError(sql.ErrNoRows)

		ctx := context.Background()
		result, err := repo.GetByID(ctx, expectedID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)
		expectedID := 1
		expectedError := errors.New("database error")

		mock.ExpectQuery(`SELECT id, user_id, amount, description, status, payment_method, created_at, updated_at 
			FROM invoices WHERE id = \$1`).
			WithArgs(expectedID).
			WillReturnError(expectedError)

		ctx := context.Background()
		result, err := repo.GetByID(ctx, expectedID)

		assert.Nil(t, result)
		assert.EqualError(t, err, expectedError.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetByUserID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)
		userID := 123

		// Mock data
		expectedInvoices := []models.Invoice{
			{
				ID:            1,
				UserID:        userID,
				Amount:        100.50,
				Description:   "Test invoice 1",
				Status:        "pending",
				PaymentMethod: "credit_card",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            2,
				UserID:        userID,
				Amount:        200.75,
				Description:   "Test invoice 2",
				Status:        "paid",
				PaymentMethod: "paypal",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}

		// Set up expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "description", "status", "payment_method", "created_at", "updated_at"})
		for _, inv := range expectedInvoices {
			rows.AddRow(
				inv.ID,
				inv.UserID,
				inv.Amount,
				inv.Description,
				inv.Status,
				inv.PaymentMethod,
				inv.CreatedAt,
				inv.UpdatedAt,
			)
		}

		mock.ExpectQuery(`SELECT id, user_id, amount, description, status, payment_method, created_at, updated_at 
			FROM invoices WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(rows)

		// Execute
		ctx := context.Background()
		result, err := repo.GetByUserID(ctx, userID)

		// Validate
		assert.NoError(t, err)
		assert.Equal(t, expectedInvoices, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)
		userID := 999

		// Set up expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "description", "status", "payment_method", "created_at", "updated_at"})

		mock.ExpectQuery(`SELECT id, user_id, amount, description, status, payment_method, created_at, updated_at 
			FROM invoices WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(rows)

		// Execute
		ctx := context.Background()
		result, err := repo.GetByUserID(ctx, userID)

		// Validate
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)
		userID := 123
		expectedError := errors.New("database error")

		mock.ExpectQuery(`SELECT id, user_id, amount, description, status, payment_method, created_at, updated_at 
			FROM invoices WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(expectedError)

		ctx := context.Background()
		result, err := repo.GetByUserID(ctx, userID)

		assert.Nil(t, result)
		assert.EqualError(t, err, expectedError.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ScanError", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)
		userID := 123

		// Set up expectations with invalid data (missing columns)
		rows := sqlmock.NewRows([]string{"id", "user_id", "amount"}).
			AddRow(1, userID, 100.50)

		mock.ExpectQuery(`SELECT id, user_id, amount, description, status, payment_method, created_at, updated_at 
			FROM invoices WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(rows)

		ctx := context.Background()
		result, err := repo.GetByUserID(ctx, userID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)

		now := time.Now()
		request := &models.CreateInvoiceRequest{
			UserID:        123,
			Amount:        100.50,
			Description:   "Test invoice",
			PaymentMethod: "credit_card",
		}

		expectedInvoice := &models.Invoice{
			ID:            1,
			UserID:        request.UserID,
			Amount:        request.Amount,
			Description:   request.Description,
			Status:        "pending",
			PaymentMethod: request.PaymentMethod,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		// Set up expectations
		mock.ExpectQuery(`INSERT INTO invoices \(user_id, amount, description, status, payment_method, created_at, updated_at\)
			VALUES \(\$1, \$2, \$3, 'pending', \$4, \$5, \$5\) 
			RETURNING id, user_id, amount, description, status, payment_method, created_at, updated_at`).
			WithArgs(
				request.UserID,
				request.Amount,
				request.Description,
				request.PaymentMethod,
				sqlmock.AnyArg(), // For timestamp
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "user_id", "amount", "description", "status", "payment_method", "created_at", "updated_at"}).
					AddRow(
						expectedInvoice.ID,
						expectedInvoice.UserID,
						expectedInvoice.Amount,
						expectedInvoice.Description,
						expectedInvoice.Status,
						expectedInvoice.PaymentMethod,
						expectedInvoice.CreatedAt,
						expectedInvoice.UpdatedAt,
					),
			)

		// Execute
		ctx := context.Background()
		result, err := repo.Create(ctx, request)

		// Validate
		assert.NoError(t, err)
		assert.Equal(t, expectedInvoice.ID, result.ID)
		assert.Equal(t, expectedInvoice.UserID, result.UserID)
		assert.Equal(t, expectedInvoice.Amount, result.Amount)
		assert.Equal(t, expectedInvoice.Description, result.Description)
		assert.Equal(t, "pending", result.Status)
		assert.Equal(t, expectedInvoice.PaymentMethod, result.PaymentMethod)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)

		request := &models.CreateInvoiceRequest{
			UserID:        123,
			Amount:        100.50,
			Description:   "Test invoice",
			PaymentMethod: "credit_card",
		}

		expectedError := errors.New("database error")

		mock.ExpectQuery(`INSERT INTO invoices \(user_id, amount, description, status, payment_method, created_at, updated_at\)
			VALUES \(\$1, \$2, \$3, 'pending', \$4, \$5, \$5\) 
			RETURNING id, user_id, amount, description, status, payment_method, created_at, updated_at`).
			WithArgs(
				request.UserID,
				request.Amount,
				request.Description,
				request.PaymentMethod,
				sqlmock.AnyArg(),
			).
			WillReturnError(expectedError)

		ctx := context.Background()
		result, err := repo.Create(ctx, request)

		assert.Nil(t, result)
		assert.EqualError(t, err, expectedError.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ScanError", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repositories.NewInvoiceRepository(db)

		request := &models.CreateInvoiceRequest{
			UserID:        123,
			Amount:        100.50,
			Description:   "Test invoice",
			PaymentMethod: "credit_card",
		}

		// Set up expectations with incomplete data
		mock.ExpectQuery(`INSERT INTO invoices \(user_id, amount, description, status, payment_method, created_at, updated_at\)
			VALUES \(\$1, \$2, \$3, 'pending', \$4, \$5, \$5\) 
			RETURNING id, user_id, amount, description, status, payment_method, created_at, updated_at`).
			WithArgs(
				request.UserID,
				request.Amount,
				request.Description,
				request.PaymentMethod,
				sqlmock.AnyArg(),
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "user_id"}). // Missing columns
										AddRow(1, request.UserID),
			)

		ctx := context.Background()
		result, err := repo.Create(ctx, request)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
