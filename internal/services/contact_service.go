package services

import (
	"database/sql"
	"time"

	"RentMe/internal/models"
)

type ContactService struct {
	db *sql.DB
}

func NewContactService(db *sql.DB) *ContactService {
	return &ContactService{db: db}
}

func (s *ContactService) ListContacts() ([]models.Contact, error) {
	rows, err := s.db.Query("SELECT id, name, role, phone, email, company, notes, created_at, updated_at FROM contacts ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := []models.Contact{}
	for rows.Next() {
		c, err := scanContact(rows)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

func (s *ContactService) GetContact(id int64) (models.Contact, error) {
	row := s.db.QueryRow("SELECT id, name, role, phone, email, company, notes, created_at, updated_at FROM contacts WHERE id = ?", id)
	return scanContactRow(row)
}

func (s *ContactService) CreateContact(c models.Contact) (models.Contact, error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	result, err := s.db.Exec(`INSERT INTO contacts (name, role, phone, email, company, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Name, c.Role, c.Phone, c.Email, c.Company, c.Notes, now, now)
	if err != nil {
		return models.Contact{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Contact{}, err
	}
	c.ID = id
	c.CreatedAt = now
	c.UpdatedAt = now
	return c, nil
}

func (s *ContactService) UpdateContact(c models.Contact) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.db.Exec(`UPDATE contacts SET name=?, role=?, phone=?, email=?, company=?, notes=?, updated_at=? WHERE id=?`,
		c.Name, c.Role, c.Phone, c.Email, c.Company, c.Notes, now, c.ID)
	return err
}

func (s *ContactService) DeleteContact(id int64) error {
	_, err := s.db.Exec("DELETE FROM contacts WHERE id = ?", id)
	return err
}

func scanContact(rows *sql.Rows) (models.Contact, error) {
	var c models.Contact
	var role, phone, email, company, notes, createdAt, updatedAt sql.NullString

	err := rows.Scan(&c.ID, &c.Name, &role, &phone, &email, &company, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.Contact{}, err
	}

	c.Role = role.String
	c.Phone = phone.String
	c.Email = email.String
	c.Company = company.String
	c.Notes = notes.String
	c.CreatedAt = createdAt.String
	c.UpdatedAt = updatedAt.String
	return c, nil
}

func scanContactRow(row *sql.Row) (models.Contact, error) {
	var c models.Contact
	var role, phone, email, company, notes, createdAt, updatedAt sql.NullString

	err := row.Scan(&c.ID, &c.Name, &role, &phone, &email, &company, &notes, &createdAt, &updatedAt)
	if err != nil {
		return models.Contact{}, err
	}

	c.Role = role.String
	c.Phone = phone.String
	c.Email = email.String
	c.Company = company.String
	c.Notes = notes.String
	c.CreatedAt = createdAt.String
	c.UpdatedAt = updatedAt.String
	return c, nil
}
