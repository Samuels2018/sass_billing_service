package services

import (
	"context"
	"sass-billing-service/src/models"
	"sass-billing-service/src/repositories"
)

type InvoiceService struct {
	repo *repositories.InvoiceRepository
}

func NewInvoiceService(repo *repositories.InvoiceRepository) *InvoiceService {
	return &InvoiceService{repo: repo}
}

func (s *InvoiceService) GetInvoiceByID(ctx context.Context, id int) (*models.Invoice, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *InvoiceService) GetInvoicesByUserID(ctx context.Context, userID int) ([]models.Invoice, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, req *models.CreateInvoiceRequest) (*models.Invoice, error) {
	return s.repo.Create(ctx, req)
}
