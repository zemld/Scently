package core

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	queries "github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/query"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

func TestGetSavepointQuery(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		i        int
		expected string
	}{
		{"savepoint with index 3", "SAVEPOINT sp_", 3, "SAVEPOINT sp_3"},
		{"rollback with index 0", "ROLLBACK TO SAVEPOINT sp_", 0, "ROLLBACK TO SAVEPOINT sp_0"},
		{"release with index 5", "RELEASE SAVEPOINT sp_", 5, "RELEASE SAVEPOINT sp_5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSavepointQuery(tt.cmd, tt.i); got != tt.expected {
				t.Fatalf("getSavepointQuery(%q, %d) = %q, want %q", tt.cmd, tt.i, got, tt.expected)
			}
		})
	}
}

type execCall struct {
	sql  string
	args []any
}

type mockTx struct {
	execCalls     []execCall
	execErrs      map[string]error
	execErrCounts map[string]int // счетчик вызовов для каждого SQL запроса
	execErrOnCall map[string]int // на каком вызове вернуть ошибку (0 = всегда)
}

func newMockTx() *mockTx {
	return &mockTx{
		execErrs:      make(map[string]error),
		execErrCounts: make(map[string]int),
		execErrOnCall: make(map[string]int),
	}
}

func (m *mockTx) setExecError(sql string, err error) {
	m.execErrs[sql] = err
	m.execErrOnCall[sql] = 0 // всегда возвращать ошибку
}

func (m *mockTx) setExecErrorOnCall(sql string, err error, callNumber int) {
	m.execErrs[sql] = err
	m.execErrOnCall[sql] = callNumber // вернуть ошибку на определенном вызове
}

func (m *mockTx) Exec(_ context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.execCalls = append(m.execCalls, execCall{sql: sql, args: arguments})
	if err, ok := m.execErrs[sql]; ok {
		m.execErrCounts[sql]++
		callNumber := m.execErrOnCall[sql]
		if callNumber == 0 || m.execErrCounts[sql] == callNumber {
			return pgconn.NewCommandTag(""), err
		}
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

func TestNewUpdateStatus(t *testing.T) {
	tests := []struct {
		name    string
		success bool
	}{
		{"success true", true},
		{"success false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := models.ProcessedState{Success: tt.success}
			if status.Success != tt.success {
				t.Fatalf("ProcessedState(%v).Success = %v, want %v", tt.success, status.Success, tt.success)
			}
			if status.SuccessfulCount != 0 {
				t.Fatalf("ProcessedState().SuccessfulCount = %d, want 0", status.SuccessfulCount)
			}
			if status.FailedCount != 0 {
				t.Fatalf("ProcessedState().FailedCount = %d, want 0", status.FailedCount)
			}
		})
	}
}

func TestUpsertSuccess(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			Type:       "Eau de Parfum",
			Family:     []string{"Floral"},
			UpperNotes: []string{"Bergamot"},
			CoreNotes:  []string{"Rose"},
			BaseNotes:  []string{"Musk"},
		},
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: "http://image1.com",
				Variants: []models.PerfumeVariant{
					{Volume: 100, Price: 5000, Link: "http://link1.com"},
				},
			},
		},
	}

	tx := newMockTx()
	ctx := context.Background()

	status := upsert(ctx, tx, []models.Perfume{perfume})

	if !status.Success {
		t.Fatalf("upsert status success = %v, want %v", status.Success, true)
	}
	if status.SuccessfulCount != 1 {
		t.Fatalf("upsert successful count = %d, want %d", status.SuccessfulCount, 1)
	}
	if status.FailedCount != 0 {
		t.Fatalf("upsert failed count = %d, want %d", status.FailedCount, 0)
	}

	// Проверяем, что были вызваны правильные SQL запросы
	expectedQueries := []string{
		getSavepointQuery(queries.Savepoint, 0),
		queries.GetOrInsertShop,
		queries.InsertVariant,
		queries.InsertFamily,
		queries.InsertUpperNote,
		queries.InsertCoreNote,
		queries.InsertBaseNote,
		queries.InsertPerfumeBaseInfo,
		getSavepointQuery(queries.ReleaseSavepoint, 0),
	}

	if len(tx.execCalls) != len(expectedQueries) {
		t.Fatalf("upsert exec calls len = %d, want %d", len(tx.execCalls), len(expectedQueries))
	}

	for i, expected := range expectedQueries {
		if tx.execCalls[i].sql != expected {
			t.Fatalf("upsert exec[%d] = %q, want %q", i, tx.execCalls[i].sql, expected)
		}
	}
}

