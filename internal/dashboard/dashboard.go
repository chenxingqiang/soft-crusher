// File: internal/dashboard/dashboard.go

package dashboard

import (
	"fmt"
	"sync"
	"time"
)

// DashboardData represents the data structure for the dashboard
type DashboardData struct {
	ActiveUsers    int
	APIRequests    int
	ErrorRate      float64
	LastUpdateTime time.Time
}

var (
	data      DashboardData
	dataMutex sync.RWMutex
)

// Init initializes the dashboard
func Init() {
	fmt.Println("Initializing dashboard...")
	updateDashboard()
	go periodicallyUpdateDashboard()
}

// GetDashboardData returns the current dashboard data
func GetDashboardData() DashboardData {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	return data
}

func updateDashboard() {
	// In a real application, this would fetch data from various sources
	dataMutex.Lock()
	defer dataMutex.Unlock()

	data.ActiveUsers = 100  // Example value
	data.APIRequests = 1000 // Example value
	data.ErrorRate = 0.01   // Example value
	data.LastUpdateTime = time.Now()
}

func periodicallyUpdateDashboard() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		updateDashboard()
	}
}
