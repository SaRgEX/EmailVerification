package model

import "time"

type ClientInput struct {
	ID               int    `json:"id"`
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
	Username         string `json:"username" binding:"required"`
	Email            string `json:"email" binding:"required"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"-"`
}

type ClientOutput struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	IsVerified bool   `json:"is_verified"`
}

type ClientVerification struct {
	Email string `json:"email" binding:"required"`
	VerificationCode
}

type VerificationCode struct {
	Code       string     `json:"code" binding:"required"`
	IsVerified bool       `json:"-"`
	Exp        *time.Time `json:"-"`
}