func TestUpsertFailureRollsBack(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			Type:       "Eau de Parfum",
			Family:     []string{"Floral"},
			UpperNotes: []string{"Bergamot"},
			CoreNotes:  []string{"Rose"},
			BaseNotes:  []string{"Musk"},
		},
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: "http://image1.com",
				Variants: []models.PerfumeVariant{
					{Volume: 100, Price: 5000, Link: "http://link1.com"},
				},
			},
		},
	}

	tx := newMockTx()
	tx.setExecError(queries.GetOrInsertShop, errors.New("database error"))
	ctx := context.Background()

	status := upsert(ctx, tx, []models.Perfume{perfume})

	if status.SuccessfulCount != 0 {
		t.Fatalf("upsert successful count = %d, want %d", status.SuccessfulCount, 0)
	}
	if status.FailedCount != 1 {
		t.Fatalf("upsert failed count = %d, want %d", status.FailedCount, 1)
	}

	// Проверяем, что был выполнен rollback
	expectedQueries := []string{
		getSavepointQuery(queries.Savepoint, 0),
		queries.GetOrInsertShop,
		getSavepointQuery(queries.RollbackSavepoint, 0),
	}

	if len(tx.execCalls) != len(expectedQueries) {
		t.Fatalf("upsert exec calls len = %d, want %d", len(tx.execCalls), len(expectedQueries))
	}

	for i, expected := range expectedQueries {
		if tx.execCalls[i].sql != expected {
			t.Fatalf("upsert exec[%d] = %q, want %q", i, tx.execCalls[i].sql, expected)
		}
	}
}

func TestUpsertMultiplePerfumes(t *testing.T) {
	perfumes := []models.Perfume{
		{
			Brand: "Brand1",
			Name:  "Name1",
			Sex:   "male",
			Properties: models.PerfumeProperties{
				Type:       "Eau de Parfum",
				Family:     []string{"Floral"},
				UpperNotes: []string{"Bergamot"},
				CoreNotes:  []string{"Rose"},
				BaseNotes:  []string{"Musk"},
			},
			Shops: []models.ShopInfo{
				{
					ShopName: "Gold Apple",
					Domain:   "goldapple.ru",
					ImageUrl: "http://image1.com",
					Variants: []models.PerfumeVariant{
						{Volume: 100, Price: 5000, Link: "http://link1.com"},
					},
				},
			},
		},
		{
			Brand: "Brand2",
			Name:  "Name2",
			Sex:   "female",
			Properties: models.PerfumeProperties{
				Type:       "Eau de Toilette",
				Family:     []string{"Woody"},
				UpperNotes: []string{"Lemon"},
				CoreNotes:  []string{"Jasmine"},
				BaseNotes:  []string{"Sandalwood"},
			},
			Shops: []models.ShopInfo{
				{
					ShopName: "Randewoo",
					Domain:   "randewoo.ru",
					ImageUrl: "http://image2.com",
					Variants: []models.PerfumeVariant{
						{Volume: 50, Price: 3000, Link: "http://link2.com"},
					},
				},
			},
		},
	}

	tx := newMockTx()
	ctx := context.Background()

	status := upsert(ctx, tx, perfumes)

	if status.SuccessfulCount != 2 {
		t.Fatalf("upsert successful count = %d, want %d", status.SuccessfulCount, 2)
	}
	if status.FailedCount != 0 {
		t.Fatalf("upsert failed count = %d, want %d", status.FailedCount, 0)
	}

	// Проверяем, что для каждого парфюма были созданы savepoints
	expectedSavepoints := []string{
		getSavepointQuery(queries.Savepoint, 0),
		getSavepointQuery(queries.ReleaseSavepoint, 0),
		getSavepointQuery(queries.Savepoint, 1),
		getSavepointQuery(queries.ReleaseSavepoint, 1),
	}

	savepointCalls := []string{}
	for _, call := range tx.execCalls {
		if call.sql == getSavepointQuery(queries.Savepoint, 0) ||
			call.sql == getSavepointQuery(queries.Savepoint, 1) ||
			call.sql == getSavepointQuery(queries.ReleaseSavepoint, 0) ||
			call.sql == getSavepointQuery(queries.ReleaseSavepoint, 1) {
			savepointCalls = append(savepointCalls, call.sql)
		}
	}

	if len(savepointCalls) != len(expectedSavepoints) {
		t.Fatalf("savepoint calls len = %d, want %d", len(savepointCalls), len(expectedSavepoints))
	}
}

