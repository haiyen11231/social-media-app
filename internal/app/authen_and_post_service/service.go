package authen_and_post_service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/haiyen11231/social-media-app.git/configs"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type AuthenAndPostService struct {
	authen_and_post.UnimplementedAuthenticateAndPostServer
	db *gorm.DB
	rdb *redis.Client
	minioClient *minio.Client
}

func NewAuthenAndPostService(cfg *configs.AuthenAndPostConfig) (*AuthenAndPostService, error) {
	db, err := gorm.Open(mysql.New(cfg.MySQL), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Println("Error connecting to MySQL DB: ", err)
		return nil, err
	}

	rdb := redis.NewClient(&cfg.Redis)
	if rdb != nil {
		log.Println("Error connecting to Redis: ", err)
		return nil, err
	}

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Println("Error connecting to Minio: ", err)
		return nil, err
	}

	return &AuthenAndPostService{
		db: db,
		rdb: rdb,
		minioClient: minioClient,
	}, nil
}

// User
func (s *AuthenAndPostService) SignUp (ctx context.Context, request *authen_and_post.SignUpRequest) (*authen_and_post.SignUpResponse, error) {
	salt := generateAlphabetSalt(16)
	hashedPassword, err := hashPassword(request.GetPassword(), salt)
	if err != nil {
		return nil, err
	}

	user := models.User{
		FirstName: request.GetFirstName(),
		LastName: request.GetLastName(),
		DateOfBirth: request.Dob.AsTime(),
		Email: request.GetEmail(),
		Username: request.GetUsername(),
		HashedPassword: hashedPassword,
		Salt: string(salt),
	}

	s.db.Create(&user)
	return &authen_and_post.SignUpResponse{Message: "User created successfully"}, nil
}

func (s *AuthenAndPostService) LogIn (ctx context.Context, request *authen_and_post.LogInRequest) (*authen_and_post.LogInResponse, error) {
	var user models.User
	result := s.db.Where(&models.User{Username: request.GetUsername()}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &authen_and_post.LogInResponse{Message: "User not found"}, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	passwordWithSalt := []byte(request.Password + user.Salt)
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), passwordWithSalt)
	if err != nil {
		return &authen_and_post.LogInResponse{Message: "Incorrect password"}, nil
	}

	accessToken, err := GenerateToken(uint64(user.ID), 15*time.Minute)
	if err != nil {
		log.Println("Failed to generate access token:", err.Error())
		return nil, err
	}

	refreshToken, err := GenerateToken(uint64(user.ID), 24*time.Hour)
	if err != nil {
		log.Println("Failed to generate refresh token:", err.Error())
		return nil, err
	}

	// Store refresh token in Redis....

	return &authen_and_post.LogInResponse{
		UserId: uint64(user.ID),
		Message: "Log in successful",
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthenAndPostService) EditUser (ctx context.Context, request *authen_and_post.EditUserRequest) (*authen_and_post.EditUserResponse, error) {
	var user models.User
	s.db.Where(&models.User{ID: uint(request.UserId)}).First(&user)
	if user.ID == 0 {
		return nil, fmt.Errorf("User not found")
	}

	if request.FirstName != nil {
		user.FirstName = request.GetFirstName()
	}

	if request.LastName != nil {
		user.LastName = request.GetLastName()
	}

	if request.Dob != nil {
		user.DateOfBirth = request.Dob.AsTime()
	}

	if request.Password != nil {
		salt := generateAlphabetSalt(16)
		hashPassword, err := hashPassword(request.GetPassword(), salt)
		if err != nil {
			return nil, err
		}

		user.HashedPassword = hashPassword
		user.Salt = string(salt)
	}

	s.db.Save(&user)

	return &authen_and_post.EditUserResponse{
		Message: "User updated successfully",
	}, nil
}

// func (s *AuthenAndPostService) AuthenticateUser (ctx context.Context, request *authen_and_post.AuthenticateUserRequest) (*authen_and_post.AuthenticateUserResponse, error) {

// }

// func (s *AuthenAndPostService) RefreshToken (ctx context.Context, request *authen_and_post.RefreshTokenRequest) (*authen_and_post.RefreshTokenResponse, error) {

// }

// Following
func (s *AuthenAndPostService) checkUserExisting (ctx context.Context, userId uint64) error {
	var user models.User
	err := s.db.Table("user").Where("id = ?", userId).First(&user).Error
	if err != nil {
		return errors.New("User not found")
	}

	return nil
}

// func (s *AuthenAndPostService) FollowUser (ctx context.Context, request *authen_and_post.FollowUserRequest) (*authen_and_post.FollowUserResponse, error) {

// }

// func (s *AuthenAndPostService) UnfollowUser (ctx context.Context, request *authen_and_post.UnfollowUserRequest) (*authen_and_post.UnfollowUserResponse, error) {

// }

// func (s *AuthenAndPostService) GetFollowerList (ctx context.Context, request *authen_and_post.GetFollowerListRequest) (*authen_and_post.GetFollowerListResponse, error) {

// }

// // Post
// func (s *AuthenAndPostService) CreatePost (ctx context.Context, request *authen_and_post.CreatePostRequest) (*authen_and_post.CreatePostResponse, error) {

// }

// func (s *AuthenAndPostService) GetPost (ctx context.Context, request *authen_and_post.GetPostRequest) (*authen_and_post.GetPostResponse, error) {

// }

// func (s *AuthenAndPostService) EditPost (ctx context.Context, request *authen_and_post.EditPostRequest) (*authen_and_post.EditPostResponse, error) {

// }

// func (s *AuthenAndPostService) DeletePost (ctx context.Context, request *authen_and_post.DeletePostRequest) (*authen_and_post.DeletePostResponse, error) {

// }

// func (s *AuthenAndPostService) CreateComment (ctx context.Context, request *authen_and_post.CreateCommentRequest) (*authen_and_post.CreateCommentResponse, error) {

// }

// func (s *AuthenAndPostService) LikePost (ctx context.Context, request *authen_and_post.LikePostRequest) (*authen_and_post.LikePostResponse, error) {

// }

func generateAlphabetSalt(length int) []byte {
	rand.Seed(time.Now().UnixNano())

	salt := make([]byte, length)
	for i := 0; i < length; i++ {
		salt[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return salt
}

func hashPassword(password string, salt []byte) (string, error) {
	passwordWithSalt := []byte(password + string(salt))
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordWithSalt, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func GenerateToken(userId uint64, expiry time.Duration) (string, error) {
    expirationTime := time.Now().Add(expiry)
	claims := jwt.MapClaims{
        "user_id": userId,
        "exp":     expirationTime.Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func getKey(userID uint64) string {
    return "user_refresh_token:" + fmt.Sprint(userID)
}

func ParseToken(tokenString, secret string) (int64, error) {
	log.Println("Token:", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		log.Println("Error parsing token:", err)
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return 0, errors.New("token expired")
		}

		// Debugging: log the type of claims["id"]
		log.Printf("Type of 'user_id' in claims: %T\n", claims["user_id"])

		switch id := claims["user_id"].(type) {
		case float64:
			return int64(id), nil
		case string:
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return 0, errors.New("invalid ID format in token")
			}
			return idInt, nil
		case int64:
			return id, nil
		default:
			return 0, errors.New("invalid ID format in token")
		}
	}
	return 0, errors.New("invalid token")
}