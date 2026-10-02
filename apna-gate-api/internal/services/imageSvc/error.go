package imagesvc

import "go-server/internal/models"

var (
	ErrImageVariant        = models.NewAppError("IMAGE_VARIANT_INVALID", "Choose list, avatar, detail, or original", 400, nil)
	ErrImageUnavailable    = models.NewAppError("IMAGE_UNAVAILABLE", "Image service is unavailable", 503, nil)
	ErrImageBusy           = models.NewAppError("IMAGE_BUSY", "Image service is busy; retry shortly", 503, nil)
	ErrImageInvalid        = models.NewAppError("IMAGE_INVALID", "Supply exactly one valid JPEG or PNG file", 400, nil)
	ErrImageTooLarge       = models.NewAppError("IMAGE_TOO_LARGE", "Image exceeds the upload or pixel limit", 413, nil)
	ErrImageNotFound       = models.NewAppError("IMAGE_NOT_FOUND", "Managed image not found", 404, nil)
	ErrImageTargetNotFound = models.NewAppError("IMAGE_TARGET_NOT_FOUND", "Image owner or entry not found", 404, nil)
	ErrImageForbidden      = models.NewAppError("IMAGE_FORBIDDEN", "Image access is not permitted", 403, nil)
	ErrImageState          = models.NewAppError("IMAGE_ENTRY_STATE", "Visitor photo cannot be uploaded in the current entry state", 409, nil)
	ErrImageProvider       = models.NewAppError("IMAGE_PROVIDER_FAILED", "Image storage operation failed", 502, nil)
	ErrImagePersistence    = models.NewAppError("IMAGE_PERSISTENCE_FAILED", "Image could not be saved; refresh before retrying", 503, nil)
)
