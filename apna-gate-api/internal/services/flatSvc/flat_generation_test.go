package flatsvc

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"go-server/internal/models"
	repository "go-server/internal/repositories"
	societysvc "go-server/internal/services/societySvc"
)

func generationInt32(value int32) *int32    { return &value }
func generationString(value string) *string { return &value }

func TestGenerateFlatInputsContinuous(t *testing.T) {
	req := &models.GenerateFlatsRequest{
		Block: "A", NumberingMode: models.FlatNumberingModeContinuous, StartFloor: generationInt32(1),
		FlatsPerFloor: 3, SequenceStart: generationString("098"), TotalFlats: generationInt32(7),
	}

	items, err := generateFlatInputs(req)
	if err != nil {
		t.Fatalf("generateFlatInputs() error = %v", err)
	}
	wantNumbers := []string{"098", "099", "100", "101", "102", "103", "104"}
	wantFloors := []string{"1", "1", "1", "2", "2", "2", "3"}
	for index, item := range items {
		if item.FlatNumber != wantNumbers[index] || item.Floor == nil || *item.Floor != wantFloors[index] {
			t.Fatalf("item %d = %q floor %v, want %q floor %q", index, item.FlatNumber, item.Floor, wantNumbers[index], wantFloors[index])
		}
	}
}

func TestGenerateFlatInputsFloorBased(t *testing.T) {
	req := &models.GenerateFlatsRequest{
		Block: "B", NumberingMode: models.FlatNumberingModeFloorBased, StartFloor: generationInt32(0),
		FlatsPerFloor: 3, UnitStart: generationString("01"), NumberOfFloors: generationInt32(3),
	}

	items, err := generateFlatInputs(req)
	if err != nil {
		t.Fatalf("generateFlatInputs() error = %v", err)
	}
	wantNumbers := []string{"001", "002", "003", "101", "102", "103", "201", "202", "203"}
	wantFloors := []string{"0", "0", "0", "1", "1", "1", "2", "2", "2"}
	gotNumbers := make([]string, 0, len(items))
	gotFloors := make([]string, 0, len(items))
	for _, item := range items {
		gotNumbers = append(gotNumbers, item.FlatNumber)
		gotFloors = append(gotFloors, *item.Floor)
	}
	if !reflect.DeepEqual(gotNumbers, wantNumbers) || !reflect.DeepEqual(gotFloors, wantFloors) {
		t.Fatalf("generated = %v floors %v, want %v floors %v", gotNumbers, gotFloors, wantNumbers, wantFloors)
	}
}

func TestGenerateFlatInputsMultiDigitFloor(t *testing.T) {
	req := &models.GenerateFlatsRequest{
		Block: "C", NumberingMode: models.FlatNumberingModeFloorBased, StartFloor: generationInt32(12),
		FlatsPerFloor: 2, UnitStart: generationString("01"), NumberOfFloors: generationInt32(1),
	}
	items, err := generateFlatInputs(req)
	if err != nil {
		t.Fatalf("generateFlatInputs() error = %v", err)
	}
	if items[0].FlatNumber != "1201" || *items[0].Floor != "12" {
		t.Fatalf("first item = %q floor %q", items[0].FlatNumber, *items[0].Floor)
	}
}

func TestGenerateFlatInputsRejectsFixedSuffixOverflow(t *testing.T) {
	req := &models.GenerateFlatsRequest{
		Block: "A", NumberingMode: models.FlatNumberingModeFloorBased,
		FlatsPerFloor: 3, UnitStart: generationString("98"), NumberOfFloors: generationInt32(1),
	}
	_, err := generateFlatInputs(req)
	if err == nil || !strings.Contains(err.Error(), "fixed unit_start width") {
		t.Fatalf("generateFlatInputs() error = %v", err)
	}
}

func TestValidateUniqueFlatInputsNormalizesBlockAndPreservesCase(t *testing.T) {
	blockA := " A "
	blockANormalized := "A"
	blockLower := "a"
	if err := validateUniqueFlatInputs([]flatCreateInput{
		{Block: &blockA, FlatNumber: "001"},
		{Block: &blockLower, FlatNumber: "001"},
	}); err != nil {
		t.Fatalf("case-sensitive blocks should not conflict: %v", err)
	}
	if err := validateUniqueFlatInputs([]flatCreateInput{
		{Block: &blockA, FlatNumber: "001"},
		{Block: &blockANormalized, FlatNumber: "001"},
	}); !errors.Is(err, ErrFlatConflict) {
		t.Fatalf("normalized duplicate error = %v, want ErrFlatConflict", err)
	}
}

