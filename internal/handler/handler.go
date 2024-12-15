package handler

import (
	"email-verification-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/sign-up", h.register)
	email := router.Group("/email")
	{
		email.POST("/verify", h.verifyEmail)
		email.POST("/refresh", h.refresh)

	}
	return router
}
