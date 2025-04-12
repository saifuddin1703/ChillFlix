package video

import (
	"chillfix/models"
	"context"
)

type VideoRepository interface {
	Create(ctx context.Context, video *models.Video) error
	FindByID(ctx context.Context, id string) (*models.Video, error)
	Update(ctx context.Context, video *models.Video) error
	Delete(ctx context.Context, id string) error
}
