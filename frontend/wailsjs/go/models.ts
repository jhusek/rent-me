export namespace models {
	
	export class ActivityLog {
	    id: number;
	    action: string;
	    entityType: string;
	    entityId: number;
	    description: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new ActivityLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.action = source["action"];
	        this.entityType = source["entityType"];
	        this.entityId = source["entityId"];
	        this.description = source["description"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class Contact {
	    id: number;
	    name: string;
	    role: string;
	    phone: string;
	    email: string;
	    company: string;
	    notes: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Contact(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.role = source["role"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.company = source["company"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class Document {
	    id: number;
	    name: string;
	    category: string;
	    filePath: string;
	    fileSize: number;
	    mimeType: string;
	    uploadDate: string;
	    notes: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Document(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.category = source["category"];
	        this.filePath = source["filePath"];
	        this.fileSize = source["fileSize"];
	        this.mimeType = source["mimeType"];
	        this.uploadDate = source["uploadDate"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class Expense {
	    id: number;
	    category: string;
	    amount: number;
	    date: string;
	    vendor: string;
	    description: string;
	    isRecurring: boolean;
	    recurringInterval: string;
	    receiptPath: string;
	    taxDeductible: boolean;
	    notes: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Expense(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.category = source["category"];
	        this.amount = source["amount"];
	        this.date = source["date"];
	        this.vendor = source["vendor"];
	        this.description = source["description"];
	        this.isRecurring = source["isRecurring"];
	        this.recurringInterval = source["recurringInterval"];
	        this.receiptPath = source["receiptPath"];
	        this.taxDeductible = source["taxDeductible"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class ExpenseSummary {
	    totalAmount: number;
	    byCategory: Record<string, number>;
	    count: number;
	    taxDeductible: number;
	
	    static createFrom(source: any = {}) {
	        return new ExpenseSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalAmount = source["totalAmount"];
	        this.byCategory = source["byCategory"];
	        this.count = source["count"];
	        this.taxDeductible = source["taxDeductible"];
	    }
	}
	export class MaintenanceRequest {
	    id: number;
	    title: string;
	    description: string;
	    priority: string;
	    status: string;
	    reportedDate: string;
	    completedDate: string;
	    cost: number;
	    vendor: string;
	    expenseId: number;
	    notes: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new MaintenanceRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.priority = source["priority"];
	        this.status = source["status"];
	        this.reportedDate = source["reportedDate"];
	        this.completedDate = source["completedDate"];
	        this.cost = source["cost"];
	        this.vendor = source["vendor"];
	        this.expenseId = source["expenseId"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class Payment {
	    id: number;
	    tenantId: number;
	    amount: number;
	    dueDate: string;
	    paidDate: string;
	    method: string;
	    status: string;
	    lateFee: number;
	    notes: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Payment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tenantId = source["tenantId"];
	        this.amount = source["amount"];
	        this.dueDate = source["dueDate"];
	        this.paidDate = source["paidDate"];
	        this.method = source["method"];
	        this.status = source["status"];
	        this.lateFee = source["lateFee"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class PaymentSummary {
	    totalDue: number;
	    totalPaid: number;
	    totalLate: number;
	    totalMissed: number;
	    totalPartial: number;
	    lateFees: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new PaymentSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalDue = source["totalDue"];
	        this.totalPaid = source["totalPaid"];
	        this.totalLate = source["totalLate"];
	        this.totalMissed = source["totalMissed"];
	        this.totalPartial = source["totalPartial"];
	        this.lateFees = source["lateFees"];
	        this.count = source["count"];
	    }
	}
	export class Property {
	    id: number;
	    name: string;
	    address: string;
	    city: string;
	    state: string;
	    zip: string;
	    purchasePrice: number;
	    purchaseDate: string;
	    beds: number;
	    baths: number;
	    sqft: number;
	    yearBuilt: number;
	    lotSize: string;
	    propertyType: string;
	    mortgagePayment: number;
	    mortgageRate: number;
	    mortgageStart: string;
	    mortgageTermYears: number;
	    insuranceMonthly: number;
	    propertyTaxAnnual: number;
	    hoaMonthly: number;
	    notes: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Property(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.city = source["city"];
	        this.state = source["state"];
	        this.zip = source["zip"];
	        this.purchasePrice = source["purchasePrice"];
	        this.purchaseDate = source["purchaseDate"];
	        this.beds = source["beds"];
	        this.baths = source["baths"];
	        this.sqft = source["sqft"];
	        this.yearBuilt = source["yearBuilt"];
	        this.lotSize = source["lotSize"];
	        this.propertyType = source["propertyType"];
	        this.mortgagePayment = source["mortgagePayment"];
	        this.mortgageRate = source["mortgageRate"];
	        this.mortgageStart = source["mortgageStart"];
	        this.mortgageTermYears = source["mortgageTermYears"];
	        this.insuranceMonthly = source["insuranceMonthly"];
	        this.propertyTaxAnnual = source["propertyTaxAnnual"];
	        this.hoaMonthly = source["hoaMonthly"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class Reminder {
	    id: number;
	    title: string;
	    description: string;
	    dueDate: string;
	    type: string;
	    isRecurring: boolean;
	    recurringInterval: string;
	    isDismissed: boolean;
	    relatedEntityType: string;
	    relatedEntityId: number;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Reminder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.dueDate = source["dueDate"];
	        this.type = source["type"];
	        this.isRecurring = source["isRecurring"];
	        this.recurringInterval = source["recurringInterval"];
	        this.isDismissed = source["isDismissed"];
	        this.relatedEntityType = source["relatedEntityType"];
	        this.relatedEntityId = source["relatedEntityId"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class Tenant {
	    id: number;
	    firstName: string;
	    lastName: string;
	    email: string;
	    phone: string;
	    emergencyContactName: string;
	    emergencyContactPhone: string;
	    leaseStart: string;
	    leaseEnd: string;
	    monthlyRent: number;
	    securityDeposit: number;
	    status: string;
	    moveInDate: string;
	    moveOutDate: string;
	    notes: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Tenant(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.firstName = source["firstName"];
	        this.lastName = source["lastName"];
	        this.email = source["email"];
	        this.phone = source["phone"];
	        this.emergencyContactName = source["emergencyContactName"];
	        this.emergencyContactPhone = source["emergencyContactPhone"];
	        this.leaseStart = source["leaseStart"];
	        this.leaseEnd = source["leaseEnd"];
	        this.monthlyRent = source["monthlyRent"];
	        this.securityDeposit = source["securityDeposit"];
	        this.status = source["status"];
	        this.moveInDate = source["moveInDate"];
	        this.moveOutDate = source["moveOutDate"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class UpdateInfo {
	    currentVersion: string;
	    latestVersion: string;
	    updateAvailable: boolean;
	    releaseURL: string;

	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.updateAvailable = source["updateAvailable"];
	        this.releaseURL = source["releaseURL"];
	    }
	}
	
	export class CategoryAmount {
	    category: string;
	    amount: number;
	
	    static createFrom(source: any = {}) {
	        return new CategoryAmount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.amount = source["amount"];
	    }
	}
	export class MonthAmount {
	    month: string;
	    amount: number;
	
	    static createFrom(source: any = {}) {
	        return new MonthAmount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.amount = source["amount"];
	    }
	}
	export class RentStatusSummary {
	    paid: number;
	    due: number;
	    late: number;
	    partial: number;
	    missed: number;
	
	    static createFrom(source: any = {}) {
	        return new RentStatusSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.paid = source["paid"];
	        this.due = source["due"];
	        this.late = source["late"];
	        this.partial = source["partial"];
	        this.missed = source["missed"];
	    }
	}
	export class DashboardData {
	    property: models.Property;
	    currentTenant?: models.Tenant;
	    monthlyIncome: number;
	    monthlyExpenses: number;
	    ytdIncome: number;
	    ytdExpenses: number;
	    netIncome: number;
	    rentStatus: RentStatusSummary;
	    incomeByMonth: MonthAmount[];
	    expensesByMonth: MonthAmount[];
	    expensesByCategory: CategoryAmount[];
	    upcomingReminders: models.Reminder[];
	    recentActivity: models.ActivityLog[];
	    openMaintenance: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.property = this.convertValues(source["property"], models.Property);
	        this.currentTenant = this.convertValues(source["currentTenant"], models.Tenant);
	        this.monthlyIncome = source["monthlyIncome"];
	        this.monthlyExpenses = source["monthlyExpenses"];
	        this.ytdIncome = source["ytdIncome"];
	        this.ytdExpenses = source["ytdExpenses"];
	        this.netIncome = source["netIncome"];
	        this.rentStatus = this.convertValues(source["rentStatus"], RentStatusSummary);
	        this.incomeByMonth = this.convertValues(source["incomeByMonth"], MonthAmount);
	        this.expensesByMonth = this.convertValues(source["expensesByMonth"], MonthAmount);
	        this.expensesByCategory = this.convertValues(source["expensesByCategory"], CategoryAmount);
	        this.upcomingReminders = this.convertValues(source["upcomingReminders"], models.Reminder);
	        this.recentActivity = this.convertValues(source["recentActivity"], models.ActivityLog);
	        this.openMaintenance = source["openMaintenance"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

