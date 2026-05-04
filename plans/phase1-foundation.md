# Phase 1: Foundation — Database, Models & Services

## Goal
Stand up the entire Go backend: SQLite database with migrations, all entity models, 
CRUD services, and Wails-bound handler methods. After this phase, the frontend can 
call any backend method.

## Deliverables
1. SQLite database initialized at startup with all tables
2. Go structs for every entity
3. Service layer with full CRUD for each entity
4. Wails-bound App methods exposing all services to the frontend
5. Verified with `wails build` or `wails dev`

---

## Step 1: Add SQLite Dependency

```bash
cd RentMe
go get modernc.org/sqlite
go get github.com/jmoiron/sqlx  # optional but helpful for named queries
```

If sqlx causes issues, use `database/sql` directly with `modernc.org/sqlite` driver.

## Step 2: Database Initialization (`internal/db/db.go`)

- Open/create SQLite file at configurable path
- Run migrations (embedded SQL or Go code)
- Return `*sql.DB` handle
- Create all tables in a single migration

### Tables

```sql
-- Property (single property, one row)
CREATE TABLE IF NOT EXISTS property (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    city TEXT NOT NULL DEFAULT '',
    state TEXT NOT NULL DEFAULT '',
    zip TEXT NOT NULL DEFAULT '',
    purchase_price REAL DEFAULT 0,
    purchase_date TEXT DEFAULT '',
    beds INTEGER DEFAULT 0,
    baths REAL DEFAULT 0,
    sqft INTEGER DEFAULT 0,
    year_built INTEGER DEFAULT 0,
    lot_size TEXT DEFAULT '',
    property_type TEXT DEFAULT 'Single Family',
    mortgage_payment REAL DEFAULT 0,
    mortgage_rate REAL DEFAULT 0,
    mortgage_start TEXT DEFAULT '',
    mortgage_term_years INTEGER DEFAULT 30,
    insurance_monthly REAL DEFAULT 0,
    property_tax_annual REAL DEFAULT 0,
    hoa_monthly REAL DEFAULT 0,
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tenants
CREATE TABLE IF NOT EXISTS tenants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT DEFAULT '',
    phone TEXT DEFAULT '',
    emergency_contact_name TEXT DEFAULT '',
    emergency_contact_phone TEXT DEFAULT '',
    lease_start TEXT DEFAULT '',
    lease_end TEXT DEFAULT '',
    monthly_rent REAL DEFAULT 0,
    security_deposit REAL DEFAULT 0,
    status TEXT DEFAULT 'active' CHECK(status IN ('active','past','pending')),
    move_in_date TEXT DEFAULT '',
    move_out_date TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Payments
CREATE TABLE IF NOT EXISTS payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id),
    amount REAL NOT NULL,
    due_date TEXT NOT NULL,
    paid_date TEXT DEFAULT '',
    method TEXT DEFAULT '' CHECK(method IN ('','cash','check','bank_transfer','venmo','zelle','paypal','other')),
    status TEXT DEFAULT 'due' CHECK(status IN ('due','paid','partial','late','missed')),
    late_fee REAL DEFAULT 0,
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Expenses
CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category TEXT NOT NULL CHECK(category IN (
        'mortgage','insurance','property_tax','repairs','maintenance',
        'utilities','hoa','management','supplies','legal',
        'advertising','travel','depreciation','other'
    )),
    amount REAL NOT NULL,
    date TEXT NOT NULL,
    vendor TEXT DEFAULT '',
    description TEXT DEFAULT '',
    is_recurring INTEGER DEFAULT 0,
    recurring_interval TEXT DEFAULT '' CHECK(recurring_interval IN ('','monthly','quarterly','annually')),
    receipt_path TEXT DEFAULT '',
    tax_deductible INTEGER DEFAULT 1,
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Maintenance Requests
CREATE TABLE IF NOT EXISTS maintenance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    priority TEXT DEFAULT 'medium' CHECK(priority IN ('low','medium','high','emergency')),
    status TEXT DEFAULT 'open' CHECK(status IN ('open','in_progress','completed','cancelled')),
    reported_date TEXT DEFAULT '',
    completed_date TEXT DEFAULT '',
    cost REAL DEFAULT 0,
    vendor TEXT DEFAULT '',
    expense_id INTEGER DEFAULT NULL REFERENCES expenses(id),
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Documents
CREATE TABLE IF NOT EXISTS documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    category TEXT DEFAULT 'other' CHECK(category IN (
        'lease','inspection','receipt','insurance','tax','photo','other'
    )),
    file_path TEXT NOT NULL,
    file_size INTEGER DEFAULT 0,
    mime_type TEXT DEFAULT '',
    upload_date TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Reminders
CREATE TABLE IF NOT EXISTS reminders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    due_date TEXT NOT NULL,
    type TEXT DEFAULT 'custom' CHECK(type IN (
        'lease_expiry','rent_due','maintenance','insurance','tax','custom'
    )),
    is_recurring INTEGER DEFAULT 0,
    recurring_interval TEXT DEFAULT '' CHECK(recurring_interval IN ('','monthly','quarterly','annually')),
    is_dismissed INTEGER DEFAULT 0,
    related_entity_type TEXT DEFAULT '',
    related_entity_id INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Contacts (vendors, agents, etc.)
CREATE TABLE IF NOT EXISTS contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    role TEXT DEFAULT '' CHECK(role IN (
        '','plumber','electrician','handyman','contractor',
        'insurance_agent','property_manager','realtor',
        'attorney','accountant','hoa','other'
    )),
    phone TEXT DEFAULT '',
    email TEXT DEFAULT '',
    company TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Activity Log
CREATE TABLE IF NOT EXISTS activity_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id INTEGER DEFAULT 0,
    description TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Step 3: Go Models (`internal/models/`)

One file per entity or a single `models.go`:
- `Property`, `Tenant`, `Payment`, `Expense`, `Maintenance`, `Document`, `Reminder`, `Contact`, `ActivityLog`
- All fields use Go types + JSON tags for Wails serialization
- Use `string` for dates (SQLite stores text), parse on frontend

## Step 4: Service Layer (`internal/services/`)

Each service gets its own file:
- `property_service.go` — Get, Update (single property, no list needed)
- `tenant_service.go` — List, Get, Create, Update, Delete
- `payment_service.go` — List (by tenant/date range), Get, Create, Update, Delete, GetSummary
- `expense_service.go` — List (by category/date range), Get, Create, Update, Delete, GetSummary
- `maintenance_service.go` — List, Get, Create, Update, Delete
- `document_service.go` — List, Get, Create, Delete (no update — re-upload)
- `reminder_service.go` — List, GetUpcoming, Create, Update, Dismiss, Delete
- `contact_service.go` — List, Get, Create, Update, Delete
- `activity_service.go` — Log, GetRecent
- `dashboard_service.go` — GetDashboardData (aggregates from all above)

## Step 5: Wails Bindings (`app.go`)

The `App` struct holds `*sql.DB` and all services. Each public method on `App` 
is auto-exposed to the frontend. Group methods with clear naming:

```go
// Tenants
func (a *App) GetTenants() ([]models.Tenant, error)
func (a *App) GetTenant(id int64) (models.Tenant, error)
func (a *App) CreateTenant(t models.Tenant) (models.Tenant, error)
func (a *App) UpdateTenant(t models.Tenant) error
func (a *App) DeleteTenant(id int64) error
// ... etc for all entities
```

## Step 6: Verify Build

```bash
wails build
```

Ensure no compile errors. Frontend doesn't need to use the bindings yet — 
just confirm the Go side compiles and the app launches.

## Acceptance Criteria
- [ ] `wails dev` launches without errors
- [ ] SQLite DB file is created on first run
- [ ] All tables exist in the DB
- [ ] Wails generates TypeScript bindings in `frontend/wailsjs/go/`
- [ ] No compile errors
