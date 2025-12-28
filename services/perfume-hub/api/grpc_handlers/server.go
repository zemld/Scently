package grpc_handlers

import (
	"context"
	"errors"

	pb "github.com/zemld/Scently/generated/proto/perfume-hub"
	protoModels "github.com/zemld/Scently/generated/proto/perfume-hub/models"
	"github.com/zemld/Scently/generated/proto/perfume-hub/requests"
	"github.com/zemld/Scently/perfume-hub/internal/db/core"
	"github.com/zemld/Scently/perfume-hub/internal/models"
)

type PerfumeStorageServer struct {
	pb.UnimplementedPerfumeStorageServer
	selectFunc core.SelectFunc
	updateFunc core.UpdateFunc
}

func NewPerfumeStorageServer(selectFunc core.SelectFunc, updateFunc core.UpdateFunc) *PerfumeStorageServer {
	return &PerfumeStorageServer{
		selectFunc: selectFunc,
		updateFunc: updateFunc,
	}
}

func (s *PerfumeStorageServer) GetPerfumes(
	ctx context.Context,
	req *requests.GetPerfumesRequest,
) (*requests.GetPerfumesResponse, error) {
	brand := req.GetBrand()
	name := req.GetName()
	sex := req.GetSex()

	params := models.NewSelectParameters().WithBrand(brand).WithName(name).WithSex(sex)

	perfumes, status := s.selectFunc(ctx, params)
	if status.Error != nil {
		return nil, errors.New(status.Error.Error())
	}

	responsePerfumes := make([]*protoModels.Perfume, len(perfumes))
	for i := range perfumes {
		responsePerfumes[i] = convertPerfumeToProto(perfumes[i])
	}

	return &requests.GetPerfumesResponse{Perfumes: responsePerfumes}, nil
}

func (s *PerfumeStorageServer) UpdatePerfumes(
	ctx context.Context,
	req *requests.UpdatePerfumesRequest,
) (*requests.UpdatePerfumesResponse, error) {
	perfumesToUpdate := make([]models.Perfume, len(req.GetPerfumes()))
	for i := range req.GetPerfumes() {
		if req.Perfumes[i] == nil {
			continue
		}
		perfumesToUpdate[i] = convertPerfumeToModel(req.Perfumes[i])
	}
	updateStatus := s.updateFunc(ctx, models.NewUpdateParameters().WithPerfumes(perfumesToUpdate))
	return &requests.UpdatePerfumesResponse{
		SuccessfulCount: int32(updateStatus.SuccessfulCount),
		FailedCount:     int32(updateStatus.FailedCount),
	}, nil
}
