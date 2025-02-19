package newsfeed_service

import (
	"fmt"

	"github.com/haiyen11231/social-media-app.git/configs"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/newsfeed"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type NewsfeedService struct {
	newsfeed.UnimplementedNewsfeedServer
	db *gorm.DB
	rdb *redis.Client
}

func NewNewsfeedService(cfg *configs.NewsfeedConfig) (*NewsfeedService, error) {
	db, err := gorm.Open(mysql.New(cfg.MySQL), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println("Error connecting to MySQL DB: ", err)
	}

	rdb := redis.NewClient(&cfg.Redis)
	if rdb == nil {
		return nil, fmt.Errorf("Error connecting to Redis: %w", err)
	}

	return &NewsfeedService{
		db: db,
		rdb: rdb,
	}, nil
}