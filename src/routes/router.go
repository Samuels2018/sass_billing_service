package router

import (
	"sass-billing-service/src/controllers"
	"sass-billing-service/src/helpers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, invoiceController *controllers.InvoiceController) {
	invoices := app.Group("/invoices")
	{
		invoices.Get("/", helpers.AuthMiddleware, invoiceController.GetInvoices)
		invoices.Post("/", helpers.AuthMiddleware, invoiceController.CreateInvoice)
		invoices.Get("/:id", helpers.AuthMiddleware, invoiceController.GetInvoice)
	}
}
