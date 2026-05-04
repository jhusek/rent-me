# RentMe — Passdown Notes

## Status as of 2026-03-25 — ALL PHASES COMPLETE ✅
All 8 phases done. 14/14 todos completed. `wails build` produces RentMe.exe.

## What's Built
- **Go Backend**: SQLite (modernc.org/sqlite), 9 tables, 10 service files, full CRUD + dashboard aggregation
- **Frontend**: Svelte 3 + TypeScript + Tailwind 3 + Chart.js + svelte-spa-router
- **11 Pages**: Dashboard, Property, Tenants, Payments, Expenses, Maintenance, Documents, Reminders, Reports, Contacts, Settings
- **8 Components**: PageHeader, Modal, Badge, EmptyState, Toast, ConfirmDialog, BarChart, DoughnutChart
- **Features**: Chart.js visualizations, sidebar nav with reminder badges, form validation, CSV export, Schedule E tax helper, depreciation tracker, first-run welcome experience

## How to Run
- Dev mode: `wails dev` (from project root)
- Build: `wails build` → produces `build/bin/RentMe.exe`
- DB file: `./rentme.db` (created on first run)

## Future Enhancements
- File picker integration for document uploads (Wails runtime.OpenFileDialog)
- Actual file copying for document storage
- Dark mode
- Data backup/restore via Wails file dialogs
- Multi-property support (if needed later)

## Environment
- Go 1.26, Wails v2.11, Windows
- Project: C:\Users\jhusek\OneDrive - Intel Corporation\Documents\Coding\VS Code\RentMe
- SQLite via modernc.org/sqlite (pure Go, no CGO)
- Frontend: Svelte 3 + TypeScript + Tailwind 3 + svelte-spa-router + Chart.js

## Key Decisions
- Single property app (no multi-property support)
- DB file: ./rentme.db (dev), %APPDATA%/RentMe/rentme.db (production later)
- Documents stored in local filesystem alongside DB
- Hash-based routing for Wails compatibility
- IRS Schedule E-aligned expense categories
