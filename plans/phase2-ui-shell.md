# Phase 2: UI Shell — Layout, Navigation & Components

## Goal
Build the frontend application shell: sidebar navigation, content area, 
Tailwind CSS styling, and reusable Svelte components. After this phase, 
clicking nav items shows placeholder pages.

## Deliverables
1. Tailwind CSS integrated and working
2. Sidebar navigation with icons and active state
3. Page routing (svelte-spa-router or simple reactive routing)
4. Reusable UI components library
5. All placeholder pages stubbed out

---

## Step 1: Install Frontend Dependencies

```bash
cd frontend
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
npm install svelte-spa-router
npm install chart.js       # for Phase 6
```

## Step 2: Configure Tailwind

`tailwind.config.js`:
```js
export default {
  content: ['./src/**/*.{html,js,svelte,ts}', './index.html'],
  theme: {
    extend: {
      colors: {
        primary: { 50: '#eff6ff', 500: '#3b82f6', 600: '#2563eb', 700: '#1d4ed8' },
        sidebar: '#1e293b',
      }
    },
  },
  plugins: [],
}
```

`src/style.css` — replace with Tailwind directives:
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

## Step 3: App Layout (`src/App.svelte`)

Two-column layout:
- **Left sidebar** (w-64, fixed): App logo/name, nav links, version
- **Main content** (flex-1, scrollable): Routed page content

### Navigation Items
| Icon | Label | Route |
|------|-------|-------|
| 🏠 | Dashboard | `/` |
| 🏡 | Property | `/property` |
| 👤 | Tenants | `/tenants` |
| 💰 | Payments | `/payments` |
| 📊 | Expenses | `/expenses` |
| 🔧 | Maintenance | `/maintenance` |
| 📁 | Documents | `/documents` |
| 🔔 | Reminders | `/reminders` |
| 📈 | Reports | `/reports` |
| 📞 | Contacts | `/contacts` |
| ⚙️ | Settings | `/settings` |

## Step 4: Page Stubs (`src/pages/`)

Create one `.svelte` file per page:
- `Dashboard.svelte`
- `Property.svelte`
- `Tenants.svelte`
- `Payments.svelte`
- `Expenses.svelte`
- `Maintenance.svelte`
- `Documents.svelte`
- `Reminders.svelte`
- `Reports.svelte`
- `Contacts.svelte`
- `Settings.svelte`

Each starts with a heading and "Coming soon" placeholder.

## Step 5: Reusable Components (`src/components/`)

### Core Components
- **PageHeader.svelte** — Title + optional action button
- **Card.svelte** — Flexible card container with title, optional footer
- **DataTable.svelte** — Sortable table with slot for rows
- **Modal.svelte** — Centered overlay modal with close button
- **FormField.svelte** — Label + input wrapper (text, number, select, date, textarea)
- **Button.svelte** — Primary, secondary, danger variants
- **Badge.svelte** — Status badges (paid, due, late, open, completed, etc.)
- **EmptyState.svelte** — "No data yet" placeholder with icon and CTA
- **Toast.svelte** — Success/error notification toast
- **ConfirmDialog.svelte** — "Are you sure?" modal

### Layout
- Clean, professional look — blue primary color, slate sidebar
- Consistent spacing (p-6 content areas, gap-4 grids)
- Responsive within the desktop window

## Step 6: Verify

```bash
cd frontend && npm run build
wails dev  # should show sidebar + placeholder pages
```

## Acceptance Criteria
- [ ] Tailwind works (utility classes render correctly)
- [ ] Sidebar shows all nav items with working routing
- [ ] Active nav item is highlighted
- [ ] All placeholder pages render
- [ ] Components compile without errors
- [ ] Clean, professional appearance
