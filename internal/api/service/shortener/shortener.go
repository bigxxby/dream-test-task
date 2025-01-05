package shortener

import (
	"time"

	"github.com/bigxxby/dream-test-task/internal/api/repo/shortener"
	"github.com/bigxxby/dream-test-task/internal/models"
	"github.com/bigxxby/dream-test-task/internal/utils"
	"github.com/google/uuid"
)

type IShortenerService interface {
	CreateShortLink(link string, userId *uuid.UUID) (*models.ShortLink, int, error)
	Redirect(shortID string) (string, int, error)
	GetLinks(userId *uuid.UUID) ([]models.ShortLink, int, error)
	GetLink(shortID string) (*models.ShortLink, int, error)
	DeleteLink(shortID string) (int, error)
}

type ShortenerService struct {
	ShortenerRepo shortener.IShortenerRepo
}

func (s *ShortenerService) GetLinks(userId *uuid.UUID) ([]models.ShortLink, int, error) {
	links, err := s.ShortenerRepo.GetLinks(userId)
	if err != nil {
		return nil, 500, err
	}
	return links, 200, nil
}
func (s *ShortenerService) GetLink(shortID string) (*models.ShortLink, int, error) {
	link, err := s.ShortenerRepo.GetShortLinkByShortID(shortID)
	if err != nil {
		return nil, 404, err
	}
	return link, 200, nil
}

func NewShortenerService(shortenerRepo shortener.IShortenerRepo) IShortenerService {
	return &ShortenerService{ShortenerRepo: shortenerRepo}
}
func (s *ShortenerService) DeleteLink(shortID string) (int, error) {
	err := s.ShortenerRepo.DeleteLink(shortID)
	if err != nil {
		return 500, err
	}
	return 200, nil
}

// CreateShortLink implements IShortenerService.
func (s *ShortenerService) CreateShortLink(link string, userId *uuid.UUID) (*models.ShortLink, int, error) {
	// Генерация уникального короткого идентификатора
	shortLink := utils.GenerateShortLink()

	// Проверка, существует ли уже короткая ссылка с таким идентификатором
	existingLink, err := s.ShortenerRepo.GetShortLinkByShortID(shortLink)
	if err != nil && err.Error() != "record not found" {
		return nil, 500, err
	}
	if existingLink != nil {
		// Если ссылка уже существует, генерируем новую
		shortLink = utils.GenerateShortLink()
	}

	expiration := time.Now().Add(30 * 24 * time.Hour)

	shortLinkModel := &models.ShortLink{
		LongLink:  link,
		ShortId:   shortLink,
		UserID:    userId, // Привязываем userId
		ExpiresAt: &expiration,
	}

	err = s.ShortenerRepo.CreateShortLink(shortLinkModel)
	if err != nil {
		return nil, 500, err
	}

	err = shortLinkModel.ValidateLongLink()
	if err != nil {
		return nil, 400, err
	}

	// shortLink = "http://localhost:" + config.AppPort + "/" + "shortener/" + shortLink
	shortLinkModel.ParseShortId()
	return shortLinkModel, 200, nil
}

func (s *ShortenerService) Redirect(shortID string) (string, int, error) {
	shortLink, err := s.ShortenerRepo.GetShortLinkByShortID(shortID)
	if err != nil {
		return "", 500, err
	}
	if shortLink == nil {
		return "", 404, nil
	}

	if shortLink.ExpiresAt != nil && shortLink.ExpiresAt.Before(time.Now()) {
		return "", 404, nil
	}

	shortLink.UpdateClicks()
	err = s.ShortenerRepo.UpdateShortLink(shortLink)
	if err != nil {
		return "", 500, err
	}

	return shortLink.LongLink, 200, nil
}
