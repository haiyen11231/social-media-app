package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

func AddFollowingRouter(r *gin.RouterGroup, svc *service.WebService) {
	followingRouter := r.Group("/following")

	followingRouter.POST("/:id", )
	followingRouter.DELETE("/:id", )
	followingRouter.GET("/:id", )
}