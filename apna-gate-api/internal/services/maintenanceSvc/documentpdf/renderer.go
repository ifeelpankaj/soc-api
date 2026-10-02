// Package documentpdf renders only permission-filtered document data. It never
// queries live configuration, creates payment requests, or stores documents.
package documentpdf

import (
	"context"
	"embed"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/signintech/gopdf"
	"go-server/internal/models"
)

// Fonts and their redistribution license travel together in the API binary.
//
//go:embed fonts/*.ttf fonts/LICENSE.txt
var fontFiles embed.FS

const (
	margin           = 44.0
	pageWidth        = 595.28
	contentWidth     = pageWidth - 2*margin
	bottom           = 776.0
	lineHeight       = 16.0
	descriptionWidth = 307.0
)

type Renderer struct{}

func New() *Renderer { return &Renderer{} }

// Money never converts the monetary value to floating point, even for int64 max.
func Money(paise int64) string {
	digits := strconv.FormatInt(paise, 10)
	sign := ""
	if strings.HasPrefix(digits, "-") {
		sign = "-"
		digits = digits[1:]
	}
	if len(digits) < 3 {
		digits = strings.Repeat("0", 3-len(digits)) + digits
	}
	whole, cents := digits[:len(digits)-2], digits[len(digits)-2:]
	for i := len(whole) - 3; i > 0; i -= 3 {
		whole = whole[:i] + "," + whole[i:]
	}
	return "INR " + sign + whole + "." + cents
}

type document struct {
	pdf              gopdf.GoPdf
	ctx              context.Context
	title, reference string
	y                float64
	err              error
}

