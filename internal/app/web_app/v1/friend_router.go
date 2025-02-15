package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

func AddFriendRouter(r *gin.RouterGroup, svc *service.WebService) {
	friendRouter := r.Group("/friend")

	friendRouter.POST("/:id", )
	friendRouter.DELETE("/:id", )
	friendRouter.GET("/:id", )
}