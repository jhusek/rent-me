package services

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"RentMe/internal/models"
)

// DocumentContent holds the base64-encoded file bytes and metadata needed for
// the in-app viewer.
type DocumentContent struct {
	Base64Data string `json:"base64Data"`
	MimeType   string `json:"mimeType"`
	Name       string `json:"name"`
}

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

	// Read file content from disk and store it in the database
	var fileData []byte
	if d.FilePath != "" {
		var err error
		fileData, err = os.ReadFile(d.FilePath)
		if err != nil {
			return models.Document{}, fmt.Errorf("failed to read file: %w", err)
		}
		if d.FileSize == 0 {
			d.FileSize = len(fileData)
		}
		if d.MimeType == "" && len(fileData) > 0 {
			d.MimeType = http.DetectContentType(fileData)
		}
	}

	result, err := s.db.Exec(`INSERT INTO documents (name, category, file_path, file_size, mime_type, upload_date, notes, file_data, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.Name, d.Category, d.FilePath, d.FileSize, d.MimeType, d.UploadDate, d.Notes, fileData, now)
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

// OpenDocument fetches the stored file bytes from the database, writes them to
// a temporary file, and opens it with the system default application.
func (s *DocumentService) OpenDocument(id int64) error {
	row := s.db.QueryRow("SELECT file_data, name FROM documents WHERE id = ?", id)
	var fileData []byte
	var name sql.NullString
	if err := row.Scan(&fileData, &name); err != nil {
		return fmt.Errorf("document not found: %w", err)
	}
	if len(fileData) == 0 {
		return fmt.Errorf("no file content stored for this document")
	}

	ext := filepath.Ext(name.String)
	tmpFile, err := os.CreateTemp("", "rentme-doc-*"+ext)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	if _, err := tmpFile.Write(fileData); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	return openWithDefaultApp(tmpPath)
}

func openWithDefaultApp(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

// GetDocumentContent retrieves the stored file bytes from the database and
// returns them base64-encoded along with the detected MIME type and document
// name so the frontend can render an in-app preview.
func (s *DocumentService) GetDocumentContent(id int64) (DocumentContent, error) {
	row := s.db.QueryRow("SELECT file_data, mime_type, name FROM documents WHERE id = ?", id)
	var fileData []byte
	var mimeType, name sql.NullString
	if err := row.Scan(&fileData, &mimeType, &name); err != nil {
		return DocumentContent{}, fmt.Errorf("document not found: %w", err)
	}
	if len(fileData) == 0 {
		return DocumentContent{}, fmt.Errorf("no file content stored for this document")
	}
	mt := mimeType.String
	if mt == "" {
		mt = http.DetectContentType(fileData)
	}
	return DocumentContent{
		Base64Data: base64.StdEncoding.EncodeToString(fileData),
		MimeType:   mt,
		Name:       name.String,
	}, nil
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
