package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-server/internal/models"
	hubsvc "go-server/internal/services/hubSvc"
	"go-server/pkg/utils"
)

type HubHandler struct{ service hubsvc.HubService }

func NewHubHandler(service hubsvc.HubService) *HubHandler { return &HubHandler{service} }
func hubRequest(c *gin.Context) (context.Context, int64, int64, bool) {
	society, user, ok := maintenanceIDs(c)
	return c.Request.Context(), society, user, ok
}
func hubResult(c *gin.Context, data any, err error) {
	var rate *hubsvc.RateError
	if errors.As(err, &rate) {
		c.Header("Retry-After", strconv.Itoa(rate.RetryAfter))
		utils.ErrorResponse(c, 429, "HUB_RATE_LIMIT", rate.Error(), nil)
		return
	}
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, 200, "Society Hub request completed", data)
}
func hubBind(c *gin.Context, v any) bool {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 256*1024)
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		hubResult(c, nil, hubsvc.ErrInvalid)
		return false
	}
	var extra any
	if !errors.Is(decoder.Decode(&extra), io.EOF) {
		hubResult(c, nil, hubsvc.ErrInvalid)
		return false
	}
	return true
}
func hubLimit(c *gin.Context) (int, bool) {
	raw := c.DefaultQuery("limit", "20")
	n, e := strconv.Atoi(raw)
	if e != nil || n < 1 || n > 100 {
		hubResult(c, nil, hubsvc.ErrInvalid)
		return 0, false
	}
	return n, true
}

// Timeout bounds every Hub request, including provider uploads.
func (h *HubHandler) Timeout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 50*time.Second)
	defer cancel()
	c.Request = c.Request.WithContext(ctx)
	c.Header("Cache-Control", "no-store")
	c.Next()
}

// Channels godoc
// @Summary List Society Hub channels and unread counts
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.HubResponse{data=[]models.HubChannelSummary}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channels [get]
func (h *HubHandler) Channels(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	v, e := h.service.Channels(ctx, s, u)
	hubResult(c, v, e)
}

// Categories godoc
// @Summary List active community categories
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.HubResponse{data=[]models.HubCategory}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channels/categories [get]
func (h *HubHandler) Categories(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	v, e := h.service.Categories(ctx, s, u)
	hubResult(c, v, e)
}

// Posts godoc
// @Summary List posts with a stable cursor and separate pinned notices
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param channelId path int true "channelId"
// @Param cursor query string false "Opaque resource-scoped cursor"
// @Param limit query int false "Page size" default(20) minimum(1) maximum(100)
// @Param category_id query int false "Community category ID"
// @Success 200 {object} models.HubResponse{data=models.HubPage}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channels/{channelId}/posts [get]
func (h *HubHandler) Posts(c *gin.Context) { h.feed(c, false) }

// Comments godoc
// @Summary List discussion comments including deleted parent placeholders
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Param cursor query string false "Opaque resource-scoped cursor"
// @Param limit query int false "Page size" default(20) minimum(1) maximum(100)
// @Success 200 {object} models.HubResponse{data=models.HubPage}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId}/comments [get]
func (h *HubHandler) Comments(c *gin.Context) { h.feed(c, true) }
func (h *HubHandler) feed(c *gin.Context, comments bool) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	name := "channelId"
	if comments {
		name = "postId"
	}
	id, ok := parsePathInt64(c, name)
	if !ok {
		return
	}
	limit, ok := hubLimit(c)
	if !ok {
		return
	}
	var category *int16
	if raw := c.Query("category_id"); raw != "" {
		n, e := strconv.ParseInt(raw, 10, 16)
		if e != nil || n <= 0 {
			hubResult(c, nil, hubsvc.ErrInvalid)
			return
		}
		v := int16(n)
		category = &v
	}
	v, e := h.service.Feed(ctx, s, u, id, c.Query("cursor"), limit, category, comments)
	hubResult(c, v, e)
}

// Post godoc
// @Summary Get a Society Hub post
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Success 200 {object} models.HubResponse{data=models.HubContent}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId} [get]
func (h *HubHandler) Post(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "postId")
	if !ok {
		return
	}
	v, e := h.service.Post(ctx, s, u, id)
	hubResult(c, v, e)
}

// CreatePost godoc
// @Summary Publish an announcement or community post
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param channelId path int true "channelId"
// @Accept json
// @Param request body models.HubPostRequest true "Request"
// @Success 200 {object} models.HubResponse{data=models.HubContent}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channels/{channelId}/posts [post]
func (h *HubHandler) CreatePost(c *gin.Context) { h.savePost(c, true) }

