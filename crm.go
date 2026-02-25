package crm

import (
	"fmt"
	"net/http"

	"github.com/digitosai/core/registry"
)

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
			tenantID = -1
		}

		fmt.Fprintf(w, `
		<div style="font-family: sans-serif; padding: 20px;">
			<h2 style="color: #2563eb;">Customer Relationship Management</h2>
			<p>Welcome to the CRM module! Your Tenant ID is: %d</p>
			
			<div style="margin-top: 20px; padding: 15px; border: 1px solid #e5e7eb; border-radius: 8px;">
				<h3>Your Leads</h3>
				<ul>
					<li>Alice Smith (New)</li>
					<li>Bob Jones (Contacted)</li>
				</ul>
			</div>
			<a href="/" style="display: inline-block; margin-top: 15px; color: #4b5563;">&larr; Back to Dashboard</a>
		</div>
		`, tenantID)
	})

	r.HandleFunc("/api/customers", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id": 1, "name": "Alice Smith"}, {"id": 2, "name": "Bob Jones"}]`))
	})
}