func TestUpsertPartialFailure(t *testing.T) {
	perfumes := []models.Perfume{
		{
			Brand: "Brand1",
			Name:  "Name1",
			Sex:   "male",
			Properties: models.PerfumeProperties{
				Type:       "Eau de Parfum",
				Family:     []string{"Floral"},
				UpperNotes: []string{"Bergamot"},
				CoreNotes:  []string{"Rose"},
				BaseNotes:  []string{"Musk"},
			},
			Shops: []models.ShopInfo{
				{
					ShopName: "Gold Apple",
					Domain:   "goldapple.ru",
					ImageUrl: "http://image1.com",
					Variants: []models.PerfumeVariant{
						{Volume: 100, Price: 5000, Link: "http://link1.com"},
					},
				},
			},
		},
		{
			Brand: "Brand2",
			Name:  "Name2",
			Sex:   "female",
			Properties: models.PerfumeProperties{
				Type:       "Eau de Toilette",
				Family:     []string{"Woody"},
				UpperNotes: []string{"Lemon"},
				CoreNotes:  []string{"Jasmine"},
				BaseNotes:  []string{"Sandalwood"},
			},
			Shops: []models.ShopInfo{
				{
					ShopName: "Randewoo",
					Domain:   "randewoo.ru",
					ImageUrl: "http://image2.com",
					Variants: []models.PerfumeVariant{
						{Volume: 50, Price: 3000, Link: "http://link2.com"},
					},
				},
			},
		},
	}

	tx := newMockTx()
	// Устанавливаем ошибку для InsertFamily на втором вызове (для второго парфюма)
	// Первый парфюм обработается успешно, второй упадет
	tx.setExecErrorOnCall(queries.InsertFamily, errors.New("database error"), 2)
	ctx := context.Background()

	status := upsert(ctx, tx, perfumes)

	// Первый парфюм должен быть успешным, второй - неудачным
	if status.SuccessfulCount != 1 {
		t.Fatalf("upsert successful count = %d, want %d", status.SuccessfulCount, 1)
	}
	if status.FailedCount != 1 {
		t.Fatalf("upsert failed count = %d, want %d", status.FailedCount, 1)
	}
}

func TestDeleteOldPerfumes(t *testing.T) {
	tx := newMockTx()
	ctx := context.Background()

	result := deleteOldPerfumes(ctx, tx)

	if !result {
		t.Fatalf("deleteOldPerfumes() = %v, want %v", result, true)
	}

	if len(tx.execCalls) != 1 {
		t.Fatalf("deleteOldPerfumes exec calls len = %d, want %d", len(tx.execCalls), 1)
	}

	if tx.execCalls[0].sql != queries.DeleteOldPerfumes {
		t.Fatalf("deleteOldPerfumes exec sql = %q, want %q", tx.execCalls[0].sql, queries.DeleteOldPerfumes)
	}
}

func TestDeleteOldPerfumesError(t *testing.T) {
	tx := newMockTx()
	tx.setExecError(queries.DeleteOldPerfumes, errors.New("database error"))
	ctx := context.Background()

	result := deleteOldPerfumes(ctx, tx)

	if result {
		t.Fatalf("deleteOldPerfumes() = %v, want %v", result, false)
	}
}

func TestRunUpdateQueries(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			Type:       "Eau de Parfum",
			Family:     []string{"Floral", "Woody"},
			UpperNotes: []string{"Bergamot", "Lemon"},
			CoreNotes:  []string{"Rose", "Jasmine"},
			BaseNotes:  []string{"Musk", "Sandalwood"},
		},
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: "http://image1.com",
				Variants: []models.PerfumeVariant{
					{Volume: 100, Price: 5000, Link: "http://link1.com"},
					{Volume: 50, Price: 3000, Link: "http://link2.com"},
				},
			},
		},
	}

	tx := newMockTx()
	ctx := context.Background()

	err := runUpdateQueries(ctx, tx, perfume)

	if err != nil {
		t.Fatalf("runUpdateQueries() error = %v, want nil", err)
	}

	// Проверяем, что все запросы были выполнены
	expectedQueries := []string{
		queries.GetOrInsertShop,
		queries.InsertVariant,
		queries.InsertVariant,
		queries.InsertFamily,
		queries.InsertFamily,
		queries.InsertUpperNote,
		queries.InsertUpperNote,
		queries.InsertCoreNote,
		queries.InsertCoreNote,
		queries.InsertBaseNote,
		queries.InsertBaseNote,
		queries.InsertPerfumeBaseInfo,
	}

	if len(tx.execCalls) != len(expectedQueries) {
		t.Fatalf("runUpdateQueries exec calls len = %d, want %d", len(tx.execCalls), len(expectedQueries))
	}

	for i, expected := range expectedQueries {
		if tx.execCalls[i].sql != expected {
			t.Fatalf("runUpdateQueries exec[%d] = %q, want %q", i, tx.execCalls[i].sql, expected)
		}
	}
}

