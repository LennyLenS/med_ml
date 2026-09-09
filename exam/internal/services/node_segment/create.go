package node_segment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"exam/internal/domain"
	nodeEntity "exam/internal/repository/node/entity"
	segmentEntity "exam/internal/repository/segment/entity"
)

const (
	avgSegmentPerNode = 3
)

func (s *service) createNodesWithSegments(
	ctx context.Context,
	mriID uuid.UUID,
	ai bool,
	arg []CreateNodesWithSegmentsArg,
	opts ...CreateNodesWithSegmentsOption,
) ([]CreateNodesWithSegmentsID, error) {
	opt := &createNodesWithSegmentsOption{}
	for _, o := range opts {
		o(opt)
	}

	nodes, segments, ids := s.createDomainNodeSegmentsFromArgs(mriID, ai, arg)

	if opt.setNodesValidation != nil {
		for i := range nodes {
			nodes[i].Validation = opt.setNodesValidation
		}
	}

	ctx, err := s.dao.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = s.dao.RollbackTx(ctx) }()

	nodeQuery := s.dao.NewNodeQuery(ctx)
	segmentQuery := s.dao.NewSegmentQuery(ctx)

	if err := nodeQuery.InsertNodes(nodeEntity.Node{}.SliceFromDomain(nodes)...); err != nil {
		return nil, fmt.Errorf("insert nodes: %w", err)
	}

	if err := segmentQuery.InsertSegments(segmentEntity.Segment{}.SliceFromDomain(segments)...); err != nil {
		return nil, fmt.Errorf("insert segments: %w", err)
	}

	if opt.newMriStatus != nil {
		if err := s.dao.NewMriQuery(ctx).UpdateMriStatus(mriID, opt.newMriStatus.String()); err != nil {
			return nil, fmt.Errorf("update mri status: %w", err)
		}
	}

	if err := s.dao.CommitTx(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return ids, nil
}

func (s *service) SaveProcessedNodesWithSegments(
	ctx context.Context,
	mriID uuid.UUID,
	arg []CreateNodesWithSegmentsArg,
) error {
	_, err := s.createNodesWithSegments(
		ctx,
		mriID,
		true,
		arg,
		WithNewMriStatus(domain.ExamStatusCompleted),
		WithSetNodesValidation(domain.NodeValidationNull),
	)

	return err
}

func (s *service) CreateManualNodesWithSegments(
	ctx context.Context,
	mriID uuid.UUID,
	arg []CreateNodesWithSegmentsArg,
) ([]CreateNodesWithSegmentsID, error) {
	return s.createNodesWithSegments(ctx, mriID, false, arg)
}

func (s *service) createDomainNodeSegmentsFromArgs(
	mriID uuid.UUID,
	ai bool,
	arg []CreateNodesWithSegmentsArg,
) (
	[]domain.Node,
	[]domain.Segment,
	[]CreateNodesWithSegmentsID,
) {
	ids := make([]CreateNodesWithSegmentsID, 0, len(arg))
	nodes := make([]domain.Node, 0, len(arg))
	segments := make([]domain.Segment, 0, avgSegmentPerNode*len(arg))

	for _, NodeAndSeg := range arg {
		nodeID := uuid.New()
		nodes = append(nodes, domain.Node{
			Id:          nodeID,
			Ai:          ai,
			MriID:       mriID,
			Knosp012:    NodeAndSeg.Node.Knosp012,
			Knosp3:      NodeAndSeg.Node.Knosp3,
			Knosp4:      NodeAndSeg.Node.Knosp4,
			Description: NodeAndSeg.Node.Description,
		})

		id := CreateNodesWithSegmentsID{NodeID: nodeID}

		for _, segment := range NodeAndSeg.Segments {
			segmentID := uuid.New()
			segments = append(segments, domain.Segment{
				Id:       segmentID,
				ImageID:  segment.ImageID,
				NodeID:   nodeID,
				Contor:   segment.Contor,
				Ai:       ai,
				Knosp012: segment.Knosp012,
				Knosp3:   segment.Knosp3,
				Knosp4:   segment.Knosp4,
			})

			id.SegmentsID = append(id.SegmentsID, segmentID)
		}

		ids = append(ids, id)
	}

	return nodes, segments, ids
}
