package web_app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
	v1 "github.com/haiyen11231/social-media-app.git/internal/app/web_app/v1"
)

type WebController struct {
	WebService service.WebService
	Port int
}

func (c *WebController) Run() {
	r := gin.Default()

	v1Router := r.Group("/v1")
	v1.AddUserRouter(v1Router, &c.WebService)
	v1.AddFriendRouter(v1Router, &c.WebService)
	v1.AddPostRouter(v1Router, &c.WebService)
	v1.AddNewsfeedRouter(v1Router, &c.WebService)

	r.Run(fmt.Sprintf(":%d", c.Port))
}

// initSwagger

// initPprof

// initPrometheus