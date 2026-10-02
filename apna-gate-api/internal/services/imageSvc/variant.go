package imagesvc

import "go-server/internal/models"

func ParseImageVariant(value string) (models.ImageVariant, error) {
	switch models.ImageVariant(value) {
	case "", models.ImageOriginal:
		return models.ImageOriginal, nil
	case models.ImageList, models.ImageAvatar, models.ImageDetail:
		return models.ImageVariant(value), nil
	default:
		return "", ErrImageVariant
	}
}
