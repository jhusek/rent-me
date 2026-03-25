package services

import (
	"database/sql"

	"RentMe/internal/models"
)

type ActivityService struct {
	db *sql.DB
}

func NewActivityService(db *sql.DB) *ActivityService {
	return &ActivityService{db: db}
}

func (s *ActivityService) LogActivity(action, entityType string, entityID int64, description string) error {
	_, err := s.db.Exec("INSERT INTO activity_log (action, entity_type, entity_id, description) VALUES (?, ?, ?, ?)",
		action, entityType, entityID, description)
	return err
}

func (s *ActivityService) GetRecentActivity(limit int) ([]models.ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := s.db.Query("SELECT id, action, entity_type, entity_id, description, created_at FROM activity_log ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []models.ActivityLog{}
	for rows.Next() {
		var a models.ActivityLog
		var description, createdAt sql.NullString
		var entityID sql.NullInt64

		err := rows.Scan(&a.ID, &a.Action, &a.EntityType, &entityID, &description, &createdAt)
		if err != nil {
			return nil, err
		}
		a.EntityID = entityID.Int64
		a.Description = description.String
		a.CreatedAt = createdAt.String
		activities = append(activities, a)
	}
	return activities, rows.Err()
}
