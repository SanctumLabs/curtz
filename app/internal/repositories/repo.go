package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/sanctumlabs/curtz/app/config"
	urlRepo "github.com/sanctumlabs/curtz/app/internal/repositories/urlrepo"
	"github.com/sanctumlabs/curtz/app/internal/repositories/userepo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repository struct {
	dbClient *mongo.Client
	userRepo *userepo.UserRepo
	urlRepo  *urlRepo.UrlRepo
}

func NewRepository(config config.DatabaseConfig) *Repository {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.User, config.Password, config.Host, config.Port)

	dbClient, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalf("DB Connection failed with err: %v", err)
	}

	ctx := context.TODO()

	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := dbClient.Database(config.Database)

	if err := dbClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("DB Connection failed with err: %v", err)
	}

	log.Println("DB Connection successful")

	return &Repository{
		dbClient: dbClient,
		userRepo: userepo.NewUserRepo(db.Collection("users"), ctx),
		urlRepo:  urlRepo.NewUrlRepo(db.Collection("urls"), ctx),
	}
}

// Disconnect disconnects from the db client database connection
func (r *Repository) Disconnect(ctx context.Context) error {
	return r.dbClient.Disconnect(ctx)
}

func (r *Repository) GetUrlRepo() *urlRepo.UrlRepo {
	return r.urlRepo
}

func (r *Repository) GetUserRepo() *userepo.UserRepo {
	return r.userRepo
}
