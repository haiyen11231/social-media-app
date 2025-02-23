package authen_and_post_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/models"
	"github.com/minio/minio-go/v7"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *AuthenAndPostService) CreatePost(ctx context.Context, request *authen_and_post.CreatePostRequest) (*authen_and_post.CreatePostResponse, error) {
	if err := s.checkUserExisting(ctx, request.UserId); err != nil {
		return nil, err
	}

	var user models.User 
	if err := s.db.WithContext(ctx).Preload("Posts").Find(&user, request.UserId).Error; err != nil {
        return nil, fmt.Errorf("failed to fetch user: %w", err)
    }

	// Upload image to MinIO if provided
	var imagePath string
	if request.ContentImagePath != "" {
		file, err := os.Open(request.ContentImagePath)
		if err != nil {
            return nil, fmt.Errorf("failed to open image file: %w", err)
        }
        defer file.Close()

		objectName := fmt.Sprintf("post/%d/%s", request.UserId, filepath.Base(request.ContentImagePath))
		_, err = s.minioClient.PutObject(ctx, "social-media-bucket", objectName, file, -1, minio.PutObjectOptions{})
        if err != nil {
            return nil, fmt.Errorf("failed to upload image to MinIO: %w", err)
        }

        imagePath = objectName
	}

	post := models.Post{
		ContentText: request.ContentText,
		ContentImagePath: imagePath,
		UserID: uint(request.UserId),
		Visible: request.Visible,
	}

	if err := s.db.WithContext(ctx).Model(&user).Association("Posts").Append(&post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return &authen_and_post.CreatePostResponse{
		Message: "Post created successfully",
		PostId: uint64(post.ID),
	}, nil
}

func (s *AuthenAndPostService) GetPost (ctx context.Context, request *authen_and_post.GetPostRequest) (*authen_and_post.GetPostResponse, error) {
	// Check Redis cache first
    cacheKey := fmt.Sprintf("post:%d", request.PostId)
    cachedPost, err := s.rdb.Get(ctx, cacheKey).Result()
    if err == nil {
        // If found in cache, unmarshal and return
        var post models.Post
        if err := json.Unmarshal([]byte(cachedPost), &post); err == nil {
            return &authen_and_post.GetPostResponse{
                Message: "Post found in cache",
                Post: &authen_and_post.Post{
                    PostId:           uint64(post.ID),
                    UserId:           uint64(post.UserID),
                    ContentText:      post.ContentText,
                    ContentImagePath: post.ContentImagePath,
                    Visible:          post.Visible,
                    CreatedAt:        timestamppb.New(post.CreatedAt),
                },
            }, nil
        }
    }

    // If not found in cache, query the database
	var post models.Post
	err = s.db.WithContext(ctx).First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &authen_and_post.GetPostResponse{Message: "Post not found"}, nil
	}
	if err != nil {
		return nil, err
	}

	// Cache the post data in Redis
    postJSON, err := json.Marshal(post)
    if err == nil {
        s.rdb.Set(ctx, cacheKey, postJSON, time.Hour) // Cache for 1 hour
    }

	// Generate pre-signed URL for the image
    var imageURL string
	if post.ContentImagePath != "" {
		presignedURL, err := s.minioClient.PresignedGetObject(ctx, "social-media-bucket", post.ContentImagePath, time.Hour, nil)
		if err == nil {
			imageURL = presignedURL.String()
		}
	}
	
	return &authen_and_post.GetPostResponse{
		Message: "Post found",
		Post: &authen_and_post.Post{
			PostId: uint64(post.ID),
			UserId: uint64(post.UserID),
			ContentText: post.ContentText,
			ContentImagePath: imageURL,
			Visible: post.Visible,
			CreatedAt: timestamppb.New(post.CreatedAt),
		},
	}, nil
}

