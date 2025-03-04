package web_app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
	v1 "github.com/haiyen11231/social-media-app.git/internal/app/web_app/v1"
)

type WebController struct {
	WebService service.WebService
	Port       int
}

func (c *WebController) Run() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the web service!")
	})

	// // Add middleware
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	// Initialize Swagger, Pprof, and Prometheus
	// initSwagger(r)
	// initPprof(r)
	// initPrometheus(r)

	// Register routes
	v1Router := r.Group("/v1")
	v1.AddUserRouter(v1Router, &c.WebService)
	v1.AddFollowingRouter(v1Router, &c.WebService)
	v1.AddPostRouter(v1Router, &c.WebService)
	v1.AddNewsfeedRouter(v1Router, &c.WebService)

	r.Run(fmt.Sprintf(":%d", c.Port))

	// // Graceful shutdown
	// server := &http.Server{
	// 	Addr:    fmt.Sprintf(":%d", c.Port),
	// 	Handler: r,
	// }

	// go func() {
	// 	log.Printf("Web server started on port %d\n", c.Port)
	// 	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("Failed to start server: %v\n", err)
	// 	}
	// }()

	// // Wait for interrupt signal
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// <-quit

	// log.Println("Shutting down server...")

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// if err := server.Shutdown(ctx); err != nil {
	// 	log.Fatalf("Server shutdown failed: %v\n", err)
	// }

	// log.Println("Server stopped.")
}	

// initSwagger

// initPprof

// initPrometheus