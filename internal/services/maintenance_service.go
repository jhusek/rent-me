package services

import (
	"database/sql"
	"time"

	"RentMe/internal/models"
)

type MaintenanceService struct {
	db *sql.DB
}

func NewMaintenanceService(db *sql.DB) *MaintenanceService {
	return &MaintenanceService{db: db}
}

func (s *MaintenanceService) ListMaintenance(status string) ([]models.MaintenanceRequest, error) {
	query := "SELECT id, title, description, priority, status, reported_date, completed_date, cost, vendor, expense_id, notes, created_at, updated_at FROM maintenance"
	var args []interface{}

	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY CASE priority WHEN 'emergency' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END, created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.MaintenanceRequest{}
	for rows.Next() {
		m, err := scanMaintenance(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (s *MaintenanceService) GetMaintenance(id int64) (models.MaintenanceRequest, error) {
	row := s.db.QueryRow("SELECT id, title, description, priority, status, reported_date, completed_date, cost, vendor, expense_id, notes, created_at, updated_at FROM maintenance WHERE id = ?", id)
	return scanMaintenanceRow(row)
}

func (s *MaintenanceService) CreateMaintenance(m models.MaintenanceRequest) (models.MaintenanceRequest, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if m.Priority == "" {
		m.Priority = "medium"
	}
	if m.Status == "" {
		m.Status = "open"
	}
	if m.ReportedDate == "" {
		m.ReportedDate = time.Now().Format("2006-01-02")
	}

	result, err := s.db.Exec(`INSERT INTO maintenance (title, description, priority, status, reported_date, completed_date, cost, vendor, expense_id, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.Title, m.Description, m.Priority, m.Status, m.ReportedDate, m.CompletedDate, m.Cost, m.Vendor, nullInt64(m.ExpenseID), m.Notes, now, now)
	if err != nil {
		return models.MaintenanceRequest{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.MaintenanceRequest{}, err
	}
	m.ID = id
	m.CreatedAt = now
	m.UpdatedAt = now
	return m, nil
}

func (s *MaintenanceService) UpdateMaintenance(m models.MaintenanceRequest) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.db.Exec(`UPDATE maintenance SET title=?, description=?, priority=?, status=?, reported_date=?, completed_date=?, cost=?, vendor=?, expense_id=?, notes=?, updated_at=? WHERE id=?`,
		m.Title, m.Description, m.Priority, m.Status, m.ReportedDate, m.CompletedDate, m.Cost, m.Vendor, nullInt64(m.ExpenseID), m.Notes, now, m.ID)
	return err
}

func (s *MaintenanceService) DeleteMaintenance(id int64) error {
	_, err := s.db.Exec("DELETE FROM maintenance WHERE id = ?", id)
	return err
}

func (s *MaintenanceService) CountOpen() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM maintenance WHERE status IN ('open', 'in_progress')").Scan(&count)
	return count, err
}

func nullInt64(v int64) interface{} {
	if v == 0 {
		return nil
	}
	return v
}

func scanMaintenance(rows *sql.Rows) (models.MaintenanceRequest, error) {
	var m models.MaintenanceRequest
	var description, priority, status, reportedDate, completedDate sql.NullString
	var vendor, notes, createdAt, updatedAt sql.NullString
	var cost sql.NullFloat64
	var expenseID sql.NullInt64

	err := rows.Scan(&m.ID, &m.Title, &description, &priority, &status, &reportedDate, &completedDate, &cost, &vendor, &expenseID, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.MaintenanceRequest{}, err
	}

	m.Description = description.String
	m.Priority = priority.String
	m.Status = status.String
	m.ReportedDate = reportedDate.String
	m.CompletedDate = completedDate.String
	m.Cost = cost.Float64
	m.Vendor = vendor.String
	m.ExpenseID = expenseID.Int64
	m.Notes = notes.String
	m.CreatedAt = createdAt.String
	m.UpdatedAt = updatedAt.String
	return m, nil
}

func scanMaintenanceRow(row *sql.Row) (models.MaintenanceRequest, error) {
	var m models.MaintenanceRequest
	var description, priority, status, reportedDate, completedDate sql.NullString
	var vendor, notes, createdAt, updatedAt sql.NullString
	var cost sql.NullFloat64
	var expenseID sql.NullInt64

	err := row.Scan(&m.ID, &m.Title, &description, &priority, &status, &reportedDate, &completedDate, &cost, &vendor, &expenseID, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.MaintenanceRequest{}, err
	}

	m.Description = description.String
	m.Priority = priority.String
	m.Status = status.String
	m.ReportedDate = reportedDate.String
	m.CompletedDate = completedDate.String
	m.Cost = cost.Float64
	m.Vendor = vendor.String
	m.ExpenseID = expenseID.Int64
	m.Notes = notes.String
	m.CreatedAt = createdAt.String
	m.UpdatedAt = updatedAt.String
	return m, nil
}