func TestRunUpdateQueriesError(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			Type:       "Eau de Parfum",
			Family:     []string{"Floral"},
			UpperNotes: []string{"Bergamot"},
			CoreNotes:  []string{"Rose"},
			BaseNotes:  []string{"Musk"},
		},
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: "http://image1.com",
				Variants: []models.PerfumeVariant{
					{Volume: 100, Price: 5000, Link: "http://link1.com"},
				},
			},
		},
	}

	tests := []struct {
		name      string
		errorSQL  string
		errorFunc func(*mockTx)
	}{
		{"shop error", queries.GetOrInsertShop, func(tx *mockTx) { tx.setExecError(queries.GetOrInsertShop, errors.New("shop error")) }},
		{"variant error", queries.InsertVariant, func(tx *mockTx) { tx.setExecError(queries.InsertVariant, errors.New("variant error")) }},
		{"family error", queries.InsertFamily, func(tx *mockTx) { tx.setExecError(queries.InsertFamily, errors.New("family error")) }},
		{"upper note error", queries.InsertUpperNote, func(tx *mockTx) { tx.setExecError(queries.InsertUpperNote, errors.New("note error")) }},
		{"core note error", queries.InsertCoreNote, func(tx *mockTx) { tx.setExecError(queries.InsertCoreNote, errors.New("note error")) }},
		{"base note error", queries.InsertBaseNote, func(tx *mockTx) { tx.setExecError(queries.InsertBaseNote, errors.New("note error")) }},
		{"perfume type error", queries.InsertPerfumeBaseInfo, func(tx *mockTx) { tx.setExecError(queries.InsertPerfumeBaseInfo, errors.New("perfume error")) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := newMockTx()
			tt.errorFunc(tx)
			ctx := context.Background()

			err := runUpdateQueries(ctx, tx, perfume)

			if err == nil {
				t.Fatalf("runUpdateQueries() error = nil, want error")
			}
		})
	}
}

func TestUpdateShopInfo(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: "http://image1.com",
				Variants: []models.PerfumeVariant{
					{Volume: 100, Price: 5000, Link: "http://link1.com"},
					{Volume: 50, Price: 3000, Link: "http://link2.com"},
				},
			},
			{
				ShopName: "Randewoo",
				Domain:   "randewoo.ru",
				ImageUrl: "http://image2.com",
				Variants: []models.PerfumeVariant{
					{Volume: 75, Price: 4000, Link: "http://link3.com"},
				},
			},
		},
	}

	tx := newMockTx()
	ctx := context.Background()

	err := updateShopInfo(ctx, tx, perfume)

	if err != nil {
		t.Fatalf("updateShopInfo() error = %v, want nil", err)
	}

	// Должно быть 2 вызова GetOrInsertShop и 3 вызова InsertVariant (2+1)
	expectedShopCalls := 2
	expectedVariantCalls := 3

	shopCalls := 0
	variantCalls := 0

	for _, call := range tx.execCalls {
		if call.sql == queries.GetOrInsertShop {
			shopCalls++
		}
		if call.sql == queries.InsertVariant {
			variantCalls++
		}
	}

	if shopCalls != expectedShopCalls {
		t.Fatalf("updateShopInfo GetOrInsertShop calls = %d, want %d", shopCalls, expectedShopCalls)
	}
	if variantCalls != expectedVariantCalls {
		t.Fatalf("updateShopInfo InsertVariant calls = %d, want %d", variantCalls, expectedVariantCalls)
	}

	// Проверяем аргументы первого варианта
	firstVariantCall := tx.execCalls[1] // После первого GetOrInsertShop
	expectedArgs := []any{"Brand", "Name", "male", "Gold Apple", 100, 5000, "http://link1.com"}
	if !reflect.DeepEqual(firstVariantCall.args, expectedArgs) {
		t.Fatalf("updateShopInfo first variant args = %v, want %v", firstVariantCall.args, expectedArgs)
	}
}

