package models

import (
	"time"
)

type Team struct {
	ID             int       `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Description    *string   `db:"description" json:"description"`
	OwnerID        *int      `db:"owner_id" json:"owner_id"`
	BudgetAmount   *float64  `db:"budget_amount" json:"budget_amount"`
	BudgetCurrency string    `db:"budget_currency" json:"budget_currency"`
	BudgetPeriod   *string   `db:"budget_period" json:"budget_period"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type TeamMember struct {
	ID        int       `db:"id" json:"id"`
	TeamID    int       `db:"team_id" json:"team_id"`
	UserID    int       `db:"user_id" json:"user_id"`
	JoinedAt  time.Time `db:"joined_at" json:"joined_at"`
}

type TeamWithMembers struct {
	Team
	Members []UserInfo `json:"members"`
	Owner   *UserInfo  `json:"owner"`
}

type UserInfo struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type CreateTeamRequest struct {
	Name           string   `json:"name"`
	Description    *string  `json:"description"`
	BudgetAmount   *float64 `json:"budget_amount"`
	BudgetCurrency string   `json:"budget_currency"`
	BudgetPeriod   *string  `json:"budget_period"`
}

type UpdateTeamRequest struct {
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	BudgetAmount   *float64 `json:"budget_amount"`
	BudgetCurrency *string  `json:"budget_currency"`
	BudgetPeriod   *string  `json:"budget_period"`
}

type AddTeamMemberRequest struct {
	UserID int `json:"user_id"`
}
