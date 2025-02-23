package authen_and_post_service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/haiyen11231/social-media-app.git/configs"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AuthenAndPostService struct {
    authen_and_post.UnimplementedAuthenticateAndPostServer
    db          *gorm.DB         // MySQL database connection
    rdb         *redis.Client    // Redis client for caching
    minioClient *minio.Client    // MinIO client for object storage
}

func NewAuthenAndPostService(cfg *configs.AuthenAndPostConfig) (*AuthenAndPostService, error) {
    // Initialize MySQL
    db, err := gorm.Open(mysql.New(mysql.Config{
        DSN:                       cfg.MySQL.DSN,
        DefaultStringSize:         cfg.MySQL.DefaultStringSize,
        DisableDatetimePrecision:  cfg.MySQL.DisableDatetimePrecision,
        DontSupportRenameIndex:    cfg.MySQL.DontSupportRenameIndex,
        SkipInitializeWithVersion: cfg.MySQL.SkipInitializeWithVersion,
    }), &gorm.Config{})
    if err != nil {
        log.Printf("Error connecting to MySQL DB: %v\n", err)
        return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
    }
    log.Println("Connected to MySQL!")

    // Initialize Redis
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    rdb := redis.NewClient(&redis.Options{
        Addr:     cfg.Redis.Addr,
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })

    if _, err := rdb.Ping(ctx).Result(); err != nil {
        log.Printf("Error connecting to Redis: %v\n", err)
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }
    log.Println("Connected to Redis!")

    // Initialize MinIO
    minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
        Secure: cfg.Minio.UseSSL,
    })
    if err != nil {
        log.Printf("Error connecting to MinIO: %v\n", err)
        return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
    }

    // Verify MinIO connection
    _, err = minioClient.ListBuckets(ctx)
    if err != nil {
        log.Printf("Error verifying MinIO connection: %v\n", err)
        return nil, fmt.Errorf("failed to verify MinIO connection: %w", err)
    }
    log.Println("Connected to MinIO!")

    return &AuthenAndPostService{
        db:          db,
        rdb:         rdb,
        minioClient: minioClient,
    }, nil
}