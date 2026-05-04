package db

import (
	"database/sql"
	"strings"

	_ "modernc.org/sqlite"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Enable WAL mode and foreign keys
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`ALTER TABLE documents ADD COLUMN file_data BLOB`,
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			if !strings.Contains(err.Error(), "duplicate column name") {
				return err
			}
		}
	}
	return nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS property (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			address TEXT,
			city TEXT,
			state TEXT,
			zip TEXT,
			purchase_price REAL,
			purchase_date TEXT,
			beds INTEGER,
			baths REAL,
			sqft INTEGER,
			year_built INTEGER,
			lot_size TEXT,
			property_type TEXT DEFAULT 'Single Family',
			mortgage_payment REAL,
			mortgage_rate REAL,
			mortgage_start TEXT,
			mortgage_term_years INTEGER DEFAULT 30,
			insurance_monthly REAL,
			property_tax_annual REAL,
			hoa_monthly REAL,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS tenants (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT,
			phone TEXT,
			emergency_contact_name TEXT,
			emergency_contact_phone TEXT,
			lease_start TEXT,
			lease_end TEXT,
			monthly_rent REAL,
			security_deposit REAL,
			status TEXT DEFAULT 'active' CHECK(status IN ('active','past','pending')),
			move_in_date TEXT,
			move_out_date TEXT,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS payments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL REFERENCES tenants(id),
			amount REAL NOT NULL,
			due_date TEXT NOT NULL,
			paid_date TEXT,
			method TEXT,
			status TEXT DEFAULT 'due' CHECK(status IN ('due','paid','partial','late','missed')),
			late_fee REAL DEFAULT 0,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS expenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT NOT NULL CHECK(category IN ('mortgage','insurance','property_tax','repairs','maintenance','utilities','hoa','management','supplies','legal','advertising','travel','depreciation','other')),
			amount REAL NOT NULL,
			date TEXT NOT NULL,
			vendor TEXT,
			description TEXT,
			is_recurring INTEGER DEFAULT 0,
			recurring_interval TEXT,
			receipt_path TEXT,
			tax_deductible INTEGER DEFAULT 1,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS maintenance (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			priority TEXT DEFAULT 'medium' CHECK(priority IN ('low','medium','high','emergency')),
			status TEXT DEFAULT 'open' CHECK(status IN ('open','in_progress','completed','cancelled')),
			reported_date TEXT,
			completed_date TEXT,
			cost REAL DEFAULT 0,
			vendor TEXT,
			expense_id INTEGER REFERENCES expenses(id),
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS documents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category TEXT DEFAULT 'other' CHECK(category IN ('lease','inspection','receipt','insurance','tax','photo','other')),
			file_path TEXT NOT NULL,
			file_size INTEGER DEFAULT 0,
			mime_type TEXT,
			upload_date TEXT,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS reminders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			due_date TEXT NOT NULL,
			type TEXT DEFAULT 'custom' CHECK(type IN ('lease_expiry','rent_due','maintenance','insurance','tax','custom')),
			is_recurring INTEGER DEFAULT 0,
			recurring_interval TEXT,
			is_dismissed INTEGER DEFAULT 0,
			related_entity_type TEXT,
			related_entity_id INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS contacts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			role TEXT,
			phone TEXT,
			email TEXT,
			company TEXT,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS activity_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			action TEXT NOT NULL,
			entity_type TEXT NOT NULL,
			entity_id INTEGER DEFAULT 0,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, t := range tables {
		if _, err := db.Exec(t); err != nil {
			return err
		}
	}

	return nil
}
