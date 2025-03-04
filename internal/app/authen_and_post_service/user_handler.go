package authen_and_post_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (s *AuthenAndPostService) SignUp(ctx context.Context, request *authen_and_post.SignUpRequest) (*authen_and_post.SignUpResponse, error) {
	log.Println("Received SignUp request")
	salt := generateAlphabetSalt(16)
	hashedPassword, err := hashPassword(request.GetPassword(), salt)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		FirstName:      request.GetFirstName(),
		LastName:       request.GetLastName(),
		DateOfBirth:    request.Dob.AsTime(),
		Email:          request.GetEmail(),
		Username:       request.GetUsername(),
		HashedPassword: hashedPassword,
		Salt:           string(salt),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &authen_and_post.SignUpResponse{Message: "User created successfully"}, nil
}

func (s *AuthenAndPostService) LogIn(ctx context.Context, request *authen_and_post.LogInRequest) (*authen_and_post.LogInResponse, error) {
	// Check Redis cache first
    cacheKey := fmt.Sprintf("user:%s", request.GetUsername())
	cachedUser, err := s.rdb.Get(ctx, cacheKey).Result()
    if err == nil {
        // Cache hit: unmarshal cached user data
        var user models.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			log.Println("User data retrieved from cache")
			return s.generateLoginResponse(ctx, &user)
		}
    }

	// Cache miss: fetch from database
    var user models.User
	result := s.db.Where(&models.User{Username: request.GetUsername()}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &authen_and_post.LogInResponse{Message: "User not found"}, nil
	} else if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", result.Error)
	}

	passwordWithSalt := []byte(request.Password + user.Salt)
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), passwordWithSalt); err != nil {
		return &authen_and_post.LogInResponse{Message: "Incorrect password"}, nil
	}

	// Cache the user data in Redis
    userJSON, err := json.Marshal(user)
	if err == nil {
		if err := s.rdb.Set(ctx, cacheKey, userJSON, 5*time.Minute).Err(); err != nil {
			log.Println("Failed to cache user data:", err)
		}
	}

	return s.generateLoginResponse(ctx, &user)
}

func (s *AuthenAndPostService) generateLoginResponse(ctx context.Context, user *models.User) (*authen_and_post.LogInResponse, error) {
    accessToken, err := GenerateToken(uint64(user.ID), 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateToken(uint64(user.ID), 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

    // Store refresh token in Redis
    refreshTokenKey := fmt.Sprintf("refresh_token:%d", user.ID)
	if err := s.rdb.Set(ctx, refreshTokenKey, refreshToken, 24*time.Hour).Err(); err != nil {
		log.Println("Failed to store refresh token in Redis:", err)
	}

	return &authen_and_post.LogInResponse{
		UserId:       uint64(user.ID),
		Message:      "Log in successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthenAndPostService) EditUser(ctx context.Context, request *authen_and_post.EditUserRequest) (*authen_and_post.EditUserResponse, error) {
	var user models.User
	s.db.Where(&models.User{ID: uint(request.UserId)}).First(&user)
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
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
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		user.HashedPassword = hashPassword
		user.Salt = string(salt)
	}

	s.db.Save(&user)

	// Invalidate cached user data
    cacheKey := fmt.Sprintf("user:%s", user.Username)
    if err := s.rdb.Del(ctx, cacheKey).Err(); err != nil {
        log.Println("Failed to invalidate cached user data:", err)
    }

	return &authen_and_post.EditUserResponse{
		Message: "User updated successfully",
	}, nil
}

func (s *AuthenAndPostService) AuthenticateUser (ctx context.Context, request *authen_and_post.AuthenticateUserRequest) (*authen_and_post.AuthenticateUserResponse, error) {
	if request.Token == "" {
		return &authen_and_post.AuthenticateUserResponse{IsValid: false, Message: "Token is required"}, errors.New("token is required")
	}

	parsedId, err := ParseToken(request.Token, os.Getenv("JWT_SECRET"))
	if err != nil {
		return &authen_and_post.AuthenticateUserResponse{IsValid: false, Message: err.Error()}, fmt.Errorf("failed to parse token: %w", err)
	}
	log.Printf("Extracted claims ID: %v", parsedId)

	// Check Redis cache first
    cacheKey := fmt.Sprintf("user:%d", parsedId)
	cachedUser, err := s.rdb.Get(ctx, cacheKey).Result()
    if err == nil {
        // Cache hit: unmarshal cached user data
        var user models.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			log.Println("User data retrieved from cache")
			return &authen_and_post.AuthenticateUserResponse{IsValid: true, Message: "Authenticated!", UserId: uint64(user.ID)}, nil
		}
    }

    // Cache miss: fetch from database
	user := &models.User{}
	if err = s.db.Model(&models.User{}).Where("id = ?", parsedId).First(user).Error; err != nil {
		return &authen_and_post.AuthenticateUserResponse{IsValid: false, Message: "Invalid Credentials!"}, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Cache the user data in Redis
    userJSON, err := json.Marshal(user)
	if err == nil {
		if err := s.rdb.Set(ctx, cacheKey, userJSON, 5*time.Minute).Err(); err != nil {
			log.Println("Failed to cache user data:", err)
		}
	}

	log.Printf("User found with ID: %v", user.ID)
	return &authen_and_post.AuthenticateUserResponse{IsValid: true, Message: "Authenticated!", UserId: uint64(user.ID)}, nil

}

func (s *AuthenAndPostService) RefreshToken (ctx context.Context, request *authen_and_post.RefreshTokenRequest) (*authen_and_post.RefreshTokenResponse, error) {
	// Validate the refresh token
    userID, err := ParseToken(request.RefreshToken, os.Getenv("JWT_SECRET"))
    if err != nil {
        log.Println("Failed to parse refresh token:", err)
        return nil, fmt.Errorf("invalid or expired refresh token")
    }

    // Get the refresh token key for the user
    refreshTokenKey := fmt.Sprintf("refresh_token:%d", userID)

    // Retrieve the stored refresh token from Redis
    storedRefreshToken, err := s.rdb.Get(ctx, refreshTokenKey).Result()
    if err == redis.Nil {
        // Refresh token not found in Redis
        return nil, fmt.Errorf("invalid or expired refresh token")
    } else if err != nil {
        // Redis error
        log.Println("Failed to retrieve refresh token from Redis:", err)
        return nil, fmt.Errorf("internal server error")
    }

    // Validate the provided refresh token against the stored token
    if storedRefreshToken != request.RefreshToken {
        return nil, fmt.Errorf("invalid refresh token")
    }

    // Generate a new access token
    newAccessToken, err := GenerateToken(uint64(userID), 15*time.Minute)
    if err != nil {
        log.Println("Failed to generate new access token:", err)
        return nil, fmt.Errorf("internal server error")
    }

    // Extend the refresh token's expiration in Redis
    if err := s.rdb.Expire(ctx, refreshTokenKey, 24*time.Hour).Err(); err != nil {
        log.Println("Failed to extend refresh token expiration in Redis:", err)
    }

    return &authen_and_post.RefreshTokenResponse{
        AccessToken: newAccessToken,
    }, nil
}

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

// func getKey(userID uint64) string {
// 	return "user_refresh_token:" + fmt.Sprint(userID)
// }

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