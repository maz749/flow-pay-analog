package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/maz749/flow-pay-analog/internal/middleware"
	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/pkg/utils"
)

type InvitationHandler struct {
	invitationRepo *repository.InvitationRepository
	teamRepo       *repository.TeamRepository
	userRepo       *repository.UserRepository
	jwtSecret      string
}

func NewInvitationHandler(
	invitationRepo *repository.InvitationRepository,
	teamRepo *repository.TeamRepository,
	userRepo *repository.UserRepository,
	jwtSecret string,
) *InvitationHandler {
	return &InvitationHandler{
		invitationRepo: invitationRepo,
		teamRepo:       teamRepo,
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (h *InvitationHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if user is admin
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if user.Role != "admin" {
		utils.SendError(w, http.StatusForbidden, "Only administrators can invite users")
		return
	}

	var req models.CreateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Role == "" {
		utils.SendError(w, http.StatusBadRequest, "Email and role are required")
		return
	}

	// Validate role
	if req.Role != "admin" && req.Role != "user" {
		utils.SendError(w, http.StatusBadRequest, "Invalid role. Must be 'admin' or 'user'")
		return
	}

	// Check if team exists (if provided)
	if req.TeamID != nil {
		_, err := h.teamRepo.GetByID(*req.TeamID)
		if err != nil {
			utils.SendError(w, http.StatusNotFound, "Team not found")
			return
		}

		// Check if user is owner of the team
		isOwner, err := h.teamRepo.IsUserOwner(*req.TeamID, userID)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to check team ownership")
			return
		}

		if !isOwner {
			utils.SendError(w, http.StatusForbidden, "Only team owner can invite to this team")
			return
		}
	}

	// Check if user with this email already exists
	existingUser, _ := h.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		utils.SendError(w, http.StatusBadRequest, "User with this email already exists")
		return
	}

	// Generate invitation token
	token, err := generateToken()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to generate invitation token")
		return
	}

	// Create invitation
	invitation := &models.Invitation{
		Email:     req.Email,
		Role:      req.Role,
		TeamID:    req.TeamID,
		Token:     token,
		Status:    "pending",
		InvitedBy: userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := h.invitationRepo.Create(invitation); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create invitation")
		return
	}

	// In production, send email here
	// For now, just return the invitation with token
	invitationURL := "http://localhost:8080/accept-invitation?token=" + token

	response := map[string]interface{}{
		"invitation": invitation,
		"url":        invitationURL,
		"message":    "Invitation created successfully. In production, an email would be sent to " + req.Email,
	}

	utils.SendJSON(w, http.StatusCreated, response)
}

func (h *InvitationHandler) GetMyInvitations(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	invitations, err := h.invitationRepo.GetByInviter(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to get invitations")
		return
	}

	utils.SendJSON(w, http.StatusOK, invitations)
}

func (h *InvitationHandler) GetInvitationByToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	if token == "" {
		utils.SendError(w, http.StatusBadRequest, "Token is required")
		return
	}

	invitation, err := h.invitationRepo.GetInvitationWithDetails(token)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Invitation not found")
		return
	}

	// Check if invitation is expired
	if invitation.Status != "pending" {
		utils.SendError(w, http.StatusBadRequest, "Invitation is no longer valid")
		return
	}

	if time.Now().After(invitation.ExpiresAt) {
		// Update status to expired
		h.invitationRepo.UpdateStatus(invitation.ID, "expired")
		utils.SendError(w, http.StatusBadRequest, "Invitation has expired")
		return
	}

	utils.SendJSON(w, http.StatusOK, invitation)
}

func (h *InvitationHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	var req models.AcceptInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" || req.Password == "" || req.Username == "" {
		utils.SendError(w, http.StatusBadRequest, "Token, password, and username are required")
		return
	}

	// Get invitation
	invitation, err := h.invitationRepo.GetByToken(req.Token)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Invitation not found")
		return
	}

	// Validate invitation
	if invitation.Status != "pending" {
		utils.SendError(w, http.StatusBadRequest, "Invitation is no longer valid")
		return
	}

	if time.Now().After(invitation.ExpiresAt) {
		h.invitationRepo.UpdateStatus(invitation.ID, "expired")
		utils.SendError(w, http.StatusBadRequest, "Invitation has expired")
		return
	}

	// Check if user already exists
	existingUser, _ := h.userRepo.GetByEmail(invitation.Email)
	if existingUser != nil {
		utils.SendError(w, http.StatusBadRequest, "User already exists")
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
		Email:        invitation.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Role:         invitation.Role,
	}

	if err := h.userRepo.Create(user); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Add user to team if specified
	if invitation.TeamID != nil {
		if err := h.teamRepo.AddMember(*invitation.TeamID, user.ID); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to add user to team")
			return
		}
	}

	// Update invitation status
	if err := h.invitationRepo.UpdateStatus(invitation.ID, "accepted"); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to update invitation status")
		return
	}

	// Generate JWT token
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

func (h *InvitationHandler) DeleteInvitation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	token := vars["token"]

	if token == "" {
		utils.SendError(w, http.StatusBadRequest, "Token is required")
		return
	}

	invitation, err := h.invitationRepo.GetByToken(token)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Invitation not found")
		return
	}

	// Only the inviter can delete the invitation
	if invitation.InvitedBy != userID {
		utils.SendError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.invitationRepo.Delete(invitation.ID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to delete invitation")
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Invitation deleted successfully"})
}
