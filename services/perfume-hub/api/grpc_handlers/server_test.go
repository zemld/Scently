package grpc_handlers

import (
	"context"
	"net"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/zemld/Scently/generated/proto/perfume-hub"
	protoModels "github.com/zemld/Scently/generated/proto/perfume-hub/models"
	"github.com/zemld/Scently/generated/proto/perfume-hub/requests"
	"github.com/zemld/Scently/perfume-hub/internal/db/core"
	"github.com/zemld/Scently/perfume-hub/internal/errors"
	"github.com/zemld/Scently/perfume-hub/internal/models"
)

const bufSize = 1024 * 1024

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func setupTestServer(
	t *testing.T,
	mockSelect core.SelectFunc,
	mockUpdate core.UpdateFunc,
) (*grpc.Server, *bufconn.Listener, *grpc.ClientConn, pb.PerfumeStorageClient) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()

	server := NewPerfumeStorageServer(mockSelect, mockUpdate)
	pb.RegisterPerfumeStorageServer(s, server)

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()

	// Создаем клиент, который подключается через буфер
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	client := pb.NewPerfumeStorageClient(conn)
	return s, lis, conn, client
}

func TestPerfumeStorageServer_GetPerfumes(t *testing.T) {
	tests := []struct {
		name           string
		request        *requests.GetPerfumesRequest
		mockPerfumes   []models.Perfume
		mockStatus     models.ProcessedState
		expectedCount  int
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "successful get with brand filter",
			request: &requests.GetPerfumesRequest{
				Brand: "Chanel",
				Sex:   "female",
			},
			mockPerfumes: []models.Perfume{
				{
					Brand: "Chanel",
					Name:  "No.5",
					Sex:   "female",
					Properties: models.PerfumeProperties{
						Type: "Eau de Parfum",
					},
				},
			},
			mockStatus:    models.NewProcessedState(),
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "successful get with name filter",
			request: &requests.GetPerfumesRequest{
				Name: "No.5",
			},
			mockPerfumes: []models.Perfume{
				{
					Brand: "Chanel",
					Name:  "No.5",
					Sex:   "female",
				},
				{
					Brand: "Chanel",
					Name:  "No.5",
					Sex:   "unisex",
				},
			},
			mockStatus:    models.NewProcessedState(),
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "empty result",
			request: &requests.GetPerfumesRequest{
				Brand: "NonExistent",
			},
			mockPerfumes:  []models.Perfume{},
			mockStatus:    models.NewProcessedState(),
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "database error",
			request: &requests.GetPerfumesRequest{
				Brand: "Chanel",
			},
			mockPerfumes: nil,
			mockStatus: models.ProcessedState{
				Error: errors.NewDBError("database connection failed", nil),
			},
			expectedCount:  0,
			expectedError:  true,
			expectedErrMsg: "database error: database connection failed",
		},
		{
			name: "successful get with all filters",
			request: &requests.GetPerfumesRequest{
				Brand: "Chanel",
				Name:  "No.5",
				Sex:   "female",
			},
			mockPerfumes: []models.Perfume{
				{
					Brand: "Chanel",
					Name:  "No.5",
					Sex:   "female",
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
							Variants: []models.PerfumeVariant{
								{Volume: 100, Price: 5000, Link: "http://example.com/link"},
							},
						},
					},
				},
			},
			mockStatus:    models.NewProcessedState(),
			expectedCount: 1,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокируем core.Select
			mockSelect := func(ctx context.Context, params *models.SelectParameters) ([]models.Perfume, models.ProcessedState) {
				// Проверяем, что параметры переданы правильно
				if tt.request.Brand != "" && params.Brand != tt.request.Brand {
					t.Errorf("Select() params.Brand = %q, want %q", params.Brand, tt.request.Brand)
				}
				if tt.request.Name != "" && params.Name != tt.request.Name {
					t.Errorf("Select() params.Name = %q, want %q", params.Name, tt.request.Name)
				}
				if tt.request.Sex != "" && params.Sex != tt.request.Sex {
					t.Errorf("Select() params.Sex = %q, want %q", params.Sex, tt.request.Sex)
				}
				return tt.mockPerfumes, tt.mockStatus
			}

			// Создаем тестовый сервер с мокированными функциями
			s, lis, conn, client := setupTestServer(t, mockSelect, nil)
			defer s.Stop()
			defer lis.Close()
			defer conn.Close()

			ctx := context.Background()
			response, err := client.GetPerfumes(ctx, tt.request)

			if tt.expectedError {
				if err == nil {
					t.Errorf("GetPerfumes() error = nil, want error")
				} else if !contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("GetPerfumes() error = %q, want to contain %q", err.Error(), tt.expectedErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("GetPerfumes() error = %v, want nil", err)
				}
				if response == nil {
					t.Fatal("GetPerfumes() response = nil, want non-nil")
				}
				if len(response.Perfumes) != tt.expectedCount {
					t.Errorf("GetPerfumes() Perfumes len = %d, want %d", len(response.Perfumes), tt.expectedCount)
				}
			}
		})
	}
}

