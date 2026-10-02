package contracts

type FinancialStore interface {
	PaymentRequestsRepository
	BillingRemindersRepository
	FinancialEventsRepository
	CollectionsRepository
	BillingReviewsRepository
	PaymentClaimsRepository
	IdempotencyRepository
	PaymentReportsRepository
	PaymentSettingsRepository
	FinancialAccessRepository
}

type PaymentStore interface {
	FinancialStore
	TransactionManager
}
