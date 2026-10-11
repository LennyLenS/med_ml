package node_segment

import (
	"context"
	"errors"
	"testing"

	"exam/internal/repository"
	mriRepo "exam/internal/repository/mri"
	daoLib "github.com/WantBeASleep/med_ml_lib/dao"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type emptyResultMri struct {
	mriRepo.Repository
	id     uuid.UUID
	status string
	err    error
}

func (r *emptyResultMri) UpdateMriStatus(id uuid.UUID, status string) error {
	r.id, r.status = id, status
	return r.err
}

type emptyResultDAO struct {
	repository.DAO
	mri       *emptyResultMri
	committed bool
}

func (d *emptyResultDAO) BeginTx(ctx context.Context, _ ...daoLib.TxOption) (context.Context, error) {
	return ctx, nil
}
func (d *emptyResultDAO) RollbackTx(context.Context) error               { return nil }
func (d *emptyResultDAO) CommitTx(context.Context) error                 { d.committed = true; return nil }
func (d *emptyResultDAO) NewMriQuery(context.Context) mriRepo.Repository { return d.mri }

func TestEmptyMLResultCompletesMri(t *testing.T) {
	id := uuid.New()
	dao := &emptyResultDAO{mri: &emptyResultMri{}}
	require.NoError(t, New(dao).SaveProcessedNodesWithSegments(context.Background(), id, nil))
	require.Equal(t, id, dao.mri.id)
	require.Equal(t, "completed", dao.mri.status)
	require.True(t, dao.committed)
}

func TestEmptyMLResultDoesNotCommitFailedStatus(t *testing.T) {
	dao := &emptyResultDAO{mri: &emptyResultMri{err: errors.New("storage error")}}
	require.ErrorContains(t, New(dao).SaveProcessedNodesWithSegments(context.Background(), uuid.New(), nil), "storage error")
	require.False(t, dao.committed)
}
