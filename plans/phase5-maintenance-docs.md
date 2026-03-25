# Phase 5: Maintenance Requests & Document Storage

## Goal
Build the maintenance request tracking system and document storage module.

## Dependencies
- Phase 3 (tenants for maintenance requests)
- Phase 4 (expenses for cost linking)

## Deliverables
1. Maintenance request CRUD with status workflow
2. Cost tracking linked to expenses
3. Document upload and categorized storage
4. File viewer/opener via Wails runtime

---

## Step 1: Maintenance Requests (`/maintenance`)

### List View
- Table: Title, Priority (color badge), Status (badge), Reported Date, Vendor, Cost
- Filters: Status, Priority
- Sort by: Date, Priority, Status
- "New Request" button

### Add/Edit Form
- Title
- Description (textarea)
- Priority: Low (green), Medium (yellow), High (orange), Emergency (red)
- Status: Open → In Progress → Completed / Cancelled
- Reported date (default today)
- Vendor (dropdown from contacts + freetext)
- Cost (when completing)
- Link to expense (auto-create expense when cost is entered)
- Notes

### Status Workflow
1. **Open** — Request logged, no action yet
2. **In Progress** — Vendor assigned or work started
3. **Completed** — Work done, cost finalized
4. **Cancelled** — Not needed

### On Complete
- Prompt for cost
- Auto-create expense entry in "repairs" or "maintenance" category
- Link expense_id to maintenance record
- Log activity

## Step 2: Document Storage (`/documents`)

### List View
- Grid or table view toggle
- Cards: File name, category badge, upload date, size
- Filters: Category dropdown
- "Upload Document" button

### Upload Flow
1. Click "Upload Document"
2. Wails `runtime.OpenFileDialog()` opens native file picker
3. User selects file
4. App copies file to `documents/{category}/{timestamp}_{filename}`
5. Metadata saved to DB
6. Confirmation toast

### Categories
- Leases & Agreements
- Inspection Reports
- Receipts & Invoices
- Insurance Documents
- Tax Documents
- Photos
- Other

### File Actions
- View/Open: `runtime.BrowserOpenURL()` for common formats
- Delete: Remove file + DB record (with confirmation)
- Download/Export: Copy to chosen location

### Storage Location
- Default: `{app_data}/RentMe/documents/`
- Dev mode: `./documents/`
- Organized into category subdirectories

## Acceptance Criteria
- [ ] Maintenance requests can be created, edited, deleted
- [ ] Status workflow transitions correctly
- [ ] Completing a request auto-creates an expense
- [ ] Documents can be uploaded via native file picker
- [ ] Documents are organized by category
- [ ] Files can be opened/viewed
- [ ] Files can be deleted (file + DB record)
- [ ] Activity log captures maintenance and document actions
