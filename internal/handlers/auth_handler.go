package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/pkg/utils"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	botToken  string
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret, botToken string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		botToken:  botToken,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.Username == "" {
		utils.SendError(w, http.StatusBadRequest, "Email, password, and username are required")
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		utils.SendError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Validate password length
	if len(req.Password) < 6 {
		utils.SendError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	// Check if user already exists
	existingUser, _ := h.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		utils.SendError(w, http.StatusConflict, "User with this email already exists")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Username:     req.Username,
	}

	if err := h.userRepo.Create(user); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	utils.SendJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		utils.SendError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		utils.SendError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	utils.SendJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SendJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) TelegramAuth(w http.ResponseWriter, r *http.Request) {
	var rawData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Convert to string map for validation
	authData := make(map[string]string)
	for key, value := range rawData {
		authData[key] = fmt.Sprintf("%v", value)
	}

	// Validate Telegram auth data
	if err := utils.ValidateTelegramAuth(authData, h.botToken); err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Invalid Telegram authentication data: "+err.Error())
		return
	}

	// Parse Telegram user data
	telegramID, _ := strconv.ParseInt(authData["id"], 10, 64)
	firstName := authData["first_name"]
	lastName := authData["last_name"]
	username := authData["username"]
	photoURL := authData["photo_url"]

	// Try to find existing user
	user, err := h.userRepo.GetByTelegramID(telegramID)
	if err != nil {
		// User doesn't exist, create new one
		displayName := firstName
		if username != "" {
			displayName = username
		}

		user = &models.User{
			Username:          displayName,
			TelegramUserID:    &telegramID,
			TelegramUsername:  strToPtr(username),
			TelegramFirstName: strToPtr(firstName),
			TelegramLastName:  strToPtr(lastName),
			TelegramPhotoURL:  strToPtr(photoURL),
			AuthType:          "telegram",
		}

		if err := h.userRepo.CreateTelegramUser(user); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
			return
		}

		// Create telegram settings with chat_id = telegram_user_id
		chatID := telegramID
		settings := &models.TelegramSettings{
			UserID:    user.ID,
			ChatID:    &chatID,
			IsEnabled: true,
		}
		if err := h.userRepo.UpsertTelegramSettings(settings); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to create telegram settings: %v\n", err)
		}
	}

	// Generate token - use telegram username or first name as "email" for JWT
	email := username
	if email == "" {
		email = fmt.Sprintf("telegram_%d", telegramID)
	}

	token, err := utils.GenerateToken(user.ID, email, h.jwtSecret)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	utils.SendJSON(w, http.StatusOK, response)
}

func strToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
