package main

import (
	"context"
	"database/sql"
	"fmt"

	"RentMe/internal/db"
	"RentMe/internal/models"
	"RentMe/internal/services"
)

type App struct {
	ctx                context.Context
	database           *sql.DB
	propertyService    *services.PropertyService
	tenantService      *services.TenantService
	paymentService     *services.PaymentService
	expenseService     *services.ExpenseService
	maintenanceService *services.MaintenanceService
	documentService    *services.DocumentService
	reminderService    *services.ReminderService
	contactService     *services.ContactService
	activityService    *services.ActivityService
	dashboardService   *services.DashboardService
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	database, err := db.InitDB("./rentme.db")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	a.database = database

	a.activityService = services.NewActivityService(database)
	a.propertyService = services.NewPropertyService(database)
	a.tenantService = services.NewTenantService(database)
	a.paymentService = services.NewPaymentService(database)
	a.expenseService = services.NewExpenseService(database)
	a.maintenanceService = services.NewMaintenanceService(database)
	a.documentService = services.NewDocumentService(database)
	a.reminderService = services.NewReminderService(database)
	a.contactService = services.NewContactService(database)
	a.dashboardService = services.NewDashboardService(
		database, a.paymentService, a.expenseService, a.reminderService,
		a.maintenanceService, a.propertyService, a.tenantService, a.activityService,
	)
}

func (a *App) shutdown(ctx context.Context) {
	if a.database != nil {
		a.database.Close()
	}
}

// ─── Property ───────────────────────────────────────────────────────────────

func (a *App) GetProperty() (models.Property, error) {
	return a.propertyService.GetProperty()
}

func (a *App) UpdateProperty(p models.Property) error {
	err := a.propertyService.UpdateProperty(p)
	if err == nil {
		a.activityService.LogActivity("updated", "property", p.ID, "Property details updated")
	}
	return err
}

// ─── Tenants ────────────────────────────────────────────────────────────────

func (a *App) GetTenants(status string) ([]models.Tenant, error) {
	return a.tenantService.ListTenants(status)
}

func (a *App) GetTenant(id int64) (models.Tenant, error) {
	return a.tenantService.GetTenant(id)
}

func (a *App) CreateTenant(t models.Tenant) (models.Tenant, error) {
	result, err := a.tenantService.CreateTenant(t)
	if err == nil {
		a.activityService.LogActivity("created", "tenant", result.ID, "Tenant "+t.FirstName+" "+t.LastName+" added")
	}
	return result, err
}

func (a *App) UpdateTenant(t models.Tenant) error {
	err := a.tenantService.UpdateTenant(t)
	if err == nil {
		a.activityService.LogActivity("updated", "tenant", t.ID, "Tenant "+t.FirstName+" "+t.LastName+" updated")
	}
	return err
}

func (a *App) DeleteTenant(id int64) error {
	err := a.tenantService.DeleteTenant(id)
	if err == nil {
		a.activityService.LogActivity("deleted", "tenant", id, "Tenant removed")
	}
	return err
}

// ─── Payments ───────────────────────────────────────────────────────────────

func (a *App) GetPayments(tenantID int64, year int, month int) ([]models.Payment, error) {
	return a.paymentService.ListPayments(tenantID, year, month)
}

func (a *App) GetPayment(id int64) (models.Payment, error) {
	return a.paymentService.GetPayment(id)
}

func (a *App) CreatePayment(p models.Payment) (models.Payment, error) {
	result, err := a.paymentService.CreatePayment(p)
	if err == nil {
		a.activityService.LogActivity("created", "payment", result.ID, fmt.Sprintf("Payment of $%.2f recorded", p.Amount))
	}
	return result, err
}

func (a *App) UpdatePayment(p models.Payment) error {
	err := a.paymentService.UpdatePayment(p)
	if err == nil {
		a.activityService.LogActivity("updated", "payment", p.ID, fmt.Sprintf("Payment of $%.2f updated", p.Amount))
	}
	return err
}

func (a *App) DeletePayment(id int64) error {
	err := a.paymentService.DeletePayment(id)
	if err == nil {
		a.activityService.LogActivity("deleted", "payment", id, "Payment removed")
	}
	return err
}

func (a *App) GetPaymentSummary(year int, month int) (models.PaymentSummary, error) {
	return a.paymentService.GetPaymentSummary(year, month)
}

// ─── Expenses ───────────────────────────────────────────────────────────────

func (a *App) GetExpenses(category string, year int) ([]models.Expense, error) {
	return a.expenseService.ListExpenses(category, year)
}

func (a *App) GetExpense(id int64) (models.Expense, error) {
	return a.expenseService.GetExpense(id)
}

func (a *App) CreateExpense(e models.Expense) (models.Expense, error) {
	result, err := a.expenseService.CreateExpense(e)
	if err == nil {
		a.activityService.LogActivity("created", "expense", result.ID, fmt.Sprintf("%s expense of $%.2f recorded", e.Category, e.Amount))
	}
	return result, err
}

func (a *App) UpdateExpense(e models.Expense) error {
	err := a.expenseService.UpdateExpense(e)
	if err == nil {
		a.activityService.LogActivity("updated", "expense", e.ID, fmt.Sprintf("%s expense of $%.2f updated", e.Category, e.Amount))
	}
	return err
}

