package dto

import "time"

// UserCreateRequest is used when creating a new user.
type UserCreateRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
	FirstName string `json:"firstName" binding:"required,max=100"`
	LastName  string `json:"lastName" binding:"required,max=100"`
	Role      string `json:"role" binding:"required,max=50"`
}

// UserUpdateRequest modifies an existing user.
type UserUpdateRequest struct {
	Email     string `json:"email" binding:"required,email,max=255"`
	FirstName string `json:"firstName" binding:"required,max=100"`
	LastName  string `json:"lastName" binding:"required,max=100"`
	Role      string `json:"role" binding:"required,max=50"`
}

// UserResponse is returned back to clients.
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserListResponse wraps a paginated list.
type UserListResponse struct {
	Total int64          `json:"total"`
	Items []UserResponse `json:"items"`
}
