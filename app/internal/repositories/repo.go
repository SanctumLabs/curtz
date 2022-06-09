package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sanctumlabs/curtz/app/config"
	urlRepo "github.com/sanctumlabs/curtz/app/internal/repositories/urlrepo"
	"github.com/sanctumlabs/curtz/app/internal/repositories/userepo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	dbClient *mongo.Client
	userRepo *userepo.UserRepo
	urlRepo  *urlRepo.UrlRepo
}

func NewRepository(config config.DatabaseConfig) *Repository {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", config.User, config.Password, config.Host, config.Port, config.Database)

	dbClient, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalf("DB Connection failed with err: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer dbClient.Disconnect(ctx)

	db := dbClient.Database(config.Database)

	return &Repository{
		dbClient: dbClient,
		userRepo: userepo.NewUserRepo(db.Collection("users"), ctx),
		urlRepo:  urlRepo.NewUrlRepo(db.Collection("urls"), ctx),
	}
}

func (r *Repository) GetUrlRepo() *urlRepo.UrlRepo {
	return r.urlRepo
}

func (r *Repository) GetUserRepo() *userepo.UserRepo {
	return r.userRepo
}
