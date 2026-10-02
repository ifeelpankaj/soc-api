package subscriptionsvc

import (
	"net/http"

	"go-server/internal/models"
)

var (
	ErrSubscriptionNotFound = models.NewAppError("SUBSCRIPTION_NOT_FOUND", "subscription not found", http.StatusNotFound, nil)
	ErrSubscriptionConflict = models.NewAppError("SUBSCRIPTION_CONFLICT", "subscription conflict", http.StatusConflict, nil)
	ErrSubscriptionRequired = models.NewAppError("SUBSCRIPTION_REQUIRED", "active subscription required", http.StatusPaymentRequired, nil)
	ErrSubscriptionExpired  = models.NewAppError("SUBSCRIPTION_EXPIRED", "subscription is expired", http.StatusPaymentRequired, nil)
	ErrFeatureDisabled      = models.NewAppError("FEATURE_DISABLED", "feature is not enabled for this subscription", http.StatusForbidden, nil)
	ErrQuotaExceeded        = models.NewAppError("SUBSCRIPTION_QUOTA_EXCEEDED", "subscription quota exceeded", http.StatusPaymentRequired, nil)
	ErrInvalidSubscription  = models.NewAppError("INVALID_SUBSCRIPTION_REQUEST", "invalid subscription request", http.StatusBadRequest, nil)
)