// UpdatePost godoc
// @Summary Edit your own post and replace its attachment set
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Accept json
// @Param request body models.HubPostRequest true "Request"
// @Success 200 {object} models.HubResponse{data=models.HubContent}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId} [patch]
func (h *HubHandler) UpdatePost(c *gin.Context) { h.savePost(c, false) }
func (h *HubHandler) savePost(c *gin.Context, create bool) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	name := "postId"
	if create {
		name = "channelId"
	}
	id, ok := parsePathInt64(c, name)
	if !ok {
		return
	}
	var req models.HubPostRequest
	if !hubBind(c, &req) {
		return
	}
	var channel int64
	if create {
		channel, id = id, 0
	}
	v, e := h.service.SavePost(ctx, s, u, channel, id, req)
	hubResult(c, v, e)
}

// CreateComment godoc
// @Summary Reply to a post or top-level comment
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Accept json
// @Param request body models.HubCommentRequest true "Request"
// @Success 200 {object} models.HubResponse{data=models.HubContent}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId}/comments [post]
func (h *HubHandler) CreateComment(c *gin.Context) { h.saveComment(c, true) }

// UpdateComment godoc
// @Summary Edit your own comment
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param commentId path int true "commentId"
// @Accept json
// @Param request body models.HubCommentRequest true "Request"
// @Success 200 {object} models.HubResponse{data=models.HubContent}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/comments/{commentId} [patch]
func (h *HubHandler) UpdateComment(c *gin.Context) { h.saveComment(c, false) }
func (h *HubHandler) saveComment(c *gin.Context, create bool) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	name := "commentId"
	if create {
		name = "postId"
	}
	id, ok := parsePathInt64(c, name)
	if !ok {
		return
	}
	var req models.HubCommentRequest
	if !hubBind(c, &req) {
		return
	}
	var post int64
	if create {
		post, id = id, 0
	}
	v, e := h.service.SaveComment(ctx, s, u, post, id, req)
	hubResult(c, v, e)
}

// DeletePost godoc
// @Summary Soft-delete own post or remove another member post
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Accept json
// @Param request body models.HubReasonRequest false "Request; removal reason is required when moderating another author"
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId} [delete]
func (h *HubHandler) DeletePost(c *gin.Context) { h.delete(c, false) }

// DeleteComment godoc
// @Summary Soft-delete own comment or moderate another member comment
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param commentId path int true "commentId"
// @Accept json
// @Param request body models.HubReasonRequest false "Request; removal reason is required when moderating another author"
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/comments/{commentId} [delete]
func (h *HubHandler) DeleteComment(c *gin.Context) { h.delete(c, true) }
func (h *HubHandler) delete(c *gin.Context, comment bool) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	name := "postId"
	if comment {
		name = "commentId"
	}
	id, ok := parsePathInt64(c, name)
	if !ok {
		return
	}
	var req models.HubReasonRequest
	if c.Request.ContentLength != 0 && !hubBind(c, &req) {
		return
	}
	hubResult(c, nil, h.service.Delete(ctx, s, u, id, comment, req.Reason))
}

// Pin godoc
// @Summary Set announcement pin state (owner/admin)
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Accept json
// @Param request body models.HubControlRequest true "Request"
// @Description Send is_pinned only; announcements only.
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId}/pin [post]
func (h *HubHandler) Pin(c *gin.Context) { h.control(c, true) }

// Lock godoc
// @Summary Set discussion lock state (owner/admin)
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Accept json
// @Param request body models.HubControlRequest true "Request"
// @Description Send locked only; locks block comments and comment edits for all roles.
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId}/lock [post]
func (h *HubHandler) Lock(c *gin.Context) { h.control(c, false) }
func (h *HubHandler) control(c *gin.Context, pin bool) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "postId")
	if !ok {
		return
	}
	var req models.HubControlRequest
	if !hubBind(c, &req) {
		return
	}
	if (pin && (req.IsPinned == nil || req.Locked != nil)) || (!pin && (req.Locked == nil || req.IsPinned != nil)) {
		hubResult(c, nil, hubsvc.ErrInvalid)
		return
	}
	hubResult(c, nil, h.service.Control(ctx, s, u, id, req.IsPinned, req.Locked))
}

// PostReaction godoc
// @Summary Add a post reaction idempotently
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Accept json
// @Param request body models.HubReactionRequest true "Request"
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId}/reactions [post]
func (h *HubHandler) PostReaction(c *gin.Context) { h.reaction(c, false, false) }

// RemovePostReaction godoc
// @Summary Remove a post reaction idempotently
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Param reactionType path string true "Reaction type" Enums(like,love,helpful)
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId}/reactions/{reactionType} [delete]
func (h *HubHandler) RemovePostReaction(c *gin.Context) { h.reaction(c, false, true) }

// CommentReaction godoc
// @Summary Add a comment reaction idempotently
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param commentId path int true "commentId"
// @Accept json
// @Param request body models.HubReactionRequest true "Request"
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/comments/{commentId}/reactions [post]
func (h *HubHandler) CommentReaction(c *gin.Context) { h.reaction(c, true, false) }

