package urlRepo

import (
	"github.com/google/uuid"
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"gorm.io/gorm"
	"sync"
)

type UrlRepo struct {
	db       *gorm.DB
	mu       sync.RWMutex
	saveChan chan entity.URL
}

func NewUrlRepo(db *gorm.DB) *UrlRepo {
	saveChan := make(chan entity.URL, 1000)

	repo := &UrlRepo{
		db:       db,
		saveChan: saveChan,
	}

	go repo.saveLoop()

	return repo
}

func (r *UrlRepo) Save(owner uuid.UUID, originalUrl, shortenedUrl string) (entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	panic("implement me")
}

func (r *UrlRepo) saveLoop() {
	for {
		urlEntity := <-r.saveChan
		urlModel := models.Url{
			BaseModel: models.BaseModel{},
			//Owner:             urlEntity.UserId.String(),
			OriginalURL:       urlEntity.OriginalUrl,
			ShortenedURLParam: urlEntity.ShortenedUrl,
			VisitCount:        nil,
		}
		r.db.Create(&urlModel)
	}
}

func (r *UrlRepo) GetByShortUrl(shortenedUrl string) (entity.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	//longUrl := r.db.Get(shortenedUrl)
	panic("implement me")
}

func (r *UrlRepo) GetByOwner(owner uuid.UUID) ([]entity.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) GetByKeyword(keyword string) ([]entity.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) GetByKeywords(keywords []string) ([]entity.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) GetByOriginalUrl(originalUrl string) ([]entity.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) GetById(id uuid.UUID) (entity.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) Delete(id uuid.UUID) error {

	panic("implement me")
}
