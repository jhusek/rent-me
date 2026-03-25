# Phase 6: Dashboard & Data Visualization

## Goal
Build the home dashboard that provides an at-a-glance overview of everything:
income, expenses, rent status, upcoming reminders, and recent activity.

## Dependencies
- Phase 4 (financial data for charts)
- Phase 5 (maintenance data for activity feed)

## Deliverables
1. Dashboard layout with summary cards
2. Income vs Expenses bar chart
3. Rent collection status for current month
4. Upcoming reminders/alerts
5. Recent activity feed
6. Property summary card

---

## Step 1: Dashboard Backend (`dashboard_service.go`)

### `GetDashboardData()` returns:
```go
type DashboardData struct {
    Property          Property
    CurrentTenant     *Tenant
    
    // Financial summary
    MonthlyIncome     float64   // rent collected this month
    MonthlyExpenses   float64   // expenses this month
    YTDIncome         float64
    YTDExpenses       float64
    NetIncome         float64   // YTD income - expenses
    
    // Rent status
    RentStatus        RentStatus // paid/due/late/partial counts
    
    // Charts data (last 12 months)
    MonthlyIncomeData   []MonthAmount
    MonthlyExpenseData  []MonthAmount
    ExpenseByCategory   []CategoryAmount
    
    // Activity & reminders
    UpcomingReminders []Reminder
    RecentActivity    []ActivityLog
    
    // Maintenance
    OpenMaintenance   int
}
```

## Step 2: Dashboard Layout (`/` route)

### Row 1: Summary Cards (4-column grid)
- **Monthly Income** — Green, dollar amount, vs last month trend arrow
- **Monthly Expenses** — Red, dollar amount, vs last month trend arrow
- **Net Income (YTD)** — Blue, dollar amount
- **Rent Status** — Badge: "Paid ✓" / "Due" / "Late ⚠"

### Row 2: Charts (2-column)
- **Left: Income vs Expenses** — Grouped bar chart (12 months)
  - Blue bars = income, Red bars = expenses
  - X-axis: months, Y-axis: dollars
  - Chart.js `Bar` type
- **Right: Expense Breakdown** — Doughnut chart
  - Top categories by amount
  - Color-coded by category

### Row 3: Activity & Alerts (2-column)
- **Left: Upcoming Reminders** — List with due dates, type badges
  - Lease expiring, rent due, maintenance follow-up
  - "View All" link to Reminders page
- **Right: Recent Activity** — Timeline of last 10 events
  - Payment recorded, tenant added, expense logged, etc.
  - Timestamp + description

### Row 4: Quick Status
- **Property Card** — Address, current tenant, lease end date
- **Open Maintenance** — Count with link to maintenance page

## Step 3: Chart.js Integration

```bash
cd frontend && npm install chart.js
```

Create `src/components/BarChart.svelte` and `DoughnutChart.svelte` wrappers.

Use reactive Svelte bindings to update charts when data changes.

## Acceptance Criteria
- [ ] Dashboard loads with real data from backend
- [ ] Summary cards show accurate monthly/YTD figures
- [ ] Bar chart renders 12 months of income vs expenses
- [ ] Doughnut chart shows expense categories
- [ ] Upcoming reminders list is accurate
- [ ] Recent activity feed updates in real-time
- [ ] All cards link to their respective detail pages
- [ ] Chart is responsive within the window
