package services

import (
	"database/sql"
	"fmt"
	"time"

	"RentMe/internal/models"
)

type ExpenseService struct {
	db *sql.DB
}

func NewExpenseService(db *sql.DB) *ExpenseService {
	return &ExpenseService{db: db}
}

func (s *ExpenseService) ListExpenses(category string, year int) ([]models.Expense, error) {
	query := "SELECT id, category, amount, date, vendor, description, is_recurring, recurring_interval, receipt_path, tax_deductible, notes, created_at, updated_at FROM expenses WHERE 1=1"
	var args []interface{}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}
	if year != 0 {
		query += " AND strftime('%Y', date) = ?"
		args = append(args, fmt.Sprintf("%04d", year))
	}
	query += " ORDER BY date DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []models.Expense{}
	for rows.Next() {
		e, err := scanExpense(rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, rows.Err()
}

func (s *ExpenseService) GetExpense(id int64) (models.Expense, error) {
	row := s.db.QueryRow("SELECT id, category, amount, date, vendor, description, is_recurring, recurring_interval, receipt_path, tax_deductible, notes, created_at, updated_at FROM expenses WHERE id = ?", id)
	return scanExpenseRow(row)
}

func (s *ExpenseService) CreateExpense(e models.Expense) (models.Expense, error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	isRecurring := 0
	if e.IsRecurring {
		isRecurring = 1
	}
	taxDeductible := 1
	if !e.TaxDeductible {
		taxDeductible = 0
	}

	result, err := s.db.Exec(`INSERT INTO expenses (category, amount, date, vendor, description, is_recurring, recurring_interval, receipt_path, tax_deductible, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.Category, e.Amount, e.Date, e.Vendor, e.Description, isRecurring, e.RecurringInterval, e.ReceiptPath, taxDeductible, e.Notes, now, now)
	if err != nil {
		return models.Expense{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Expense{}, err
	}
	e.ID = id
	e.CreatedAt = now
	e.UpdatedAt = now
	return e, nil
}

func (s *ExpenseService) UpdateExpense(e models.Expense) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	isRecurring := 0
	if e.IsRecurring {
		isRecurring = 1
	}
	taxDeductible := 1
	if !e.TaxDeductible {
		taxDeductible = 0
	}

	_, err := s.db.Exec(`UPDATE expenses SET category=?, amount=?, date=?, vendor=?, description=?, is_recurring=?, recurring_interval=?, receipt_path=?, tax_deductible=?, notes=?, updated_at=? WHERE id=?`,
		e.Category, e.Amount, e.Date, e.Vendor, e.Description, isRecurring, e.RecurringInterval, e.ReceiptPath, taxDeductible, e.Notes, now, e.ID)
	return err
}

func (s *ExpenseService) DeleteExpense(id int64) error {
	_, err := s.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}

func (s *ExpenseService) GetExpenseSummary(year int) (models.ExpenseSummary, error) {
	query := "SELECT category, COALESCE(SUM(amount), 0), COUNT(*), COALESCE(SUM(CASE WHEN tax_deductible = 1 THEN amount ELSE 0 END), 0) FROM expenses WHERE 1=1"
	var args []interface{}

	if year != 0 {
		query += " AND strftime('%Y', date) = ?"
		args = append(args, fmt.Sprintf("%04d", year))
	}
	query += " GROUP BY category"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return models.ExpenseSummary{}, err
	}
	defer rows.Close()

	summary := models.ExpenseSummary{
		ByCategory: make(map[string]float64),
	}
	for rows.Next() {
		var cat string
		var amount float64
		var count int
		var taxDed float64
		if err := rows.Scan(&cat, &amount, &count, &taxDed); err != nil {
			return models.ExpenseSummary{}, err
		}
		summary.ByCategory[cat] = amount
		summary.TotalAmount += amount
		summary.Count += count
		summary.TaxDeductible += taxDed
	}
	return summary, rows.Err()
}

func scanExpense(rows *sql.Rows) (models.Expense, error) {
	var e models.Expense
	var vendor, description, recurringInterval, receiptPath, notes, createdAt, updatedAt sql.NullString
	var isRecurring, taxDeductible sql.NullInt64

	err := rows.Scan(&e.ID, &e.Category, &e.Amount, &e.Date, &vendor, &description, &isRecurring, &recurringInterval, &receiptPath, &taxDeductible, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.Expense{}, err
	}

	e.Vendor = vendor.String
	e.Description = description.String
	e.IsRecurring = isRecurring.Int64 == 1
	e.RecurringInterval = recurringInterval.String
	e.ReceiptPath = receiptPath.String
	e.TaxDeductible = taxDeductible.Int64 == 1
	e.Notes = notes.String
	e.CreatedAt = createdAt.String
	e.UpdatedAt = updatedAt.String
	return e, nil
}

func scanExpenseRow(row *sql.Row) (models.Expense, error) {
	var e models.Expense
	var vendor, description, recurringInterval, receiptPath, notes, createdAt, updatedAt sql.NullString
	var isRecurring, taxDeductible sql.NullInt64

	err := row.Scan(&e.ID, &e.Category, &e.Amount, &e.Date, &vendor, &description, &isRecurring, &recurringInterval, &receiptPath, &taxDeductible, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.Expense{}, err
	}

	e.Vendor = vendor.String
	e.Description = description.String
	e.IsRecurring = isRecurring.Int64 == 1
	e.RecurringInterval = recurringInterval.String
	e.ReceiptPath = receiptPath.String
	e.TaxDeductible = taxDeductible.Int64 == 1
	e.Notes = notes.String
	e.CreatedAt = createdAt.String
	e.UpdatedAt = updatedAt.String
	return e, nil
}
