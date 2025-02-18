package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/newsfeed"
	"github.com/haiyen11231/social-media-app.git/internal/models"
)

func (svc *WebService) GetNewsfeed(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")

	response, err := svc.newsfeedClient.GetNewsfeed(ctx, &newsfeed.GetNewsfeedRequest{
		UserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}