package services

import (
	"database/sql"
	"fmt"
	"time"

	"RentMe/internal/models"
)

type MonthAmount struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

type CategoryAmount struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

type RentStatusSummary struct {
	Paid    int `json:"paid"`
	Due     int `json:"due"`
	Late    int `json:"late"`
	Partial int `json:"partial"`
	Missed  int `json:"missed"`
}

type DashboardData struct {
	Property           models.Property    `json:"property"`
	CurrentTenant      *models.Tenant     `json:"currentTenant"`
	MonthlyIncome      float64            `json:"monthlyIncome"`
	MonthlyExpenses    float64            `json:"monthlyExpenses"`
	YTDIncome          float64            `json:"ytdIncome"`
	YTDExpenses        float64            `json:"ytdExpenses"`
	NetIncome          float64            `json:"netIncome"`
	RentStatus         RentStatusSummary  `json:"rentStatus"`
	IncomeByMonth      []MonthAmount      `json:"incomeByMonth"`
	ExpensesByMonth    []MonthAmount      `json:"expensesByMonth"`
	ExpensesByCategory []CategoryAmount   `json:"expensesByCategory"`
	UpcomingReminders  []models.Reminder  `json:"upcomingReminders"`
	RecentActivity     []models.ActivityLog `json:"recentActivity"`
	OpenMaintenance    int                `json:"openMaintenance"`
}

type DashboardService struct {
	db                 *sql.DB
	paymentService     *PaymentService
	expenseService     *ExpenseService
	reminderService    *ReminderService
	maintenanceService *MaintenanceService
	propertyService    *PropertyService
	tenantService      *TenantService
	activityService    *ActivityService
}

func NewDashboardService(db *sql.DB, ps *PaymentService, es *ExpenseService, rs *ReminderService, ms *MaintenanceService, props *PropertyService, ts *TenantService, as *ActivityService) *DashboardService {
	return &DashboardService{
		db:                 db,
		paymentService:     ps,
		expenseService:     es,
		reminderService:    rs,
		maintenanceService: ms,
		propertyService:    props,
		tenantService:      ts,
		activityService:    as,
	}
}

func (s *DashboardService) GetDashboardData() (DashboardData, error) {
	var data DashboardData
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	property, err := s.propertyService.GetProperty()
	if err != nil {
		return data, err
	}
	data.Property = property

	tenant, err := s.tenantService.GetActiveTenant()
	if err != nil {
		return data, err
	}
	data.CurrentTenant = tenant

	// Monthly income (paid payments this month)
	var monthlyIncome sql.NullFloat64
	err = s.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM payments WHERE status = 'paid' AND strftime('%Y', paid_date) = ? AND CAST(strftime('%m', paid_date) AS INTEGER) = ?",
		fmt.Sprintf("%04d", currentYear), currentMonth).Scan(&monthlyIncome)
	if err != nil {
		return data, err
	}
	data.MonthlyIncome = monthlyIncome.Float64

	// Monthly expenses
	var monthlyExpenses sql.NullFloat64
	err = s.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE strftime('%Y', date) = ? AND CAST(strftime('%m', date) AS INTEGER) = ?",
		fmt.Sprintf("%04d", currentYear), currentMonth).Scan(&monthlyExpenses)
	if err != nil {
		return data, err
	}
	data.MonthlyExpenses = monthlyExpenses.Float64

	// YTD income
	var ytdIncome sql.NullFloat64
	err = s.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM payments WHERE status = 'paid' AND strftime('%Y', paid_date) = ?",
		fmt.Sprintf("%04d", currentYear)).Scan(&ytdIncome)
	if err != nil {
		return data, err
	}
	data.YTDIncome = ytdIncome.Float64

	// YTD expenses
	var ytdExpenses sql.NullFloat64
	err = s.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE strftime('%Y', date) = ?",
		fmt.Sprintf("%04d", currentYear)).Scan(&ytdExpenses)
	if err != nil {
		return data, err
	}
	data.YTDExpenses = ytdExpenses.Float64
	data.NetIncome = data.YTDIncome - data.YTDExpenses

	// Rent status for current month
	err = s.db.QueryRow(`SELECT 
		COALESCE(SUM(CASE WHEN status='paid' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status='due' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status='late' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status='partial' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status='missed' THEN 1 ELSE 0 END), 0)
		FROM payments WHERE strftime('%Y', due_date) = ? AND CAST(strftime('%m', due_date) AS INTEGER) = ?`,
		fmt.Sprintf("%04d", currentYear), currentMonth).Scan(
		&data.RentStatus.Paid, &data.RentStatus.Due, &data.RentStatus.Late,
		&data.RentStatus.Partial, &data.RentStatus.Missed)
	if err != nil {
		return data, err
	}

	// Income by month (last 12 months)
	data.IncomeByMonth, err = s.getIncomeByMonth(now)
	if err != nil {
		return data, err
	}

	// Expenses by month (last 12 months)
	data.ExpensesByMonth, err = s.getExpensesByMonth(now)
	if err != nil {
		return data, err
	}

	// Expenses by category (current year)
	data.ExpensesByCategory, err = s.getExpensesByCategory(currentYear)
	if err != nil {
		return data, err
	}

	// Upcoming reminders (next 30 days)
	reminders, err := s.reminderService.GetUpcomingReminders(30)
	if err != nil {
		return data, err
	}
	if reminders == nil {
		reminders = []models.Reminder{}
	}
	data.UpcomingReminders = reminders

	// Recent activity
	activity, err := s.activityService.GetRecentActivity(10)
	if err != nil {
		return data, err
	}
	if activity == nil {
		activity = []models.ActivityLog{}
	}
	data.RecentActivity = activity

	// Open maintenance count
	openCount, err := s.maintenanceService.CountOpen()
	if err != nil {
		return data, err
	}
	data.OpenMaintenance = openCount

	return data, nil
}

