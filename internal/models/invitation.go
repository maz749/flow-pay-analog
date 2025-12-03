package models

import (
	"time"
)

type Invitation struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Role      string    `db:"role" json:"role"`
	TeamID    *int      `db:"team_id" json:"team_id"`
	Token     string    `db:"token" json:"token"`
	Status    string    `db:"status" json:"status"` // pending, accepted, expired
	InvitedBy int       `db:"invited_by" json:"invited_by"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type InvitationWithDetails struct {
	Invitation
	TeamName      *string `json:"team_name"`
	InviterEmail  string  `json:"inviter_email"`
	InviterName   string  `json:"inviter_name"`
}

type CreateInvitationRequest struct {
	Email  string `json:"email"`
	Role   string `json:"role"`
	TeamID *int   `json:"team_id"`
}

type AcceptInvitationRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
	Username string `json:"username"`
}
