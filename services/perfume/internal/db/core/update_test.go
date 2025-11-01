package core

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/constants"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

func TestGetSavepointQuery(t *testing.T) {
	if got := getSavepointQuery("SAVEPOINT sp_", 3); got != "SAVEPOINT sp_3" {
		t.Fatalf("getSavepointQuery() = %q, want %q", got, "SAVEPOINT sp_3")
	}
	if got := getSavepointQuery("ROLLBACK TO SAVEPOINT sp_", 0); got != "ROLLBACK TO SAVEPOINT sp_0" {
		t.Fatalf("getSavepointQuery() = %q, want %q", got, "ROLLBACK TO SAVEPOINT sp_0")
	}
}

type execCall struct {
	sql  string
	args []any
}

type mockTx struct {
	execCalls []execCall
	execErrs  map[string]error
}

func newMockTx() *mockTx {
	return &mockTx{execErrs: make(map[string]error)}
}

func (m *mockTx) setExecError(sql string, err error) {
	m.execErrs[sql] = err
}

func (m *mockTx) Exec(_ context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.execCalls = append(m.execCalls, execCall{sql: sql, args: arguments})
	if err, ok := m.execErrs[sql]; ok {
		return pgconn.NewCommandTag(""), err
	}
	return pgconn.NewCommandTag("EXEC 0"), nil
}

func (m *mockTx) Begin(context.Context) (pgx.Tx, error) { panic("not implemented") }
func (m *mockTx) Commit(context.Context) error          { panic("not implemented") }
func (m *mockTx) Rollback(context.Context) error        { panic("not implemented") }
func (m *mockTx) CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error) {
	panic("not implemented")
}
func (m *mockTx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults { panic("not implemented") }
func (m *mockTx) LargeObjects() pgx.LargeObjects                         { panic("not implemented") }
func (m *mockTx) Prepare(context.Context, string, string) (*pgconn.StatementDescription, error) {
	panic("not implemented")
}
func (m *mockTx) Query(context.Context, string, ...any) (pgx.Rows, error) { panic("not implemented") }
func (m *mockTx) QueryRow(context.Context, string, ...any) pgx.Row        { panic("not implemented") }
func (m *mockTx) Conn() *pgx.Conn                                         { return nil }

func TestTruncateSuccess(t *testing.T) {
	tx := newMockTx()
	ctx := context.Background()

	if ok := truncate(ctx, tx); !ok {
		t.Fatalf("truncate() = %v, want %v", ok, true)
	}

	if len(tx.execCalls) != 1 {
		t.Fatalf("truncate execCalls len = %d, want %d", len(tx.execCalls), 1)
	}
	if tx.execCalls[0].sql != constants.Truncate {
		t.Fatalf("truncate executed %q, want %q", tx.execCalls[0].sql, constants.Truncate)
	}
}

func TestTruncateError(t *testing.T) {
	tx := newMockTx()
	tx.setExecError(constants.Truncate, errors.New("fail"))
	ctx := context.Background()

	if ok := truncate(ctx, tx); ok {
		t.Fatalf("truncate() = %v, want %v", ok, false)
	}
	if len(tx.execCalls) != 1 {
		t.Fatalf("truncate execCalls len = %d, want %d", len(tx.execCalls), 1)
	}
}

func TestUpsertSuccess(t *testing.T) {
	perfume := models.Perfume{
		Brand:       "Brand",
		Name:        "Name",
		Sex:         "male",
		Type:        "type",
		Family:      []string{"family"},
		UpperNotes:  []string{"upper"},
		MiddleNotes: []string{"middle"},
		BaseNotes:   []string{"base"},
		Link:        "http://link",
		Volume:      100,
		ImageUrl:    "http://image",
	}

	tx := newMockTx()
	ctx := context.Background()

	status := upsert(ctx, tx, []models.Perfume{perfume})

	if !status.State.Success {
		t.Fatalf("upsert status success = %v, want %v", status.State.Success, true)
	}
	if status.State.SuccessfulCount != 1 {
		t.Fatalf("upsert successful count = %d, want %d", status.State.SuccessfulCount, 1)
	}
	if len(status.SuccessfulPerfumes) != 1 || !reflect.DeepEqual(status.SuccessfulPerfumes[0], perfume) {
		t.Fatalf("upsert successful perfumes = %#v, want %#v", status.SuccessfulPerfumes, []models.Perfume{perfume})
	}
	if len(status.FailedPerfumes) != 0 {
		t.Fatalf("upsert failed perfumes len = %d, want %d", len(status.FailedPerfumes), 0)
	}

	executed := make([]string, len(tx.execCalls))
	for i, call := range tx.execCalls {
		executed[i] = call.sql
	}

	expected := []string{
		getSavepointQuery(constants.Savepoint, 0),
		constants.UpdatePerfumes,
		constants.UpdatePerfumeLinks,
		getSavepointQuery(constants.ReleaseSavepoint, 0),
	}

	if len(executed) != len(expected) {
		t.Fatalf("upsert exec len = %d, want %d", len(executed), len(expected))
	}
	for i := range expected {
		if executed[i] != expected[i] {
			t.Fatalf("upsert exec[%d] = %q, want %q", i, executed[i], expected[i])
		}
	}
}

func TestUpsertFailureRollsBack(t *testing.T) {
	perfume := models.Perfume{
		Brand:       "Brand",
		Name:        "Name",
		Sex:         "male",
		Type:        "type",
		Family:      []string{"family"},
		UpperNotes:  []string{"upper"},
		MiddleNotes: []string{"middle"},
		BaseNotes:   []string{"base"},
		Link:        "http://link",
		Volume:      100,
		ImageUrl:    "http://image",
	}

	tx := newMockTx()
	tx.setExecError(constants.UpdatePerfumes, errors.New("boom"))
	ctx := context.Background()

	status := upsert(ctx, tx, []models.Perfume{perfume})

	if status.State.SuccessfulCount != 0 {
		t.Fatalf("upsert successful count = %d, want %d", status.State.SuccessfulCount, 0)
	}
	if status.State.FailedCount != 1 {
		t.Fatalf("upsert failed count = %d, want %d", status.State.FailedCount, 1)
	}
	if len(status.FailedPerfumes) != 1 || !reflect.DeepEqual(status.FailedPerfumes[0], perfume) {
		t.Fatalf("upsert failed perfumes = %#v, want %#v", status.FailedPerfumes, []models.Perfume{perfume})
	}
	if len(status.SuccessfulPerfumes) != 0 {
		t.Fatalf("upsert successful perfumes len = %d, want %d", len(status.SuccessfulPerfumes), 0)
	}

	executed := make([]string, len(tx.execCalls))
	for i, call := range tx.execCalls {
		executed[i] = call.sql
	}

	expected := []string{
		getSavepointQuery(constants.Savepoint, 0),
		constants.UpdatePerfumes,
		constants.UpdatePerfumeLinks,
		getSavepointQuery(constants.RollbackSavepoint, 0),
	}

	if len(executed) != len(expected) {
		t.Fatalf("upsert exec len = %d, want %d", len(executed), len(expected))
	}
	for i := range expected {
		if executed[i] != expected[i] {
			t.Fatalf("upsert exec[%d] = %q, want %q", i, executed[i], expected[i])
		}
	}
}