func (s *DashboardService) getIncomeByMonth(now time.Time) ([]MonthAmount, error) {
	// Last 12 months including current
	startDate := now.AddDate(0, -11, 0)
	startStr := fmt.Sprintf("%04d-%02d-01", startDate.Year(), startDate.Month())

	rows, err := s.db.Query(`SELECT strftime('%Y-%m', paid_date) as month, COALESCE(SUM(amount), 0) as total 
		FROM payments WHERE status = 'paid' AND paid_date >= ? 
		GROUP BY strftime('%Y-%m', paid_date) ORDER BY month`, startStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monthMap := make(map[string]float64)
	for rows.Next() {
		var month string
		var amount float64
		if err := rows.Scan(&month, &amount); err != nil {
			return nil, err
		}
		monthMap[month] = amount
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := []MonthAmount{}
	for i := 11; i >= 0; i-- {
		d := now.AddDate(0, -i, 0)
		key := fmt.Sprintf("%04d-%02d", d.Year(), d.Month())
		result = append(result, MonthAmount{
			Month:  key,
			Amount: monthMap[key],
		})
	}
	return result, nil
}

func (s *DashboardService) getExpensesByMonth(now time.Time) ([]MonthAmount, error) {
	startDate := now.AddDate(0, -11, 0)
	startStr := fmt.Sprintf("%04d-%02d-01", startDate.Year(), startDate.Month())

	rows, err := s.db.Query(`SELECT strftime('%Y-%m', date) as month, COALESCE(SUM(amount), 0) as total 
		FROM expenses WHERE date >= ? 
		GROUP BY strftime('%Y-%m', date) ORDER BY month`, startStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monthMap := make(map[string]float64)
	for rows.Next() {
		var month string
		var amount float64
		if err := rows.Scan(&month, &amount); err != nil {
			return nil, err
		}
		monthMap[month] = amount
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := []MonthAmount{}
	for i := 11; i >= 0; i-- {
		d := now.AddDate(0, -i, 0)
		key := fmt.Sprintf("%04d-%02d", d.Year(), d.Month())
		result = append(result, MonthAmount{
			Month:  key,
			Amount: monthMap[key],
		})
	}
	return result, nil
}

func (s *DashboardService) getExpensesByCategory(year int) ([]CategoryAmount, error) {
	rows, err := s.db.Query(`SELECT category, COALESCE(SUM(amount), 0) as total 
		FROM expenses WHERE strftime('%Y', date) = ? 
		GROUP BY category ORDER BY total DESC`,
		fmt.Sprintf("%04d", year))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []CategoryAmount{}
	for rows.Next() {
		var ca CategoryAmount
		if err := rows.Scan(&ca.Category, &ca.Amount); err != nil {
			return nil, err
		}
		result = append(result, ca)
	}
	return result, rows.Err()
}
