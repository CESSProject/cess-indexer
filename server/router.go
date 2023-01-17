package server

import (
	resp "cess-indexer/server/response"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())
	router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		resp.RespError(c, resp.NewError(http.StatusInternalServerError, err.(error)))
	}))
	return router
}