func (a *App) DeleteExpense(id int64) error {
	err := a.expenseService.DeleteExpense(id)
	if err == nil {
		a.activityService.LogActivity("deleted", "expense", id, "Expense removed")
	}
	return err
}

func (a *App) GetExpenseSummary(year int) (models.ExpenseSummary, error) {
	return a.expenseService.GetExpenseSummary(year)
}

// ─── Maintenance ────────────────────────────────────────────────────────────

func (a *App) GetMaintenanceRequests(status string) ([]models.MaintenanceRequest, error) {
	return a.maintenanceService.ListMaintenance(status)
}

func (a *App) GetMaintenanceRequest(id int64) (models.MaintenanceRequest, error) {
	return a.maintenanceService.GetMaintenance(id)
}

func (a *App) CreateMaintenanceRequest(m models.MaintenanceRequest) (models.MaintenanceRequest, error) {
	result, err := a.maintenanceService.CreateMaintenance(m)
	if err == nil {
		a.activityService.LogActivity("created", "maintenance", result.ID, "Maintenance request: "+m.Title)
	}
	return result, err
}

func (a *App) UpdateMaintenanceRequest(m models.MaintenanceRequest) error {
	err := a.maintenanceService.UpdateMaintenance(m)
	if err == nil {
		a.activityService.LogActivity("updated", "maintenance", m.ID, "Maintenance request updated: "+m.Title)
	}
	return err
}

func (a *App) DeleteMaintenanceRequest(id int64) error {
	err := a.maintenanceService.DeleteMaintenance(id)
	if err == nil {
		a.activityService.LogActivity("deleted", "maintenance", id, "Maintenance request removed")
	}
	return err
}

// ─── Documents ──────────────────────────────────────────────────────────────

func (a *App) GetDocuments(category string) ([]models.Document, error) {
	return a.documentService.ListDocuments(category)
}

func (a *App) GetDocument(id int64) (models.Document, error) {
	return a.documentService.GetDocument(id)
}

func (a *App) CreateDocument(d models.Document) (models.Document, error) {
	result, err := a.documentService.CreateDocument(d)
	if err == nil {
		a.activityService.LogActivity("created", "document", result.ID, "Document uploaded: "+d.Name)
	}
	return result, err
}

func (a *App) DeleteDocument(id int64) error {
	err := a.documentService.DeleteDocument(id)
	if err == nil {
		a.activityService.LogActivity("deleted", "document", id, "Document removed")
	}
	return err
}

// ─── Reminders ──────────────────────────────────────────────────────────────

func (a *App) GetReminders(showDismissed bool) ([]models.Reminder, error) {
	return a.reminderService.ListReminders(showDismissed)
}

func (a *App) GetUpcomingReminders(days int) ([]models.Reminder, error) {
	return a.reminderService.GetUpcomingReminders(days)
}

func (a *App) CreateReminder(r models.Reminder) (models.Reminder, error) {
	result, err := a.reminderService.CreateReminder(r)
	if err == nil {
		a.activityService.LogActivity("created", "reminder", result.ID, "Reminder created: "+r.Title)
	}
	return result, err
}

func (a *App) UpdateReminder(r models.Reminder) error {
	err := a.reminderService.UpdateReminder(r)
	if err == nil {
		a.activityService.LogActivity("updated", "reminder", r.ID, "Reminder updated: "+r.Title)
	}
	return err
}

func (a *App) DismissReminder(id int64) error {
	err := a.reminderService.DismissReminder(id)
	if err == nil {
		a.activityService.LogActivity("dismissed", "reminder", id, "Reminder dismissed")
	}
	return err
}

func (a *App) DeleteReminder(id int64) error {
	err := a.reminderService.DeleteReminder(id)
	if err == nil {
		a.activityService.LogActivity("deleted", "reminder", id, "Reminder removed")
	}
	return err
}

// ─── Contacts ───────────────────────────────────────────────────────────────

func (a *App) GetContacts() ([]models.Contact, error) {
	return a.contactService.ListContacts()
}

func (a *App) GetContact(id int64) (models.Contact, error) {
	return a.contactService.GetContact(id)
}

func (a *App) CreateContact(c models.Contact) (models.Contact, error) {
	result, err := a.contactService.CreateContact(c)
	if err == nil {
		a.activityService.LogActivity("created", "contact", result.ID, "Contact added: "+c.Name)
	}
	return result, err
}

func (a *App) UpdateContact(c models.Contact) error {
	err := a.contactService.UpdateContact(c)
	if err == nil {
		a.activityService.LogActivity("updated", "contact", c.ID, "Contact updated: "+c.Name)
	}
	return err
}

func (a *App) DeleteContact(id int64) error {
	err := a.contactService.DeleteContact(id)
	if err == nil {
		a.activityService.LogActivity("deleted", "contact", id, "Contact removed")
	}
	return err
}

// ─── Dashboard & Activity ───────────────────────────────────────────────────

func (a *App) GetDashboardData() (services.DashboardData, error) {
	return a.dashboardService.GetDashboardData()
}

func (a *App) GetRecentActivity(limit int) ([]models.ActivityLog, error) {
	return a.activityService.GetRecentActivity(limit)
}
