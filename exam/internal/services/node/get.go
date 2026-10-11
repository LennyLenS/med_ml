package node

import (
	"context"
	"errors"

	"exam/internal/domain"
	daoEntity "exam/internal/repository/entity"
	nodeEntity "exam/internal/repository/node/entity"

	"github.com/google/uuid"
)

func (s *service) GetNodesByMriID(ctx context.Context, id uuid.UUID) ([]domain.Node, error) {
	nodes, err := s.dao.NewNodeQuery(ctx).GetNodesByMriID(id)
	if errors.Is(err, daoEntity.ErrNotFound) {
		exists, checkErr := s.dao.NewMriQuery(ctx).CheckExist(id)
		if checkErr != nil {
			return nil, checkErr
		}
		if exists {
			return []domain.Node{}, nil
		}
	}
	if err != nil {
		return nil, err
	}

	return nodeEntity.Node{}.SliceToDomain(nodes), nil
}
