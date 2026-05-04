# Phase 3: Property Profile & Tenant Management

## Goal
Build the first two functional pages: Property Profile (view/edit the single 
rental property) and Tenants (full CRUD for tenant management with lease tracking).

## Dependencies
- Phase 1 (backend services must be working)
- Phase 2 (UI shell and components must be in place)

## Deliverables
1. Property page: view and edit all property details
2. Tenant list page with status indicators
3. Tenant add/edit forms
4. Tenant detail view with lease history
5. Contacts page (vendors, agents)

---

## Step 1: Property Page

### View Mode
- Property summary card: address, beds/baths/sqft, year built
- Financial card: purchase price, mortgage payment, insurance, taxes, HOA
- Notes section
- Edit button → switches to edit mode

### Edit Mode
- Form with all property fields
- Save/Cancel buttons
- Auto-creates property row on first save if none exists
- Logs activity on save

## Step 2: Tenant Management

### List View (`/tenants`)
- Table: Name, Status (badge), Lease Period, Monthly Rent, Phone
- Filter: Active / Past / All
- "Add Tenant" button → opens modal or form
- Click row → detail view

### Add/Edit Tenant Form
- Personal: First/Last name, email, phone, emergency contact
- Lease: Start date, end date, monthly rent, security deposit
- Status: Active / Past / Pending
- Move-in/out dates
- Notes

### Detail View (`/tenants/:id`)
- Contact info card
- Lease details card
- Payment history summary (links to Payments page with filter)
- Maintenance requests related to this tenant
- Status transition buttons (Activate, Move Out)

## Step 3: Contacts Page (`/contacts`)

- List of vendor/service contacts
- Add/edit with: Name, Role (dropdown), Phone, Email, Company, Notes
- Quick-call/email links (if applicable)
- Used later by maintenance module for vendor assignment

## Acceptance Criteria
- [ ] Property can be viewed and edited
- [ ] Tenants can be listed, filtered, created, edited, deleted
- [ ] Tenant detail view shows comprehensive info
- [ ] Contacts can be managed
- [ ] All changes are persisted to SQLite
- [ ] Activity log captures key actions