type creationTransactionKey struct{}

type creationTransactionManager struct {
	called    bool
	committed bool
}

func (m *creationTransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	m.called = true
	err := fn(context.WithValue(ctx, creationTransactionKey{}, true))
	m.committed = err == nil
	return err
}

type creationSocietyService struct {
	societysvc.SocietyService
}

func (s *creationSocietyService) EnsureRole(context.Context, int64, int64, ...string) error {
	return nil
}

type creationSubscriptionGuard struct {
	lockedInsideTransaction bool
	adding                  int64
}

func (g *creationSubscriptionGuard) EnsureSocietyOperational(context.Context, int64) error {
	return nil
}
func (g *creationSubscriptionGuard) CanAddResident(context.Context, int64, int64) error {
	return nil
}
func (g *creationSubscriptionGuard) CanAddResidentWithLock(context.Context, int64, int64) error {
	return nil
}
func (g *creationSubscriptionGuard) CanAddFlatWithLock(ctx context.Context, _ int64, adding int64) error {
	g.lockedInsideTransaction, _ = ctx.Value(creationTransactionKey{}).(bool)
	g.adding = adding
	return nil
}

type creationFlatRepository struct {
	repository.FlatRepository
	createdInsideTransaction bool
	items                    map[int64]*models.Flat
	nextID                   int64
}

func (r *creationFlatRepository) Create(ctx context.Context, flat *models.Flat) error {
	r.createdInsideTransaction, _ = ctx.Value(creationTransactionKey{}).(bool)
	r.nextID++
	flat.ID = r.nextID
	r.items[flat.ID] = flat
	return nil
}

func (r *creationFlatRepository) Get(_ context.Context, filter *models.FlatFilter) (*models.Flat, error) {
	return r.items[*filter.ID], nil
}

type creationVisitorSettings struct {
	insideTransaction bool
	err               error
}

func (s *creationVisitorSettings) CreateDefaultFlatSettings(ctx context.Context, _ int64, _ int64, _ int64) error {
	s.insideTransaction, _ = ctx.Value(creationTransactionKey{}).(bool)
	return s.err
}

func TestCreateManyFlatsKeepsQuotaAndWritesInOneTransaction(t *testing.T) {
	tx := &creationTransactionManager{}
	quota := &creationSubscriptionGuard{}
	repo := &creationFlatRepository{items: make(map[int64]*models.Flat)}
	visitorSettings := &creationVisitorSettings{}
	svc := &FlatSvc{
		flatRepo: repo, txManager: tx, societySvc: &creationSocietyService{},
		subscriptionSvc: quota, visitorSettingSvc: visitorSettings,
	}
	block := "A"
	floor := "1"

	result, err := svc.createManyFlats(context.Background(), 42, 7, []flatCreateInput{
		{Block: &block, Floor: &floor, FlatNumber: "001"},
		{Block: &block, Floor: &floor, FlatNumber: "002"},
	})
	if err != nil {
		t.Fatalf("createManyFlats() error = %v", err)
	}
	if result.Total != 2 || !tx.called || !tx.committed {
		t.Fatalf("result/transaction = total %d, called %v, committed %v", result.Total, tx.called, tx.committed)
	}
	if !quota.lockedInsideTransaction || quota.adding != 2 || !repo.createdInsideTransaction || !visitorSettings.insideTransaction {
		t.Fatalf("transaction participation: quota=%v adding=%d repo=%v settings=%v", quota.lockedInsideTransaction, quota.adding, repo.createdInsideTransaction, visitorSettings.insideTransaction)
	}
}

func TestCreateManyFlatsRollsBackOnVisitorSettingsFailure(t *testing.T) {
	tx := &creationTransactionManager{}
	repo := &creationFlatRepository{items: make(map[int64]*models.Flat)}
	visitorErr := errors.New("visitor defaults failed")
	svc := &FlatSvc{
		flatRepo: repo, txManager: tx, societySvc: &creationSocietyService{},
		subscriptionSvc: &creationSubscriptionGuard{}, visitorSettingSvc: &creationVisitorSettings{err: visitorErr},
	}
	block := "A"
	floor := "1"

	_, err := svc.createManyFlats(context.Background(), 42, 7, []flatCreateInput{{Block: &block, Floor: &floor, FlatNumber: "001"}})
	if !errors.Is(err, visitorErr) || tx.committed {
		t.Fatalf("error/commit = %v/%v, want visitor error and rollback", err, tx.committed)
	}
}
