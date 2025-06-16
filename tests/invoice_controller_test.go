package controllers

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"sass-billing-service/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInvoiceService es un mock del servicio para testing
type MockInvoiceService struct {
	mock.Mock
}

func (m *MockInvoiceService) GetInvoiceByID(ctx context.Context, id int) (*models.Invoice, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Invoice), args.Error(1)
}

func (m *MockInvoiceService) GetInvoicesByUserID(ctx context.Context, userID int) ([]models.Invoice, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Invoice), args.Error(1)
}

func (m *MockInvoiceService) CreateInvoice(ctx context.Context, req *models.CreateInvoiceRequest) (*models.Invoice, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.Invoice), args.Error(1)
}

func TestNewInvoiceController(t *testing.T) {
	mockService := new(MockInvoiceService)
	controller := NewInvoiceController(mockService)

	assert.NotNil(t, controller)
	assert.Equal(t, mockService, controller.service)
}

func TestInvoiceController_GetInvoices(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		// Configurar el mock
		expectedUserID := 123
		expectedInvoices := []models.Invoice{
			{ID: 1, UserID: expectedUserID, Amount: 100.50},
			{ID: 2, UserID: expectedUserID, Amount: 200.75},
		}

		mockService.On("GetInvoicesByUserID", mock.Anything, expectedUserID).
			Return(expectedInvoices, nil)

		// Crear contexto de prueba
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.Request().URI().SetQueryString("user_id=" + strconv.Itoa(expectedUserID))

		// Ejecutar
		err := controller.GetInvoices(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		// Crear contexto de prueba con user_id inválido
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.Request().URI().SetQueryString("user_id=abc")

		// Ejecutar
		err := controller.GetInvoices(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())
		mockService.AssertNotCalled(t, "GetInvoicesByUserID")
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		expectedUserID := 123
		expectedError := errors.New("service error")

		mockService.On("GetInvoicesByUserID", mock.Anything, expectedUserID).
			Return([]models.Invoice{}, expectedError)

		// Crear contexto de prueba
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.Request().URI().SetQueryString("user_id=" + strconv.Itoa(expectedUserID))

		// Ejecutar
		err := controller.GetInvoices(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())
		mockService.AssertExpectations(t)
	})
}

func TestInvoiceController_GetInvoice(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		expectedID := 1
		expectedInvoice := &models.Invoice{
			ID:     expectedID,
			UserID: 123,
			Amount: 100.50,
		}

		mockService.On("GetInvoiceByID", mock.Anything, expectedID).
			Return(expectedInvoice, nil)

		// Crear contexto de prueba
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.SetParams("id", strconv.Itoa(expectedID))

		// Ejecutar
		err := controller.GetInvoice(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidInvoiceID", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		// Crear contexto de prueba con ID inválido
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.SetParams("id", "abc")

		// Ejecutar
		err := controller.GetInvoice(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())
		mockService.AssertNotCalled(t, "GetInvoiceByID")
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		expectedID := 999
		mockService.On("GetInvoiceByID", mock.Anything, expectedID).
			Return(&models.Invoice{}, errors.New("not found"))

		// Crear contexto de prueba
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.SetParams("id", strconv.Itoa(expectedID))

		// Ejecutar
		err := controller.GetInvoice(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, ctx.Response().StatusCode())
		mockService.AssertExpectations(t)
	})
}

func TestInvoiceController_CreateInvoice(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		req := &models.CreateInvoiceRequest{
			UserID:        123,
			Amount:        100.50,
			Description:   "Test invoice",
			PaymentMethod: "credit_card",
		}

		expectedInvoice := &models.Invoice{
			ID:            1,
			UserID:        req.UserID,
			Amount:        req.Amount,
			Description:   req.Description,
			PaymentMethod: req.PaymentMethod,
			Status:        "pending",
		}

		mockService.On("CreateInvoice", mock.Anything, req).
			Return(expectedInvoice, nil)

		// Crear contexto de prueba con body
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.Request().Header.SetContentType("application/json")
		ctx.Request().SetBody([]byte(`{
			"user_id": 123,
			"amount": 100.50,
			"description": "Test invoice",
			"payment_method": "credit_card"
		}`))

		// Ejecutar
		err := controller.CreateInvoice(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, ctx.Response().StatusCode())
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		// Crear contexto de prueba con body inválido
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.Request().Header.SetContentType("application/json")
		ctx.Request().SetBody([]byte(`{ invalid json }`))

		// Ejecutar
		err := controller.CreateInvoice(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())
		mockService.AssertNotCalled(t, "CreateInvoice")
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		// Crear contexto de prueba con campos faltantes
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.Request().Header.SetContentType("application/json")
		ctx.Request().SetBody([]byte(`{
			"user_id": 123,
			"amount": 100.50
		}`)) // Faltan description y payment_method

		// Ejecutar
		err := controller.CreateInvoice(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())
		mockService.AssertNotCalled(t, "CreateInvoice")
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockInvoiceService)
		controller := NewInvoiceController(mockService)

		req := &models.CreateInvoiceRequest{
			UserID:        123,
			Amount:        100.50,
			Description:   "Test invoice",
			PaymentMethod: "credit_card",
		}

		expectedError := errors.New("service error")

		mockService.On("CreateInvoice", mock.Anything, req).
			Return(&models.Invoice{}, expectedError)

		// Crear contexto de prueba
		app := fiber.New()
		ctx := app.AcquireCtx(&fiber.Ctx{})
		defer app.ReleaseCtx(ctx)
		ctx.Request().Header.SetContentType("application/json")
		ctx.Request().SetBody([]byte(`{
			"user_id": 123,
			"amount": 100.50,
			"description": "Test invoice",
			"payment_method": "credit_card"
		}`))

		// Ejecutar
		err := controller.CreateInvoice(ctx)

		// Validar
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())
		mockService.AssertExpectations(t)
	})
}
