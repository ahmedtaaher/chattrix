package handler

import (
	"chattrix/dto"
	"chattrix/service"
	"chattrix/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const MaxFileSize = 2 << 20

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) SignUp(context *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := context.ShouldBindJSON(&registerRequest); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Register(registerRequest); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusCreated, "user registered successfully", nil)
}

func (h *AuthHandler) Login(context *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := context.ShouldBindJSON(&loginRequest); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Login(loginRequest)
	if err != nil {
		utils.ErrorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "login successful", gin.H{
		"token": token,
		"type":  "Bearer",
	})
}

func (h *AuthHandler) GetProfile(context *gin.Context) {
	userIDRaw, exists := context.Get("user_id")
	if !exists {
		utils.ErrorResponse(context, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(context, http.StatusUnauthorized, "invalid user id")
		return
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "profile fetched successfully", gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"avatar_url": user.AvatarURL,
		"created_at": user.CreatedAt,
	})
}

func (h *AuthHandler) UpdateProfile(context *gin.Context) {
	var updateRequest dto.UpdateProfileRequest

	if err := context.ShouldBindJSON(&updateRequest); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	userIDRaw, exists := context.Get("user_id")
	if !exists {
		utils.ErrorResponse(context, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(context, http.StatusUnauthorized, "invalid user id")
		return
	}

	err := h.service.UpdateProfile(userID, updateRequest)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "profile updated successfully", nil)
}

func (h *AuthHandler) ChangePassword(context *gin.Context) {
	var changePasswordRequest dto.ChangePasswordRequest

	if err := context.ShouldBindJSON(&changePasswordRequest); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	userIDRaw, exists := context.Get("user_id")
	if !exists {
		utils.ErrorResponse(context, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(context, http.StatusUnauthorized, "invalid user id")
		return
	}

	if err := h.service.ChangePassword(userID, changePasswordRequest); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "password changed successfully", nil)
}

func (h *AuthHandler) UploadAvatar(context *gin.Context) {
	userIDRaw, exists := context.Get("user_id")
	if !exists {
		utils.ErrorResponse(context, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(context, http.StatusUnauthorized, "invalid user id")
		return
	}

	file, err := context.FormFile("avatar")
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "avatar file is required")
		return
	}

	if file.Size > MaxFileSize {
		utils.ErrorResponse(context, http.StatusBadRequest, "file too large (max 2MB)")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		utils.ErrorResponse(context, http.StatusBadRequest, "only .jpg, .jpeg, .png allowed")
		return
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "failed to create upload directory")
		return
	}

	filename := fmt.Sprintf("%s_%d%s", userID.String(), time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := context.SaveUploadedFile(file, filePath); err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "failed to save file")
		return
	}

	avatarPath := "/uploads/avatars/" + filename

	if user.AvatarURL != nil && *user.AvatarURL != "" {
		oldPath := "." + *user.AvatarURL
		if err := os.Remove(oldPath); err != nil {
			fmt.Println("Delete failed:", err) 
		}
	}

	if err = h.service.UpdateAvatar(userID, avatarPath); err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "failed to update avatar")
		return
	}

	fullURL := fmt.Sprintf("http://localhost:8080%s", avatarPath)
	utils.SuccessResponse(context, http.StatusOK, "avatar uploaded successfully", gin.H{
		"avatar_url": fullURL,
	})
}

func (h *AuthHandler) SearchUsers(context *gin.Context) {
	query := context.Query("q")

	users, err := h.service.SearchUsers(query)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "users found", users)
}

