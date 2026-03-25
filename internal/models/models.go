package models

type Property struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Address           string  `json:"address"`
	City              string  `json:"city"`
	State             string  `json:"state"`
	Zip               string  `json:"zip"`
	PurchasePrice     float64 `json:"purchasePrice"`
	PurchaseDate      string  `json:"purchaseDate"`
	Beds              int     `json:"beds"`
	Baths             float64 `json:"baths"`
	Sqft              int     `json:"sqft"`
	YearBuilt         int     `json:"yearBuilt"`
	LotSize           string  `json:"lotSize"`
	PropertyType      string  `json:"propertyType"`
	MortgagePayment   float64 `json:"mortgagePayment"`
	MortgageRate      float64 `json:"mortgageRate"`
	MortgageStart     string  `json:"mortgageStart"`
	MortgageTermYears int     `json:"mortgageTermYears"`
	InsuranceMonthly  float64 `json:"insuranceMonthly"`
	PropertyTaxAnnual float64 `json:"propertyTaxAnnual"`
	HOAMonthly        float64 `json:"hoaMonthly"`
	Notes             string  `json:"notes"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

type Tenant struct {
	ID                    int64   `json:"id"`
	FirstName             string  `json:"firstName"`
	LastName              string  `json:"lastName"`
	Email                 string  `json:"email"`
	Phone                 string  `json:"phone"`
	EmergencyContactName  string  `json:"emergencyContactName"`
	EmergencyContactPhone string  `json:"emergencyContactPhone"`
	LeaseStart            string  `json:"leaseStart"`
	LeaseEnd              string  `json:"leaseEnd"`
	MonthlyRent           float64 `json:"monthlyRent"`
	SecurityDeposit       float64 `json:"securityDeposit"`
	Status                string  `json:"status"`
	MoveInDate            string  `json:"moveInDate"`
	MoveOutDate           string  `json:"moveOutDate"`
	Notes                 string  `json:"notes"`
	CreatedAt             string  `json:"createdAt"`
	UpdatedAt             string  `json:"updatedAt"`
}

type Payment struct {
	ID       int64   `json:"id"`
	TenantID int64   `json:"tenantId"`
	Amount   float64 `json:"amount"`
	DueDate  string  `json:"dueDate"`
	PaidDate string  `json:"paidDate"`
	Method   string  `json:"method"`
	Status   string  `json:"status"`
	LateFee  float64 `json:"lateFee"`
	Notes    string  `json:"notes"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Expense struct {
	ID                int64   `json:"id"`
	Category          string  `json:"category"`
	Amount            float64 `json:"amount"`
	Date              string  `json:"date"`
	Vendor            string  `json:"vendor"`
	Description       string  `json:"description"`
	IsRecurring       bool    `json:"isRecurring"`
	RecurringInterval string  `json:"recurringInterval"`
	ReceiptPath       string  `json:"receiptPath"`
	TaxDeductible     bool    `json:"taxDeductible"`
	Notes             string  `json:"notes"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

type MaintenanceRequest struct {
	ID            int64   `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Priority      string  `json:"priority"`
	Status        string  `json:"status"`
	ReportedDate  string  `json:"reportedDate"`
	CompletedDate string  `json:"completedDate"`
	Cost          float64 `json:"cost"`
	Vendor        string  `json:"vendor"`
	ExpenseID     int64   `json:"expenseId"`
	Notes         string  `json:"notes"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type Document struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	FilePath   string `json:"filePath"`
	FileSize   int    `json:"fileSize"`
	MimeType   string `json:"mimeType"`
	UploadDate string `json:"uploadDate"`
	Notes      string `json:"notes"`
	CreatedAt  string `json:"createdAt"`
}

type Reminder struct {
	ID                int64  `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	DueDate           string `json:"dueDate"`
	Type              string `json:"type"`
	IsRecurring       bool   `json:"isRecurring"`
	RecurringInterval string `json:"recurringInterval"`
	IsDismissed       bool   `json:"isDismissed"`
	RelatedEntityType string `json:"relatedEntityType"`
	RelatedEntityID   int64  `json:"relatedEntityId"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type Contact struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ActivityLog struct {
	ID          int64  `json:"id"`
	Action      string `json:"action"`
	EntityType  string `json:"entityType"`
	EntityID    int64  `json:"entityId"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type PaymentSummary struct {
	TotalDue     float64 `json:"totalDue"`
	TotalPaid    float64 `json:"totalPaid"`
	TotalLate    float64 `json:"totalLate"`
	TotalMissed  float64 `json:"totalMissed"`
	TotalPartial float64 `json:"totalPartial"`
	LateFees     float64 `json:"lateFees"`
	Count        int     `json:"count"`
}

type ExpenseSummary struct {
	TotalAmount    float64          `json:"totalAmount"`
	ByCategory     map[string]float64 `json:"byCategory"`
	Count          int              `json:"count"`
	TaxDeductible  float64          `json:"taxDeductible"`
}
