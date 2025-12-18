package handlers

import (
    "net/http"
    "log"

    "paydeya-backend/internal/repositories"
    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type ProfileHandler struct {
    authService *services.AuthService
    userRepo    *repositories.UserRepository
    fileService *services.FileService
}

func NewProfileHandler(authService *services.AuthService, userRepo *repositories.UserRepository, fileService *services.FileService) *ProfileHandler {
    return &ProfileHandler{
        authService: authService,
        userRepo:    userRepo,
        fileService: fileService,
    }
}

// GetProfile godoc
// @Summary Получить профиль пользователя
// @Description Возвращает данные профиля текущего пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} ProfileResponse "Данные профиля"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
    userID := c.GetInt("userID")

    // Получаем пользователя из БД
    user, err := h.userRepo.GetUserProfile(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
        return
    }

    if user == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Получаем специализации из БД
    specializations, err := h.userRepo.GetUserSpecializations(c.Request.Context(), userID)
    if err != nil {
        // Логируем ошибку но продолжаем (специализации не критичны)
        log.Printf("Warning: failed to get specializations for user %d: %v", userID, err)
        specializations = []string{}
    }

    c.JSON(http.StatusOK, gin.H{
        "id":               user.ID,
        "email":            user.Email,
        "fullName":         user.FullName,
        "role":             user.Role,
        "avatarUrl":        user.AvatarURL,
        "isVerified":       user.IsVerified,
        "specializations":  specializations,
        "createdAt":        user.CreatedAt,
        "updatedAt":        user.UpdatedAt,
    })
}

// UpdateProfile godoc
// @Summary Обновить профиль
// @Description Обновляет данные профиля пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body UpdateProfileRequest true "Данные для обновления"
// @Success 200 {object} UpdateProfileResponse "Профиль обновлен"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
    userID := c.GetInt("userID")

    var req UpdateProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Обновляем данные в БД
    err := h.userRepo.UpdateUserProfile(c.Request.Context(), userID, req.FullName, req.Specializations)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Profile updated successfully",
        "userID":  userID,
        "data":    req,
    })
}

// UploadAvatar godoc
// @Summary Загрузить аватар
// @Description Загружает аватар пользователя
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param avatar formData file true "Файл аватара (макс. 5MB)"
// @Success 200 {object} UploadAvatarResponse "Аватар загружен"
// @Failure 400 {object} ErrorResponse "Неверный файл"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /profile/avatar [post]
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
    userID := c.GetInt("userID")

    // Получаем файл из формы
    file, err := c.FormFile("avatar")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
        return
    }

    // Проверяем размер файла (макс 5MB)
    if file.Size > 5*1024*1024 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 5MB"})
        return
    }

    // Получаем текущего пользователя чтобы удалить старый аватар
    user, err := h.userRepo.GetUserByID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
        return
    }

    // Сохраняем новый аватар
    avatarURL, err := h.fileService.SaveAvatar(userID, file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save avatar: " + err.Error()})
        return
    }

    // Удаляем старый аватар если он был
    if user.AvatarURL != "" {
        h.fileService.DeleteAvatar(user.AvatarURL)
    }

    // Обновляем аватар в БД
    err = h.userRepo.UpdateUserAvatar(c.Request.Context(), userID, avatarURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar in database"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":   "Avatar uploaded successfully",
        "avatarUrl": avatarURL,
    })
}

// Request/Response models for Swagger

// ProfileResponse represents user profile response
// @Description Ответ с данными профиля пользователя
type ProfileResponse struct {
    ID              int       `json:"id" example:"123"`
    Email           string    `json:"email" example:"user@example.com"`
    FullName        string    `json:"fullName" example:"Иван Иванов"`
    Role            string    `json:"role" example:"teacher"`
    AvatarURL       string    `json:"avatarUrl" example:"https://example.com/avatars/123.jpg"`
    IsVerified      bool      `json:"isVerified" example:"true"`
    Specializations []string  `json:"specializations" example:"math,physics"`
    CreatedAt       string    `json:"createdAt" example:"2023-01-15T10:30:00Z"`
    UpdatedAt       string    `json:"updatedAt" example:"2023-01-15T10:30:00Z"`
}

// UpdateProfileRequest represents update profile request
// @Description Запрос на обновление профиля
type UpdateProfileRequest struct {
    FullName        string   `json:"fullName" example:"Иван Иванов"`
    Specializations []string `json:"specializations" example:"math,physics"`
}

// UpdateProfileResponse represents update profile response
// @Description Ответ на обновление профиля
type UpdateProfileResponse struct {
    Message string               `json:"message" example:"Profile updated successfully"`
    UserID  int                  `json:"userID" example:"123"`
    Data    UpdateProfileRequest `json:"data"`
}

// UploadAvatarResponse represents upload avatar response
// @Description Ответ на загрузку аватара
type UploadAvatarResponse struct {
    Message   string `json:"message" example:"Avatar uploaded successfully"`
    AvatarURL string `json:"avatarUrl" example:"https://example.com/avatars/123.jpg"`
}

