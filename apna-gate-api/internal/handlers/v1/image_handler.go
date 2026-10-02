package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go-server/internal/models"
	imagesvc "go-server/internal/services/imageSvc"
	"go-server/pkg/utils"
	"net/http"
)

type imageExecutor interface {
	Execute(context.Context, string, models.ImageTarget, http.ResponseWriter, *http.Request) (*models.ImageView, error)
}
type ImageHandler struct{ service imageExecutor }

func NewImageHandler(service imageExecutor) *ImageHandler { return &ImageHandler{service: service} }

func (h *ImageHandler) execute(c *gin.Context) {
	actor, ok := currentUserID(c)
	if !ok {
		return
	}
	target := models.ImageTarget{ActorID: actor}
	if c.Param("entryId") != "" {
		var valid bool
		target.SocietyID, target.EntryID, valid = visitorEntryPath(c)
		if !valid {
			return
		}
	}
	if c.Param("flatId") != "" {
		var valid bool
		target.FlatID, valid = parsePathInt64(c, "flatId")
		if !valid {
			return
		}
	}
	c.Request.Header.Set("X-Request-ID", c.GetString("request_id"))
	view, err := h.service.Execute(c.Request.Context(), c.Request.Method, target, c.Writer, c.Request)
	if err != nil {
		if errors.Is(err, imagesvc.ErrImageBusy) {
			c.Header("Retry-After", "2")
		}
		handleServiceError(c, err)
		return
	}
	if c.Request.Method == http.MethodDelete {
		c.Status(http.StatusNoContent)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Image ready", view)
}

// PutAvatar godoc
// @Summary PutAvatar
// @Description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
// @Tags Images
// @Produce json
// @Security AccessToken
// @Accept multipart/form-data
// @Param file formData file true "JPEG or PNG, maximum 5 MiB and 25 megapixels"
// @Success 200 {object} models.ImageViewAPIResponse
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 413 {object} models.ErrorResponseDoc
// @Failure 502 {object} models.ErrorResponseDoc
// @Failure 503 {object} models.ErrorResponseDoc
// @Router /v1/auth/profile/avatar [put]
func (h *ImageHandler) PutAvatar(c *gin.Context) { h.execute(c) }

// GetAvatar godoc
// @Summary GetAvatar
// @Param variant query string false "Server-selected image view" Enums(list,avatar,detail,original) default(original)
// @Description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
// @Tags Images
// @Produce json
// @Security AccessToken
// @Success 200 {object} models.ImageViewAPIResponse
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 413 {object} models.ErrorResponseDoc
// @Failure 502 {object} models.ErrorResponseDoc
// @Failure 503 {object} models.ErrorResponseDoc
// @Router /v1/auth/profile/avatar [get]
func (h *ImageHandler) GetAvatar(c *gin.Context) { h.execute(c) }

// DeleteAvatar godoc
// @Summary DeleteAvatar
// @Description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
// @Tags Images
// @Produce json
// @Security AccessToken
// @Success 204 "Image reference removed"
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 413 {object} models.ErrorResponseDoc
// @Failure 502 {object} models.ErrorResponseDoc
// @Failure 503 {object} models.ErrorResponseDoc
// @Router /v1/auth/profile/avatar [delete]
func (h *ImageHandler) DeleteAvatar(c *gin.Context) { h.execute(c) }

// PutVisitorPhoto godoc
// @Summary PutVisitorPhoto
// @Description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
// @Tags Images
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Accept multipart/form-data
// @Param file formData file true "JPEG or PNG, maximum 5 MiB and 25 megapixels"
// @Success 200 {object} models.ImageViewAPIResponse
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 413 {object} models.ErrorResponseDoc
// @Failure 502 {object} models.ErrorResponseDoc
// @Failure 503 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/photo [put]
func (h *ImageHandler) PutVisitorPhoto(c *gin.Context) { h.execute(c) }

// GetVisitorPhoto godoc
// @Summary GetVisitorPhoto
// @Param variant query string false "Server-selected image view" Enums(list,avatar,detail,original) default(original)
// @Description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
// @Tags Images
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Success 200 {object} models.ImageViewAPIResponse
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 413 {object} models.ErrorResponseDoc
// @Failure 502 {object} models.ErrorResponseDoc
// @Failure 503 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/photo [get]
func (h *ImageHandler) GetVisitorPhoto(c *gin.Context) { h.execute(c) }

// DeleteVisitorPhoto godoc
// @Summary DeleteVisitorPhoto
// @Description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
// @Tags Images
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Success 204 "Image reference removed"
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 413 {object} models.ErrorResponseDoc
// @Failure 502 {object} models.ErrorResponseDoc
// @Failure 503 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/photo [delete]
func (h *ImageHandler) DeleteVisitorPhoto(c *gin.Context) { h.execute(c) }

// GetResidentVisitorPhoto godoc
// @Summary GetResidentVisitorPhoto
// @Param variant query string false "Server-selected image view" Enums(list,avatar,detail,original) default(original)
// @Description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
// @Tags Images
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.ImageViewAPIResponse
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 413 {object} models.ErrorResponseDoc
// @Failure 502 {object} models.ErrorResponseDoc
// @Failure 503 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-entries/{entryId}/photo [get]
func (h *ImageHandler) GetResidentVisitorPhoto(c *gin.Context) { h.execute(c) }
