package node

import (
	"context"

	"exam/internal/domain"
	nodeEntity "exam/internal/repository/node/entity"

	"github.com/google/uuid"
)

func (s *service) GetNodesByMriID(ctx context.Context, id uuid.UUID) ([]domain.Node, error) {
	nodes, err := s.dao.NewNodeQuery(ctx).GetNodesByMriID(id)
	if err != nil {
		return nil, err
	}

	return nodeEntity.Node{}.SliceToDomain(nodes), nil
}
