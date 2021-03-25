package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mazxaxz/donut-batcher/pkg/rest"
)

func setupRouting(handlers ...rest.SetupRouterer) http.Handler {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	/* no need for CORS right now */

	v1 := router.Group("v1")
	for _, handler := range handlers {
		handler.SetupRouter(v1)
	}

	/* In a production API we should add requestI middleware, but for this exercise it does not matter */

	return router
}
