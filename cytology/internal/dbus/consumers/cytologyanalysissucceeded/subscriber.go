package cytologyanalysissucceeded

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/WantBeASleep/med_ml_lib/dbus"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"

	"cytology/internal/domain"
	pb "cytology/internal/generated/dbus/consume/cytologyanalysissucceeded"
	repositoryentity "cytology/internal/repository/entity"
	"cytology/internal/services"
	"cytology/internal/services/cytology_analysis"
	"cytology/internal/services/segmentation"
	"cytology/internal/services/segmentation_group"
)

type subscriber struct {
	services analysisServices
}

func New(services *services.Services) dbus.Consumer[*pb.AnalysisSucceeded] {
	return &subscriber{services: serviceAdapter{services: services}}
}

type analysisServices interface {
	GetAnalysis(context.Context, uuid.UUID) (cytology_analysis.Analysis, error)
	BeginTx(context.Context) (context.Context, error)
	RollbackTx(context.Context) error
	CommitTx(context.Context) error
	GetSegmentationGroups(context.Context, uuid.UUID) ([]domain.SegmentationGroup, error)
	DeleteSegmentationGroup(context.Context, int) error
	CreateSegmentationGroup(context.Context, segmentation_group.CreateSegmentationGroupArg) (int, error)
	CreateSegmentation(context.Context, segmentation.CreateSegmentationArg) (int, error)
	CompleteAnalysis(context.Context, uuid.UUID, []byte) error
}

type serviceAdapter struct {
	services *services.Services
}

func (a serviceAdapter) GetAnalysis(ctx context.Context, id uuid.UUID) (cytology_analysis.Analysis, error) {
	return a.services.CytologyAnalysis.Get(ctx, id)
}

func (a serviceAdapter) BeginTx(ctx context.Context) (context.Context, error) {
	return a.services.BeginTx(ctx)
}

func (a serviceAdapter) RollbackTx(ctx context.Context) error {
	return a.services.RollbackTx(ctx)
}

func (a serviceAdapter) CommitTx(ctx context.Context) error {
	return a.services.CommitTx(ctx)
}

func (a serviceAdapter) GetSegmentationGroups(ctx context.Context, id uuid.UUID) ([]domain.SegmentationGroup, error) {
	return a.services.SegmentationGroup.GetSegmentationGroupsByCytologyID(ctx, id)
}

func (a serviceAdapter) DeleteSegmentationGroup(ctx context.Context, id int) error {
	return a.services.SegmentationGroup.DeleteSegmentationGroup(ctx, id)
}

func (a serviceAdapter) CreateSegmentationGroup(ctx context.Context, arg segmentation_group.CreateSegmentationGroupArg) (int, error) {
	return a.services.SegmentationGroup.CreateSegmentationGroup(ctx, arg)
}

func (a serviceAdapter) CreateSegmentation(ctx context.Context, arg segmentation.CreateSegmentationArg) (int, error) {
	return a.services.Segmentation.CreateSegmentation(ctx, arg)
}

func (a serviceAdapter) CompleteAnalysis(ctx context.Context, id uuid.UUID, result []byte) error {
	return a.services.CytologyAnalysis.Complete(ctx, id, result)
}

func (h *subscriber) Consume(ctx context.Context, message *pb.AnalysisSucceeded) error {
	analysisID, err := uuid.Parse(message.GetAnalysisId())
	if err != nil {
		return fmt.Errorf("analysis id is not uuid: %s", message.GetAnalysisId())
	}
	if message.GetResult() == nil {
		return errors.New("result is nil")
	}

	analysis, err := h.services.GetAnalysis(ctx, analysisID)
	if err != nil {
		return fmt.Errorf("get cytology analysis: %w", err)
	}

	txCtx, err := h.services.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin analysis result transaction: %w", err)
	}
	defer h.services.RollbackTx(txCtx)

	groups, err := h.services.GetSegmentationGroups(txCtx, analysis.CytologyID)
	if err != nil && !errors.Is(err, repositoryentity.ErrNotFound) {
		return fmt.Errorf("get previous segmentation groups: %w", err)
	}
	for _, group := range groups {
		if group.IsAI {
			if err := h.services.DeleteSegmentationGroup(txCtx, group.Id); err != nil {
				return fmt.Errorf("delete previous AI segmentation group: %w", err)
			}
		}
	}

	for _, feature := range message.GetResult().GetFeatures() {
		if err := h.createFeature(txCtx, analysis.CytologyID, feature); err != nil {
			return err
		}
	}

	result, err := protojson.Marshal(message.GetResult())
	if err != nil {
		return fmt.Errorf("marshal analysis result: %w", err)
	}
	if err := h.services.CompleteAnalysis(txCtx, analysisID, result); err != nil {
		return fmt.Errorf("complete cytology analysis: %w", err)
	}
	if err := h.services.CommitTx(txCtx); err != nil {
		return fmt.Errorf("commit analysis result transaction: %w", err)
	}
	return nil
}

func (h *subscriber) createFeature(ctx context.Context, cytologyID uuid.UUID, feature *pb.Feature) error {
	if feature == nil || feature.GetProperties() == nil || feature.GetProperties().GetClassification() == nil {
		return nil
	}
	classification := feature.GetProperties().GetClassification()

	details, err := json.Marshal(map[string]any{"classification": classification})
	if err != nil {
		return fmt.Errorf("marshal segmentation group details: %w", err)
	}
	groupID, err := h.services.CreateSegmentationGroup(ctx, segmentation_group.CreateSegmentationGroupArg{
		CytologyID: cytologyID,
		SegType:    domain.SegTypeNIL,
		GroupType:  domain.GroupTypeCE,
		IsAI:       true,
		Details:    details,
	})
	if err != nil {
		return fmt.Errorf("create segmentation group: %w", err)
	}

	points := geometryPoints(feature.GetGeometry())
	if len(points) == 0 {
		return nil
	}
	_, err = h.services.CreateSegmentation(ctx, segmentation.CreateSegmentationArg{
		SegmentationGroupID: groupID,
		Points:              points,
	})
	if err != nil {
		return fmt.Errorf("create segmentation: %w", err)
	}
	return nil
}

func geometryPoints(geometry *pb.Geometry) []domain.SegmentationPoint {
	if geometry == nil {
		return nil
	}
	var protoPoints []*pb.Point
	switch value := geometry.GetGeometryType().(type) {
	case *pb.Geometry_Point:
		protoPoints = []*pb.Point{value.Point}
	case *pb.Geometry_Polygon:
		if len(value.Polygon.GetRings()) > 0 {
			protoPoints = value.Polygon.GetRings()[0].GetPoints()
		}
	}
	points := make([]domain.SegmentationPoint, 0, len(protoPoints))
	for _, point := range protoPoints {
		if point != nil {
			points = append(points, domain.SegmentationPoint{X: int(point.GetX()), Y: int(point.GetY())})
		}
	}
	return points
}