func unsupportedText() error {
	return models.NewAppError("DOCUMENT_UNSUPPORTED_TEXT", "PDF documents currently support English/Latin-script text only; a document field contains unsupported text", 422, nil)
}
func validateText(text string) error {
	if !utf8.ValidString(text) {
		return unsupportedText()
	}
	for _, r := range text {
		if r == '\n' || r == '\t' || r == '\r' || (r >= 32 && r <= 126) || unicode.In(r, unicode.Latin) || unicode.IsMark(r) || r == '\u00a0' || strings.ContainsRune("‘’“”–—•…₹", r) {
			continue
		}
		return unsupportedText()
	}
	return nil
}
func newDocument(ctx context.Context, title, reference string, created time.Time) *document {
	d := &document{ctx: ctx, title: title, reference: reference}
	regularFont, err := fontFiles.ReadFile("fonts/NotoSans-Regular.ttf")
	if err != nil {
		d.err = err
		return d
	}
	boldFont, err := fontFiles.ReadFile("fonts/NotoSans-Bold.ttf")
	if err != nil {
		d.err = err
		return d
	}
	d.pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	missing := func(r rune) { d.err = unsupportedText() }
	if err := d.pdf.AddTTFFontDataWithOption("regular", regularFont, gopdf.TtfOption{OnGlyphNotFound: missing}); err != nil {
		d.err = err
		return d
	}
	if err := d.pdf.AddTTFFontDataWithOption("bold", boldFont, gopdf.TtfOption{OnGlyphNotFound: missing}); err != nil {
		d.err = err
		return d
	}
	// Metadata contains no personal data. Invoice metadata remains stable after payment.
	d.pdf.SetInfo(gopdf.PdfInfo{Title: title, Author: "", Creator: "Apna Gate", Producer: "Apna Gate", CreationDate: created})
	d.newPage()
	return d
}
func (d *document) check() bool {
	if d.err == nil {
		d.err = d.ctx.Err()
	}
	return d.err == nil
}
func (d *document) font(bold bool, size float64) {
	if !d.check() {
		return
	}
	family := "regular"
	if bold {
		family = "bold"
	}
	if err := d.pdf.SetFont(family, "", size); err != nil {
		d.err = err
	}
}
func (d *document) text(x, y, width float64, text string, bold bool, size float64, right bool) {
	if !d.check() {
		return
	}
	if err := validateText(text); err != nil {
		d.err = err
		return
	}
	d.font(bold, size)
	if !d.check() {
		return
	}
	if right {
		w, err := d.pdf.MeasureTextWidth(text)
		if err != nil {
			d.err = err
			return
		}
		x += width - w
	}
	d.pdf.SetXY(x, y)
	if err := d.pdf.Cell(&gopdf.Rect{W: width, H: lineHeight}, text); err != nil {
		d.err = err
	}
}
func (d *document) lines(text string, width float64, bold bool, size float64) []string {
	if !d.check() {
		return nil
	}
	if err := validateText(text); err != nil {
		d.err = err
		return nil
	}
	d.font(bold, size)
	var result []string
	for _, part := range strings.Split(strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n"), "\n") {
		part = strings.ReplaceAll(part, "\t", " ")
		if part == "" {
			result = append(result, "")
			continue
		}
		lines, err := d.pdf.SplitTextWithWordWrap(part, width)
		if err != nil {
			d.err = err
			return nil
		}
		result = append(result, lines...)
	}
	return result
}
func (d *document) newPage() {
	if !d.check() {
		return
	}
	if d.pdf.GetNumberOfPages() >= 100 {
		d.err = models.NewAppError("DOCUMENT_TOO_LARGE", "Document exceeds the supported page count", 422, nil)
		return
	}
	d.pdf.AddPage()
	d.pdf.SetFillColor(235, 72, 8)
	d.pdf.RectFromUpperLeftWithStyle(margin, 34, 32, 32, "F")
	d.pdf.SetTextColor(255, 255, 255)
	d.text(margin+7, 39, 20, "AG", true, 13, false)
	d.pdf.SetTextColor(17, 28, 51)
	d.text(margin+43, 32, contentWidth-43, "Apna Gate", true, 22, false)
	d.pdf.SetTextColor(100, 112, 132)
	d.text(margin+44, 59, contentWidth-44, "SOCIETY MANAGEMENT", false, 8, false)
	d.pdf.SetTextColor(235, 72, 8)
	d.text(margin, 98, contentWidth, strings.ToUpper(d.title), true, 23, false)
	d.pdf.SetTextColor(100, 112, 132)
	d.text(margin, 132, contentWidth, d.reference, false, 10, false)
	d.pdf.SetStrokeColor(201, 210, 220)
	d.pdf.SetLineWidth(0.7)
	d.pdf.Line(margin, 157, pageWidth-margin, 157)
	d.y = 175
}
func (d *document) section(title string) {
	d.ensure(34)
	d.pdf.SetFillColor(250, 239, 231)
	d.pdf.RectFromUpperLeftWithStyle(margin, d.y, contentWidth, 29, "F")
	d.pdf.SetTextColor(219, 63, 5)
	d.text(margin+12, d.y+5, contentWidth-24, strings.ToUpper(title), true, 11, false)
	d.pdf.SetTextColor(28, 40, 57)
	d.y += 35
}
func (d *document) ensure(height float64) {
	if d.y+height > bottom {
		d.newPage()
	}
}
func (d *document) paragraph(text string, bold bool, size float64) {
	for _, line := range d.lines(text, contentWidth, bold, size) {
		d.ensure(lineHeight)
		d.text(margin, d.y, contentWidth, line, bold, size, false)
		d.y += lineHeight
	}
}
func (d *document) field(label, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	d.paragraph(label+": "+value, false, 10)
}
func (d *document) issuer(i models.MaintenanceIssuer, zone *time.Location, notice bool) {
	d.paragraph(i.Name, true, 14)
	d.field("Society", i.SocietyCode)
	for _, v := range []string{i.AddressLine1, i.AddressLine2, i.Landmark} {
		if strings.TrimSpace(v) != "" {
			d.paragraph(v, false, 10)
		}
	}
	var locality []string
	for _, v := range []string{i.City, i.State, i.Pincode, i.Country} {
		if strings.TrimSpace(v) != "" {
			locality = append(locality, v)
		}
	}
	if len(locality) > 0 {
		d.paragraph(strings.Join(locality, ", "), false, 10)
	}
	if notice && i.CaptureSource == "rollout" {
		d.y += 5
		d.paragraph("Society details captured on "+i.CapturedAt.In(zone).Format("02 Jan 2006")+" after bill issuance; historical issuer details were unavailable.", false, 9)
	}
	d.y += 14
}
func (d *document) flat(b models.MaintenanceBill) {
	value := b.Flat.FlatNumber
	if b.Flat.Block != nil && *b.Flat.Block != "" {
		value = *b.Flat.Block + " / " + value
	}
	d.field("Flat", value)
}
func (d *document) tableHeader() {
	d.ensure(30)
	d.pdf.SetFillColor(250, 228, 214)
	d.pdf.RectFromUpperLeftWithStyle(margin, d.y, contentWidth, 25, "F")
	d.text(margin+8, d.y+4, descriptionWidth-16, "DESCRIPTION", true, 10, false)
	d.text(margin+descriptionWidth, d.y+4, contentWidth-descriptionWidth-8, "AMOUNT (INR)", true, 10, true)
	d.y += 31
}
func (d *document) items(items []models.MaintenanceItem, total int64) {
	d.tableHeader()
	for _, item := range items {
		lines := d.lines(item.Description, descriptionWidth-16, false, 10)
		if len(lines) == 0 {
			lines = []string{""}
		}
		for i, line := range lines {
			if d.y+lineHeight+6 > bottom {
				d.newPage()
				d.tableHeader()
			}
			d.text(margin+8, d.y, descriptionWidth-16, line, false, 10, false)
			if i == 0 {
				d.text(margin+descriptionWidth, d.y, contentWidth-descriptionWidth-8, Money(item.AmountPaise), false, 10, true)
			}
			d.y += lineHeight
		}
		d.y += 6
		d.pdf.SetStrokeColor(225, 230, 236)
		d.pdf.Line(margin, d.y, pageWidth-margin, d.y)
		d.y += 7
	}
	d.ensure(40)
	d.y += 8
	d.pdf.SetFillColor(253, 240, 232)
	d.pdf.RectFromUpperLeftWithStyle(margin, d.y-5, contentWidth, 31, "F")
	d.pdf.SetTextColor(210, 63, 7)
	d.text(margin+8, d.y, descriptionWidth-16, "TOTAL", true, 13, false)
	amountSize := 13.0
	if len(Money(total)) > 18 {
		amountSize = 9
	}
	d.text(margin+descriptionWidth, d.y, contentWidth-descriptionWidth-8, Money(total), true, amountSize, true)
	d.pdf.SetTextColor(28, 40, 57)
	d.y += 35
}
func (d *document) finish() ([]byte, error) {
	if !d.check() {
		return nil, d.err
	}
	pages := d.pdf.GetNumberOfPages()
	for p := 1; p <= pages; p++ {
		if err := d.pdf.SetPage(p); err != nil {
			return nil, err
		}
		d.pdf.SetTextColor(95, 105, 119)
		d.text(margin, 802, contentWidth, "Generated by Apna Gate", false, 8, false)
		d.text(margin, 802, contentWidth, fmt.Sprintf("Page %d of %d", p, pages), false, 8, true)
	}
	if !d.check() {
		return nil, d.err
	}
	data, err := d.pdf.GetBytesPdfReturnErr()
	if d.err != nil {
		return nil, d.err
	}
	return data, err
}

