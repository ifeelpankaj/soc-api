package imagesvc

import (
	"context"
	"go-server/internal/models"
	"net/http"
)

type ImageService interface {
	Execute(ctx context.Context, method string, target models.ImageTarget, writer http.ResponseWriter, request *http.Request) (view *models.ImageView, err error)
}

var _ ImageService = (*Service)(nil)
