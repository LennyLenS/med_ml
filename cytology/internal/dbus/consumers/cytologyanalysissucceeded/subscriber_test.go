package cytologyanalysissucceeded

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"cytology/internal/domain"
	pb "cytology/internal/generated/dbus/consume/cytologyanalysissucceeded"
	"cytology/internal/services/cytology_analysis"
	"cytology/internal/services/segmentation"
	"cytology/internal/services/segmentation_group"
)

type analysisServicesSpy struct {
	analysis          cytology_analysis.Analysis
	groups            []domain.SegmentationGroup
	deletedGroupIDs   []int
	createdGroups     []segmentation_group.CreateSegmentationGroupArg
	createdSegments   []segmentation.CreateSegmentationArg
	completedID       uuid.UUID
	completedResult   []byte
	transactionCommit bool
}

func (s *analysisServicesSpy) GetAnalysis(context.Context, uuid.UUID) (cytology_analysis.Analysis, error) {
	return s.analysis, nil
}

func (s *analysisServicesSpy) BeginTx(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (s *analysisServicesSpy) RollbackTx(context.Context) error {
	return nil
}

func (s *analysisServicesSpy) CommitTx(context.Context) error {
	s.transactionCommit = true
	return nil
}

func (s *analysisServicesSpy) GetSegmentationGroups(context.Context, uuid.UUID) ([]domain.SegmentationGroup, error) {
	return s.groups, nil
}

func (s *analysisServicesSpy) DeleteSegmentationGroup(_ context.Context, id int) error {
	s.deletedGroupIDs = append(s.deletedGroupIDs, id)
	return nil
}

func (s *analysisServicesSpy) CreateSegmentationGroup(_ context.Context, arg segmentation_group.CreateSegmentationGroupArg) (int, error) {
	s.createdGroups = append(s.createdGroups, arg)
	return 101, nil
}

func (s *analysisServicesSpy) CreateSegmentation(_ context.Context, arg segmentation.CreateSegmentationArg) (int, error) {
	s.createdSegments = append(s.createdSegments, arg)
	return 202, nil
}

func (s *analysisServicesSpy) CompleteAnalysis(_ context.Context, id uuid.UUID, result []byte) error {
	s.completedID = id
	s.completedResult = result
	return nil
}

func TestConsume_ProjectsLatestResultForReadHandlers(t *testing.T) {
	analysisID := uuid.New()
	cytologyID := uuid.New()
	services := &analysisServicesSpy{
		analysis: cytology_analysis.Analysis{ID: analysisID, CytologyID: cytologyID},
		groups: []domain.SegmentationGroup{
			{Id: 10, IsAI: true},
			{Id: 11, IsAI: false},
		},
	}
	subscriber := &subscriber{services: services}
	comment := "Проверено моделью"

	err := subscriber.Consume(context.Background(), &pb.AnalysisSucceeded{
		AnalysisId: analysisID.String(),
		Result: &pb.Result{
			Features: []*pb.Feature{{
				Geometry: &pb.Geometry{GeometryType: &pb.Geometry_Polygon{Polygon: &pb.Polygon{Rings: []*pb.Ring{{Points: []*pb.Point{
					{X: 10.4, Y: 20.6},
					{X: 30.1, Y: 40.9},
				}}}}}},
				Properties: &pb.Properties{Classification: &pb.Classification{
					Name:  "Скопление папиллярное",
					Color: &pb.Color{R: 1, G: 2, B: 3},
				}},
			}},
			Conclusion: &pb.Conclusion{
				Category:   5,
				ShortLabel: "V",
				CommentRu:  &comment,
				Counts:     &pb.Counts{ClusterTotal: 1},
			},
		},
	})

	require.NoError(t, err)
	require.Equal(t, []int{10}, services.deletedGroupIDs)
	require.Len(t, services.createdGroups, 1)
	require.Equal(t, cytologyID, services.createdGroups[0].CytologyID)
	require.Equal(t, domain.SegTypeCPS, services.createdGroups[0].SegType)
	require.Equal(t, domain.GroupTypeCL, services.createdGroups[0].GroupType)
	require.True(t, services.createdGroups[0].IsAI)
	require.JSONEq(t, `{"classification":{"name":"Скопление папиллярное","color":{"r":1,"g":2,"b":3}}}`, string(services.createdGroups[0].Details))
	require.Len(t, services.createdSegments, 1)
	require.Equal(t, 101, services.createdSegments[0].SegmentationGroupID)
	require.Equal(t, []domain.SegmentationPoint{{X: 10, Y: 20}, {X: 30, Y: 40}}, services.createdSegments[0].Points)
	require.Equal(t, analysisID, services.completedID)
	require.JSONEq(t, `{"features":[{"geometry":{"polygon":{"rings":[{"points":[{"x":10.4,"y":20.6},{"x":30.1,"y":40.9}]}]}},"properties":{"classification":{"name":"Скопление папиллярное","color":{"r":1,"g":2,"b":3}}}}],"conclusion":{"category":5,"shortLabel":"V","commentRu":"Проверено моделью","counts":{"clusterTotal":"1"}}}`, string(services.completedResult))
	require.True(t, services.transactionCommit)
}
