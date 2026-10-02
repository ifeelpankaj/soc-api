package models

import (
	"fmt"
	"time"
)

// PopulatePresentation keeps calendar arithmetic on the server using the
// frozen bill timezone. Clients render these values without inferring payment.
func (b *MaintenanceBill) PopulatePresentation(now time.Time) error {
	loc, err := time.LoadLocation(b.Timezone)
	if err != nil {
		return err
	}
	today, err := time.Parse("2006-01-02", now.In(loc).Format("2006-01-02"))
	if err != nil {
		return err
	}
	due, err := time.Parse("2006-01-02", b.DueDate)
	if err != nil {
		return err
	}
	days := int(due.Sub(today).Hours() / 24)
	b.Status = "unpaid"
	if b.OutstandingAmountPaise == 0 && b.PaidAmountPaise == b.TotalPaise {
		b.Status = "paid"
	} else if days < 0 {
		b.Status = "overdue"
	}
	b.DisplayStatus = b.Status
	if b.Status != "paid" {
		switch b.PaymentClaimStatus {
		case "pending":
			b.DisplayStatus = "pending_review"
		case "rejected":
			b.DisplayStatus = "rejected"
		}
	}
	switch {
	case b.Status == "paid":
		b.DueMessage = "Payment verified"
	case days == 0:
		b.DueMessage = "Due today"
	case days == 1:
		b.DueMessage = "1 day left"
	case days > 0:
		b.DueMessage = fmt.Sprintf("%d days left", days)
	case days == -1:
		b.DueMessage = "1 day overdue"
	default:
		b.DueMessage = fmt.Sprintf("%d days overdue", -days)
	}
	return nil
}
