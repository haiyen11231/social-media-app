package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/models"
)

func (svc *WebService) FollowUser(ctx *gin.Context) {
	friendID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: "Invalid user ID"})
		return
	}


}

func (svc *WebService) UnfollowUser(ctx *gin.Context) {

}

func (svc *WebService) GetFollowerList(ctx *gin.Context) {
	
}