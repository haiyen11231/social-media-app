package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/models"
)

func (svc *WebService) FollowUser(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")
	followingId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: "Invalid user ID"})
		return
	}

	response, err := svc.authenAndPostClient.FollowUser(ctx, &authen_and_post.FollowUserRequest{
		UserId:      userId,
		FollowingId: uint64(followingId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (svc *WebService) UnfollowUser(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")
	followingId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: "Invalid user ID"})
		return
	}

	response, err := svc.authenAndPostClient.UnfollowUser(ctx, &authen_and_post.UnfollowUserRequest{
		UserId:      userId,
		FollowingId: uint64(followingId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (svc *WebService) GetFollowerList(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")

	response, err := svc.authenAndPostClient.GetFollowerList(ctx, &authen_and_post.GetFollowerListRequest{
		UserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}