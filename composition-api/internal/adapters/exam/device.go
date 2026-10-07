package exam

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"composition-api/internal/adapters/exam/mappers"
	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

func (a *adapter) CreateDevice(ctx context.Context, name string) (int, error) {
	res, err := a.client.CreateDevice(ctx, &pb.CreateDeviceIn{Name: name})
	if err != nil {
		return 0, err
	}

	return int(res.Id), nil
}

func (a *adapter) GetDeviceList(ctx context.Context) ([]domain.Device, error) {
	res, err := a.client.GetDeviceList(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return mappers.Device{}.SliceDomain(res.Devices), nil
}


