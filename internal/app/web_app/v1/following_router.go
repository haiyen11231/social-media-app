package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

func AddFollowingRouter(r *gin.RouterGroup, svc *service.WebService) {
	followingRouter := r.Group("/following")

	followingRouter.Use(svc.AuthenticateUser)
	followingRouter.POST("/:id", svc.FollowUser)
	followingRouter.DELETE("/:id", svc.UnfollowUser)
	followingRouter.GET("/:id", svc.GetFollowerList)
}