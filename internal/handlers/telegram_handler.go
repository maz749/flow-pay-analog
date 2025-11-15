package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/maz749/flow-pay-analog/internal/middleware"
	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/pkg/utils"
)

type TelegramHandler struct {
	userRepo *repository.UserRepository
}

func NewTelegramHandler(userRepo *repository.UserRepository) *TelegramHandler {
	return &TelegramHandler{
		userRepo: userRepo,
	}
}

type TelegramSettingsRequest struct {
	ChatID    *int64 `json:"chat_id"`
	IsEnabled bool   `json:"is_enabled"`
}

func (h *TelegramHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	settings, err := h.userRepo.GetTelegramSettings(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch settings")
		return
	}

	if settings == nil {
		settings = &models.TelegramSettings{
			UserID:    userID,
			IsEnabled: false,
		}
	}

	utils.SendJSON(w, http.StatusOK, settings)
}

func (h *TelegramHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req TelegramSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	settings := &models.TelegramSettings{
		UserID:    userID,
		ChatID:    req.ChatID,
		IsEnabled: req.IsEnabled,
	}

	if err := h.userRepo.UpsertTelegramSettings(settings); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to update settings")
		return
	}

	utils.SendJSON(w, http.StatusOK, settings)
}
