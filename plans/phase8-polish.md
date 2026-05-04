# Phase 8: Polish & Settings

## Goal
Final quality pass: global search, data validation, error handling, 
app settings, and overall UX polish.

## Dependencies
- All prior phases

## Deliverables
1. Global search across all entities
2. Comprehensive input validation
3. Error handling with user-friendly messages
4. Settings page with configurable options
5. Data backup/restore
6. Final UI polish

---

## Step 1: Global Search

- Search bar in the top of the sidebar or content header
- Searches across: tenants, payments, expenses, maintenance, documents, contacts
- Results grouped by entity type
- Click result → navigate to that entity

### Implementation
- Backend: `Search(query string)` method
- SQLite `LIKE` queries across key text fields
- Returns `[]SearchResult{Type, ID, Title, Subtitle, URL}`

## Step 2: Input Validation

- All forms: required field indicators (*)
- Email format validation
- Phone number formatting
- Date range validation (lease end > lease start)
- Amount > 0 for payments/expenses
- Prevent duplicate tenant creation
- Frontend inline validation + backend validation

## Step 3: Error Handling

- Wails method errors → toast notifications
- Network/DB errors → friendly error messages
- Loading states on all data fetches
- Empty states with helpful CTAs
- Confirmation dialogs on destructive actions

## Step 4: Settings Page (`/settings`)

### Late Fee Configuration
- Grace period (days after due date)
- Late fee type: Flat amount or Percentage
- Late fee amount/percentage

### Display Preferences
- Currency format (USD default)
- Date format (MM/DD/YYYY default)

### Data Management
- Database location display
- Documents folder location
- Backup database (copy DB file to chosen location)
- Restore from backup
- Export all data as JSON

### About
- App version
- License info
- GitHub link (if applicable)

## Step 5: Final UI Polish

- Consistent spacing and alignment across all pages
- Hover effects on interactive elements
- Focus states for accessibility
- Loading spinners during data operations
- Smooth transitions between pages
- Keyboard shortcuts (Ctrl+N for new, Esc for close modal)
- Favicon and app icon

## Step 6: Testing & QA

- Test all CRUD operations across entities
- Test edge cases: empty DB, missing data, long strings
- Test file operations: upload, open, delete
- Test charts with various data amounts
- Verify CSV exports open correctly in Excel
- Window resize behavior
- First-run experience (empty state)

## Acceptance Criteria
- [ ] Global search returns relevant results across all entities
- [ ] All forms validate input before submission
- [ ] Errors show user-friendly messages
- [ ] Settings page allows configuration of key options
- [ ] Database backup/restore works
- [ ] No visual glitches or broken layouts
- [ ] First-run experience is welcoming (empty states guide the user)
- [ ] App feels professional and intuitive
