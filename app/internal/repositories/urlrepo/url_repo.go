package urlRepo

import (
	"sync"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"gorm.io/gorm"
)

type UrlRepo struct {
	db       *gorm.DB
	mu       sync.RWMutex
	saveChan chan entities.URL
}

func NewUrlRepo(db *gorm.DB) *UrlRepo {
	saveChan := make(chan entities.URL, 1000)

	repo := &UrlRepo{
		db:       db,
		saveChan: saveChan,
	}

	go repo.saveLoop()

	return repo
}

func (r *UrlRepo) Save(url entities.URL) (entities.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	panic("implement me")
}

func (r *UrlRepo) saveLoop() {
	for {
		urlentities := <-r.saveChan
		urlModel := models.Url{
			BaseModel: models.BaseModel{},
			//Owner:             urlentities.UserId.String(),
			OriginalURL:       urlentities.OriginalUrl,
			ShortenedURLParam: urlentities.ShortCode,
			VisitCount:        nil,
		}
		r.db.Create(&urlModel)
	}
}

func (r *UrlRepo) GetByShortUrl(shortenedUrl string) (entities.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	//longUrl := r.db.Get(shortenedUrl)
	panic("implement me")
}

func (r *UrlRepo) GetByOwner(owner string) ([]entities.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) GetByKeyword(keyword string) ([]entities.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) GetByKeywords(keywords []string) ([]entities.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) GetByOriginalUrl(originalUrl string) ([]entities.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) GetById(id string) (entities.URL, error) {

	panic("implement me")
}

func (r *UrlRepo) Delete(id string) error {

	panic("implement me")
}
