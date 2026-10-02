package jobs

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestBeginRunRecordsSuccessMetrics(t *testing.T) {
	const name = "metrics-test"

	complete := BeginRun(name)
	complete(OutcomeSuccess)

	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	for _, family := range metricFamilies {
		if family.GetName() != "apna_gate_jobs_runs_total" {
			continue
		}
		for _, metric := range family.Metric {
			var job, outcome string
			for _, label := range metric.Label {
				switch label.GetName() {
				case "job_name":
					job = label.GetValue()
				case "outcome":
					outcome = label.GetValue()
				}
			}
			if job == name && outcome == OutcomeSuccess && metric.GetCounter().GetValue() > 0 {
				return
			}
		}
	}
	t.Fatal("did not find the successful metrics-test job run")
}
