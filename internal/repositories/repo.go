package repositories

import (
	"fmt"
	"github.com/sanctumlabs/curtz/config"
	"github.com/sanctumlabs/curtz/internal/repositories/entities"
	urlRepo "github.com/sanctumlabs/curtz/internal/repositories/urlrepo"
	"github.com/sanctumlabs/curtz/internal/repositories/userepo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type repo struct {
	db       *gorm.DB
	userRepo *userepo.UserRepo
	urlRepo  *urlRepo.UrlRepo
}

func NewRepo(config config.DatabaseConfig) *repo {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host, config.User, config.Password, config.Database, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("DB Connection failed with err: %v", err)
	}

	if err = db.AutoMigrate(&entities.UserModel{}, &entities.UrlModel{}); err != nil {
		log.Fatalf("AutoMigration failed with err: %v", err)
	}

	return &repo{
		db:       db,
		userRepo: userepo.NewUserRepo(db),
		urlRepo:  urlRepo.NewUrlRepo(db),
	}
}
