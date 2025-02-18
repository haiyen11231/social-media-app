package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/models"
)

func (svc *WebService) CreatePost(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")
	var jsonRequest models.CreatePostRequest
	if err := ctx.ShouldBindJSON(&jsonRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: err.Error()})
		return
	}

	response, err := svc.authenAndPostClient.CreatePost(ctx, &authen_and_post.CreatePostRequest{
		UserId: userId,
		ContentText: jsonRequest.ContentText,
		ContentImagePath: jsonRequest.ContentImagePath,
		Visible: jsonRequest.Visible,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (svc *WebService) GetPost(ctx *gin.Context) {
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: "Invalid post ID"})
		return
	}

	reponse, err := svc.authenAndPostClient.GetPost(ctx, &authen_and_post.GetPostRequest{
		PostId: uint64(postId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.PostDetailResponse{
		PostID: uint(reponse.Post.PostId),
		UserID: uint(reponse.Post.UserId),
		ContentText: reponse.Post.ContentText,
		ContentImagePath: reponse.Post.ContentImagePath,
		Visible: reponse.Post.Visible,
		CreatedAt: reponse.Post.CreatedAt.AsTime(),
	})
}

func (svc *WebService) EditPost(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: "Invalid post ID"})
		return
	}

	var jsonRequest models.EditPostRequest
	if err := ctx.ShouldBindJSON(&jsonRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: err.Error()})
		return
	}

	var contentText *string
	if jsonRequest.ContentText != nil {
		contentText = jsonRequest.ContentText
	}

	var contentImagePath *string
	if jsonRequest.ContentImagePath != nil {
		contentImagePath = jsonRequest.ContentImagePath
	}

	var visible *bool 
	if jsonRequest.Visible != nil {
		visible = jsonRequest.Visible
	}

	response, err := svc.authenAndPostClient.EditPost(ctx, &authen_and_post.EditPostRequest{
		PostId: uint64(postId),
		UserId: userId,
		ContentText: contentText,
		ContentImagePath: contentImagePath,
		Visible: visible,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (svc *WebService) DeletePost(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: "Invalid post ID"})
		return
	}

	response, err := svc.authenAndPostClient.DeletePost(ctx, &authen_and_post.DeletePostRequest{
		PostId: uint64(postId),
		UserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (svc *WebService) CreateComment(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: "Invalid post ID"})
		return
	}

	var jsonRequest models.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&jsonRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: err.Error()})
		return
	}

	response, err := svc.authenAndPostClient.CreateComment(ctx, &authen_and_post.CreateCommentRequest{
		UserId: userId,
		PostId: uint64(postId),
		ContentText: jsonRequest.ContentText,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (svc *WebService) LikePost(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: "Invalid post ID"})
		return
	}

	response, err := svc.authenAndPostClient.LikePost(ctx, &authen_and_post.LikePostRequest{
		UserId: userId,
		PostId: uint64(postId),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}