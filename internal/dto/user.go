package dto

import (
	"time"

	"gin-learning/internal/model"
)

type CreateUserRequest struct {
	Name  string `json:"name"  binding:"required,min=1,max=100"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"  binding:"required,min=1,max=100"`
	Email string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToUserResponseList(users []model.User) []UserResponse {
	resp := make([]UserResponse, len(users))
	for i := range users {
		resp[i] = ToUserResponse(&users[i])
	}
	return resp
}
