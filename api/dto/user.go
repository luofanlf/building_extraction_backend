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

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type UserResponse struct {
	Username        string `json:"username"`
	ExtractionCount int    `json:"extraction_count"`
	RemainingCount  int    `json:"remaining_count"`
	UserRole        int    `json:"is_admin"`
}

func UserToResponse(user *model.User) UserResponse {
	return UserResponse{
		Username:        user.Username,
		ExtractionCount: user.ExtractionCount,
		RemainingCount:  user.RemainingCount,
		UserRole:        user.UserRole,
	}
}
