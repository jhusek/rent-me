package services

import (
	"database/sql"
	"time"

	"RentMe/internal/models"
)

type TenantService struct {
	db *sql.DB
}

func NewTenantService(db *sql.DB) *TenantService {
	return &TenantService{db: db}
}

func (s *TenantService) ListTenants(statusFilter string) ([]models.Tenant, error) {
	query := "SELECT id, first_name, last_name, email, phone, emergency_contact_name, emergency_contact_phone, lease_start, lease_end, monthly_rent, security_deposit, status, move_in_date, move_out_date, notes, created_at, updated_at FROM tenants"
	var args []interface{}

	if statusFilter != "" {
		query += " WHERE status = ?"
		args = append(args, statusFilter)
	}
	query += " ORDER BY last_name, first_name"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []models.Tenant{}
	for rows.Next() {
		t, err := scanTenant(rows)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, rows.Err()
}

func (s *TenantService) GetTenant(id int64) (models.Tenant, error) {
	row := s.db.QueryRow("SELECT id, first_name, last_name, email, phone, emergency_contact_name, emergency_contact_phone, lease_start, lease_end, monthly_rent, security_deposit, status, move_in_date, move_out_date, notes, created_at, updated_at FROM tenants WHERE id = ?", id)
	return scanTenantRow(row)
}

func (s *TenantService) CreateTenant(t models.Tenant) (models.Tenant, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if t.Status == "" {
		t.Status = "active"
	}

	result, err := s.db.Exec(`INSERT INTO tenants (first_name, last_name, email, phone, emergency_contact_name, emergency_contact_phone, lease_start, lease_end, monthly_rent, security_deposit, status, move_in_date, move_out_date, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.FirstName, t.LastName, t.Email, t.Phone, t.EmergencyContactName, t.EmergencyContactPhone, t.LeaseStart, t.LeaseEnd, t.MonthlyRent, t.SecurityDeposit, t.Status, t.MoveInDate, t.MoveOutDate, t.Notes, now, now)
	if err != nil {
		return models.Tenant{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Tenant{}, err
	}
	t.ID = id
	t.CreatedAt = now
	t.UpdatedAt = now
	return t, nil
}

func (s *TenantService) UpdateTenant(t models.Tenant) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.db.Exec(`UPDATE tenants SET first_name=?, last_name=?, email=?, phone=?, emergency_contact_name=?, emergency_contact_phone=?, lease_start=?, lease_end=?, monthly_rent=?, security_deposit=?, status=?, move_in_date=?, move_out_date=?, notes=?, updated_at=? WHERE id=?`,
		t.FirstName, t.LastName, t.Email, t.Phone, t.EmergencyContactName, t.EmergencyContactPhone, t.LeaseStart, t.LeaseEnd, t.MonthlyRent, t.SecurityDeposit, t.Status, t.MoveInDate, t.MoveOutDate, t.Notes, now, t.ID)
	return err
}

func (s *TenantService) DeleteTenant(id int64) error {
	_, err := s.db.Exec("DELETE FROM tenants WHERE id = ?", id)
	return err
}

func (s *TenantService) GetActiveTenant() (*models.Tenant, error) {
	row := s.db.QueryRow("SELECT id, first_name, last_name, email, phone, emergency_contact_name, emergency_contact_phone, lease_start, lease_end, monthly_rent, security_deposit, status, move_in_date, move_out_date, notes, created_at, updated_at FROM tenants WHERE status = 'active' LIMIT 1")
	t, err := scanTenantRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func scanTenant(rows *sql.Rows) (models.Tenant, error) {
	var t models.Tenant
	var email, phone, ecName, ecPhone, leaseStart, leaseEnd sql.NullString
	var moveIn, moveOut, notes, createdAt, updatedAt, status sql.NullString
	var monthlyRent, securityDeposit sql.NullFloat64

	err := rows.Scan(&t.ID, &t.FirstName, &t.LastName, &email, &phone, &ecName, &ecPhone, &leaseStart, &leaseEnd, &monthlyRent, &securityDeposit, &status, &moveIn, &moveOut, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.Tenant{}, err
	}

	t.Email = email.String
	t.Phone = phone.String
	t.EmergencyContactName = ecName.String
	t.EmergencyContactPhone = ecPhone.String
	t.LeaseStart = leaseStart.String
	t.LeaseEnd = leaseEnd.String
	t.MonthlyRent = monthlyRent.Float64
	t.SecurityDeposit = securityDeposit.Float64
	t.Status = status.String
	t.MoveInDate = moveIn.String
	t.MoveOutDate = moveOut.String
	t.Notes = notes.String
	t.CreatedAt = createdAt.String
	t.UpdatedAt = updatedAt.String
	return t, nil
}

func scanTenantRow(row *sql.Row) (models.Tenant, error) {
	var t models.Tenant
	var email, phone, ecName, ecPhone, leaseStart, leaseEnd sql.NullString
	var moveIn, moveOut, notes, createdAt, updatedAt, status sql.NullString
	var monthlyRent, securityDeposit sql.NullFloat64

	err := row.Scan(&t.ID, &t.FirstName, &t.LastName, &email, &phone, &ecName, &ecPhone, &leaseStart, &leaseEnd, &monthlyRent, &securityDeposit, &status, &moveIn, &moveOut, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.Tenant{}, err
	}

	t.Email = email.String
	t.Phone = phone.String
	t.EmergencyContactName = ecName.String
	t.EmergencyContactPhone = ecPhone.String
	t.LeaseStart = leaseStart.String
	t.LeaseEnd = leaseEnd.String
	t.MonthlyRent = monthlyRent.Float64
	t.SecurityDeposit = securityDeposit.Float64
	t.Status = status.String
	t.MoveInDate = moveIn.String
	t.MoveOutDate = moveOut.String
	t.Notes = notes.String
	t.CreatedAt = createdAt.String
	t.UpdatedAt = updatedAt.String
	return t, nil
}