func documentZone(b models.MaintenanceBill) (*time.Location, error) {
	if b.ID <= 0 || b.Issuer.Name == "" || b.Issuer.CapturedAt.IsZero() || b.CreatedAt.IsZero() {
		return nil, fmt.Errorf("issued bill document snapshot is incomplete")
	}
	return time.LoadLocation(b.Timezone)
}
func (r *Renderer) Invoice(ctx context.Context, b models.MaintenanceBill) ([]byte, error) {
	zone, err := documentZone(b)
	if err != nil {
		return nil, err
	}
	d := newDocument(ctx, "Maintenance Invoice", "Invoice no.  "+b.BillNumber, b.CreatedAt)
	if !d.check() {
		return nil, d.err
	}
	d.section("Society")
	d.issuer(b.Issuer, zone, false)
	d.section("Invoice details")
	d.flat(b)
	if month, err := time.Parse("2006-01", b.BillingMonth); err == nil {
		d.field("Billing month", month.Format("January 2006"))
		d.field("Billing period", month.Format("02 Jan 2006")+" - "+month.AddDate(0, 1, -1).Format("02 Jan 2006"))
	}
	d.field("Issue date", b.CreatedAt.In(zone).Format("02 Jan 2006"))
	d.field("Due date", b.DueDate)
	d.y += 15
	d.section("Charges")
	d.items(b.Items, b.TotalPaise)
	d.section("Payment instructions")
	d.paragraph("Pay this amount using the payment details shared by your society. After payment, submit the reference in the app.", false, 10)
	d.y += 10
	d.pdf.SetTextColor(219, 63, 5)
	d.paragraph("IMPORTANT NOTES", true, 10)
	d.pdf.SetTextColor(28, 40, 57)
	d.paragraph("This is a system generated invoice. Please pay on or before the due date. For questions, contact your society management.", false, 9)
	d.paragraph("Not a payment receipt. This invoice is void once payment is verified and a receipt is issued.", false, 9)
	if b.IsCatchUp {
		d.paragraph("Missed-month bill: issued using current pricing and flat information.", false, 10)
	}
	if b.Outstanding != nil {
		d.y += 12
		d.paragraph("Outstanding across bills (information only)", true, 12)
		d.field("Calculated at", b.Outstanding.CalculatedAt.In(zone).Format("02 Jan 2006 15:04:05 MST"))
		d.field("Current month outstanding", Money(b.Outstanding.CurrentMonthPaise))
		d.field("Previous months outstanding", Money(b.Outstanding.PreviousOutstandingPaise))
		d.field("Total outstanding across bills", Money(b.Outstanding.TotalOutstandingPaise))
		d.paragraph("Previous dues are not included again in this bill's charges. Pay each outstanding bill separately.", false, 10)
	}
	return d.finish()
}
func (r *Renderer) Receipt(ctx context.Context, b models.MaintenanceBill, p models.UPIPayment) ([]byte, error) {
	zone, err := documentZone(b)
	if err != nil {
		return nil, err
	}
	if p.BillID != b.ID || p.SocietyID != b.SocietyID || p.AmountPaise != b.TotalPaise || (p.Status != "verified" && p.Status != "reversed") {
		return nil, fmt.Errorf("invalid verified payment document")
	}
	d := newDocument(ctx, "Payment Receipt", "Receipt ID  "+p.ReceiptNumber, p.VerifiedAt)
	if !d.check() {
		return nil, d.err
	}
	if p.Status == "reversed" {
		d.pdf.SetTextColor(170, 30, 30)
		d.paragraph("REVERSED", true, 16)
		d.pdf.SetTextColor(28, 40, 57)
		d.paragraph("This verification was reversed. This document is not evidence of an active settlement. Reversal does not constitute a refund or move money.", false, 10)
		d.y += 12
	}
	d.section("Society information")
	d.issuer(b.Issuer, zone, true)
	d.section("Bill details")
	d.flat(b)
	d.field("Bill reference", p.BillNumber)
	d.field("Billing month", b.BillingMonth)
	d.y += 8
	d.section("Payment amount and status")
	d.pdf.SetTextColor(219, 63, 5)
	d.paragraph("Amount: "+Money(p.AmountPaise), true, 14)
	d.pdf.SetTextColor(28, 40, 57)
	d.field("Status", strings.ToUpper(p.Status))
	d.y += 8
	d.section("Payment verification details")
	d.field("Bank credit date", p.CreditDate)
	d.field("Verified at", p.VerifiedAt.In(zone).Format("02 Jan 2006 15:04:05 MST"))
	d.field("Payee", p.Destination.PayeeName)
	d.field("UPI destination", p.Destination.UPIID)
	d.field("Destination version", strconv.FormatInt(p.SettingsVersion, 10))
	if p.Reference != nil {
		d.field("Bank reference", *p.Reference)
	} else {
		d.field("Bank reference", "Restricted to the payer and society administrators")
	}
	d.field("Verified by", p.VerifiedByName)
	if p.ReversedAt != nil {
		d.field("Reversed at", p.ReversedAt.In(zone).Format("02 Jan 2006 15:04:05 MST"))
	}
	if p.ReversalReason != nil {
		d.field("Reversal reason", *p.ReversalReason)
	}
	d.y += 12
	d.paragraph("Records society-admin verification of bank credit. Apna Gate does not collect or transfer these funds.", false, 9)
	return d.finish()
}
