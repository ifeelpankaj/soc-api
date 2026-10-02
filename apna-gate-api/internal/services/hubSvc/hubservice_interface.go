package hubsvc

import (
	"context"
	"go-server/internal/models"
	"net/http"
)

type HubService interface {
	HubServiceCommands
	HubServiceQueries
	HubServiceJobs
}

type HubServiceCommands interface {
	Reaction(ctx context.Context, society, user, id int64, comment, remove bool, kind string) error
	SavePost(ctx context.Context, society, user, channel, id int64, req models.HubPostRequest) (models.HubContent, error)
	SaveComment(ctx context.Context, society, user, post, id int64, req models.HubCommentRequest) (models.HubContent, error)
	Delete(ctx context.Context, society, user, id int64, comment bool, reason string) error
	Control(ctx context.Context, society, user, id int64, pin, locked *bool) error
	Read(ctx context.Context, society, user, channel, post int64) error
	Report(ctx context.Context, society, user, post int64, reason string) (models.HubReport, error)
	Resolve(ctx context.Context, society, user, id int64, req models.HubResolveRequest) error
	Upload(ctx context.Context, society, user int64, w http.ResponseWriter, r *http.Request) (models.HubUploadView, error)
	DeleteUpload(ctx context.Context, society, user, id int64) error
}

type HubServiceQueries interface {
	Reports(ctx context.Context, society, user int64, status, raw string, limit int) (models.HubReportPage, error)
	Channels(ctx context.Context, society, user int64) ([]models.HubChannelSummary, error)
	Categories(ctx context.Context, society, user int64) ([]models.HubCategory, error)
	Post(ctx context.Context, society, user, id int64) (models.HubContent, error)
	Feed(ctx context.Context, society, user, target int64, raw string, limit int, category *int16, comments bool) (models.HubPage, error)
	Attachment(ctx context.Context, society, user, id int64) (models.HubAttachmentView, error)
}

type HubServiceJobs interface {
	CleanupUploads(ctx context.Context) error
}

var _ HubService = (*Service)(nil)
