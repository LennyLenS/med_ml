package exam

import (
	"context"

	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

type Adapter interface {
	// DEVICE
	CreateDevice(ctx context.Context, name string) (int, error)
	GetDeviceList(ctx context.Context) ([]domain.Device, error)
	// MRI
	CreateMri(ctx context.Context, in CreateMriIn) (uuid.UUID, error)
	GetMriById(ctx context.Context, id uuid.UUID) (domain.Mri, error)
	GetMrisByExternalId(ctx context.Context, id uuid.UUID) ([]domain.Mri, error)
	GetMrisByAuthor(ctx context.Context, id uuid.UUID) ([]domain.Mri, error)
	GetEchographicByMriId(ctx context.Context, id uuid.UUID) (domain.Echographic, error)
	UpdateMri(ctx context.Context, in UpdateMriIn) (domain.Mri, error)
	UpdateEchographic(ctx context.Context, in domain.Echographic) (domain.Echographic, error)
	DeleteMri(ctx context.Context, id uuid.UUID) error
	// KT
	CreateKt(ctx context.Context, in CreateKtIn) (uuid.UUID, error)
	GetKtById(ctx context.Context, id uuid.UUID) (domain.KT, error)
	GetKtsByAuthor(ctx context.Context, id uuid.UUID) ([]domain.KT, error)
	UpdateKt(ctx context.Context, in UpdateKtIn) (domain.KT, error)
	DeleteKt(ctx context.Context, id uuid.UUID) error
	// IMAGE
	GetImagesByMriId(ctx context.Context, id uuid.UUID) ([]domain.Image, error)
	// NODE
	GetNodesByMriId(ctx context.Context, id uuid.UUID) ([]domain.Node, error)
	UpdateNode(ctx context.Context, in UpdateNodeIn) (domain.Node, error)
	// SEGMENT
	CreateSegment(ctx context.Context, in CreateSegmentIn) (uuid.UUID, error)
	GetSegmentsByNodeId(ctx context.Context, id uuid.UUID) ([]domain.Segment, error)
	UpdateSegment(ctx context.Context, in UpdateSegmentIn) (domain.Segment, error)
	// доменные области слишком сильно пересекаются, вынесено в одну надобласть
	// NODE-SEGMENT
	CreateNodeWithSegments(ctx context.Context, in CreateNodeWithSegmentsIn) (uuid.UUID, []uuid.UUID, error)
	GetNodesWithSegmentsByImageId(ctx context.Context, id uuid.UUID) ([]domain.Node, []domain.Segment, error)
	DeleteNode(ctx context.Context, id uuid.UUID) error
	DeleteSegment(ctx context.Context, id uuid.UUID) error
}

type adapter struct {
	client pb.ExamSrvClient
}

func NewAdapter(client pb.ExamSrvClient) Adapter {
	return &adapter{client: client}
}


