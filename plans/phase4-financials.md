# Phase 4: Financial Tracking — Payments & Expenses

## Goal
Build the rent payment tracking and expense management modules. This is the 
financial backbone of the app.

## Dependencies
- Phase 3 (tenants must exist for payment tracking)

## Deliverables
1. Payment logging with status tracking
2. Payment history with filtering
3. Expense tracking with categories
4. Recurring expense support
5. Monthly/annual financial summaries

---

## Step 1: Rent Payments (`/payments`)

### List View
- Table: Due Date, Tenant, Amount, Status (badge), Paid Date, Method, Late Fee
- Filters: Month picker, status filter, tenant filter
- "Record Payment" button

### Record Payment Form
- Tenant (dropdown of active tenants)
- Amount (pre-filled from tenant's monthly_rent)
- Due date (default: 1st of current month)
- Paid date
- Payment method (Cash, Check, Bank Transfer, Venmo, Zelle, PayPal, Other)
- Late fee (auto-calculated if paid_date > due_date + grace period)
- Notes

### Auto-Generate Monthly Dues
- Backend method: `GenerateMonthlyPayments(month, year)` 
- Creates "due" records for all active tenants for the given month
- Idempotent (won't duplicate if already generated)
- Called from dashboard or payments page

### Payment Summary Card
- Total collected this month vs total expected
- Late payments count
- Year-to-date income

## Step 2: Expense Tracking (`/expenses`)

### List View
- Table: Date, Category (badge), Amount, Vendor, Description, Recurring?
- Filters: Category dropdown, date range, vendor search
- "Add Expense" button

### Add/Edit Expense Form
- Category (dropdown of IRS-aligned categories)
- Amount
- Date
- Vendor (text input with autocomplete from existing vendors)
- Description
- Is Recurring + Interval (monthly/quarterly/annually)
- Tax deductible checkbox
- Receipt file upload (via Wails file dialog)
- Notes

### Expense Categories (IRS Schedule E aligned)
- Mortgage Interest
- Insurance
- Property Tax
- Repairs
- Maintenance
- Utilities
- HOA Fees
- Management Fees
- Supplies
- Legal & Professional
- Advertising
- Travel
- Depreciation
- Other

### Expense Summary Card
- Total expenses this month
- Top 3 categories this month
- Year-to-date expenses
- Monthly average

## Step 3: Recurring Expense Automation

- Settings: Configure which expenses recur
- On app launch or monthly trigger: auto-create recurring expenses
- User can edit/adjust auto-created entries
- Visual indicator on recurring expenses

## Acceptance Criteria
- [ ] Payments can be recorded, edited, deleted
- [ ] Payment status is tracked (due → paid/partial/late/missed)
- [ ] Monthly payment generation works
- [ ] Late fees auto-calculate
- [ ] Expenses can be logged with all categories
- [ ] Recurring expenses auto-populate
- [ ] Receipt file attachment works via Wails file dialog
- [ ] Summary cards show accurate calculations
- [ ] Activity log captures all financial actions
