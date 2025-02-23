package authen_and_post_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/models"
)

func (s *AuthenAndPostService) checkUserExisting(ctx context.Context, userID uint64) error {
    var user models.User
    if err := s.db.WithContext(ctx).Table("user").Where("id = ?", userID).First(&user).Error; err != nil {
        return fmt.Errorf("user with ID %d not found: %w", userID, err)
    }
    return nil
}

func (s *AuthenAndPostService) FollowUser (ctx context.Context, request *authen_and_post.FollowUserRequest) (*authen_and_post.FollowUserResponse, error) {
	if request.UserId == request.FollowingId {
        return nil, errors.New("a user cannot follow themselves")
    }

    if err := s.checkUserExisting(ctx, request.UserId); err != nil {
        return nil, err
    }

    if err := s.checkUserExisting(ctx, request.FollowingId); err != nil {
        return nil, err
    }

	var user models.User
    if err := s.db.WithContext(ctx).Preload("Following").First(&user, request.UserId).Error; err != nil {
        return nil, fmt.Errorf("failed to fetch user: %w", err)
    }

    var alreadyFollowed bool
	for _, following := range user.Following {
        if following.ID == uint(request.FollowingId) {
            alreadyFollowed = true
            break
        }
    }

	if !alreadyFollowed {
		var followingUser models.User
        if err := s.db.WithContext(ctx).First(&followingUser, request.FollowingId).Error; err != nil {
            return nil, fmt.Errorf("failed to fetch following user: %w", err)
        }

		if err := s.db.WithContext(ctx).Model(&user).Association("Following").Append(&followingUser); err != nil {
			return nil, fmt.Errorf("failed to add into Following: %w", err)
		}

		if err := s.db.WithContext(ctx).Model(&followingUser).Association("Followers").Append(&user); err != nil {
			return nil, fmt.Errorf("failed to add into Follower: %w", err)
		}
	}

	return &authen_and_post.FollowUserResponse{Message: "Followed"}, nil
}

func (s *AuthenAndPostService) UnfollowUser (ctx context.Context, request *authen_and_post.UnfollowUserRequest) (*authen_and_post.UnfollowUserResponse, error) {
	if err := s.checkUserExisting(ctx, request.UserId); err != nil {
        return nil, err
    }

    if err := s.checkUserExisting(ctx, request.FollowingId); err != nil {
        return nil, err
    }

	var user models.User
	if err := s.db.WithContext(ctx).Preload("Following").First(&user, request.UserId).Error; err != nil {
        return nil, fmt.Errorf("failed to fetch user: %w", err)
    }

    var alreadyFollowed bool
    for _, following := range user.Following {
        if following.ID == uint(request.FollowingId) {
            alreadyFollowed = true
            break
        }
    }

	if alreadyFollowed {
		var followingUser models.User
		if err := s.db.WithContext(ctx).First(&followingUser, request.FollowingId).Error; err != nil {
            return nil, fmt.Errorf("failed to fetch following user: %w", err)
        }
		
		if err := s.db.WithContext(ctx).Model(&user).Association("Following").Delete(&followingUser); err != nil {
			return nil, fmt.Errorf("failed to remove from Following: %w", err)
		}

		if err := s.db.WithContext(ctx).Model(&followingUser).Association("Followers").Delete(&user); err != nil {
			return nil, fmt.Errorf("failed to remove from Follower: %w", err)
		}
	}

	return &authen_and_post.UnfollowUserResponse{Message: "Unfollowed"}, nil
}

func (s *AuthenAndPostService) GetFollowerList (ctx context.Context, request *authen_and_post.GetFollowerListRequest) (*authen_and_post.GetFollowerListResponse, error) {
	if err := s.checkUserExisting(ctx, request.UserId); err != nil {
        return nil, err
    }

	var user models.User
	if err := s.db.WithContext(ctx).Preload("Followers").First(&user, request.UserId).Error; err != nil {
        return nil, fmt.Errorf("failed to fetch user: %w", err)
    }
	
	var followerList []*authen_and_post.GetFollowerListResponse_FollowerInfo
	for _, follower := range user.Followers {
		followerList = append(followerList, &authen_and_post.GetFollowerListResponse_FollowerInfo{
			UserId: uint64(follower.ID),
			FirstName: follower.FirstName,
			LastName: follower.LastName,
			Username: follower.Username,
		})
	}

	return &authen_and_post.GetFollowerListResponse{
		Message: "Success",
		Followers: followerList,
	}, nil
}

