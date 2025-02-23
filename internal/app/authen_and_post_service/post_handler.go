package authen_and_post_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/models"
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

	post := models.Post{
		ContentText: request.ContentText,
		ContentImagePath: request.ContentImagePath,
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
	var post models.Post
	err := s.db.WithContext(ctx).First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &authen_and_post.GetPostResponse{Message: "Post not found"}, nil
	}

	if err != nil {
		return nil, err
	}

	return &authen_and_post.GetPostResponse{
		Message: "Post found",
		Post: &authen_and_post.Post{
			PostId: uint64(post.ID),
			UserId: uint64(post.UserID),
			ContentText: post.ContentText,
			ContentImagePath: post.ContentImagePath,
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

	if request.ContentImagePath != nil {
		post.ContentImagePath = request.GetContentImagePath()
	}

	if request.Visible != nil {
		post.Visible = request.GetVisible()
	}

	if err := s.db.WithContext(ctx).Save(&post).Error; err != nil {
		return nil, fmt.Errorf("failed to edit post: %w", err)
	}

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

	if err := s.db.WithContext(ctx).Delete(&post).Error; err != nil {
		return nil, fmt.Errorf("failed to delete post: %w", err)
	}

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