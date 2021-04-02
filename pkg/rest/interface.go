package rest

import "github.com/gin-gonic/gin"

type SetupRouterer interface {
	SetupRouter(group *gin.RouterGroup)
}
