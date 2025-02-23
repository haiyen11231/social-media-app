package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

func AddPostRouter(r *gin.RouterGroup, svc *service.WebService) {
	postRouter := r.Group("/posts")
	
	postRouter.Use(svc.AuthenticateUser)
	postRouter.POST("/", )
	postRouter.GET("/:id", )
	postRouter.PUT("/:id", )
	postRouter.DELETE("/:id", )

	postRouter.POST("/:id/comments", )
	postRouter.POST("/:id/likes", )
}