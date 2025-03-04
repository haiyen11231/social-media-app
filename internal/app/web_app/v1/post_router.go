package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

func AddPostRouter(r *gin.RouterGroup, svc *service.WebService) {
	postRouter := r.Group("/posts")

	postRouter.Use(svc.AuthenticateUser)
	postRouter.POST("/", svc.CreatePost)
	postRouter.GET("/:id", svc.GetPost)
	postRouter.PUT("/:id", svc.EditPost)
	postRouter.DELETE("/:id", svc.DeletePost)

	postRouter.POST("/:id/comments", svc.CreateComment)
	postRouter.POST("/:id/likes", svc.LikePost)
}