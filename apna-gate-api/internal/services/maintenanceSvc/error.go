package maintenancesvc

import (
	"errors"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
)

func invalid(message string) error {
	return models.NewAppError("MAINTENANCE_INVALID", message, 400, nil)
}

func paymentConflict(code, message string) error { return models.NewAppError(code, message, 409, nil) }

func paymentError(err error) error {
	if errors.Is(err, contracts.ErrNotFound) {
		return models.NewAppError("PAYMENT_NOT_FOUND", "Payment resource not found", 404, nil)
	}
	if errors.Is(err, contracts.ErrAlreadyExists) {
		return paymentConflict("PAYMENT_CONFLICT", "Reference or financial operation is already reserved")
	}
	return err
}