func TestUpdateShopInfoError(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: "http://image1.com",
				Variants: []models.PerfumeVariant{
					{Volume: 100, Price: 5000, Link: "http://link1.com"},
				},
			},
		},
	}

	tests := []struct {
		name     string
		errorSQL string
	}{
		{"shop error", queries.GetOrInsertShop},
		{"variant error", queries.InsertVariant},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := newMockTx()
			tx.setExecError(tt.errorSQL, errors.New("database error"))
			ctx := context.Background()

			err := updateShopInfo(ctx, tx, perfume)

			if err == nil {
				t.Fatalf("updateShopInfo() error = nil, want error")
			}
		})
	}
}

func TestUpdateFamilies(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			Family: []string{"Floral", "Woody", "Oriental"},
		},
	}

	tx := newMockTx()
	ctx := context.Background()

	err := updateFamilies(ctx, tx, perfume)

	if err != nil {
		t.Fatalf("updateFamilies() error = %v, want nil", err)
	}

	if len(tx.execCalls) != 3 {
		t.Fatalf("updateFamilies exec calls len = %d, want %d", len(tx.execCalls), 3)
	}

	for i, call := range tx.execCalls {
		if call.sql != queries.InsertFamily {
			t.Fatalf("updateFamilies exec[%d] sql = %q, want %q", i, call.sql, queries.InsertFamily)
		}
		expectedArgs := []any{"Brand", "Name", "male", perfume.Properties.Family[i]}
		if !reflect.DeepEqual(call.args, expectedArgs) {
			t.Fatalf("updateFamilies exec[%d] args = %v, want %v", i, call.args, expectedArgs)
		}
	}
}

func TestUpdateFamiliesError(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			Family: []string{"Floral"},
		},
	}

	tx := newMockTx()
	tx.setExecError(queries.InsertFamily, errors.New("database error"))
	ctx := context.Background()

	err := updateFamilies(ctx, tx, perfume)

	if err == nil {
		t.Fatalf("updateFamilies() error = nil, want error")
	}
}

func TestUpdateNotes(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			UpperNotes: []string{"Bergamot", "Lemon"},
			CoreNotes:  []string{"Rose", "Jasmine"},
			BaseNotes:  []string{"Musk", "Sandalwood"},
		},
	}

	tests := []struct {
		name     string
		query    string
		notes    []string
		noteType string
	}{
		{"upper notes", queries.InsertUpperNote, perfume.Properties.UpperNotes, "upper"},
		{"core notes", queries.InsertCoreNote, perfume.Properties.CoreNotes, "core"},
		{"base notes", queries.InsertBaseNote, perfume.Properties.BaseNotes, "base"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := newMockTx()
			ctx := context.Background()

			err := updateNotes(ctx, tx, tt.query, perfume, tt.notes)

			if err != nil {
				t.Fatalf("updateNotes(%s) error = %v, want nil", tt.noteType, err)
			}

			if len(tx.execCalls) != len(tt.notes) {
				t.Fatalf("updateNotes(%s) exec calls len = %d, want %d", tt.noteType, len(tx.execCalls), len(tt.notes))
			}

			for i, call := range tx.execCalls {
				if call.sql != tt.query {
					t.Fatalf("updateNotes(%s) exec[%d] sql = %q, want %q", tt.noteType, i, call.sql, tt.query)
				}
				expectedArgs := []any{"Brand", "Name", "male", tt.notes[i]}
				if !reflect.DeepEqual(call.args, expectedArgs) {
					t.Fatalf("updateNotes(%s) exec[%d] args = %v, want %v", tt.noteType, i, call.args, expectedArgs)
				}
			}
		})
	}
}

func TestUpdateNotesError(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			UpperNotes: []string{"Bergamot"},
		},
	}

	tx := newMockTx()
	tx.setExecError(queries.InsertUpperNote, errors.New("database error"))
	ctx := context.Background()

	err := updateNotes(ctx, tx, queries.InsertUpperNote, perfume, perfume.Properties.UpperNotes)

	if err == nil {
		t.Fatalf("updateNotes() error = nil, want error")
	}
}