func (s *AuthenAndPostService) EditPost (ctx context.Context, request *authen_and_post.EditPostRequest) (*authen_and_post.EditPostResponse, error) {
	var post models.Post
	err := s.db.WithContext(ctx).First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &authen_and_post.EditPostResponse{Message: "Post not found"}, nil
	}

	if err != nil {
		return nil, err
	}

	if request.ContentText != nil {
		post.ContentText = request.GetContentText()
	}

	// Check if the image is being updated
	if request.ContentImagePath != nil && *request.ContentImagePath != post.ContentImagePath {
		// Delete the old image from MinIO if it exists
        if post.ContentImagePath != "" {
            err := s.minioClient.RemoveObject(ctx, "social-media-bucket", post.ContentImagePath, minio.RemoveObjectOptions{})
            if err != nil {
                return nil, fmt.Errorf("failed to delete old image from MinIO: %w", err)
            }
        }

		// Upload the new image to MinIO
		file, err := os.Open(*request.ContentImagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open new image file: %w", err)
		}
		defer file.Close()
 
		objectName := fmt.Sprintf("posts/%d/%s", post.UserID, filepath.Base(*request.ContentImagePath))
		_, err = s.minioClient.PutObject(ctx, "social-media-bucket", objectName, file, -1, minio.PutObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to upload new image to MinIO: %w", err)
		}
 
		// Update the ContentImagePath in the post
		post.ContentImagePath = objectName
	}

	if request.Visible != nil {
		post.Visible = request.GetVisible()
	}

	// Save the updated post to the database
	if err := s.db.WithContext(ctx).Save(&post).Error; err != nil {
		return nil, fmt.Errorf("failed to edit post: %w", err)
	}

	// Invalidate Redis cache
    cacheKey := fmt.Sprintf("post:%d", request.PostId)
    s.rdb.Del(ctx, cacheKey)

	return &authen_and_post.EditPostResponse{Message: "Post edited successfully"}, nil
}

func (s *AuthenAndPostService) DeletePost (ctx context.Context, request *authen_and_post.DeletePostRequest) (*authen_and_post.DeletePostResponse, error) {
	var post models.Post
	err := s.db.WithContext(ctx).First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &authen_and_post.DeletePostResponse{Message: "Post not found"}, nil
	}

	if err != nil {
		return nil, err
	}

	// Delete the associated image from MinIO if it exists
    if post.ContentImagePath != "" {
        err := s.minioClient.RemoveObject(ctx, "social-media-bucket", post.ContentImagePath, minio.RemoveObjectOptions{})
        if err != nil {
            return nil, fmt.Errorf("failed to delete image from MinIO: %w", err)
        }
    }

	// Delete the post from the database
	if err := s.db.WithContext(ctx).Delete(&post).Error; err != nil {
		return nil, fmt.Errorf("failed to delete post: %w", err)
	}

	// Invalidate Redis cache
    cacheKey := fmt.Sprintf("post:%d", request.PostId)
    s.rdb.Del(ctx, cacheKey)

	return &authen_and_post.DeletePostResponse{Message: "Post deleted successfully"}, nil
}

func (s *AuthenAndPostService) CreateComment (ctx context.Context, request *authen_and_post.CreateCommentRequest) (*authen_and_post.CreateCommentResponse, error) {
	if err := s.checkUserExisting(ctx, request.UserId); err != nil {
        return nil, err
    }

	var post models.Post
	err := s.db.WithContext(ctx).First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &authen_and_post.CreateCommentResponse{Message: "Post not found"}, nil
	}

	if err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Preload("Comments").First(&post, request.PostId).Error; err != nil {
        return nil, fmt.Errorf("failed to fetch post: %w", err)
    }

	comment := models.Comment{
		UserID: uint(request.UserId),
		PostID: uint(request.PostId),
		ContentText: request.ContentText,
	}

	if err := s.db.WithContext(ctx).Model(&post).Association("Comments").Append(&comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &authen_and_post.CreateCommentResponse{
		Message: "Comment created successfully", 
		CommentId: uint64(comment.ID),
	}, nil
}

func (s *AuthenAndPostService) LikePost (ctx context.Context, request *authen_and_post.LikePostRequest) (*authen_and_post.LikePostResponse, error) {
	var user models.User
	err := s.db.WithContext(ctx).First(&user, request.UserId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &authen_and_post.LikePostResponse{Message: "User not found"}, nil
	}

	var post models.Post
	err = s.db.WithContext(ctx).First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &authen_and_post.LikePostResponse{Message: "Post not found"}, nil
	}

	if err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Preload("LikedUsers").First(&post, request.PostId).Error; err != nil {
        return nil, fmt.Errorf("failed to fetch post: %w", err)
    }

	if err := s.db.WithContext(ctx).Model(&post).Association("LikedUsers").Append(&user); err != nil {
		return nil, fmt.Errorf("failed to like post: %w", err)
	}

	return &authen_and_post.LikePostResponse{Message: "Post liked successfully"}, nil
}