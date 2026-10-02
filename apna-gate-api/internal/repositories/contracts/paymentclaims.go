package contracts

import (
	"context"
)

type PaymentClaimsRepository interface {
	GetPendingUPIClaim(ctx context.Context, arg GetPendingUPIClaimInput) (MaintenancePaymentClaim, error)
	GetUPIClaim(ctx context.Context, arg GetUPIClaimInput) (MaintenancePaymentClaim, error)
	InsertUPIClaim(ctx context.Context, arg InsertUPIClaimInput) (MaintenancePaymentClaim, error)
	ListUPIClaims(ctx context.Context, arg ListUPIClaimsInput) ([]MaintenancePaymentClaim, error)
	MaintenancePendingClaimsCount(ctx context.Context) (int64, error)
	ReleaseUPIClaimReference(ctx context.Context, arg ReleaseUPIClaimReferenceInput) error
	ReviewUPIClaim(ctx context.Context, arg ReviewUPIClaimInput) (MaintenancePaymentClaim, error)
}
