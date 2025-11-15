package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/maz749/flow-pay-analog/internal/middleware"
	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/pkg/utils"
)

type SubscriptionHandler struct {
	subRepo    *repository.SubscriptionRepository
	notifRepo  *repository.NotificationRepository
}

func NewSubscriptionHandler(subRepo *repository.SubscriptionRepository, notifRepo *repository.NotificationRepository) *SubscriptionHandler {
	return &SubscriptionHandler{
		subRepo:   subRepo,
		notifRepo: notifRepo,
	}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" || req.Amount <= 0 || req.BillingPeriod == "" || req.StartDate == "" {
		utils.SendError(w, http.StatusBadRequest, "Name, amount, billing period, and start date are required")
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid start date format (expected YYYY-MM-DD)")
		return
	}

	// Validate billing period
	validPeriods := map[string]bool{"monthly": true, "yearly": true, "weekly": true, "custom": true}
	if !validPeriods[req.BillingPeriod] {
		utils.SendError(w, http.StatusBadRequest, "Invalid billing period")
		return
	}

	// Calculate next billing date
	nextBillingDate := utils.CalculateNextBillingDate(startDate, req.BillingPeriod, req.BillingDay, req.CustomPeriodDays)

	// Create subscription
	subscription := &models.Subscription{
		UserID:           userID,
		Name:             req.Name,
		Description:      req.Description,
		Amount:           req.Amount,
		Currency:         req.Currency,
		BillingPeriod:    req.BillingPeriod,
		BillingDay:       req.BillingDay,
		CustomPeriodDays: req.CustomPeriodDays,
		StartDate:        startDate,
		NextBillingDate:  nextBillingDate,
		IsActive:         true,
		Category:         req.Category,
		Icon:             req.Icon,
		Color:            req.Color,
	}

	if subscription.Currency == "" {
		subscription.Currency = "RUB"
	}

	if err := h.subRepo.Create(subscription); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create subscription")
		return
	}

	// Create notifications
	if len(req.NotifyDaysBefore) > 0 {
		for _, days := range req.NotifyDaysBefore {
			notif := &models.Notification{
				SubscriptionID:   subscription.ID,
				NotifyDaysBefore: days,
				IsEnabled:        true,
			}
			if err := h.notifRepo.Create(notif); err != nil {
				// Log error but don't fail the request
				continue
			}
		}
	}

	// Load notifications for response
	notifications, _ := h.notifRepo.GetBySubscription(subscription.ID)
	subscription.Notifications = notifications

	utils.SendJSON(w, http.StatusCreated, subscription)
}

func (h *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	subscriptions, err := h.subRepo.GetAll(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch subscriptions")
		return
	}

	// Load notifications for each subscription
	for i := range subscriptions {
		notifications, _ := h.notifRepo.GetBySubscription(subscriptions[i].ID)
		subscriptions[i].Notifications = notifications
	}

	utils.SendJSON(w, http.StatusOK, subscriptions)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	subscription, err := h.subRepo.GetByID(id, userID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Subscription not found")
		return
	}

	// Load notifications
	notifications, _ := h.notifRepo.GetBySubscription(subscription.ID)
	subscription.Notifications = notifications

	utils.SendJSON(w, http.StatusOK, subscription)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	// Get existing subscription
	subscription, err := h.subRepo.GetByID(id, userID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Subscription not found")
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		subscription.Name = *req.Name
	}
	if req.Description != nil {
		subscription.Description = req.Description
	}
	if req.Amount != nil {
		subscription.Amount = *req.Amount
	}
	if req.Currency != nil {
		subscription.Currency = *req.Currency
	}
	if req.BillingPeriod != nil {
		subscription.BillingPeriod = *req.BillingPeriod
		// Recalculate next billing date if period changed
		subscription.NextBillingDate = utils.CalculateNextBillingDate(
			subscription.StartDate,
			subscription.BillingPeriod,
			subscription.BillingDay,
			subscription.CustomPeriodDays,
		)
	}
	if req.BillingDay != nil {
		subscription.BillingDay = req.BillingDay
	}
	if req.CustomPeriodDays != nil {
		subscription.CustomPeriodDays = req.CustomPeriodDays
	}
	if req.IsActive != nil {
		subscription.IsActive = *req.IsActive
	}
	if req.Category != nil {
		subscription.Category = req.Category
	}
	if req.Icon != nil {
		subscription.Icon = req.Icon
	}
	if req.Color != nil {
		subscription.Color = req.Color
	}

	if err := h.subRepo.Update(subscription); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to update subscription")
		return
	}

	// Load notifications
	notifications, _ := h.notifRepo.GetBySubscription(subscription.ID)
	subscription.Notifications = notifications

	utils.SendJSON(w, http.StatusOK, subscription)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	if err := h.subRepo.Delete(id, userID); err != nil {
		utils.SendError(w, http.StatusNotFound, "Subscription not found")
		return
	}

	utils.SendSuccess(w, "Subscription deleted successfully")
}

func (h *SubscriptionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	stats, err := h.subRepo.GetStats(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch stats")
		return
	}

	utils.SendJSON(w, http.StatusOK, stats)
}
