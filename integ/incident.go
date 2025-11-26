//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/opsorch/opsorch-adapter/incident"
	"github.com/opsorch/opsorch-core/schema"
)

func main() {
	// Test statistics
	var totalTests, passedTests, failedTests int
	startTime := time.Now()

	testResult := func(name string, err error) {
		totalTests++
		if err != nil {
			failedTests++
			log.Printf("❌ %s: %v", name, err)
		} else {
			passedTests++
			fmt.Printf("✅ %s passed\n", name)
		}
	}

	fmt.Println("=================================")
	fmt.Println("Example Adapter Integration Test")
	fmt.Println("=================================")
	fmt.Printf("Started: %s\n\n", startTime.Format("2006-01-02 15:04:05"))

	ctx := context.Background()

	// Create the incident provider with custom config
	config := map[string]any{
		"source":          "integration-test",
		"defaultSeverity": "sev2",
	}

	provider, err := incident.New(config)
	if err != nil {
		log.Fatalf("Failed to create example incident provider: %v", err)
	}

	// Test 1: Query incidents (should be empty initially)
	fmt.Println("\n=== Test 1: Query Empty Incidents ===")
	incidents, err := provider.Query(ctx, schema.IncidentQuery{})
	if err != nil {
		testResult("Query empty incidents", err)
	} else {
		if len(incidents) != 0 {
			testResult("Validate empty list", fmt.Errorf("expected 0 incidents, got %d", len(incidents)))
		} else {
			fmt.Println("Correctly returned empty list")
			testResult("Query empty incidents", nil)
		}
	}

	// Test 2: Create a new incident
	fmt.Println("\n=== Test 2: Create New Incident ===")
	newIncident, err := provider.Create(ctx, schema.CreateIncidentInput{
		Title:    "Integration Test Incident",
		Status:   "open",
		Severity: "critical",
		Service:  "test-service",
		Fields: map[string]any{
			"description": "This is a test incident created by the integration test.",
		},
		Metadata: map[string]any{
			"environment": "staging",
			"team":        "platform",
		},
	})
	if err != nil {
		testResult("Create incident", err)
	} else {
		fmt.Printf("Successfully created incident:\n")
		fmt.Printf("  ID: %s\n", newIncident.ID)
		fmt.Printf("  Title: %s\n", newIncident.Title)
		fmt.Printf("  Status: %s\n", newIncident.Status)
		fmt.Printf("  Severity: %s\n", newIncident.Severity)
		fmt.Printf("  Service: %s\n", newIncident.Service)
		fmt.Printf("  Created: %s\n", newIncident.CreatedAt.Format("2006-01-02 15:04:05"))

		// Validate created incident
		if newIncident.Title != "Integration Test Incident" {
			testResult("Validate created incident title", fmt.Errorf("title mismatch"))
		} else if newIncident.Status != "open" {
			testResult("Validate created incident status", fmt.Errorf("status mismatch"))
		} else if newIncident.Metadata["source"] != "integration-test" {
			testResult("Validate source metadata", fmt.Errorf("expected source=integration-test, got %v", newIncident.Metadata["source"]))
		} else {
			testResult("Create incident", nil)
		}
	}

	var incident1ID string
	if err == nil {
		incident1ID = newIncident.ID
	}

	// Test 3: Create incident with default severity
	fmt.Println("\n=== Test 3: Create Incident with Default Severity ===")
	incident2, err := provider.Create(ctx, schema.CreateIncidentInput{
		Title:   "Test Incident 2",
		Status:  "acknowledged",
		Service: "api-service",
	})
	if err != nil {
		testResult("Create incident with default severity", err)
	} else {
		if incident2.Severity != "sev2" {
			testResult("Validate default severity", fmt.Errorf("expected sev2, got %s", incident2.Severity))
		} else {
			fmt.Printf("Correctly applied default severity: %s\n", incident2.Severity)
			testResult("Create incident with default severity", nil)
		}
	}

	// Test 4: Create more incidents for query testing
	fmt.Println("\n=== Test 4: Create Additional Test Incidents ===")
	testIncidents := []schema.CreateIncidentInput{
		{
			Title:    "High Priority Alert",
			Status:   "open",
			Severity: "high",
			Service:  "database-service",
			Metadata: map[string]any{"environment": "production", "team": "database"},
		},
		{
			Title:    "Low Priority Warning",
			Status:   "resolved",
			Severity: "low",
			Service:  "monitoring-service",
			Metadata: map[string]any{"environment": "staging", "team": "platform"},
		},
		{
			Title:    "Critical Production Issue",
			Status:   "open",
			Severity: "critical",
			Service:  "api-service",
			Metadata: map[string]any{"environment": "production", "team": "backend"},
		},
	}

	createdCount := 0
	for _, input := range testIncidents {
		_, err := provider.Create(ctx, input)
		if err != nil {
			testResult("Create test incident", err)
		} else {
			createdCount++
		}
	}
	if createdCount == len(testIncidents) {
		fmt.Printf("Created %d test incidents\n", createdCount)
		testResult("Create test incidents batch", nil)
	}

	// Test 5: Query all incidents
	fmt.Println("\n=== Test 5: Query All Incidents ===")
	allIncidents, err := provider.Query(ctx, schema.IncidentQuery{})
	if err != nil {
		testResult("Query all incidents", err)
	} else {
		fmt.Printf("Found %d incidents\n", len(allIncidents))
		for i, inc := range allIncidents {
			fmt.Printf("  [%d] ID: %s, Title: %s, Status: %s, Severity: %s\n",
				i+1, inc.ID, inc.Title, inc.Status, inc.Severity)
		}
		if len(allIncidents) >= 5 {
			testResult("Query all incidents", nil)
		} else {
			testResult("Query all incidents", fmt.Errorf("expected at least 5 incidents, got %d", len(allIncidents)))
		}
	}

	// Test 6: Get specific incident
	if incident1ID != "" {
		fmt.Println("\n=== Test 6: Get Specific Incident ===")
		inc, err := provider.Get(ctx, incident1ID)
		if err != nil {
			testResult("Get incident by ID", err)
		} else {
			fmt.Printf("Retrieved incident:\n")
			fmt.Printf("  ID: %s\n", inc.ID)
			fmt.Printf("  Title: %s\n", inc.Title)
			fmt.Printf("  Status: %s\n", inc.Status)
			fmt.Printf("  Severity: %s\n", inc.Severity)

			if inc.ID != incident1ID {
				testResult("Validate retrieved incident ID", fmt.Errorf("ID mismatch"))
			} else {
				testResult("Get incident by ID", nil)
			}
		}
	}

	// Test 7: Error handling - invalid incident ID
	fmt.Println("\n=== Test 7: Error Handling - Invalid ID ===")
	_, err = provider.Get(ctx, "INVALID_ID_9999")
	if err != nil {
		fmt.Printf("Correctly handled invalid ID: %v\n", err)
		testResult("Error handling for invalid ID", nil)
	} else {
		testResult("Error handling for invalid ID", fmt.Errorf("should have returned error for invalid ID"))
	}

	// Test 8: Query all incidents using Query method
	fmt.Println("\n=== Test 8: Query All Incidents ===")
	queryIncidents, err := provider.Query(ctx, schema.IncidentQuery{})
	if err != nil {
		testResult("Query all incidents", err)
	} else {
		fmt.Printf("Found %d incidents via Query\n", len(queryIncidents))
		testResult("Query all incidents", nil)
	}

	// Test 9: Query by status
	fmt.Println("\n=== Test 9: Query Open Incidents ===")
	openIncidents, err := provider.Query(ctx, schema.IncidentQuery{
		Statuses: []string{"open"},
	})
	if err != nil {
		testResult("Query open incidents", err)
	} else {
		fmt.Printf("Found %d open incidents\n", len(openIncidents))
		allOpen := true
		for _, inc := range openIncidents {
			if inc.Status != "open" {
				allOpen = false
				testResult("Validate status filter", fmt.Errorf("expected open status, got %s", inc.Status))
				break
			}
		}
		if allOpen {
			testResult("Query open incidents", nil)
		}
	}

	// Test 10: Query by severity
	fmt.Println("\n=== Test 10: Query Critical Incidents ===")
	criticalIncidents, err := provider.Query(ctx, schema.IncidentQuery{
		Severities: []string{"critical"},
	})
	if err != nil {
		testResult("Query by severity", err)
	} else {
		fmt.Printf("Found %d critical incidents\n", len(criticalIncidents))
		allCritical := true
		for _, inc := range criticalIncidents {
			if inc.Severity != "critical" {
				allCritical = false
				testResult("Validate severity filter", fmt.Errorf("expected critical, got %s", inc.Severity))
				break
			}
		}
		if allCritical {
			testResult("Query by severity", nil)
		}
	}

	// Test 11: Query by service scope
	fmt.Println("\n=== Test 11: Query by Service Scope ===")
	serviceIncidents, err := provider.Query(ctx, schema.IncidentQuery{
		Scope: schema.QueryScope{Service: "api-service"},
	})
	if err != nil {
		testResult("Query by service scope", err)
	} else {
		fmt.Printf("Found %d incidents for service 'api-service'\n", len(serviceIncidents))
		allMatch := true
		for _, inc := range serviceIncidents {
			if inc.Service != "api-service" {
				allMatch = false
				testResult("Validate service filter", fmt.Errorf("expected api-service, got %s", inc.Service))
				break
			}
		}
		if allMatch {
			testResult("Query by service scope", nil)
		}
	}

	// Test 12: Query by environment metadata
	fmt.Println("\n=== Test 12: Query by Environment Scope ===")
	prodIncidents, err := provider.Query(ctx, schema.IncidentQuery{
		Scope: schema.QueryScope{Environment: "production"},
	})
	if err != nil {
		testResult("Query by environment scope", err)
	} else {
		fmt.Printf("Found %d incidents in production environment\n", len(prodIncidents))
		testResult("Query by environment scope", nil)
	}

	// Test 13: Query by team metadata
	fmt.Println("\n=== Test 13: Query by Team Scope ===")
	platformIncidents, err := provider.Query(ctx, schema.IncidentQuery{
		Scope: schema.QueryScope{Team: "platform"},
	})
	if err != nil {
		testResult("Query by team scope", err)
	} else {
		fmt.Printf("Found %d incidents for platform team\n", len(platformIncidents))
		testResult("Query by team scope", nil)
	}

	// Test 14: Query with combined filters
	fmt.Println("\n=== Test 14: Query with Combined Filters ===")
	combinedIncidents, err := provider.Query(ctx, schema.IncidentQuery{
		Statuses:   []string{"open", "acknowledged"},
		Severities: []string{"critical", "high"},
	})
	if err != nil {
		testResult("Query with combined filters", err)
	} else {
		fmt.Printf("Found %d incidents (open/acknowledged + critical/high severity)\n", len(combinedIncidents))
		allMatch := true
		for _, inc := range combinedIncidents {
			statusMatch := inc.Status == "open" || inc.Status == "acknowledged"
			severityMatch := inc.Severity == "critical" || inc.Severity == "high"
			if !statusMatch || !severityMatch {
				allMatch = false
				testResult("Validate combined filters", fmt.Errorf("incident %s doesn't match filters", inc.ID))
				break
			}
		}
		if allMatch {
			testResult("Query with combined filters", nil)
		}
	}

	// Test 15: Query with limit
	fmt.Println("\n=== Test 15: Query with Limit ===")
	limitedIncidents, err := provider.Query(ctx, schema.IncidentQuery{
		Limit: 2,
	})
	if err != nil {
		testResult("Query with limit", err)
	} else {
		if len(limitedIncidents) <= 2 {
			fmt.Printf("Correctly limited to %d incidents\n", len(limitedIncidents))
			testResult("Query with limit", nil)
		} else {
			testResult("Query with limit", fmt.Errorf("expected max 2 incidents, got %d", len(limitedIncidents)))
		}
	}

	// Test 16: Update incident status
	if incident1ID != "" {
		fmt.Println("\n=== Test 16: Update Incident Status ===")
		newStatus := "acknowledged"
		updatedIncident, err := provider.Update(ctx, incident1ID, schema.UpdateIncidentInput{
			Status: &newStatus,
		})
		if err != nil {
			testResult("Update incident status", err)
		} else {
			if updatedIncident.Status != "acknowledged" {
				testResult("Validate updated status", fmt.Errorf("expected acknowledged, got %s", updatedIncident.Status))
			} else {
				fmt.Printf("✅ Incident status updated to: %s\n", updatedIncident.Status)
				testResult("Update incident status", nil)
			}
		}

		// Test 17: Update incident title and severity
		fmt.Println("\n=== Test 17: Update Multiple Fields ===")
		newTitle := "Updated Integration Test Incident"
		newSeverity := "high"
		updatedIncident, err = provider.Update(ctx, incident1ID, schema.UpdateIncidentInput{
			Title:    &newTitle,
			Severity: &newSeverity,
		})
		if err != nil {
			testResult("Update multiple fields", err)
		} else {
			if updatedIncident.Title != newTitle || updatedIncident.Severity != newSeverity {
				testResult("Validate multiple field updates", fmt.Errorf("field update mismatch"))
			} else {
				fmt.Printf("✅ Updated title to: %s\n", updatedIncident.Title)
				fmt.Printf("✅ Updated severity to: %s\n", updatedIncident.Severity)
				testResult("Update multiple fields", nil)
			}
		}
	}

	// Test 18: GetTimeline for incident (should be empty initially)
	if incident1ID != "" {
		fmt.Println("\n=== Test 18: Get Incident Timeline ===")
		timeline, err := provider.GetTimeline(ctx, incident1ID)
		if err != nil {
			testResult("Get incident timeline", err)
		} else {
			fmt.Printf("Timeline has %d entries\n", len(timeline))
			testResult("Get incident timeline", nil)
		}

		// Test 19: Append timeline entry
		fmt.Println("\n=== Test 19: Append Timeline Entry ===")
		err = provider.AppendTimeline(ctx, incident1ID, schema.TimelineAppendInput{
			At:   time.Now(),
			Kind: "note",
			Body: "This is a test note added via the integration test.",
			Actor: map[string]any{
				"name": "integration-test",
			},
		})
		testResult("Append timeline note", err)

		// Test 20: Verify timeline was appended
		fmt.Println("\n=== Test 20: Verify Timeline Entry ===")
		timeline, err = provider.GetTimeline(ctx, incident1ID)
		if err != nil {
			testResult("Get timeline after append", err)
		} else if len(timeline) != 1 {
			testResult("Validate timeline entry count", fmt.Errorf("expected 1 entry, got %d", len(timeline)))
		} else {
			entry := timeline[0]
			fmt.Printf("Timeline entry:\n")
			fmt.Printf("  ID: %s\n", entry.ID)
			fmt.Printf("  Kind: %s\n", entry.Kind)
			fmt.Printf("  Body: %s\n", entry.Body)
			fmt.Printf("  At: %s\n", entry.At.Format("2006-01-02 15:04:05"))

			if entry.Body != "This is a test note added via the integration test." {
				testResult("Validate timeline entry body", fmt.Errorf("body mismatch"))
			} else {
				testResult("Verify timeline entry", nil)
			}
		}
	}

	// Test 21: Query by text search
	fmt.Println("\n=== Test 21: Query by Text Search ===")
	searchIncidents, err := provider.Query(ctx, schema.IncidentQuery{
		Query: "critical",
	})
	if err != nil {
		testResult("Query by text search", err)
	} else {
		fmt.Printf("Found %d incidents matching 'critical'\n", len(searchIncidents))
		testResult("Query by text search", nil)
	}

	// Print summary
	duration := time.Since(startTime)
	fmt.Println("\n=================================")
	fmt.Println("Test Summary")
	fmt.Println("=================================")
	fmt.Printf("Total Tests: %d\n", totalTests)
	fmt.Printf("Passed: %d ✅\n", passedTests)
	fmt.Printf("Failed: %d ❌\n", failedTests)
	fmt.Printf("Duration: %v\n", duration.Round(time.Millisecond))
	if totalTests > 0 {
		fmt.Printf("Success Rate: %.1f%%\n", float64(passedTests)/float64(totalTests)*100)
	}

	if failedTests == 0 {
		fmt.Println("\n✅ All tests passed successfully!")
	} else {
		fmt.Printf("\n⚠️  %d test(s) failed. Please review the output above.\n", failedTests)
	}
}
