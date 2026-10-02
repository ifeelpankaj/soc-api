package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-server/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type documentHandlerFake struct {
	society, user, id int64
	resident, receipt bool
	calls             int
	err               error
}

func (s *documentHandlerFake) Invoice(_ context.Context, society, user, id int64, resident bool) (models.MaintenancePDF, error) {
	s.society, s.user, s.id, s.resident = society, user, id, resident
	s.calls++
	return models.MaintenancePDF{Bytes: []byte("%PDF-test"), Filename: fmt.Sprintf("maintenance-invoice-%d.pdf", id)}, s.err
}
func (s *documentHandlerFake) Receipt(ctx context.Context, society, user, id int64, resident bool) (models.MaintenancePDF, error) {
	s.receipt = true
	return s.Invoice(ctx, society, user, id, resident)
}
func documentRouter(s *documentHandlerFake, authenticated bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if authenticated {
		r.Use(func(c *gin.Context) { c.Set("user_id", int64(7)) })
	}
	h := NewMaintenanceDocumentHandler(s)
	g := r.Group("/societies/:societyId/maintenance")
	g.GET("/bills/:id/invoice", h.InvoicePDF)
	g.GET("/my/bills/:id/invoice", h.MyInvoicePDF)
	g.GET("/payments/:id/receipt/pdf", h.ReceiptPDF)
	g.GET("/my/payments/:id/receipt/pdf", h.MyReceiptPDF)
	return r
}
func TestMaintenancePDFHandlers(t *testing.T) {
	for _, path := range []string{"bills/123/invoice", "my/bills/123/invoice", "payments/123/receipt/pdf", "my/payments/123/receipt/pdf"} {
		t.Run(path, func(t *testing.T) {
			s := &documentHandlerFake{}
			r := documentRouter(s, true)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/"+path, nil))
			if w.Code != 200 || w.Header().Get("Content-Type") != "application/pdf" || w.Header().Get("Cache-Control") != "private, no-store" || w.Header().Get("X-Content-Type-Options") != "nosniff" || !strings.Contains(w.Header().Get("Content-Disposition"), "attachment; filename=maintenance-invoice-123.pdf") {
				t.Fatalf("headers: %d %+v", w.Code, w.Header())
			}
			if s.society != 42 || s.user != 7 || s.id != 123 || s.resident != strings.HasPrefix(path, "my/") || s.receipt != strings.Contains(path, "receipt") {
				t.Fatalf("scope: %+v", s)
			}
			s.err = models.NewAppError("DOCUMENT_UNSUPPORTED_TEXT", "Unsupported text", 422, nil)
			w = httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/"+path, nil))
			if w.Code != 422 || !strings.Contains(w.Body.String(), "DOCUMENT_UNSUPPORTED_TEXT") || w.Header().Get("Content-Disposition") != "" || strings.Contains(w.Body.String(), "%PDF") {
				t.Fatalf("partial PDF on error: %d %s", w.Code, w.Body)
			}
			unauth := &documentHandlerFake{}
			w = httptest.NewRecorder()
			documentRouter(unauth, false).ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/"+path, nil))
			if w.Code != 401 || unauth.calls != 0 {
				t.Fatal("unauthenticated download")
			}
		})
	}
}