// RemoveCommentReaction godoc
// @Summary Remove a comment reaction idempotently
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param commentId path int true "commentId"
// @Param reactionType path string true "Reaction type" Enums(like,love,helpful)
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/comments/{commentId}/reactions/{reactionType} [delete]
func (h *HubHandler) RemoveCommentReaction(c *gin.Context) { h.reaction(c, true, true) }
func (h *HubHandler) reaction(c *gin.Context, comment, remove bool) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	name := "postId"
	if comment {
		name = "commentId"
	}
	id, ok := parsePathInt64(c, name)
	if !ok {
		return
	}
	req := models.HubReactionRequest{ReactionType: c.Param("reactionType")}
	if !remove && !hubBind(c, &req) {
		return
	}
	hubResult(c, nil, h.service.Reaction(ctx, s, u, id, comment, remove, req.ReactionType))
}

// Read godoc
// @Summary Advance the channel read watermark monotonically
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param channelId path int true "channelId"
// @Accept json
// @Param request body models.HubReadRequest true "Request"
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channels/{channelId}/read [post]
func (h *HubHandler) Read(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "channelId")
	if !ok {
		return
	}
	var req models.HubReadRequest
	if !hubBind(c, &req) {
		return
	}
	if req.PostID <= 0 {
		hubResult(c, nil, hubsvc.ErrInvalid)
		return
	}
	hubResult(c, nil, h.service.Read(ctx, s, u, id, req.PostID))
}

// Report godoc
// @Summary Report a post (one report per member and post)
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param postId path int true "postId"
// @Accept json
// @Param request body models.HubReasonRequest true "Request"
// @Success 200 {object} models.HubResponse{data=models.HubReport}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/posts/{postId}/reports [post]
func (h *HubHandler) Report(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "postId")
	if !ok {
		return
	}
	var req models.HubReasonRequest
	if !hubBind(c, &req) {
		return
	}
	v, e := h.service.Report(ctx, s, u, id, req.Reason)
	hubResult(c, v, e)
}

// Reports godoc
// @Summary List moderation reports (owner/admin)
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param cursor query string false "Opaque resource-scoped cursor"
// @Param limit query int false "Page size" default(20) minimum(1) maximum(100)
// @Param status query string false "Report status" Enums(pending,dismissed,actioned)
// @Success 200 {object} models.HubResponse{data=models.HubReportPage}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channel-reports [get]
func (h *HubHandler) Reports(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	limit, ok := hubLimit(c)
	if !ok {
		return
	}
	v, e := h.service.Reports(ctx, s, u, c.Query("status"), c.Query("cursor"), limit)
	hubResult(c, v, e)
}

// Resolve godoc
// @Summary Dismiss a report or remove the reported post (owner/admin)
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param reportId path int true "reportId"
// @Accept json
// @Param request body models.HubResolveRequest true "Request"
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channel-reports/{reportId} [patch]
func (h *HubHandler) Resolve(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "reportId")
	if !ok {
		return
	}
	var req models.HubResolveRequest
	if !hubBind(c, &req) {
		return
	}
	hubResult(c, nil, h.service.Resolve(ctx, s, u, id, req))
}

// Upload godoc
// @Summary Upload a private image or PDF before publishing
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Accept multipart/form-data
// @Param file formData file true "One JPEG/PNG (5 MiB) or PDF (10 MiB); expires unclaimed after 24 hours"
// @Success 200 {object} models.HubResponse{data=models.HubUploadView}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Failure 413 {object} models.ErrorResponseDoc
// @Failure 502 {object} models.ErrorResponseDoc
// @Failure 503 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channel-uploads [post]
func (h *HubHandler) Upload(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	v, e := h.service.Upload(ctx, s, u, c.Writer, c.Request)
	hubResult(c, v, e)
}

// DeleteUpload godoc
// @Summary Queue an unclaimed upload for deletion
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param uploadId path int true "uploadId"
// @Success 200 {object} models.HubResponse{data=object}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channel-uploads/{uploadId} [delete]
func (h *HubHandler) DeleteUpload(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "uploadId")
	if !ok {
		return
	}
	hubResult(c, nil, h.service.DeleteUpload(ctx, s, u, id))
}

// Attachment godoc
// @Summary Get an authorized signed attachment URL
// @Description Requires active owner/admin/resident membership and an operational society. Staff and platform roles without membership have no access.
// @Tags Society Hub
// @Produce json
// @Security AccessToken
// @Param societyId path int true "Society ID"
// @Param uploadId path int true "uploadId"
// @Description Attachment IDs are upload IDs. Pending uploads are visible only to the uploader. Signed URLs expire after 15 minutes.
// @Success 200 {object} models.HubResponse{data=models.HubAttachmentView}
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Failure 409 {object} models.ErrorResponseDoc
// @Failure 429 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Router /v1/societies/{societyId}/channel-attachments/{uploadId} [get]
func (h *HubHandler) Attachment(c *gin.Context) {
	ctx, s, u, ok := hubRequest(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "uploadId")
	if !ok {
		return
	}
	v, e := h.service.Attachment(ctx, s, u, id)
	hubResult(c, v, e)
}
