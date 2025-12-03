package repository

import (
	"fmt"

	"github.com/maz749/flow-pay-analog/internal/models"
	"github.com/maz749/flow-pay-analog/pkg/database"
)

type CancellationRepository struct {
	db *database.Database
}

func NewCancellationRepository(db *database.Database) *CancellationRepository {
	return &CancellationRepository{db: db}
}

func (r *CancellationRepository) GetAll() ([]models.CancellationInstruction, error) {
	var instructions []models.CancellationInstruction
	query := `
		SELECT id, service_name, instructions, note, url, category
		FROM cancellation_instructions
		ORDER BY service_name ASC
	`
	err := r.db.Select(&instructions, query)
	if err != nil {
		return nil, err
	}
	return instructions, nil
}

func (r *CancellationRepository) GetByServiceName(serviceName string) (*models.CancellationInstruction, error) {
	var instruction models.CancellationInstruction
	query := `
		SELECT id, service_name, instructions, note, url, category
		FROM cancellation_instructions
		WHERE LOWER(service_name) = LOWER($1)
	`
	err := r.db.Get(&instruction, query, serviceName)
	if err != nil {
		return nil, fmt.Errorf("cancellation instruction not found")
	}
	return &instruction, nil
}

func (r *CancellationRepository) Search(searchTerm string) ([]models.CancellationInstruction, error) {
	var instructions []models.CancellationInstruction
	query := `
		SELECT id, service_name, instructions, note, url, category
		FROM cancellation_instructions
		WHERE LOWER(service_name) LIKE LOWER($1)
		ORDER BY service_name ASC
	`
	err := r.db.Select(&instructions, query, "%"+searchTerm+"%")
	if err != nil {
		return nil, err
	}
	return instructions, nil
}
