package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/maz749/flow-pay-analog/internal/middleware"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/pkg/utils"
)

type CancellationHandler struct {
	cancelRepo *repository.CancellationRepository
	subRepo    *repository.SubscriptionRepository
}

func NewCancellationHandler(cancelRepo *repository.CancellationRepository, subRepo *repository.SubscriptionRepository) *CancellationHandler {
	return &CancellationHandler{
		cancelRepo: cancelRepo,
		subRepo:    subRepo,
	}
}

// GetAllInstructions returns all cancellation instructions
func (h *CancellationHandler) GetAllInstructions(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("search")

	var instructions interface{}
	var err error

	if searchTerm != "" {
		instructions, err = h.cancelRepo.Search(searchTerm)
	} else {
		instructions, err = h.cancelRepo.GetAll()
	}

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch cancellation instructions")
		return
	}

	utils.SendJSON(w, http.StatusOK, instructions)
}

// GetInstructionByService returns cancellation instruction for a specific service
func (h *CancellationHandler) GetInstructionByService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["serviceName"]

	instruction, err := h.cancelRepo.GetByServiceName(serviceName)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Cancellation instruction not found for this service")
		return
	}

	utils.SendJSON(w, http.StatusOK, instruction)
}

// CancelSubscription marks a subscription as cancelled
func (h *CancellationHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	subID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	// Check if subscription exists and belongs to user
	sub, err := h.subRepo.GetByID(subID, userID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Subscription not found")
		return
	}

	// Cancel the subscription
	if err := h.subRepo.Cancel(subID, userID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to cancel subscription")
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Subscription cancelled successfully",
		"subscription": sub,
	})
}
