package jobs

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	OutcomeSuccess     = "success"
	OutcomeFailure     = "failure"
	OutcomeSkipped     = "skipped"
	OutcomeInterrupted = "interrupted"
)

var (
	jobRunsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "apna_gate", Subsystem: "jobs", Name: "runs_total",
		Help: "Total completed background job runs by outcome.",
	}, []string{"job_name", "outcome"})
	jobDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "apna_gate", Subsystem: "jobs", Name: "duration_seconds",
		Help: "Background job run duration in seconds.", Buckets: prometheus.DefBuckets,
	}, []string{"job_name", "outcome"})
	jobInProgress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "apna_gate", Subsystem: "jobs", Name: "in_progress",
		Help: "Background jobs currently executing.",
	}, []string{"job_name"})
	jobLastCompletion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "apna_gate", Subsystem: "jobs", Name: "last_completion_timestamp_seconds",
		Help: "Unix timestamp of the most recent completed background job run.",
	}, []string{"job_name"})
	jobLastSuccess = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "apna_gate", Subsystem: "jobs", Name: "last_success_timestamp_seconds",
		Help: "Unix timestamp of the most recent successful background job run.",
	}, []string{"job_name"})
)

func init() {
	prometheus.MustRegister(jobRunsTotal, jobDurationSeconds, jobInProgress, jobLastCompletion, jobLastSuccess)
	// Known jobs exist before their first invocation. Zero timestamps mean
	// "never observed in this process", not an epoch-age success.
	for _, name := range []string{"expiry", JobCleanup, JobMonthlyVisitorReport, JobMaintenanceBilling, JobMaintenanceReminders} {
		jobInProgress.WithLabelValues(name).Set(0)
		jobLastCompletion.WithLabelValues(name).Set(0)
		jobLastSuccess.WithLabelValues(name).Set(0)
		for _, outcome := range []string{OutcomeSuccess, OutcomeFailure, OutcomeSkipped, OutcomeInterrupted} {
			jobRunsTotal.WithLabelValues(name, outcome).Add(0)
		}
	}
}

// BeginRun marks a job as active and returns a completion callback.
func BeginRun(name string) func(outcome string) {
	started := time.Now()
	jobInProgress.WithLabelValues(name).Inc()
	return func(outcome string) {
		jobInProgress.WithLabelValues(name).Dec()
		jobRunsTotal.WithLabelValues(name, outcome).Inc()
		jobDurationSeconds.WithLabelValues(name, outcome).Observe(time.Since(started).Seconds())
		jobLastCompletion.WithLabelValues(name).SetToCurrentTime()
		if outcome == OutcomeSuccess {
			jobLastSuccess.WithLabelValues(name).SetToCurrentTime()
		}
	}
}