func TestPerfumeStorageServer_UpdatePerfumes(t *testing.T) {
	imageUrl := "http://example.com/image.jpg"
	tests := []struct {
		name               string
		request            *requests.UpdatePerfumesRequest
		mockStatus         models.ProcessedState
		expectedSuccessful int32
		expectedFailed     int32
		verifyParams       func(t *testing.T, params *models.UpdateParameters)
	}{
		{
			name: "successful update single perfume",
			request: &requests.UpdatePerfumesRequest{
				Perfumes: []*protoModels.Perfume{
					{
						Brand:    "Chanel",
						Name:     "No.5",
						Sex:      "female",
						ImageUrl: &imageUrl,
						Properties: &protoModels.Perfume_Properties{
							PerfumeType: "Eau de Parfum",
							Family:      []string{"Floral"},
						},
					},
				},
			},
			mockStatus: models.ProcessedState{
				SuccessfulCount: 1,
				FailedCount:     0,
			},
			expectedSuccessful: 1,
			expectedFailed:     0,
			verifyParams: func(t *testing.T, params *models.UpdateParameters) {
				if len(params.Perfumes) != 1 {
					t.Errorf("Update() params.Perfumes len = %d, want 1", len(params.Perfumes))
				}
				if params.Perfumes[0].Brand != "Chanel" {
					t.Errorf("Update() params.Perfumes[0].Brand = %q, want %q", params.Perfumes[0].Brand, "Chanel")
				}
			},
		},
		{
			name: "successful update multiple perfumes",
			request: &requests.UpdatePerfumesRequest{
				Perfumes: []*protoModels.Perfume{
					{
						Brand: "Chanel",
						Name:  "No.5",
						Sex:   "female",
					},
					{
						Brand: "Dior",
						Name:  "Sauvage",
						Sex:   "male",
					},
				},
			},
			mockStatus: models.ProcessedState{
				SuccessfulCount: 2,
				FailedCount:     0,
			},
			expectedSuccessful: 2,
			expectedFailed:     0,
			verifyParams: func(t *testing.T, params *models.UpdateParameters) {
				if len(params.Perfumes) != 2 {
					t.Errorf("Update() params.Perfumes len = %d, want 2", len(params.Perfumes))
				}
			},
		},
		{
			name: "update with partial failures",
			request: &requests.UpdatePerfumesRequest{
				Perfumes: []*protoModels.Perfume{
					{
						Brand: "Chanel",
						Name:  "No.5",
						Sex:   "female",
					},
					{
						Brand: "Invalid",
						Name:  "Perfume",
						Sex:   "unknown",
					},
				},
			},
			mockStatus: models.ProcessedState{
				SuccessfulCount: 1,
				FailedCount:     1,
			},
			expectedSuccessful: 1,
			expectedFailed:     1,
		},
		{
			name: "update with nil perfume in request",
			request: &requests.UpdatePerfumesRequest{
				Perfumes: []*protoModels.Perfume{
					{
						Brand: "Chanel",
						Name:  "No.5",
						Sex:   "female",
					},
					nil,
					{
						Brand: "Dior",
						Name:  "Sauvage",
						Sex:   "male",
					},
				},
			},
			mockStatus: models.ProcessedState{
				SuccessfulCount: 2,
				FailedCount:     0,
			},
			expectedSuccessful: 2,
			expectedFailed:     0,
			verifyParams: func(t *testing.T, params *models.UpdateParameters) {
				// nil perfumes должны быть пропущены
				if len(params.Perfumes) != 3 {
					t.Errorf("Update() params.Perfumes len = %d, want 3", len(params.Perfumes))
				}
			},
		},
		{
			name: "empty request",
			request: &requests.UpdatePerfumesRequest{
				Perfumes: []*protoModels.Perfume{},
			},
			mockStatus: models.ProcessedState{
				SuccessfulCount: 0,
				FailedCount:     0,
			},
			expectedSuccessful: 0,
			expectedFailed:     0,
		},
		{
			name: "update with complex perfume data",
			request: &requests.UpdatePerfumesRequest{
				Perfumes: []*protoModels.Perfume{
					{
						Brand:    "Chanel",
						Name:     "No.5",
						Sex:      "female",
						ImageUrl: &imageUrl,
						Properties: &protoModels.Perfume_Properties{
							PerfumeType: "Eau de Parfum",
							Family:      []string{"Floral", "Woody"},
							UpperNotes:  []string{"Bergamot", "Lemon"},
							CoreNotes:   []string{"Rose", "Jasmine"},
							BaseNotes:   []string{"Musk", "Sandalwood"},
						},
						Shops: []*protoModels.Perfume_ShopInfo{
							{
								ShopName: "Gold Apple",
								Domain:   "goldapple.ru",
								ImageUrl: &imageUrl,
								Variants: []*protoModels.Perfume_ShopInfo_Variant{
									{Volume: 100, Price: 5000, Link: "http://example.com/link1"},
									{Volume: 50, Price: 3000, Link: "http://example.com/link2"},
								},
							},
						},
					},
				},
			},
			mockStatus: models.ProcessedState{
				SuccessfulCount: 1,
				FailedCount:     0,
			},
			expectedSuccessful: 1,
			expectedFailed:     0,
			verifyParams: func(t *testing.T, params *models.UpdateParameters) {
				if len(params.Perfumes) != 1 {
					t.Fatalf("Update() params.Perfumes len = %d, want 1", len(params.Perfumes))
				}
				perfume := params.Perfumes[0]
				if len(perfume.Properties.Family) != 2 {
					t.Errorf("Update() perfume.Properties.Family len = %d, want 2", len(perfume.Properties.Family))
				}
				if len(perfume.Shops) != 1 {
					t.Errorf("Update() perfume.Shops len = %d, want 1", len(perfume.Shops))
				}
				if len(perfume.Shops[0].Variants) != 2 {
					t.Errorf("Update() perfume.Shops[0].Variants len = %d, want 2", len(perfume.Shops[0].Variants))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокируем core.Update
			mockUpdate := func(ctx context.Context, params *models.UpdateParameters) models.ProcessedState {
				if tt.verifyParams != nil {
					tt.verifyParams(t, params)
				}
				return tt.mockStatus
			}

			// Создаем тестовый сервер с мокированными функциями
			s, lis, conn, client := setupTestServer(t, nil, mockUpdate)
			defer s.Stop()
			defer lis.Close()
			defer conn.Close()

			ctx := context.Background()
			response, err := client.UpdatePerfumes(ctx, tt.request)

			if err != nil {
				t.Fatalf("UpdatePerfumes() error = %v, want nil", err)
			}
			if response == nil {
				t.Fatal("UpdatePerfumes() response = nil, want non-nil")
			}
			if response.SuccessfulCount != tt.expectedSuccessful {
				t.Errorf("UpdatePerfumes() SuccessfulCount = %d, want %d", response.SuccessfulCount, tt.expectedSuccessful)
			}
			if response.FailedCount != tt.expectedFailed {
				t.Errorf("UpdatePerfumes() FailedCount = %d, want %d", response.FailedCount, tt.expectedFailed)
			}
		})
	}
}

func TestGetPerfumes_ParametersMapping(t *testing.T) {
	tests := []struct {
		name          string
		request       *requests.GetPerfumesRequest
		expectedBrand string
		expectedName  string
		expectedSex   string
	}{
		{
			name: "all parameters set",
			request: &requests.GetPerfumesRequest{
				Brand: "Chanel",
				Name:  "No.5",
				Sex:   "female",
			},
			expectedBrand: "Chanel",
			expectedName:  "No.5",
			expectedSex:   "female",
		},
		{
			name: "only brand",
			request: &requests.GetPerfumesRequest{
				Brand: "Chanel",
			},
			expectedBrand: "Chanel",
			expectedName:  "",
			expectedSex:   "",
		},
		{
			name: "only name",
			request: &requests.GetPerfumesRequest{
				Name: "No.5",
			},
			expectedBrand: "",
			expectedName:  "No.5",
			expectedSex:   "",
		},
		{
			name: "only sex",
			request: &requests.GetPerfumesRequest{
				Sex: "male",
			},
			expectedBrand: "",
			expectedName:  "",
			expectedSex:   "male",
		},
		{
			name:          "empty request",
			request:       &requests.GetPerfumesRequest{},
			expectedBrand: "",
			expectedName:  "",
			expectedSex:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокируем core.Select для проверки параметров
			mockSelect := func(ctx context.Context, params *models.SelectParameters) ([]models.Perfume, models.ProcessedState) {
				if params.Brand != tt.expectedBrand {
					t.Errorf("SelectParameters.Brand = %q, want %q", params.Brand, tt.expectedBrand)
				}
				if params.Name != tt.expectedName {
					t.Errorf("SelectParameters.Name = %q, want %q", params.Name, tt.expectedName)
				}
				if params.Sex != tt.expectedSex {
					t.Errorf("SelectParameters.Sex = %q, want %q", params.Sex, tt.expectedSex)
				}
				return []models.Perfume{}, models.NewProcessedState()
			}

			s, lis, conn, client := setupTestServer(t, mockSelect, nil)
			defer s.Stop()
			defer lis.Close()
			defer conn.Close()

			ctx := context.Background()
			_, err := client.GetPerfumes(ctx, tt.request)
			if err != nil {
				t.Errorf("GetPerfumes() error = %v, want nil", err)
			}
		})
	}
}

