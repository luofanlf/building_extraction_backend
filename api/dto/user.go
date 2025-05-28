package dto

import "building-extraction/internal/model"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type UserResponse struct {
	Username        string `json:"username"`
	ExtractionCount int    `json:"extraction_count"`
	RemainingCount  int    `json:"remaining_count"`
}

func UserToResponse(user *model.User) UserResponse {
	return UserResponse{
		Username:        user.Username,
		ExtractionCount: user.ExtractionCount,
		RemainingCount:  user.RemainingCount,
	}
}
