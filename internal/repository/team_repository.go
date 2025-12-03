package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/pkg/database"
)

type TeamRepository struct {
	db *database.Database
}

func NewTeamRepository(db *database.Database) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *models.Team) error {
	query := `
		INSERT INTO teams (name, description, owner_id, budget_amount, budget_currency, budget_period)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		team.Name, team.Description, team.OwnerID, team.BudgetAmount, team.BudgetCurrency, team.BudgetPeriod,
	).Scan(&team.ID, &team.CreatedAt, &team.UpdatedAt)
}

func (r *TeamRepository) GetByID(id int) (*models.Team, error) {
	var team models.Team
	query := `SELECT * FROM teams WHERE id = $1`
	err := r.db.Get(&team, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("team not found")
		}
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) GetByOwnerID(ownerID int) ([]models.Team, error) {
	var teams []models.Team
	query := `SELECT * FROM teams WHERE owner_id = $1 ORDER BY created_at DESC`
	err := r.db.Select(&teams, query, ownerID)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *TeamRepository) GetTeamsByUserID(userID int) ([]models.Team, error) {
	var teams []models.Team
	query := `
		SELECT t.* FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY t.created_at DESC
	`
	err := r.db.Select(&teams, query, userID)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *TeamRepository) Update(team *models.Team) error {
	query := `
		UPDATE teams
		SET name = $1, description = $2, budget_amount = $3, budget_currency = $4, budget_period = $5, updated_at = $6
		WHERE id = $7
	`
	result, err := r.db.Exec(query, team.Name, team.Description, team.BudgetAmount, team.BudgetCurrency, team.BudgetPeriod, time.Now(), team.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("team not found")
	}
	return nil
}

func (r *TeamRepository) Delete(id int) error {
	query := `DELETE FROM teams WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("team not found")
	}
	return nil
}

func (r *TeamRepository) AddMember(teamID, userID int) error {
	query := `
		INSERT INTO team_members (team_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (team_id, user_id) DO NOTHING
	`
	_, err := r.db.Exec(query, teamID, userID)
	return err
}

func (r *TeamRepository) RemoveMember(teamID, userID int) error {
	query := `DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, teamID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("member not found in team")
	}
	return nil
}

func (r *TeamRepository) GetTeamMembers(teamID int) ([]models.UserInfo, error) {
	var members []models.UserInfo
	query := `
		SELECT u.id, u.email, u.username, u.role
		FROM users u
		INNER JOIN team_members tm ON u.id = tm.user_id
		WHERE tm.team_id = $1
		ORDER BY u.username
	`
	err := r.db.Select(&members, query, teamID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *TeamRepository) IsUserMember(teamID, userID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM team_members WHERE team_id = $1 AND user_id = $2`
	err := r.db.QueryRow(query, teamID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TeamRepository) IsUserOwner(teamID, userID int) (bool, error) {
	var ownerID *int
	query := `SELECT owner_id FROM teams WHERE id = $1`
	err := r.db.QueryRow(query, teamID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("team not found")
		}
		return false, err
	}
	if ownerID == nil {
		return false, nil
	}
	return *ownerID == userID, nil
}

func (r *TeamRepository) GetTeamWithMembers(teamID int) (*models.TeamWithMembers, error) {
	team, err := r.GetByID(teamID)
	if err != nil {
		return nil, err
	}

	members, err := r.GetTeamMembers(teamID)
	if err != nil {
		return nil, err
	}

	var owner *models.UserInfo
	if team.OwnerID != nil {
		query := `SELECT id, email, username, role FROM users WHERE id = $1`
		var ownerInfo models.UserInfo
		err := r.db.Get(&ownerInfo, query, *team.OwnerID)
		if err == nil {
			owner = &ownerInfo
		}
	}

	return &models.TeamWithMembers{
		Team:    *team,
		Members: members,
		Owner:   owner,
	}, nil
}
