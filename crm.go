package crm

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/digitosai/core/registry"
)

type Customer struct {
	ID   int
	Name string
}

type CRMModule struct{}

func init() {
	// Register the module with the core platform upon initialization
	registry.Register(&CRMModule{})
}

func (m *CRMModule) Info() registry.ModuleInfo {
	return registry.ModuleInfo{
		ID:          "crm",
		Name:        "Customer CRM",
		Description: "Manage your customers, leads, and interactions.",
		Icon:        "👤", // Using emoji as icon for simplicity
		Path:        "/crm",
	}
}

func (m *CRMModule) SetupRoutes(r *http.ServeMux) {
	// Root route for the CRM module (mounted at /crm in the core router)
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		tenantID, ok := req.Context().Value("tenantID").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieve the isolated DB connection mapped to this Tenant
		tenantDBRaw := req.Context().Value("tenantDB")
		if tenantDBRaw == nil {
			http.Error(w, "Internal Server Error: No Tenant DB Connection", http.StatusInternalServerError)
			return
		}
		tenantDB := tenantDBRaw.(*sql.DB)

		// Create table if not exists
		_, err := tenantDB.Exec(`CREATE TABLE IF NOT EXISTS customers (id SERIAL PRIMARY KEY, name TEXT NOT NULL)`)
		if err != nil {
			log.Printf("Failed to create customers table: %v", err)
		} else {
			// Seed dummy data if empty
			var count int
			tenantDB.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
			if count == 0 {
				tenantDB.Exec("INSERT INTO customers (name) VALUES ('Alice Smith (New Lead)'), ('Bob Jones (Contacted)')")
			}
		}

		// Query from the Tenant Database
		rows, err := tenantDB.Query("SELECT id, name FROM customers")

		customerListHTML := ""
		if err != nil {
			customerListHTML = fmt.Sprintf("<li>Error loading customers: %v</li>", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				var customer Customer
				rows.Scan(&customer.ID, &customer.Name)
				customerListHTML += fmt.Sprintf("<li>#%d - %s</li>", customer.ID, customer.Name)
			}
		}

		if customerListHTML == "" {
			customerListHTML = "<li>No customers found.</li>"
		}

		fmt.Fprintf(w, `
		<div style="font-family: sans-serif; padding: 20px;">
			<h2 style="color: #2563eb;">Customer Relationship Management</h2>
			<p>Welcome to the CRM module! Your isolated Tenant ID is: <strong>%d</strong></p>
			
			<div style="margin-top: 20px; padding: 15px; border: 1px solid #e5e7eb; border-radius: 8px;">
				<h3>Your Leads (Fetched from isolated database)</h3>
				<ul style="line-height: 1.6;">
					%s
				</ul>
			</div>
			<a href="/" style="display: inline-block; margin-top: 15px; color: #4b5563; text-decoration: none;">&larr; Back to Dashboard</a>
		</div>
		`, tenantID, customerListHTML)
	})

	r.HandleFunc("/api/customers", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id": 1, "name": "API Route Placeholder"}]`))
	})
}
