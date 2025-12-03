package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/pkg/database"
)

type InvitationRepository struct {
	db *database.Database
}

func NewInvitationRepository(db *database.Database) *InvitationRepository {
	return &InvitationRepository{db: db}
}

func (r *InvitationRepository) Create(invitation *models.Invitation) error {
	query := `
		INSERT INTO invitations (email, role, team_id, token, status, invited_by, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		invitation.Email, invitation.Role, invitation.TeamID, invitation.Token,
		invitation.Status, invitation.InvitedBy, invitation.ExpiresAt,
	).Scan(&invitation.ID, &invitation.CreatedAt)
}

func (r *InvitationRepository) GetByToken(token string) (*models.Invitation, error) {
	var invitation models.Invitation
	query := `SELECT * FROM invitations WHERE token = $1`
	err := r.db.Get(&invitation, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, err
	}
	return &invitation, nil
}

func (r *InvitationRepository) GetByID(id int) (*models.Invitation, error) {
	var invitation models.Invitation
	query := `SELECT * FROM invitations WHERE id = $1`
	err := r.db.Get(&invitation, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, err
	}
	return &invitation, nil
}

func (r *InvitationRepository) GetByInviter(inviterID int) ([]models.InvitationWithDetails, error) {
	var invitations []models.InvitationWithDetails
	query := `
		SELECT
			i.*,
			t.name as team_name,
			u.email as inviter_email,
			u.username as inviter_name
		FROM invitations i
		LEFT JOIN teams t ON i.team_id = t.id
		INNER JOIN users u ON i.invited_by = u.id
		WHERE i.invited_by = $1
		ORDER BY i.created_at DESC
	`
	err := r.db.Select(&invitations, query, inviterID)
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (r *InvitationRepository) GetPendingByEmail(email string) ([]models.InvitationWithDetails, error) {
	var invitations []models.InvitationWithDetails
	query := `
		SELECT
			i.*,
			t.name as team_name,
			u.email as inviter_email,
			u.username as inviter_name
		FROM invitations i
		LEFT JOIN teams t ON i.team_id = t.id
		INNER JOIN users u ON i.invited_by = u.id
		WHERE i.email = $1 AND i.status = 'pending' AND i.expires_at > NOW()
		ORDER BY i.created_at DESC
	`
	err := r.db.Select(&invitations, query, email)
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (r *InvitationRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE invitations SET status = $1 WHERE id = $2`
	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}
	return nil
}

func (r *InvitationRepository) Delete(id int) error {
	query := `DELETE FROM invitations WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}
	return nil
}

func (r *InvitationRepository) ExpireOldInvitations() error {
	query := `
		UPDATE invitations
		SET status = 'expired'
		WHERE status = 'pending' AND expires_at < NOW()
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *InvitationRepository) GetInvitationWithDetails(token string) (*models.InvitationWithDetails, error) {
	var invitation models.InvitationWithDetails
	query := `
		SELECT
			i.*,
			t.name as team_name,
			u.email as inviter_email,
			u.username as inviter_name
		FROM invitations i
		LEFT JOIN teams t ON i.team_id = t.id
		INNER JOIN users u ON i.invited_by = u.id
		WHERE i.token = $1
	`
	err := r.db.Get(&invitation, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, err
	}
	return &invitation, nil
}
