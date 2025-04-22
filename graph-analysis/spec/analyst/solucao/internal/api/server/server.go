package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()
	s := &Server{router: r}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}

func (s *Server) Run() error {
	return s.router.Run(":8080")
}
