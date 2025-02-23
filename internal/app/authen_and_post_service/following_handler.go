package authen_and_post_service

import (
	"context"
	"errors"

	"github.com/haiyen11231/social-media-app.git/internal/models"
)

// Following
func (s *AuthenAndPostService) checkUserExisting(ctx context.Context, userId uint64) error {
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

