package services

import (
	"database/sql"
	"time"

	"RentMe/internal/models"
)

type DocumentService struct {
	db *sql.DB
}

func NewDocumentService(db *sql.DB) *DocumentService {
	return &DocumentService{db: db}
}

func (s *DocumentService) ListDocuments(category string) ([]models.Document, error) {
	query := "SELECT id, name, category, file_path, file_size, mime_type, upload_date, notes, created_at FROM documents"
	var args []interface{}

	if category != "" {
		query += " WHERE category = ?"
		args = append(args, category)
	}
	query += " ORDER BY created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := []models.Document{}
	for rows.Next() {
		d, err := scanDocument(rows)
		if err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

func (s *DocumentService) GetDocument(id int64) (models.Document, error) {
	row := s.db.QueryRow("SELECT id, name, category, file_path, file_size, mime_type, upload_date, notes, created_at FROM documents WHERE id = ?", id)
	return scanDocumentRow(row)
}

func (s *DocumentService) CreateDocument(d models.Document) (models.Document, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if d.Category == "" {
		d.Category = "other"
	}
	if d.UploadDate == "" {
		d.UploadDate = time.Now().Format("2006-01-02")
	}

	result, err := s.db.Exec(`INSERT INTO documents (name, category, file_path, file_size, mime_type, upload_date, notes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.Name, d.Category, d.FilePath, d.FileSize, d.MimeType, d.UploadDate, d.Notes, now)
	if err != nil {
		return models.Document{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Document{}, err
	}
	d.ID = id
	d.CreatedAt = now
	return d, nil
}

func (s *DocumentService) DeleteDocument(id int64) error {
	_, err := s.db.Exec("DELETE FROM documents WHERE id = ?", id)
	return err
}

func scanDocument(rows *sql.Rows) (models.Document, error) {
	var d models.Document
	var category, mimeType, uploadDate, notes, createdAt sql.NullString
	var fileSize sql.NullInt64

	err := rows.Scan(&d.ID, &d.Name, &category, &d.FilePath, &fileSize, &mimeType, &uploadDate, &notes, &createdAt)
	if err != nil {
		return models.Document{}, err
	}

	d.Category = category.String
	d.FileSize = int(fileSize.Int64)
	d.MimeType = mimeType.String
	d.UploadDate = uploadDate.String
	d.Notes = notes.String
	d.CreatedAt = createdAt.String
	return d, nil
}

func scanDocumentRow(row *sql.Row) (models.Document, error) {
	var d models.Document
	var category, mimeType, uploadDate, notes, createdAt sql.NullString
	var fileSize sql.NullInt64

	err := row.Scan(&d.ID, &d.Name, &category, &d.FilePath, &fileSize, &mimeType, &uploadDate, &notes, &createdAt)
	if err != nil {
		return models.Document{}, err
	}

	d.Category = category.String
	d.FileSize = int(fileSize.Int64)
	d.MimeType = mimeType.String
	d.UploadDate = uploadDate.String
	d.Notes = notes.String
	d.CreatedAt = createdAt.String
	return d, nil
}
