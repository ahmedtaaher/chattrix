package dto

type RegisterRequest struct {
  Username string `json:"username" binding:"required,min=3,max=50"`
  Nickname string `json:"nickname" binding:"required,min=3,max=50"`
  Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
  Username string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
  Nickname string `json:"nickname" binding:"required,min=3,max=50"`
}

type ChangePasswordRequest struct {
  CurrentPassword string `json:"current_password" binding:"required"`
  NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UserResponse struct {
  Token    string `json:"token"`
}