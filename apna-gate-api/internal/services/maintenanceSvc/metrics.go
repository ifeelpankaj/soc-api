package maintenancesvc

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

var financialOperations = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "apna_gate_maintenance_operations_total", Help: "Maintenance financial commands and billing runs by bounded operation and outcome."}, []string{"operation", "outcome"})
var pendingClaims = prometheus.NewGauge(prometheus.GaugeOpts{Name: "apna_gate_maintenance_pending_claims", Help: "Pending claims observed at the latest job execution."})

func init() { prometheus.MustRegister(financialOperations, pendingClaims) }
func observeFinancial(operation string, err error) {
	parts := strings.Split(operation, "/")
	if len(parts) > 2 {
		parts = parts[:2]
	}
	operation = strings.Join(parts, "/")
	outcome := "success"
	if err != nil {
		outcome = "failure"
	}
	financialOperations.WithLabelValues(operation, outcome).Inc()
}
