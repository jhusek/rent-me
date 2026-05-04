package services

import (
	"database/sql"
	"time"

	"RentMe/internal/models"
)

type ReminderService struct {
	db *sql.DB
}

func NewReminderService(db *sql.DB) *ReminderService {
	return &ReminderService{db: db}
}

func (s *ReminderService) ListReminders(showDismissed bool) ([]models.Reminder, error) {
	query := "SELECT id, title, description, due_date, type, is_recurring, recurring_interval, is_dismissed, related_entity_type, related_entity_id, created_at, updated_at FROM reminders"
	if !showDismissed {
		query += " WHERE is_dismissed = 0"
	}
	query += " ORDER BY due_date ASC"

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := []models.Reminder{}
	for rows.Next() {
		r, err := scanReminder(rows)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, r)
	}
	return reminders, rows.Err()
}

func (s *ReminderService) GetUpcomingReminders(days int) ([]models.Reminder, error) {
	query := "SELECT id, title, description, due_date, type, is_recurring, recurring_interval, is_dismissed, related_entity_type, related_entity_id, created_at, updated_at FROM reminders WHERE is_dismissed = 0 AND due_date <= date('now', '+' || ? || ' days') ORDER BY due_date ASC"

	rows, err := s.db.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := []models.Reminder{}
	for rows.Next() {
		r, err := scanReminder(rows)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, r)
	}
	return reminders, rows.Err()
}

func (s *ReminderService) CreateReminder(r models.Reminder) (models.Reminder, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if r.Type == "" {
		r.Type = "custom"
	}

	isRecurring := 0
	if r.IsRecurring {
		isRecurring = 1
	}
	isDismissed := 0
	if r.IsDismissed {
		isDismissed = 1
	}

	result, err := s.db.Exec(`INSERT INTO reminders (title, description, due_date, type, is_recurring, recurring_interval, is_dismissed, related_entity_type, related_entity_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.Title, r.Description, r.DueDate, r.Type, isRecurring, r.RecurringInterval, isDismissed, r.RelatedEntityType, r.RelatedEntityID, now, now)
	if err != nil {
		return models.Reminder{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Reminder{}, err
	}
	r.ID = id
	r.CreatedAt = now
	r.UpdatedAt = now
	return r, nil
}

func (s *ReminderService) UpdateReminder(r models.Reminder) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	isRecurring := 0
	if r.IsRecurring {
		isRecurring = 1
	}
	isDismissed := 0
	if r.IsDismissed {
		isDismissed = 1
	}

	_, err := s.db.Exec(`UPDATE reminders SET title=?, description=?, due_date=?, type=?, is_recurring=?, recurring_interval=?, is_dismissed=?, related_entity_type=?, related_entity_id=?, updated_at=? WHERE id=?`,
		r.Title, r.Description, r.DueDate, r.Type, isRecurring, r.RecurringInterval, isDismissed, r.RelatedEntityType, r.RelatedEntityID, now, r.ID)
	return err
}

func (s *ReminderService) DismissReminder(id int64) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.db.Exec("UPDATE reminders SET is_dismissed = 1, updated_at = ? WHERE id = ?", now, id)
	return err
}

func (s *ReminderService) DeleteReminder(id int64) error {
	_, err := s.db.Exec("DELETE FROM reminders WHERE id = ?", id)
	return err
}

func scanReminder(rows *sql.Rows) (models.Reminder, error) {
	var r models.Reminder
	var description, rType, recurringInterval, relatedEntityType sql.NullString
	var createdAt, updatedAt sql.NullString
	var isRecurring, isDismissed, relatedEntityID sql.NullInt64

	err := rows.Scan(&r.ID, &r.Title, &description, &r.DueDate, &rType, &isRecurring, &recurringInterval, &isDismissed, &relatedEntityType, &relatedEntityID, &createdAt, &updatedAt)
	if err != nil {
		return models.Reminder{}, err
	}

	r.Description = description.String
	r.Type = rType.String
	r.IsRecurring = isRecurring.Int64 == 1
	r.RecurringInterval = recurringInterval.String
	r.IsDismissed = isDismissed.Int64 == 1
	r.RelatedEntityType = relatedEntityType.String
	r.RelatedEntityID = relatedEntityID.Int64
	r.CreatedAt = createdAt.String
	r.UpdatedAt = updatedAt.String
	return r, nil
}
