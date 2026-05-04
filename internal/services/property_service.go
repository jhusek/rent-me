package services

import (
	"database/sql"
	"time"

	"RentMe/internal/models"
)

type PropertyService struct {
	db *sql.DB
}

func NewPropertyService(db *sql.DB) *PropertyService {
	return &PropertyService{db: db}
}

func (s *PropertyService) GetProperty() (models.Property, error) {
	var p models.Property
	var name, address, city, state, zip, purchaseDate, lotSize, propertyType sql.NullString
	var mortgageStart, notes, createdAt, updatedAt sql.NullString
	var purchasePrice, baths, mortgagePayment, mortgageRate sql.NullFloat64
	var insuranceMonthly, propertyTaxAnnual, hoaMonthly sql.NullFloat64
	var beds, sqft, yearBuilt, mortgageTermYears sql.NullInt64

	row := s.db.QueryRow("SELECT id, name, address, city, state, zip, purchase_price, purchase_date, beds, baths, sqft, year_built, lot_size, property_type, mortgage_payment, mortgage_rate, mortgage_start, mortgage_term_years, insurance_monthly, property_tax_annual, hoa_monthly, notes, created_at, updated_at FROM property LIMIT 1")
	err := row.Scan(&p.ID, &name, &address, &city, &state, &zip, &purchasePrice, &purchaseDate, &beds, &baths, &sqft, &yearBuilt, &lotSize, &propertyType, &mortgagePayment, &mortgageRate, &mortgageStart, &mortgageTermYears, &insuranceMonthly, &propertyTaxAnnual, &hoaMonthly, &notes, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return models.Property{}, nil
	}
	if err != nil {
		return models.Property{}, err
	}

	p.Name = name.String
	p.Address = address.String
	p.City = city.String
	p.State = state.String
	p.Zip = zip.String
	p.PurchasePrice = purchasePrice.Float64
	p.PurchaseDate = purchaseDate.String
	p.Beds = int(beds.Int64)
	p.Baths = baths.Float64
	p.Sqft = int(sqft.Int64)
	p.YearBuilt = int(yearBuilt.Int64)
	p.LotSize = lotSize.String
	p.PropertyType = propertyType.String
	p.MortgagePayment = mortgagePayment.Float64
	p.MortgageRate = mortgageRate.Float64
	p.MortgageStart = mortgageStart.String
	p.MortgageTermYears = int(mortgageTermYears.Int64)
	p.InsuranceMonthly = insuranceMonthly.Float64
	p.PropertyTaxAnnual = propertyTaxAnnual.Float64
	p.HOAMonthly = hoaMonthly.Float64
	p.Notes = notes.String
	p.CreatedAt = createdAt.String
	p.UpdatedAt = updatedAt.String

	return p, nil
}

func (s *PropertyService) UpdateProperty(p models.Property) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	existing, err := s.GetProperty()
	if err != nil {
		return err
	}

	if existing.ID == 0 {
		_, err = s.db.Exec(`INSERT INTO property (name, address, city, state, zip, purchase_price, purchase_date, beds, baths, sqft, year_built, lot_size, property_type, mortgage_payment, mortgage_rate, mortgage_start, mortgage_term_years, insurance_monthly, property_tax_annual, hoa_monthly, notes, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			p.Name, p.Address, p.City, p.State, p.Zip, p.PurchasePrice, p.PurchaseDate, p.Beds, p.Baths, p.Sqft, p.YearBuilt, p.LotSize, p.PropertyType, p.MortgagePayment, p.MortgageRate, p.MortgageStart, p.MortgageTermYears, p.InsuranceMonthly, p.PropertyTaxAnnual, p.HOAMonthly, p.Notes, now)
		return err
	}

	_, err = s.db.Exec(`UPDATE property SET name=?, address=?, city=?, state=?, zip=?, purchase_price=?, purchase_date=?, beds=?, baths=?, sqft=?, year_built=?, lot_size=?, property_type=?, mortgage_payment=?, mortgage_rate=?, mortgage_start=?, mortgage_term_years=?, insurance_monthly=?, property_tax_annual=?, hoa_monthly=?, notes=?, updated_at=? WHERE id=?`,
		p.Name, p.Address, p.City, p.State, p.Zip, p.PurchasePrice, p.PurchaseDate, p.Beds, p.Baths, p.Sqft, p.YearBuilt, p.LotSize, p.PropertyType, p.MortgagePayment, p.MortgageRate, p.MortgageStart, p.MortgageTermYears, p.InsuranceMonthly, p.PropertyTaxAnnual, p.HOAMonthly, p.Notes, now, existing.ID)
	return err
}
