package repository

import "go-server/pkg/database"

// financialRepository composes persistence capabilities that share a database
// and transaction context. Each capability owns its queries and row mappings.
type financialRepository struct {
	*billingRemindersRepository
	*billingReviewsRepository
	*collectionsRepository
	*financialAccessRepository
	*financialEventsRepository
	*idempotencyRepository
	*paymentClaimsRepository
	*paymentReportsRepository
	*paymentRequestsRepository
	*paymentSettingsRepository
}

func newFinancialRepository(d *database.Database) *financialRepository {
	return &financialRepository{
		billingRemindersRepository: &billingRemindersRepository{database: d},
		billingReviewsRepository:   &billingReviewsRepository{database: d},
		collectionsRepository:      &collectionsRepository{database: d},
		financialAccessRepository:  &financialAccessRepository{database: d},
		financialEventsRepository:  &financialEventsRepository{database: d},
		idempotencyRepository:      &idempotencyRepository{database: d},
		paymentClaimsRepository:    &paymentClaimsRepository{database: d},
		paymentReportsRepository:   &paymentReportsRepository{database: d},
		paymentRequestsRepository:  &paymentRequestsRepository{database: d},
		paymentSettingsRepository:  &paymentSettingsRepository{database: d},
	}
}
