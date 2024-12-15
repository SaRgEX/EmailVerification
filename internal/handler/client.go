package handler

import (
	"context"
	"email-verification-service/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (h *Handler) register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var input model.ClientInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.s.Register(ctx, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, statusResponse{
		Status: "registered",
	})
}

func (h *Handler) verifyEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var input model.ClientVerification
	err := c.ShouldBindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.s.Verify(ctx, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, statusResponse{
		Status: "verified",
	})
}

func (h *Handler) refresh(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var input model.ClientVerification
	err := c.ShouldBindJSON(&input)

	_, err = h.s.Refresh(ctx, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "refreshed",
	})
}
