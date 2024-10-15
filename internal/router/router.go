package router

import (
	"log"

	"github.com/bulgil/blog-rest-api/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func NewRouter(env string, logger *logrus.Logger, storage *pgxpool.Pool) *gin.Engine {
	engine := gin.New()

	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	routes(engine, logger, storage)

	log.Println("router initialized")
	return engine
}

func routes(router *gin.Engine, logger *logrus.Logger, storage *pgxpool.Pool) {
	router.POST("/posts", handlers.CreatePost(storage, logger))
	router.PUT("/posts/:id", handlers.UpdatePost(storage, logger))
	router.DELETE("/posts/:id", handlers.DeletePost(storage, logger))
	router.GET("posts/:id", handlers.GetPost(storage, logger))
	router.GET("/posts", handlers.GetAllPosts(storage, logger))
}
