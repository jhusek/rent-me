# RentMe — Master Plan

## Overview
Desktop rental property management app built with Go + Wails v2 + Svelte + SQLite.

## Phase Index

| Phase | File | Description | Dependencies |
|-------|------|-------------|--------------|
| 1 | [phase1-foundation.md](phase1-foundation.md) | Database, models, services, Wails bindings | None |
| 2 | [phase2-ui-shell.md](phase2-ui-shell.md) | App layout, navigation, Tailwind, reusable components | Phase 1 |
| 3 | [phase3-property-tenants.md](phase3-property-tenants.md) | Property profile + tenant management pages | Phase 1 + 2 |
| 4 | [phase4-financials.md](phase4-financials.md) | Rent payments + expense tracking | Phase 3 |
| 5 | [phase5-maintenance-docs.md](phase5-maintenance-docs.md) | Maintenance requests + document storage | Phase 3 + 4 |
| 6 | [phase6-dashboard.md](phase6-dashboard.md) | Dashboard with charts and activity feed | Phase 4 + 5 |
| 7 | [phase7-reminders-reports.md](phase7-reminders-reports.md) | Reminder system + tax reports + CSV export | Phase 4 + 6 |
| 8 | [phase8-polish.md](phase8-polish.md) | Search, validation, settings, final QA | All prior |

## Execution Strategy
- Phases 1 & 2 can be built in parallel (Go backend + frontend shell)
- Phase 3 requires both 1 & 2
- Phases 4 & 5 have partial overlap but should be sequential
- Phase 6 consumes data from 4 & 5
- Phase 7 builds on 6
- Phase 8 is the final sweep

## Tech Decisions
- **SQLite:** `modernc.org/sqlite` (pure Go, no CGO needed)
- **Frontend:** Svelte 3 + TypeScript + Tailwind CSS
- **Charts:** Chart.js via CDN or npm
- **File Storage:** Local `documents/` folder alongside SQLite DB
- **DB Location:** `%APPDATA%/RentMe/rentme.db` (production) or `./rentme.db` (dev)
