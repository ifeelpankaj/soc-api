package testutil

import (
	"time"

	"go-server/internal/models"
)

func Ptr[T any](value T) *T {
	return &value
}

func ActiveSociety(id int64) *models.Society {
	return &models.Society{
		ID:          id,
		Name:        "Apna Gate Test Society",
		SocietyCode: "AGTEST",
		City:        Ptr("Pune"),
		State:       Ptr("Maharashtra"),
		Pincode:     Ptr("411045"),
		Country:     "India",
		Status:      models.SocietyStatusActive,
		CreatedBy:   1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func OccupiedFlat(societyID int64, flatID int64) *models.Flat {
	return &models.Flat{
		ID:         flatID,
		SocietyID:  societyID,
		FlatNumber: "A-101",
		Block:      Ptr("A"),
		Floor:      Ptr("1"),
		Status:     models.FlatStatusOccupied,
		IsActive:   true,
		CreatedBy:  Ptr[int64](1),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func VisitorEntry(societyID int64, flatID int64, status models.VisitorStatus) *models.VisitorEntry {
	return &models.VisitorEntry{
		ID:          100,
		SocietyID:   societyID,
		FlatID:      flatID,
		VisitorID:   200,
		Source:      models.VisitorEntrySourcePublicQR,
		Purpose:     models.VisitorPurposeGuest,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		QRExpiresAt: Ptr(time.Now().Add(time.Hour)),
		Visitor:     &models.VisitorSummary{FullName: "Test Visitor", PhoneNumber: Ptr("+919876543210")},
		Flat:        &models.VisitorFlatSummary{ID: flatID, FlatNumber: "A-101", Block: Ptr("A"), Floor: Ptr("1")},
	}
}
