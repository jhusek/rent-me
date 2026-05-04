# Phase 7: Reminders & Tax Reports

## Goal
Build the reminder/notification system and tax reporting with CSV export.

## Dependencies
- Phase 4 (financial data for reports)
- Phase 6 (dashboard integration for reminder badges)

## Deliverables
1. Reminder CRUD with dismissal
2. Auto-generated reminders (lease expiry, etc.)
3. Notification badge in sidebar
4. Tax report (Schedule E format)
5. CSV export
6. Year-over-year comparison

---

## Step 1: Reminder System (`/reminders`)

### List View
- Tabs: Upcoming | All | Dismissed
- Table: Title, Type (badge), Due Date, Recurring?, Actions
- "Add Reminder" button

### Types
- **Lease Expiry** — Auto-generated 90/60/30 days before lease end
- **Rent Due** — Monthly, tied to tenant payment schedule  
- **Maintenance Follow-up** — When maintenance has been open > 7 days
- **Insurance Renewal** — Annual
- **Property Tax** — Based on tax due dates
- **Custom** — User-created

### Auto-Generation Logic
On app startup or manual refresh:
1. Check active tenant lease_end → create reminders at 90/60/30 days
2. Check open maintenance > 7 days old → create follow-up reminder
3. Skip if reminder already exists (dedup by type + entity + date)

### Actions
- Dismiss (hides from active list)
- Snooze (push due_date forward by 1/7/30 days)
- Delete

### Sidebar Badge
- Count of active (non-dismissed) reminders due within 7 days
- Red dot or number badge on 🔔 nav item
- Updates when reminders are dismissed/created

## Step 2: Tax Reports (`/reports`)

### Annual Income Summary
- Total rental income by year
- Total expenses by IRS category
- Net rental income (income - expenses)
- Effective tax rate helper (informational)

### Schedule E Helper View
Organized by IRS Schedule E line items:
- Line 3: Rents received
- Line 5: Advertising
- Line 6: Auto and travel
- Line 7: Cleaning and maintenance
- Line 8: Commissions
- Line 9: Insurance
- Line 10: Legal and professional fees
- Line 11: Management fees
- Line 12: Mortgage interest
- Line 13: Other interest
- Line 14: Repairs
- Line 16: Taxes
- Line 17: Utilities
- Line 19: Other
- Line 20: Total expenses
- Line 21: Net income/loss

### Depreciation Tracker
- Property cost basis (purchase price + improvements)
- Depreciation period: 27.5 years (residential)
- Annual depreciation amount = cost basis / 27.5
- Accumulated depreciation
- Remaining depreciable basis

### Year-over-Year Comparison
- Side-by-side table: current year vs previous year
- Income, expenses by category, net
- Percentage change indicators

## Step 3: CSV Export

- Export buttons on reports page
- Formats:
  - Annual income/expense summary
  - All payments for a year
  - All expenses for a year
  - Schedule E data
- Standard CSV format, opens in Excel

### Implementation
```go
func (a *App) ExportPaymentsCSV(year int) (string, error) {
    // Use runtime.SaveFileDialog for save location
    // Write CSV with encoding/csv
    // Return file path for confirmation
}
```

## Acceptance Criteria
- [ ] Reminders can be created, dismissed, snoozed, deleted
- [ ] Auto-reminders generate for lease expiry
- [ ] Sidebar badge shows upcoming reminder count
- [ ] Schedule E report displays correct line items
- [ ] Depreciation tracker calculates correctly
- [ ] Year-over-year comparison works
- [ ] CSV export produces valid files
- [ ] Reports handle edge cases (no data, partial year)
