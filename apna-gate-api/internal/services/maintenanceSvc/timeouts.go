package maintenancesvc

import "time"

const (
	FinancialTimeout = time.Minute
	PushTimeout      = 30 * time.Second
	PreviewLifetime  = 15 * time.Minute
)
