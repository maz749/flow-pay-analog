package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/pkg/utils"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
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
