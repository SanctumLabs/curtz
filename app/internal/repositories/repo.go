package repositories

import (
	"fmt"
	"github.com/sanctumlabs/curtz/app/config"
	models2 "github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"github.com/sanctumlabs/curtz/app/internal/repositories/urlrepo"
	"github.com/sanctumlabs/curtz/app/internal/repositories/userepo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type repository struct {
	db       *gorm.DB
	userRepo *userepo.UserRepo
	urlRepo  *urlRepo.UrlRepo
}

func NewRepository(config config.DatabaseConfig) *repository {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host, config.User, config.Password, config.Database, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("DB Connection failed with err: %v", err)
	}

	if err = db.AutoMigrate(&models2.User{}, &models2.Url{}); err != nil {
		log.Fatalf("AutoMigration failed with err: %v", err)
	}

	return &repository{
		db:       db,
		userRepo: userepo.NewUserRepo(db),
		urlRepo:  urlRepo.NewUrlRepo(db),
	}
}

func (r repository) GetUrlRepo() *urlRepo.UrlRepo {
	return r.urlRepo
}

func (r repository) GetUserRepo() *userepo.UserRepo {
	return r.userRepo
}
