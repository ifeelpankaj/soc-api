package maintenancesvc

import (
	"context"
	"fmt"
	"go-server/internal/models"
	service "go-server/internal/services"
)

type DocumentBillReader interface {
	Get(context.Context, models.MaintenanceBillFilter) (models.MaintenanceBill, error)
}
type DocumentPaymentReader interface {
	Payment(context.Context, int64, int64, int64, bool) (models.UPIPayment, error)
}
type DocumentRenderer interface {
	Invoice(context.Context, models.MaintenanceBill) ([]byte, error)
	Receipt(context.Context, models.MaintenanceBill, models.UPIPayment) ([]byte, error)
}
type DocumentService struct {
	bills       DocumentBillReader
	payments    DocumentPaymentReader
	renderer    DocumentRenderer
	operational OperationalGuard
}

func NewDocumentService(b DocumentBillReader, p DocumentPaymentReader, r DocumentRenderer, g OperationalGuard) *DocumentService {
	if b == nil || p == nil || r == nil || g == nil {
		panic("maintenanceSvc: required dependency is nil")
	}
	return &DocumentService{b, p, r, g}
}
func (s *DocumentService) Invoice(ctx context.Context, society, user, id int64, resident bool) (models.MaintenancePDF, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	var result models.MaintenancePDF
	if err := s.operational.EnsureSocietyOperational(ctx, society); err != nil {
		return result, err
	}
	bill, err := s.bills.Get(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: user, ID: id, Resident: resident})
	if err != nil {
		return result, err
	}
	if reader, ok := s.bills.(interface {
		Outstanding(context.Context, int64, int64, int64) (models.MaintenanceOutstanding, error)
	}); ok {
		summary, err := reader.Outstanding(ctx, society, user, bill.FlatID)
		if err != nil {
			return result, err
		}
		bill.Outstanding = &summary
	}
	data, err := s.renderer.Invoice(ctx, bill)
	if err != nil {
		return result, err
	}
	return models.MaintenancePDF{Bytes: data, Filename: fmt.Sprintf("maintenance-invoice-%d.pdf", bill.ID)}, nil
}
func (s *DocumentService) Receipt(ctx context.Context, society, user, id int64, resident bool) (models.MaintenancePDF, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	var result models.MaintenancePDF
	if err := s.operational.EnsureSocietyOperational(ctx, society); err != nil {
		return result, err
	}
	payment, err := s.payments.Payment(ctx, society, user, id, resident)
	if err != nil {
		return result, err
	}
	bill, err := s.bills.Get(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: user, ID: payment.BillID, Resident: resident})
	if err != nil {
		return result, err
	}
	data, err := s.renderer.Receipt(ctx, bill, payment)
	if err != nil {
		return result, err
	}
	return models.MaintenancePDF{Bytes: data, Filename: fmt.Sprintf("maintenance-receipt-%d.pdf", payment.ID)}, nil
}
