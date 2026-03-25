package services

import (
	"database/sql"
	"fmt"
	"time"

	"RentMe/internal/models"
)

type PaymentService struct {
	db *sql.DB
}

func NewPaymentService(db *sql.DB) *PaymentService {
	return &PaymentService{db: db}
}

func (s *PaymentService) ListPayments(tenantID int64, year int, month int) ([]models.Payment, error) {
	query := "SELECT id, tenant_id, amount, due_date, paid_date, method, status, late_fee, notes, created_at, updated_at FROM payments WHERE 1=1"
	var args []interface{}

	if tenantID != 0 {
		query += " AND tenant_id = ?"
		args = append(args, tenantID)
	}
	if year != 0 {
		query += " AND strftime('%Y', due_date) = ?"
		args = append(args, fmt.Sprintf("%04d", year))
	}
	if month != 0 {
		query += " AND CAST(strftime('%m', due_date) AS INTEGER) = ?"
		args = append(args, month)
	}
	query += " ORDER BY due_date DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []models.Payment{}
	for rows.Next() {
		p, err := scanPayment(rows)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, rows.Err()
}

func (s *PaymentService) GetPayment(id int64) (models.Payment, error) {
	row := s.db.QueryRow("SELECT id, tenant_id, amount, due_date, paid_date, method, status, late_fee, notes, created_at, updated_at FROM payments WHERE id = ?", id)
	return scanPaymentRow(row)
}

func (s *PaymentService) CreatePayment(p models.Payment) (models.Payment, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if p.Status == "" {
		p.Status = "due"
	}

	result, err := s.db.Exec(`INSERT INTO payments (tenant_id, amount, due_date, paid_date, method, status, late_fee, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.TenantID, p.Amount, p.DueDate, p.PaidDate, p.Method, p.Status, p.LateFee, p.Notes, now, now)
	if err != nil {
		return models.Payment{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Payment{}, err
	}
	p.ID = id
	p.CreatedAt = now
	p.UpdatedAt = now
	return p, nil
}

func (s *PaymentService) UpdatePayment(p models.Payment) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.db.Exec(`UPDATE payments SET tenant_id=?, amount=?, due_date=?, paid_date=?, method=?, status=?, late_fee=?, notes=?, updated_at=? WHERE id=?`,
		p.TenantID, p.Amount, p.DueDate, p.PaidDate, p.Method, p.Status, p.LateFee, p.Notes, now, p.ID)
	return err
}

func (s *PaymentService) DeletePayment(id int64) error {
	_, err := s.db.Exec("DELETE FROM payments WHERE id = ?", id)
	return err
}

func (s *PaymentService) GetPaymentSummary(year int, month int) (models.PaymentSummary, error) {
	query := "SELECT COALESCE(SUM(CASE WHEN status='due' THEN amount ELSE 0 END), 0), COALESCE(SUM(CASE WHEN status='paid' THEN amount ELSE 0 END), 0), COALESCE(SUM(CASE WHEN status='late' THEN amount ELSE 0 END), 0), COALESCE(SUM(CASE WHEN status='missed' THEN amount ELSE 0 END), 0), COALESCE(SUM(CASE WHEN status='partial' THEN amount ELSE 0 END), 0), COALESCE(SUM(late_fee), 0), COUNT(*) FROM payments WHERE 1=1"
	var args []interface{}

	if year != 0 {
		query += " AND strftime('%Y', due_date) = ?"
		args = append(args, fmt.Sprintf("%04d", year))
	}
	if month != 0 {
		query += " AND CAST(strftime('%m', due_date) AS INTEGER) = ?"
		args = append(args, month)
	}

	var summary models.PaymentSummary
	err := s.db.QueryRow(query, args...).Scan(&summary.TotalDue, &summary.TotalPaid, &summary.TotalLate, &summary.TotalMissed, &summary.TotalPartial, &summary.LateFees, &summary.Count)
	if err != nil {
		return models.PaymentSummary{}, err
	}
	return summary, nil
}

func scanPayment(rows *sql.Rows) (models.Payment, error) {
	var p models.Payment
	var paidDate, method, status, notes, createdAt, updatedAt sql.NullString
	var lateFee sql.NullFloat64

	err := rows.Scan(&p.ID, &p.TenantID, &p.Amount, &p.DueDate, &paidDate, &method, &status, &lateFee, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.Payment{}, err
	}

	p.PaidDate = paidDate.String
	p.Method = method.String
	p.Status = status.String
	p.LateFee = lateFee.Float64
	p.Notes = notes.String
	p.CreatedAt = createdAt.String
	p.UpdatedAt = updatedAt.String
	return p, nil
}

func scanPaymentRow(row *sql.Row) (models.Payment, error) {
	var p models.Payment
	var paidDate, method, status, notes, createdAt, updatedAt sql.NullString
	var lateFee sql.NullFloat64

	err := row.Scan(&p.ID, &p.TenantID, &p.Amount, &p.DueDate, &paidDate, &method, &status, &lateFee, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.Payment{}, err
	}

	p.PaidDate = paidDate.String
	p.Method = method.String
	p.Status = status.String
	p.LateFee = lateFee.Float64
	p.Notes = notes.String
	p.CreatedAt = createdAt.String
	p.UpdatedAt = updatedAt.String
	return p, nil
}
