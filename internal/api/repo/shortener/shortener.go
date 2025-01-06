package shortener

import (
	"github.com/bigxxby/dream-test-task/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IShortenerRepo interface {
	CreateShortLink(link *models.ShortLink) error
	UpdateShortLink(link *models.ShortLink) error
	GetShortLinkByShortID(shortID string) (*models.ShortLink, error)
	GetLinkStat(shortID string) (int, error) // Возвращает количество кликов для короткой ссылки
	DeleteLink(shortID string) error         // Удаляет короткую ссылку
	GetLinks(userId *uuid.UUID) ([]models.ShortLink, error)
}

type ShortenerRepo struct {
	Db *gorm.DB
}

// NewShortenerRepo создаёт новый экземпляр репозитория с подключением к базе данных.
func NewShortenerRepo(db *gorm.DB) IShortenerRepo {
	return &ShortenerRepo{Db: db}
}

func (sr *ShortenerRepo) GetLinks(userId *uuid.UUID) ([]models.ShortLink, error) {
	var links []models.ShortLink
	err := sr.Db.Where("user_id = ?", userId).Find(&links).Error
	if err != nil {
		return nil, err
	}
	return links, nil
}
func (sr *ShortenerRepo) UpdateShortLink(link *models.ShortLink) error {
	result := sr.Db.Save(link)
	if result.Error != nil {
		return result.Error
	}
	return result.Error
}

func (sr *ShortenerRepo) CreateShortLink(link *models.ShortLink) error {
	result := sr.Db.Create(link)
	if result.Error != nil {
		return result.Error
	}
	return result.Error
}

// GetShortLinkByShortID находит короткую ссылку по короткому идентификатору.
func (sr *ShortenerRepo) GetShortLinkByShortID(shortID string) (*models.ShortLink, error) {
	var link models.ShortLink
	err := sr.Db.Where("short_id = ?", shortID).First(&link).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Если запись не найдена, возвращаем nil
		}
		return nil, err // Возвращаем ошибку для других случаев
	}
	return &link, nil
}

// GetLinkStat возвращает количество кликов по короткой ссылке.
func (sr *ShortenerRepo) GetLinkStat(shortID string) (int, error) {
	var link models.ShortLink
	err := sr.Db.Where("short_id = ?", shortID).First(&link).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // Если запись не найдена, возвращаем 0 кликов
		}
		return 0, err // Возвращаем ошибку для других случаев
	}
	return link.Clicks, nil
}

// DeleteLink удаляет короткую ссылку по её короткому идентификатору.
func (sr *ShortenerRepo) DeleteLink(shortID string) error {
	var link models.ShortLink
	// Проверяем, существует ли такая ссылка
	err := sr.Db.Where("short_id = ?", shortID).First(&link).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // Если запись не найдена, ничего не делаем
		}
		return err // Возвращаем ошибку для других случаев
	}
	// Удаляем ссылку
	return sr.Db.Delete(&link).Error
}
