package visitorentrysvc

import (
	"context"
	"time"

	"go-server/internal/models"
	service "go-server/internal/services"
)

const visitorStatsTimezone = "Asia/Kolkata"

func (s *VisitorEntrySvc) GetDailyEntryStats(ctx context.Context, societyID int64, days int32) (*models.VisitorEntryDailyStatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if days <= 0 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	rows, err := s.entryRepo.GetDailyStatsCreated(ctx, societyID, days)
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64, len(rows))
	for _, row := range rows {
		counts[row.Date] = row.Count
	}

	loc, err := time.LoadLocation(visitorStatsTimezone)
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	daily := make([]models.VisitorDailyCountResponse, 0, days)
	var total int64
	for offset := int(days) - 1; offset >= 0; offset-- {
		date := now.AddDate(0, 0, -offset).Format("2006-01-02")
		count := counts[date]
		total += count
		daily = append(daily, models.VisitorDailyCountResponse{
			Date:  date,
			Count: count,
		})
	}

	return &models.VisitorEntryDailyStatsResponse{
		Days:     days,
		Metric:   "created",
		Timezone: visitorStatsTimezone,
		Daily:    daily,
		Total:    total,
	}, nil
}
