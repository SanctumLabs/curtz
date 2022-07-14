package repositories

import (
	"context"
	"fmt"

	"github.com/sanctumlabs/curtz/app/config"
	urlRepo "github.com/sanctumlabs/curtz/app/internal/repositories/urlrepo"
	"github.com/sanctumlabs/curtz/app/internal/repositories/userepo"
	"github.com/sanctumlabs/curtz/app/tools/logger"
	"github.com/sanctumlabs/curtz/app/tools/monitoring"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var log = logger.NewLogger("repository")

// Repository represents a database repository
type Repository struct {
	dbClient *mongo.Client
	userRepo *userepo.UserRepo
	urlRepo  *urlRepo.UrlRepo
}

// NewRepository creates a new repository with provided config
func NewRepository(config config.DatabaseConfig) *Repository {
	defer monitoring.ErrorHandler()
	var uri string

	if config.IsSRV {
		uri = fmt.Sprintf("mongodb+srv://%s:%s@%s", config.User, config.Password, config.Host)	
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", config.User, config.Password, config.Host, config.Port)
	}

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetRetryWrites(true)

	dbClient, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatalf("DB Client configuration failed with err: %v", err)
	}

	ctx := context.Background()

	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database %s", err)
	}

	db := dbClient.Database(config.Database)

	if err := dbClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("DB Connection failed with err: %v", err)
	}

	log.Info("DB Connection successful")

	return &Repository{
		dbClient: dbClient,
		userRepo: userepo.NewUserRepo(db.Collection("users"), ctx),
		urlRepo:  urlRepo.NewUrlRepo(db.Collection("urls"), ctx),
	}
}

// Disconnect disconnects from the db client database connection
func (r *Repository) Disconnect(ctx context.Context) error {
	defer monitoring.ErrorHandler()
	return r.dbClient.Disconnect(ctx)
}

// GetUrlRepo returns the Url repository
func (r *Repository) GetUrlRepo() *urlRepo.UrlRepo {
	return r.urlRepo
}

// GetUserRepo returns configured user repository
func (r *Repository) GetUserRepo() *userepo.UserRepo {
	return r.userRepo
}