func TestUpdatePerfumes_Conversion(t *testing.T) {
	imageUrl := "http://example.com/image.jpg"
	protoPerfume := &protoModels.Perfume{
		Brand:    "Chanel",
		Name:     "No.5",
		Sex:      "female",
		ImageUrl: &imageUrl,
		Properties: &protoModels.Perfume_Properties{
			PerfumeType: "Eau de Parfum",
			Family:      []string{"Floral"},
			UpperNotes:  []string{"Bergamot"},
			CoreNotes:   []string{"Rose"},
			BaseNotes:   []string{"Musk"},
		},
		Shops: []*protoModels.Perfume_ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: &imageUrl,
				Variants: []*protoModels.Perfume_ShopInfo_Variant{
					{Volume: 100, Price: 5000, Link: "http://example.com/link"},
				},
			},
		},
	}

	modelPerfume := convertPerfumeToModel(protoPerfume)

	if modelPerfume.Brand != protoPerfume.Brand {
		t.Errorf("convertPerfumeToModel() Brand = %q, want %q", modelPerfume.Brand, protoPerfume.Brand)
	}
	if modelPerfume.Name != protoPerfume.Name {
		t.Errorf("convertPerfumeToModel() Name = %q, want %q", modelPerfume.Name, protoPerfume.Name)
	}
	if modelPerfume.Sex != protoPerfume.Sex {
		t.Errorf("convertPerfumeToModel() Sex = %q, want %q", modelPerfume.Sex, protoPerfume.Sex)
	}
	if modelPerfume.ImageUrl != imageUrl {
		t.Errorf("convertPerfumeToModel() ImageUrl = %q, want %q", modelPerfume.ImageUrl, imageUrl)
	}
	if modelPerfume.Properties.Type != protoPerfume.Properties.PerfumeType {
		t.Errorf("convertPerfumeToModel() Properties.Type = %q, want %q", modelPerfume.Properties.Type, protoPerfume.Properties.PerfumeType)
	}
	if !reflect.DeepEqual(modelPerfume.Properties.Family, protoPerfume.Properties.Family) {
		t.Errorf("convertPerfumeToModel() Properties.Family = %v, want %v", modelPerfume.Properties.Family, protoPerfume.Properties.Family)
	}
	if len(modelPerfume.Shops) != len(protoPerfume.Shops) {
		t.Errorf("convertPerfumeToModel() Shops len = %d, want %d", len(modelPerfume.Shops), len(protoPerfume.Shops))
	}
}
