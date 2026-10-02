package contracts

import (
	"context"
	"go-server/internal/models"
)

type ImageRepository interface {
	Read(context.Context, models.ImageTarget) (models.StoredImage, error)
	Replace(context.Context, models.ImageTarget, models.StoredImage, func(context.Context) error) (models.StoredImage, error)
	Resolve(context.Context, models.ImageTarget) (models.StoredImage, error)
}