func TestUpdatePerfumeType(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			Type: "Eau de Parfum",
		},
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: "http://image1.com",
				Variants: []models.PerfumeVariant{},
			},
		},
	}

	tx := newMockTx()
	ctx := context.Background()

	err := updatePerfumeType(ctx, tx, perfume)

	if err != nil {
		t.Fatalf("updatePerfumeType() error = %v, want nil", err)
	}

	if len(tx.execCalls) != 1 {
		t.Fatalf("updatePerfumeType exec calls len = %d, want %d", len(tx.execCalls), 1)
	}

	call := tx.execCalls[0]
	if call.sql != queries.InsertPerfumeBaseInfo {
		t.Fatalf("updatePerfumeType exec sql = %q, want %q", call.sql, queries.InsertPerfumeBaseInfo)
	}

	expectedArgs := []any{"Brand", "Name", "male", "Eau de Parfum", "http://image1.com"}
	if !reflect.DeepEqual(call.args, expectedArgs) {
		t.Fatalf("updatePerfumeType exec args = %v, want %v", call.args, expectedArgs)
	}
}

func TestUpdatePerfumeTypeError(t *testing.T) {
	perfume := models.Perfume{
		Brand: "Brand",
		Name:  "Name",
		Sex:   "male",
		Properties: models.PerfumeProperties{
			Type: "Eau de Parfum",
		},
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: "http://image1.com",
				Variants: []models.PerfumeVariant{},
			},
		},
	}

	tx := newMockTx()
	tx.setExecError(queries.InsertPerfumeBaseInfo, errors.New("database error"))
	ctx := context.Background()

	err := updatePerfumeType(ctx, tx, perfume)

	if err == nil {
		t.Fatalf("updatePerfumeType() error = nil, want error")
	}
}

func TestGetPreferredImageUrl(t *testing.T) {
	tests := []struct {
		name     string
		perfume  models.Perfume
		expected string
	}{
		{
			"Gold Apple priority",
			models.Perfume{
				Shops: []models.ShopInfo{
					{ShopName: "Randewoo", ImageUrl: "http://randewoo.com/image"},
					{ShopName: "Gold Apple", ImageUrl: "http://goldapple.com/image"},
					{ShopName: "Letu", ImageUrl: "http://letu.com/image"},
				},
			},
			"http://goldapple.com/image",
		},
		{
			"Randewoo when Gold Apple missing",
			models.Perfume{
				Shops: []models.ShopInfo{
					{ShopName: "Letu", ImageUrl: "http://letu.com/image"},
					{ShopName: "Randewoo", ImageUrl: "http://randewoo.com/image"},
				},
			},
			"http://randewoo.com/image",
		},
		{
			"Letu when others missing",
			models.Perfume{
				Shops: []models.ShopInfo{
					{ShopName: "Letu", ImageUrl: "http://letu.com/image"},
				},
			},
			"http://letu.com/image",
		},
		{
			"empty when no shops",
			models.Perfume{
				Shops: []models.ShopInfo{},
			},
			"",
		},
		{
			"unknown shop returns its image",
			models.Perfume{
				Shops: []models.ShopInfo{
					{ShopName: "Unknown Shop", ImageUrl: "http://unknown.com/image"},
				},
			},
			"http://unknown.com/image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPreferredImageUrl(tt.perfume)

			if result != tt.expected {
				t.Fatalf("getPreferredImageUrl() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestUpdateSavepointStatus(t *testing.T) {
	tx := newMockTx()
	ctx := context.Background()

	updateSavepointStatus(ctx, tx, queries.Savepoint, 0)
	updateSavepointStatus(ctx, tx, queries.RollbackSavepoint, 1)
	updateSavepointStatus(ctx, tx, queries.ReleaseSavepoint, 2)

	if len(tx.execCalls) != 3 {
		t.Fatalf("updateSavepointStatus exec calls len = %d, want %d", len(tx.execCalls), 3)
	}

	expected := []string{
		getSavepointQuery(queries.Savepoint, 0),
		getSavepointQuery(queries.RollbackSavepoint, 1),
		getSavepointQuery(queries.ReleaseSavepoint, 2),
	}

	for i, expectedSQL := range expected {
		if tx.execCalls[i].sql != expectedSQL {
			t.Fatalf("updateSavepointStatus exec[%d] = %q, want %q", i, tx.execCalls[i].sql, expectedSQL)
		}
	}
}
