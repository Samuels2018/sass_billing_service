package controllers

import (
	"sass-billing-service/src/models"
	"sass-billing-service/src/services"
	"sass-billing-service/src/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type InvoiceController struct {
	service *services.InvoiceService
}

func NewInvoiceController(service *services.InvoiceService) *InvoiceController {
	return &InvoiceController{service: service}
}

func (c *InvoiceController) GetInvoices(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Query("user_id"))
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid user ID")
	}

	invoices, err := c.service.GetInvoicesByUserID(ctx.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, invoices)
}

func (c *InvoiceController) GetInvoice(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid invoice ID")
	}

	invoice, err := c.service.GetInvoiceByID(ctx.Context(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusNotFound, "Invoice not found")
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, invoice)
}

func (c *InvoiceController) CreateInvoice(ctx *fiber.Ctx) error {
	var req models.CreateInvoiceRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validar campos requeridos
	if req.UserID == 0 || req.Amount <= 0 || req.Description == "" || req.PaymentMethod == "" {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Missing required fields")
	}

	invoice, err := c.service.CreateInvoice(ctx.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(ctx, fiber.StatusCreated, invoice)
}
