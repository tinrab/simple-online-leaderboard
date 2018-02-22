package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	api := router.Group("/api")
	{
		api.POST("/scores", postScoreEndpoint)
		api.GET("/scores", getScoresEndpoint)
	}

	http.Handle("/", cors.AllowAll().Handler(router))
}
