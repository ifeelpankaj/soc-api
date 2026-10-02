package imagesvc

import "github.com/prometheus/client_golang/prometheus"

var imageDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "apna_gate_image_operation_seconds", Help: "Image API operation duration", Buckets: []float64{.05, .1, .5, 1, 2, 5, 10, 20, 45},
}, []string{"operation"})
var imageFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "apna_gate_image_failures_total", Help: "Image API failures",
}, []string{"operation"})
var cleanupFailures = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "apna_gate_image_cleanup_failures_total", Help: "Image cleanup operations requiring reconciliation",
})
var imageStageDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "apna_gate_image_stage_seconds", Help: "Image operation stage duration", Buckets: []float64{.001, .005, .01, .05, .1, .5, 1, 2, 5, 10, 20},
}, []string{"stage"})
var imageNormalization = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "apna_gate_image_normalization_total", Help: "Visitor image normalization path",
}, []string{"path"})

func init() {
	prometheus.MustRegister(imageDuration, imageFailures, cleanupFailures, imageStageDuration, imageNormalization)
}
