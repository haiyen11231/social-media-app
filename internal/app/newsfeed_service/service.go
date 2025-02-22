package newsfeed_service

import (
	"context"
	"log"

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
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: cfg.MySQL.DSN,
		DefaultStringSize: cfg.MySQL.DefaultStringSize,
		DisableDatetimePrecision: cfg.MySQL.DisableDatetimePrecision,
		DontSupportRenameIndex: cfg.MySQL.DontSupportRenameIndex,
		SkipInitializeWithVersion: cfg.MySQL.SkipInitializeWithVersion,
	}), &gorm.Config{})
	if err != nil {
		log.Println("Error connecting to MySQL DB: ", err)
		return nil, err
	}
	log.Println("Connected to MySQL!")

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB: cfg.Redis.DB,
	})
	if _, err = rdb.Ping(context.Background()).Result(); err != nil {
		log.Println("Error connecting to Redis:", err)
		return nil, err
	}
	log.Println("Connected to Redis!")

	return &NewsfeedService{
		db: db,
		rdb: rdb,
	}, nil
}