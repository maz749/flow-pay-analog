package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/maz749/flow-pay-analog/internal/middleware"
	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/pkg/utils"
)

type TeamHandler struct {
	teamRepo *repository.TeamRepository
	userRepo *repository.UserRepository
}

func NewTeamHandler(teamRepo *repository.TeamRepository, userRepo *repository.UserRepository) *TeamHandler {
	return &TeamHandler{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		utils.SendError(w, http.StatusBadRequest, "Team name is required")
		return
	}

	team := &models.Team{
		Name:           req.Name,
		Description:    req.Description,
		OwnerID:        &userID,
		BudgetAmount:   req.BudgetAmount,
		BudgetCurrency: req.BudgetCurrency,
		BudgetPeriod:   req.BudgetPeriod,
	}

	if err := h.teamRepo.Create(team); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create team")
		return
	}

	// Add creator as team member
	if err := h.teamRepo.AddMember(team.ID, userID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to add creator to team")
		return
	}

	utils.SendJSON(w, http.StatusCreated, team)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	// Check if user is member or owner
	isMember, err := h.teamRepo.IsUserMember(teamID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to check team membership")
		return
	}

	isOwner, err := h.teamRepo.IsUserOwner(teamID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to check team ownership")
		return
	}

	if !isMember && !isOwner {
		utils.SendError(w, http.StatusForbidden, "Access denied")
		return
	}

	teamWithMembers, err := h.teamRepo.GetTeamWithMembers(teamID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Team not found")
		return
	}

	utils.SendJSON(w, http.StatusOK, teamWithMembers)
}

func (h *TeamHandler) GetMyTeams(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get teams where user is a member
	memberTeams, err := h.teamRepo.GetTeamsByUserID(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to get teams")
		return
	}

	utils.SendJSON(w, http.StatusOK, memberTeams)
}

func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	// Check if user is owner
	isOwner, err := h.teamRepo.IsUserOwner(teamID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to check team ownership")
		return
	}

	if !isOwner {
		utils.SendError(w, http.StatusForbidden, "Only team owner can update the team")
		return
	}

	var req models.UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing team
	team, err := h.teamRepo.GetByID(teamID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Team not found")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		team.Name = *req.Name
	}
	if req.Description != nil {
		team.Description = req.Description
	}
	if req.BudgetAmount != nil {
		team.BudgetAmount = req.BudgetAmount
	}
	if req.BudgetCurrency != nil {
		team.BudgetCurrency = *req.BudgetCurrency
	}
	if req.BudgetPeriod != nil {
		team.BudgetPeriod = req.BudgetPeriod
	}

	if err := h.teamRepo.Update(team); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to update team")
		return
	}

	utils.SendJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	// Check if user is owner
	isOwner, err := h.teamRepo.IsUserOwner(teamID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to check team ownership")
		return
	}

	if !isOwner {
		utils.SendError(w, http.StatusForbidden, "Only team owner can delete the team")
		return
	}

	if err := h.teamRepo.Delete(teamID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to delete team")
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Team deleted successfully"})
}

func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	// Check if user is owner
	isOwner, err := h.teamRepo.IsUserOwner(teamID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to check team ownership")
		return
	}

	if !isOwner {
		utils.SendError(w, http.StatusForbidden, "Only team owner can add members")
		return
	}

	var req models.AddTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if user exists
	_, err = h.userRepo.GetByID(req.UserID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "User not found")
		return
	}

	if err := h.teamRepo.AddMember(teamID, req.UserID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to add member")
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Member added successfully"})
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	memberID, err := strconv.Atoi(vars["memberId"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid member ID")
		return
	}

	// Check if user is owner
	isOwner, err := h.teamRepo.IsUserOwner(teamID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to check team ownership")
		return
	}

	if !isOwner {
		utils.SendError(w, http.StatusForbidden, "Only team owner can remove members")
		return
	}

	// Don't allow removing the owner
	if memberID == userID {
		utils.SendError(w, http.StatusBadRequest, "Cannot remove team owner")
		return
	}

	if err := h.teamRepo.RemoveMember(teamID, memberID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to remove member")
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Member removed successfully"})
}

func (h *TeamHandler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	// Check if user is member or owner
	isMember, err := h.teamRepo.IsUserMember(teamID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to check team membership")
		return
	}

	isOwner, err := h.teamRepo.IsUserOwner(teamID, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to check team ownership")
		return
	}

	if !isMember && !isOwner {
		utils.SendError(w, http.StatusForbidden, "Access denied")
		return
	}

	members, err := h.teamRepo.GetTeamMembers(teamID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to get team members")
		return
	}

	utils.SendJSON(w, http.StatusOK, members)
}
